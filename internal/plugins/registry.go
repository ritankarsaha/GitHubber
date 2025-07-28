/*
 * GitHubber - Plugin Registry Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Plugin registry and lifecycle management
 */

package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DefaultRegistry is the global plugin registry
var DefaultRegistry = NewRegistry()

// Registry implements PluginRegistry
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
	configs map[string]*PluginConfig
	running map[string]bool
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
		configs: make(map[string]*PluginConfig),
		running: make(map[string]bool),
	}
}

// Register registers a plugin
func (r *Registry) Register(plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	name := plugin.Name()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already registered", name)
	}

	r.plugins[name] = plugin
	r.running[name] = false
	return nil
}

// Unregister unregisters a plugin
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Stop plugin if running
	if r.running[name] {
		if err := plugin.Stop(); err != nil {
			return fmt.Errorf("failed to stop plugin %s: %w", name, err)
		}
	}

	delete(r.plugins, name)
	delete(r.configs, name)
	delete(r.running, name)
	return nil
}

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List returns all registered plugins
func (r *Registry) List() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugins := make([]Plugin, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// ListByType returns plugins filtered by type
func (r *Registry) ListByType(pluginType PluginType) []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var plugins []Plugin
	for _, plugin := range r.plugins {
		if plugin.Type() == pluginType {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// StartPlugin starts a plugin
func (r *Registry) StartPlugin(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if r.running[name] {
		return fmt.Errorf("plugin %s is already running", name)
	}

	if err := plugin.Start(); err != nil {
		return fmt.Errorf("failed to start plugin %s: %w", name, err)
	}

	r.running[name] = true
	return nil
}

// StopPlugin stops a plugin
func (r *Registry) StopPlugin(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if !r.running[name] {
		return fmt.Errorf("plugin %s is not running", name)
	}

	if err := plugin.Stop(); err != nil {
		return fmt.Errorf("failed to stop plugin %s: %w", name, err)
	}

	r.running[name] = false
	return nil
}

// RestartPlugin restarts a plugin
func (r *Registry) RestartPlugin(name string) error {
	if err := r.StopPlugin(name); err != nil {
		return err
	}
	return r.StartPlugin(name)
}

// StartAll starts all registered plugins
func (r *Registry) StartAll() error {
	r.mu.RLock()
	plugins := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		plugins = append(plugins, name)
	}
	r.mu.RUnlock()

	var errors []string
	for _, name := range plugins {
		if err := r.StartPlugin(name); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to start some plugins: %v", errors)
	}

	return nil
}

// StopAll stops all running plugins
func (r *Registry) StopAll() error {
	r.mu.RLock()
	plugins := make([]string, 0, len(r.running))
	for name, running := range r.running {
		if running {
			plugins = append(plugins, name)
		}
	}
	r.mu.RUnlock()

	var errors []string
	for _, name := range plugins {
		if err := r.StopPlugin(name); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to stop some plugins: %v", errors)
	}

	return nil
}

// LoadConfig loads plugin configuration from file
func (r *Registry) LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configs map[string]*PluginConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for name, config := range configs {
		r.configs[name] = config
	}

	return nil
}

// SaveConfig saves plugin configuration to file
func (r *Registry) SaveConfig(path string) error {
	r.mu.RLock()
	configs := make(map[string]*PluginConfig, len(r.configs))
	for name, config := range r.configs {
		configs[name] = config
	}
	r.mu.RUnlock()

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// UpdatePluginConfig updates plugin configuration
func (r *Registry) UpdatePluginConfig(name string, config *PluginConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Validate configuration
	if err := plugin.Validate(config); err != nil {
		return fmt.Errorf("invalid config for plugin %s: %w", name, err)
	}

	r.configs[name] = config
	return nil
}

// DiscoverPlugins discovers plugins in specified directories
func (r *Registry) DiscoverPlugins(paths []string) error {
	loader := NewPluginLoader()

	for _, path := range paths {
		err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if d.IsDir() {
				return nil
			}

			// Look for plugin files (e.g., .so, .dll, or executable files)
			if r.isPluginFile(filePath) {
				plugin, err := loader.LoadPlugin(filePath)
				if err != nil {
					// Log error but continue discovery
					fmt.Printf("Warning: failed to load plugin %s: %v\n", filePath, err)
					return nil
				}

				if err := r.Register(plugin); err != nil {
					fmt.Printf("Warning: failed to register plugin %s: %v\n", plugin.Name(), err)
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to discover plugins in %s: %w", path, err)
		}
	}

	return nil
}

// InstallPlugin installs a plugin from a source
func (r *Registry) InstallPlugin(source string) error {
	// This is a simplified implementation
	// In production, you'd want to handle different sources (URLs, files, etc.)
	loader := NewPluginLoader()

	plugin, err := loader.LoadPlugin(source)
	if err != nil {
		return fmt.Errorf("failed to load plugin from %s: %w", source, err)
	}

	return r.Register(plugin)
}

// UninstallPlugin uninstalls a plugin
func (r *Registry) UninstallPlugin(name string) error {
	return r.Unregister(name)
}

// Helper method to check if a file is a plugin
func (r *Registry) isPluginFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".so", ".dll", ".dylib":
		return true
	case ".exe":
		return true
	case "":
		// Executable without extension (Unix)
		if info, err := os.Stat(filePath); err == nil {
			return info.Mode().Perm()&0111 != 0
		}
	}
	return false
}

// Manager implements PluginManager
type Manager struct {
	registry PluginRegistry
	health   map[string]*PluginHealth
	metrics  map[string]*PluginMetrics
	mu       sync.RWMutex
}

// NewManager creates a new plugin manager
func NewManager(registry PluginRegistry) *Manager {
	if registry == nil {
		registry = DefaultRegistry
	}

	return &Manager{
		registry: registry,
		health:   make(map[string]*PluginHealth),
		metrics:  make(map[string]*PluginMetrics),
	}
}

// ExecuteCommand executes a command on a plugin
func (m *Manager) ExecuteCommand(pluginName, command string, args []string) error {
	plugin, err := m.registry.Get(pluginName)
	if err != nil {
		return err
	}

	commandPlugin, ok := plugin.(CommandPlugin)
	if !ok {
		return fmt.Errorf("plugin %s does not support commands", pluginName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return commandPlugin.ExecuteCommand(ctx, command, args)
}

// HandleWebhook handles a webhook event
func (m *Manager) HandleWebhook(event *WebhookEvent) error {
	plugins := m.registry.ListByType(PluginTypeWebhook)
	
	var errors []string
	for _, plugin := range plugins {
		webhookPlugin, ok := plugin.(WebhookPlugin)
		if !ok {
			continue
		}

		// Check if plugin supports this event type
		supported := false
		for _, eventType := range webhookPlugin.GetSupportedEvents() {
			if eventType == event.Type {
				supported = true
				break
			}
		}

		if !supported {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := webhookPlugin.HandleWebhook(ctx, event); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", plugin.Name(), err))
		}
		cancel()
	}

	if len(errors) > 0 {
		return fmt.Errorf("webhook handling failed for some plugins: %v", errors)
	}

	return nil
}

// TriggerBuild triggers a build using a CI plugin
func (m *Manager) TriggerBuild(pluginName string, config *BuildConfig) (*BuildResult, error) {
	plugin, err := m.registry.Get(pluginName)
	if err != nil {
		return nil, err
	}

	ciPlugin, ok := plugin.(CIPlugin)
	if !ok {
		return nil, fmt.Errorf("plugin %s does not support CI/CD", pluginName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	return ciPlugin.TriggerBuild(ctx, config)
}

// SendNotification sends a notification using a notifier plugin
func (m *Manager) SendNotification(pluginName string, notification *Notification) error {
	plugin, err := m.registry.Get(pluginName)
	if err != nil {
		return err
	}

	notifierPlugin, ok := plugin.(NotifierPlugin)
	if !ok {
		return fmt.Errorf("plugin %s does not support notifications", pluginName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return notifierPlugin.SendNotification(ctx, notification)
}

// GetPluginHealth returns plugin health status
func (m *Manager) GetPluginHealth(name string) (*PluginHealth, error) {
	plugin, err := m.registry.Get(name)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	health, exists := m.health[name]
	m.mu.RUnlock()

	if !exists {
		// Create initial health status
		health = &PluginHealth{
			Name:      name,
			Status:    "unknown",
			LastCheck: time.Now(),
		}

		m.mu.Lock()
		m.health[name] = health
		m.mu.Unlock()
	}

	// Update health status
	if plugin.IsRunning() {
		health.Status = "healthy"
	} else {
		health.Status = "stopped"
	}
	health.LastCheck = time.Now()

	return health, nil
}

// GetPluginMetrics returns plugin performance metrics
func (m *Manager) GetPluginMetrics(name string) (*PluginMetrics, error) {
	_, err := m.registry.Get(name)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	metrics, exists := m.metrics[name]
	m.mu.RUnlock()

	if !exists {
		// Create initial metrics
		metrics = &PluginMetrics{
			Name:        name,
			CollectedAt: time.Now(),
		}

		m.mu.Lock()
		m.metrics[name] = metrics
		m.mu.Unlock()
	}

	// In a real implementation, you'd collect actual metrics here
	metrics.CollectedAt = time.Now()

	return metrics, nil
}

// SendMessage sends a message to a plugin
func (m *Manager) SendMessage(pluginName string, message *PluginMessage) (*PluginMessage, error) {
	_, err := m.registry.Get(pluginName)
	if err != nil {
		return nil, err
	}

	// In a real implementation, you'd have a message queue/communication system
	// For now, just return a placeholder response
	response := &PluginMessage{
		ID:        fmt.Sprintf("response-%d", time.Now().Unix()),
		From:      pluginName,
		To:        message.From,
		Type:      "response",
		Data:      map[string]interface{}{"status": "received"},
		Timestamp: time.Now(),
	}

	return response, nil
}

// BroadcastMessage broadcasts a message to all plugins
func (m *Manager) BroadcastMessage(message *PluginMessage) error {
	plugins := m.registry.List()
	
	var errors []string
	for _, plugin := range plugins {
		if plugin.IsRunning() {
			if _, err := m.SendMessage(plugin.Name(), message); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", plugin.Name(), err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("broadcast failed for some plugins: %v", errors)
	}

	return nil
}

// StartHealthChecker starts a health checker goroutine
func (m *Manager) StartHealthChecker(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.checkAllPluginHealth()
			}
		}
	}()
}

// checkAllPluginHealth checks health of all plugins
func (m *Manager) checkAllPluginHealth() {
	plugins := m.registry.List()
	
	for _, plugin := range plugins {
		if _, err := m.GetPluginHealth(plugin.Name()); err != nil {
			// Log error
			fmt.Printf("Health check failed for plugin %s: %v\n", plugin.Name(), err)
		}
	}
}