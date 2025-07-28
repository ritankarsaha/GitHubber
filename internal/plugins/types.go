/*
 * GitHubber - Plugin System Types and Interfaces
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Extensible plugin architecture for GitHubber
 */

package plugins

import (
	"context"
	"time"

	"github.com/ritankarsaha/git-tool/internal/providers"
)

// PluginType represents the type of plugin
type PluginType string

const (
	PluginTypeProvider   PluginType = "provider"
	PluginTypeCommand    PluginType = "command"
	PluginTypeWebhook    PluginType = "webhook"
	PluginTypeCI         PluginType = "ci"
	PluginTypeNotifier   PluginType = "notifier"
	PluginTypeIntegration PluginType = "integration"
)

// Plugin represents the base plugin interface
type Plugin interface {
	// Metadata
	Name() string
	Version() string
	Type() PluginType
	Description() string
	Author() string
	
	// Lifecycle
	Initialize(config *PluginConfig) error
	Start() error
	Stop() error
	IsRunning() bool
	
	// Configuration
	GetConfigSchema() *ConfigSchema
	Validate(config *PluginConfig) error
}

// CommandPlugin extends Plugin for command-based plugins
type CommandPlugin interface {
	Plugin
	GetCommands() []CommandDefinition
	ExecuteCommand(ctx context.Context, cmd string, args []string) error
}

// WebhookPlugin extends Plugin for webhook handling
type WebhookPlugin interface {
	Plugin
	HandleWebhook(ctx context.Context, event *WebhookEvent) error
	GetSupportedEvents() []string
}

// CIPlugin extends Plugin for CI/CD integration
type CIPlugin interface {
	Plugin
	TriggerBuild(ctx context.Context, config *BuildConfig) (*BuildResult, error)
	GetBuildStatus(ctx context.Context, buildID string) (*BuildStatus, error)
	CancelBuild(ctx context.Context, buildID string) error
}

// NotifierPlugin extends Plugin for notifications
type NotifierPlugin interface {
	Plugin
	SendNotification(ctx context.Context, notification *Notification) error
	GetSupportedChannels() []string
}

// IntegrationPlugin extends Plugin for external integrations
type IntegrationPlugin interface {
	Plugin
	Connect(ctx context.Context, credentials map[string]string) error
	Disconnect(ctx context.Context) error
	IsConnected() bool
	SyncData(ctx context.Context) error
}

// ProviderPlugin extends Plugin for custom providers
type ProviderPlugin interface {
	Plugin
	CreateProvider(config *providers.ProviderConfig) (providers.Provider, error)
	GetSupportedProviderTypes() []providers.ProviderType
}

// PluginConfig represents plugin configuration
type PluginConfig struct {
	Name     string                 `json:"name"`
	Type     PluginType            `json:"type"`
	Version  string                `json:"version"`
	Enabled  bool                  `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
	
	// Plugin-specific configuration
	Binary   string                 `json:"binary,omitempty"`
	Args     []string              `json:"args,omitempty"`
	Env      map[string]string     `json:"env,omitempty"`
	
	// Security settings
	Sandboxed    bool     `json:"sandboxed"`
	Permissions  []string `json:"permissions"`
	AllowedHosts []string `json:"allowed_hosts,omitempty"`
	
	// Resource limits
	MaxMemory  int64         `json:"max_memory,omitempty"`
	MaxCPU     float64       `json:"max_cpu,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
}

// ConfigSchema defines the configuration schema for a plugin
type ConfigSchema struct {
	Properties map[string]*PropertySchema `json:"properties"`
	Required   []string                   `json:"required"`
}

// PropertySchema defines a configuration property
type PropertySchema struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Min         *float64    `json:"min,omitempty"`
	Max         *float64    `json:"max,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
}

// CommandDefinition defines a plugin command
type CommandDefinition struct {
	Name        string            `json:"name"`
	Usage       string            `json:"usage"`
	Description string            `json:"description"`
	Flags       []FlagDefinition  `json:"flags"`
	Subcommands []CommandDefinition `json:"subcommands,omitempty"`
}

// FlagDefinition defines a command flag
type FlagDefinition struct {
	Name        string `json:"name"`
	Short       string `json:"short,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
}

// WebhookEvent represents a webhook event
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Headers   map[string]string      `json:"headers"`
	
	// Provider context
	Provider *providers.ProviderConfig `json:"provider,omitempty"`
	Repository struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
	} `json:"repository,omitempty"`
}

// BuildConfig represents CI/CD build configuration
type BuildConfig struct {
	Repository struct {
		Owner string `json:"owner"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"repository"`
	
	Branch    string            `json:"branch"`
	Commit    string            `json:"commit"`
	Variables map[string]string `json:"variables,omitempty"`
	
	// Build settings
	BuildFile string   `json:"build_file,omitempty"`
	Commands  []string `json:"commands,omitempty"`
	Image     string   `json:"image,omitempty"`
	Timeout   time.Duration `json:"timeout,omitempty"`
}

// BuildResult represents the result of a build trigger
type BuildResult struct {
	BuildID   string    `json:"build_id"`
	Status    string    `json:"status"`
	URL       string    `json:"url,omitempty"`
	StartedAt time.Time `json:"started_at"`
}

// BuildStatus represents the current status of a build
type BuildStatus struct {
	BuildID     string            `json:"build_id"`
	Status      string            `json:"status"`
	Phase       string            `json:"phase,omitempty"`
	Progress    float64           `json:"progress,omitempty"`
	StartedAt   time.Time         `json:"started_at"`
	FinishedAt  *time.Time        `json:"finished_at,omitempty"`
	Duration    time.Duration     `json:"duration,omitempty"`
	URL         string            `json:"url,omitempty"`
	Logs        string            `json:"logs,omitempty"`
	Artifacts   []Artifact        `json:"artifacts,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// Artifact represents a build artifact
type Artifact struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	URL      string `json:"url,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}

// Notification represents a notification to be sent
type Notification struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Level       NotificationLevel      `json:"level"`
	Channels    []string               `json:"channels"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Attachments []NotificationAttachment `json:"attachments,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NotificationLevel represents the severity level
type NotificationLevel string

const (
	NotificationLevelInfo    NotificationLevel = "info"
	NotificationLevelWarning NotificationLevel = "warning"
	NotificationLevelError   NotificationLevel = "error"
	NotificationLevelSuccess NotificationLevel = "success"
)

// NotificationAttachment represents an attachment
type NotificationAttachment struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Size     int64  `json:"size,omitempty"`
	Preview  string `json:"preview,omitempty"`
}

// PluginRegistry manages plugin registration and lifecycle
type PluginRegistry interface {
	// Registration
	Register(plugin Plugin) error
	Unregister(name string) error
	Get(name string) (Plugin, error)
	List() []Plugin
	ListByType(pluginType PluginType) []Plugin
	
	// Lifecycle management
	StartPlugin(name string) error
	StopPlugin(name string) error
	RestartPlugin(name string) error
	StartAll() error
	StopAll() error
	
	// Configuration
	LoadConfig(path string) error
	SaveConfig(path string) error
	UpdatePluginConfig(name string, config *PluginConfig) error
	
	// Discovery
	DiscoverPlugins(paths []string) error
	InstallPlugin(source string) error
	UninstallPlugin(name string) error
}

// PluginManager handles plugin execution and communication
type PluginManager interface {
	// Execution
	ExecuteCommand(pluginName, command string, args []string) error
	HandleWebhook(event *WebhookEvent) error
	TriggerBuild(pluginName string, config *BuildConfig) (*BuildResult, error)
	SendNotification(pluginName string, notification *Notification) error
	
	// Health and monitoring
	GetPluginHealth(name string) (*PluginHealth, error)
	GetPluginMetrics(name string) (*PluginMetrics, error)
	
	// Communication
	SendMessage(pluginName string, message *PluginMessage) (*PluginMessage, error)
	BroadcastMessage(message *PluginMessage) error
}

// PluginHealth represents plugin health status
type PluginHealth struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Uptime    time.Duration `json:"uptime"`
	LastCheck time.Time `json:"last_check"`
	Errors    []string  `json:"errors,omitempty"`
}

// PluginMetrics represents plugin performance metrics
type PluginMetrics struct {
	Name           string        `json:"name"`
	CPUUsage       float64       `json:"cpu_usage"`
	MemoryUsage    int64         `json:"memory_usage"`
	RequestCount   int64         `json:"request_count"`
	ErrorCount     int64         `json:"error_count"`
	AverageLatency time.Duration `json:"average_latency"`
	CollectedAt    time.Time     `json:"collected_at"`
}

// PluginMessage represents inter-plugin communication
type PluginMessage struct {
	ID        string                 `json:"id"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// PluginLoader handles plugin loading and discovery
type PluginLoader interface {
	LoadPlugin(path string) (Plugin, error)
	LoadFromConfig(config *PluginConfig) (Plugin, error)
	ValidatePlugin(plugin Plugin) error
	GetPluginInfo(path string) (*PluginInfo, error)
}

// PluginInfo represents plugin metadata
type PluginInfo struct {
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Type        PluginType `json:"type"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	License     string     `json:"license"`
	Homepage    string     `json:"homepage"`
	Repository  string     `json:"repository"`
	
	// Requirements
	MinVersion    string   `json:"min_version"`
	Dependencies  []string `json:"dependencies"`
	Permissions   []string `json:"permissions"`
	
	// Build info
	BuildDate     time.Time `json:"build_date"`
	CommitHash    string    `json:"commit_hash"`
	Architecture  string    `json:"architecture"`
	OS           string    `json:"os"`
}

// PluginContext provides context and utilities to plugins
type PluginContext interface {
	// Logging
	Log(level string, message string, fields map[string]interface{})
	
	// Configuration
	GetConfig() *PluginConfig
	UpdateConfig(config *PluginConfig) error
	
	// Provider access
	GetProvider(name string) (providers.Provider, error)
	GetProviderManager() interface{} // Returns the actual provider manager
	
	// Events
	EmitEvent(event *PluginEvent) error
	SubscribeToEvent(eventType string, handler EventHandler) error
	
	// Storage
	GetStorage() PluginStorage
	
	// HTTP utilities
	MakeHTTPRequest(method, url string, headers map[string]string, body []byte) (*HTTPResponse, error)
	
	// Filesystem access (sandboxed)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	ListFiles(dir string) ([]string, error)
}

// PluginEvent represents an event emitted by plugins
type PluginEvent struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventHandler handles plugin events
type EventHandler func(event *PluginEvent) error

// PluginStorage provides sandboxed storage for plugins
type PluginStorage interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
	List(prefix string) ([]string, error)
	Clear() error
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
}

// PluginExecutor handles different plugin execution methods
type PluginExecutor interface {
	Execute(plugin Plugin, method string, args ...interface{}) (interface{}, error)
	ExecuteAsync(plugin Plugin, method string, args ...interface{}) (<-chan interface{}, <-chan error)
}

// SecurityManager handles plugin security and sandboxing
type SecurityManager interface {
	ValidatePlugin(plugin Plugin) error
	CreateSandbox(plugin Plugin) (Sandbox, error)
	CheckPermissions(plugin Plugin, permission string) bool
	AuditPluginActivity(plugin Plugin, activity string, details map[string]interface{})
}

// Sandbox represents a plugin execution environment
type Sandbox interface {
	Start() error
	Stop() error
	Execute(command string, args ...string) ([]byte, error)
	SetResourceLimits(memory int64, cpu float64) error
	GetResourceUsage() (*ResourceUsage, error)
}

// ResourceUsage represents current resource usage
type ResourceUsage struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryBytes   int64   `json:"memory_bytes"`
	DiskBytes     int64   `json:"disk_bytes"`
	NetworkRx     int64   `json:"network_rx"`
	NetworkTx     int64   `json:"network_tx"`
}