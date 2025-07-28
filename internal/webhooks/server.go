/*
 * GitHubber - Webhook Server Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Real-time webhook integration and event handling
 */

package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/ritankarsaha/git-tool/internal/plugins"
	"github.com/ritankarsaha/git-tool/internal/providers"
)

// WebhookServer handles incoming webhooks from various providers
type WebhookServer struct {
	port          int
	router        *mux.Router
	server        *http.Server
	handlers      map[string]WebhookHandler
	eventQueue    chan *WebhookEvent
	subscribers   map[string][]EventSubscriber
	pluginManager plugins.PluginManager
	mu            sync.RWMutex
	running       bool
}

// WebhookHandler processes webhooks for a specific provider
type WebhookHandler interface {
	HandleWebhook(ctx context.Context, event *WebhookEvent) error
	ValidateSignature(payload []byte, signature string, secret string) bool
	ParseEvent(headers http.Header, body []byte) (*WebhookEvent, error)
	GetSupportedEvents() []string
}

// WebhookEvent represents a normalized webhook event
type WebhookEvent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Action      string                 `json:"action"`
	Provider    string                 `json:"provider"`
	Repository  *RepositoryInfo        `json:"repository"`
	Sender      *UserInfo              `json:"sender"`
	Data        map[string]interface{} `json:"data"`
	Headers     map[string]string      `json:"headers"`
	Timestamp   time.Time              `json:"timestamp"`
	Signature   string                 `json:"signature"`
	DeliveryID  string                 `json:"delivery_id"`
}

// RepositoryInfo contains repository information from webhook
type RepositoryInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    string `json:"owner"`
	URL      string `json:"url"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
}

// UserInfo contains user information from webhook
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// EventSubscriber processes webhook events
type EventSubscriber interface {
	ProcessEvent(ctx context.Context, event *WebhookEvent) error
	GetEventTypes() []string
}

// WebhookConfig represents webhook server configuration
type WebhookConfig struct {
	Port          int               `json:"port"`
	Path          string            `json:"path"`
	Secret        string            `json:"secret"`
	EnableLogging bool              `json:"enable_logging"`
	QueueSize     int               `json:"queue_size"`
	Workers       int               `json:"workers"`
	Timeout       time.Duration     `json:"timeout"`
	TLS           *TLSConfig        `json:"tls,omitempty"`
	Providers     map[string]*ProviderWebhookConfig `json:"providers"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

// ProviderWebhookConfig contains provider-specific webhook configuration
type ProviderWebhookConfig struct {
	Secret    string            `json:"secret"`
	Events    []string          `json:"events"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
}

// NewWebhookServer creates a new webhook server
func NewWebhookServer(config *WebhookConfig, pluginManager plugins.PluginManager) *WebhookServer {
	if config.QueueSize == 0 {
		config.QueueSize = 1000
	}
	if config.Workers == 0 {
		config.Workers = 10
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	server := &WebhookServer{
		port:          config.Port,
		router:        mux.NewRouter(),
		eventQueue:    make(chan *WebhookEvent, config.QueueSize),
		handlers:      make(map[string]WebhookHandler),
		subscribers:   make(map[string][]EventSubscriber),
		pluginManager: pluginManager,
	}

	// Register built-in handlers
	server.RegisterHandler("github", NewGitHubWebhookHandler())
	server.RegisterHandler("gitlab", NewGitLabWebhookHandler())
	server.RegisterHandler("bitbucket", NewBitbucketWebhookHandler())

	server.setupRoutes()
	server.startWorkers(config.Workers)

	return server
}

// RegisterHandler registers a webhook handler for a provider
func (ws *WebhookServer) RegisterHandler(provider string, handler WebhookHandler) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.handlers[provider] = handler
}

// Subscribe adds an event subscriber
func (ws *WebhookServer) Subscribe(eventType string, subscriber EventSubscriber) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	
	if ws.subscribers[eventType] == nil {
		ws.subscribers[eventType] = make([]EventSubscriber, 0)
	}
	ws.subscribers[eventType] = append(ws.subscribers[eventType], subscriber)
}

// Start starts the webhook server
func (ws *WebhookServer) Start() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.running {
		return fmt.Errorf("webhook server is already running")
	}

	ws.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", ws.port),
		Handler: ws.router,
	}

	ws.running = true

	go func() {
		log.Printf("Starting webhook server on port %d", ws.port)
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Webhook server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the webhook server
func (ws *WebhookServer) Stop() error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.running {
		return fmt.Errorf("webhook server is not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := ws.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown webhook server: %w", err)
	}

	close(ws.eventQueue)
	ws.running = false

	return nil
}

// setupRoutes configures the HTTP routes
func (ws *WebhookServer) setupRoutes() {
	// Generic webhook endpoint
	ws.router.HandleFunc("/webhook/{provider}", ws.handleWebhook).Methods("POST")
	
	// Provider-specific endpoints
	ws.router.HandleFunc("/github", ws.handleGitHubWebhook).Methods("POST")
	ws.router.HandleFunc("/gitlab", ws.handleGitLabWebhook).Methods("POST")
	ws.router.HandleFunc("/bitbucket", ws.handleBitbucketWebhook).Methods("POST")
	
	// Health check endpoint
	ws.router.HandleFunc("/health", ws.handleHealth).Methods("GET")
	
	// Metrics endpoint
	ws.router.HandleFunc("/metrics", ws.handleMetrics).Methods("GET")
}

// handleWebhook handles generic webhook requests
func (ws *WebhookServer) handleWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	ws.processWebhook(w, r, provider)
}

// handleGitHubWebhook handles GitHub-specific webhooks
func (ws *WebhookServer) handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	ws.processWebhook(w, r, "github")
}

// handleGitLabWebhook handles GitLab-specific webhooks
func (ws *WebhookServer) handleGitLabWebhook(w http.ResponseWriter, r *http.Request) {
	ws.processWebhook(w, r, "gitlab")
}

// handleBitbucketWebhook handles Bitbucket-specific webhooks
func (ws *WebhookServer) handleBitbucketWebhook(w http.ResponseWriter, r *http.Request) {
	ws.processWebhook(w, r, "bitbucket")
}

// processWebhook processes a webhook request
func (ws *WebhookServer) processWebhook(w http.ResponseWriter, r *http.Request, provider string) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get handler for provider
	ws.mu.RLock()
	handler, exists := ws.handlers[provider]
	ws.mu.RUnlock()

	if !exists {
		http.Error(w, fmt.Sprintf("No handler for provider: %s", provider), http.StatusBadRequest)
		return
	}

	// Parse event
	event, err := handler.ParseEvent(r.Header, body)
	if err != nil {
		log.Printf("Failed to parse webhook event: %v", err)
		http.Error(w, "Failed to parse webhook event", http.StatusBadRequest)
		return
	}

	// Validate signature if provided
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		signature = r.Header.Get("X-GitLab-Token")
	}
	if signature == "" {
		signature = r.Header.Get("X-Hook-UUID")
	}

	if signature != "" {
		// This would typically use a configured secret
		if !handler.ValidateSignature(body, signature, "") {
			log.Printf("Invalid webhook signature for provider: %s", provider)
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Set additional event data
	event.Provider = provider
	event.Timestamp = time.Now()
	event.DeliveryID = r.Header.Get("X-GitHub-Delivery")
	if event.DeliveryID == "" {
		event.DeliveryID = r.Header.Get("X-GitLab-Event-UUID")
	}

	// Queue event for processing
	select {
	case ws.eventQueue <- event:
		log.Printf("Queued webhook event: %s/%s from %s", event.Type, event.Action, provider)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		log.Printf("Webhook event queue is full")
		http.Error(w, "Event queue is full", http.StatusServiceUnavailable)
	}
}

// handleHealth handles health check requests
func (ws *WebhookServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"queue_size": len(ws.eventQueue),
		"running":   ws.running,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// handleMetrics handles metrics requests
func (ws *WebhookServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"queue_size":     len(ws.eventQueue),
		"queue_capacity": cap(ws.eventQueue),
		"handlers":       len(ws.handlers),
		"subscribers":    len(ws.subscribers),
		"uptime":         time.Since(time.Now()), // Would track actual uptime
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// startWorkers starts background workers to process events
func (ws *WebhookServer) startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go ws.worker(i)
	}
}

// worker processes events from the queue
func (ws *WebhookServer) worker(id int) {
	log.Printf("Starting webhook worker %d", id)

	for event := range ws.eventQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		
		if err := ws.processEvent(ctx, event); err != nil {
			log.Printf("Worker %d failed to process event %s: %v", id, event.ID, err)
		}
		
		cancel()
	}

	log.Printf("Webhook worker %d stopped", id)
}

// processEvent processes a webhook event
func (ws *WebhookServer) processEvent(ctx context.Context, event *WebhookEvent) error {
	// Notify plugin manager
	if ws.pluginManager != nil {
		pluginEvent := &plugins.WebhookEvent{
			ID:        event.ID,
			Type:      event.Type,
			Source:    event.Provider,
			Timestamp: event.Timestamp,
			Data:      event.Data,
			Headers:   event.Headers,
		}

		if err := ws.pluginManager.HandleWebhook(pluginEvent); err != nil {
			log.Printf("Plugin manager failed to handle webhook: %v", err)
		}
	}

	// Notify subscribers
	ws.mu.RLock()
	subscribers := ws.subscribers[event.Type]
	if subscribers == nil {
		subscribers = ws.subscribers["*"] // Wildcard subscribers
	}
	ws.mu.RUnlock()

	for _, subscriber := range subscribers {
		if err := subscriber.ProcessEvent(ctx, event); err != nil {
			log.Printf("Subscriber failed to process event: %v", err)
		}
	}

	return nil
}

// GitHub Webhook Handler
type GitHubWebhookHandler struct{}

func NewGitHubWebhookHandler() *GitHubWebhookHandler {
	return &GitHubWebhookHandler{}
}

func (gh *GitHubWebhookHandler) HandleWebhook(ctx context.Context, event *WebhookEvent) error {
	// Process GitHub-specific logic
	return nil
}

func (gh *GitHubWebhookHandler) ValidateSignature(payload []byte, signature string, secret string) bool {
	if secret == "" {
		return true // Skip validation if no secret configured
	}

	// GitHub uses HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (gh *GitHubWebhookHandler) ParseEvent(headers http.Header, body []byte) (*WebhookEvent, error) {
	eventType := headers.Get("X-GitHub-Event")
	deliveryID := headers.Get("X-GitHub-Delivery")

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	event := &WebhookEvent{
		ID:         deliveryID,
		Type:       eventType,
		Data:       payload,
		Headers:    convertHeaders(headers),
		Timestamp:  time.Now(),
		DeliveryID: deliveryID,
	}

	// Extract action if present
	if action, ok := payload["action"].(string); ok {
		event.Action = action
	}

	// Extract repository information
	if repo, ok := payload["repository"].(map[string]interface{}); ok {
		event.Repository = parseGitHubRepository(repo)
	}

	// Extract sender information
	if sender, ok := payload["sender"].(map[string]interface{}); ok {
		event.Sender = parseGitHubUser(sender)
	}

	return event, nil
}

func (gh *GitHubWebhookHandler) GetSupportedEvents() []string {
	return []string{
		"push", "pull_request", "issues", "issue_comment",
		"pull_request_review", "pull_request_review_comment",
		"create", "delete", "fork", "watch", "star",
		"release", "deployment", "deployment_status",
		"check_run", "check_suite", "workflow_run",
	}
}

// GitLab Webhook Handler
type GitLabWebhookHandler struct{}

func NewGitLabWebhookHandler() *GitLabWebhookHandler {
	return &GitLabWebhookHandler{}
}

func (gl *GitLabWebhookHandler) HandleWebhook(ctx context.Context, event *WebhookEvent) error {
	return nil
}

func (gl *GitLabWebhookHandler) ValidateSignature(payload []byte, signature string, secret string) bool {
	return signature == secret // GitLab uses token-based authentication
}

func (gl *GitLabWebhookHandler) ParseEvent(headers http.Header, body []byte) (*WebhookEvent, error) {
	eventType := headers.Get("X-GitLab-Event")

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	event := &WebhookEvent{
		ID:        headers.Get("X-GitLab-Event-UUID"),
		Type:      eventType,
		Data:      payload,
		Headers:   convertHeaders(headers),
		Timestamp: time.Now(),
	}

	// GitLab has different event structure
	if project, ok := payload["project"].(map[string]interface{}); ok {
		event.Repository = parseGitLabProject(project)
	}

	if user, ok := payload["user"].(map[string]interface{}); ok {
		event.Sender = parseGitLabUser(user)
	}

	return event, nil
}

func (gl *GitLabWebhookHandler) GetSupportedEvents() []string {
	return []string{
		"Push Hook", "Tag Push Hook", "Issue Hook", "Note Hook",
		"Merge Request Hook", "Wiki Page Hook", "Deployment Hook",
		"Job Hook", "Pipeline Hook", "Build Hook",
	}
}

// Bitbucket Webhook Handler
type BitbucketWebhookHandler struct{}

func NewBitbucketWebhookHandler() *BitbucketWebhookHandler {
	return &BitbucketWebhookHandler{}
}

func (bb *BitbucketWebhookHandler) HandleWebhook(ctx context.Context, event *WebhookEvent) error {
	return nil
}

func (bb *BitbucketWebhookHandler) ValidateSignature(payload []byte, signature string, secret string) bool {
	if secret == "" {
		return true
	}

	// Bitbucket uses HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (bb *BitbucketWebhookHandler) ParseEvent(headers http.Header, body []byte) (*WebhookEvent, error) {
	eventType := headers.Get("X-Event-Key")

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	event := &WebhookEvent{
		ID:        headers.Get("X-Hook-UUID"),
		Type:      eventType,
		Data:      payload,
		Headers:   convertHeaders(headers),
		Timestamp: time.Now(),
	}

	// Bitbucket structure
	if repository, ok := payload["repository"].(map[string]interface{}); ok {
		event.Repository = parseBitbucketRepository(repository)
	}

	if actor, ok := payload["actor"].(map[string]interface{}); ok {
		event.Sender = parseBitbucketUser(actor)
	}

	return event, nil
}

func (bb *BitbucketWebhookHandler) GetSupportedEvents() []string {
	return []string{
		"repo:push", "pullrequest:created", "pullrequest:updated",
		"pullrequest:approved", "pullrequest:unapproved",
		"pullrequest:fulfilled", "pullrequest:rejected",
		"issue:created", "issue:updated", "issue:comment_created",
	}
}

// Helper functions
func convertHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

func parseGitHubRepository(repo map[string]interface{}) *RepositoryInfo {
	info := &RepositoryInfo{}
	
	if id, ok := repo["id"].(float64); ok {
		info.ID = int64(id)
	}
	if name, ok := repo["name"].(string); ok {
		info.Name = name
	}
	if fullName, ok := repo["full_name"].(string); ok {
		info.FullName = fullName
		parts := strings.Split(fullName, "/")
		if len(parts) == 2 {
			info.Owner = parts[0]
		}
	}
	if url, ok := repo["html_url"].(string); ok {
		info.URL = url
	}
	if cloneURL, ok := repo["clone_url"].(string); ok {
		info.CloneURL = cloneURL
	}
	if private, ok := repo["private"].(bool); ok {
		info.Private = private
	}
	
	return info
}

func parseGitHubUser(user map[string]interface{}) *UserInfo {
	info := &UserInfo{}
	
	if id, ok := user["id"].(float64); ok {
		info.ID = int64(id)
	}
	if login, ok := user["login"].(string); ok {
		info.Username = login
	}
	if name, ok := user["name"].(string); ok {
		info.Name = name
	}
	if email, ok := user["email"].(string); ok {
		info.Email = email
	}
	if avatar, ok := user["avatar_url"].(string); ok {
		info.Avatar = avatar
	}
	
	return info
}

func parseGitLabProject(project map[string]interface{}) *RepositoryInfo {
	info := &RepositoryInfo{}
	
	if id, ok := project["id"].(float64); ok {
		info.ID = int64(id)
	}
	if name, ok := project["name"].(string); ok {
		info.Name = name
	}
	if pathWithNamespace, ok := project["path_with_namespace"].(string); ok {
		info.FullName = pathWithNamespace
		parts := strings.Split(pathWithNamespace, "/")
		if len(parts) >= 2 {
			info.Owner = strings.Join(parts[:len(parts)-1], "/")
		}
	}
	if url, ok := project["web_url"].(string); ok {
		info.URL = url
	}
	if cloneURL, ok := project["git_http_url"].(string); ok {
		info.CloneURL = cloneURL
	}
	
	return info
}

func parseGitLabUser(user map[string]interface{}) *UserInfo {
	info := &UserInfo{}
	
	if id, ok := user["id"].(float64); ok {
		info.ID = int64(id)
	}
	if username, ok := user["username"].(string); ok {
		info.Username = username
	}
	if name, ok := user["name"].(string); ok {
		info.Name = name
	}
	if email, ok := user["email"].(string); ok {
		info.Email = email
	}
	if avatar, ok := user["avatar_url"].(string); ok {
		info.Avatar = avatar
	}
	
	return info
}

func parseBitbucketRepository(repository map[string]interface{}) *RepositoryInfo {
	info := &RepositoryInfo{}
	
	if name, ok := repository["name"].(string); ok {
		info.Name = name
	}
	if fullName, ok := repository["full_name"].(string); ok {
		info.FullName = fullName
		parts := strings.Split(fullName, "/")
		if len(parts) == 2 {
			info.Owner = parts[0]
		}
	}
	if links, ok := repository["links"].(map[string]interface{}); ok {
		if html, ok := links["html"].(map[string]interface{}); ok {
			if href, ok := html["href"].(string); ok {
				info.URL = href
			}
		}
		if clone, ok := links["clone"].([]interface{}); ok {
			for _, c := range clone {
				if cloneLink, ok := c.(map[string]interface{}); ok {
					if name, ok := cloneLink["name"].(string); ok && name == "https" {
						if href, ok := cloneLink["href"].(string); ok {
							info.CloneURL = href
						}
					}
				}
			}
		}
	}
	if isPrivate, ok := repository["is_private"].(bool); ok {
		info.Private = isPrivate
	}
	
	return info
}

func parseBitbucketUser(user map[string]interface{}) *UserInfo {
	info := &UserInfo{}
	
	if username, ok := user["username"].(string); ok {
		info.Username = username
	}
	if displayName, ok := user["display_name"].(string); ok {
		info.Name = displayName
	}
	if uuid, ok := user["uuid"].(string); ok {
		// Convert UUID to numeric ID (simplified)
		if id, err := strconv.ParseInt(strings.ReplaceAll(uuid, "-", "")[:8], 16, 64); err == nil {
			info.ID = id
		}
	}
	if links, ok := user["links"].(map[string]interface{}); ok {
		if avatar, ok := links["avatar"].(map[string]interface{}); ok {
			if href, ok := avatar["href"].(string); ok {
				info.Avatar = href
			}
		}
	}
	
	return info
}