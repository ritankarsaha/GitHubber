/*
 * GitHubber - Plugin Loader Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Plugin loading and discovery system
 */

package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"time"
)

// Loader implements PluginLoader
type Loader struct {
	loadedPlugins map[string]*plugin.Plugin
}

// NewPluginLoader creates a new plugin loader
func NewPluginLoader() *Loader {
	return &Loader{
		loadedPlugins: make(map[string]*plugin.Plugin),
	}
}

// LoadPlugin loads a plugin from a file path
func (l *Loader) LoadPlugin(path string) (Plugin, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin file not found: %s", path)
	}

	// For Go plugins (.so files on Linux/macOS, .dll on Windows)
	if isGoPlugin(path) {
		return l.loadGoPlugin(path)
	}

	// For executable plugins
	if isExecutable(path) {
		return l.loadExecutablePlugin(path)
	}

	// For configuration-based plugins
	if isConfigPlugin(path) {
		return l.loadConfigPlugin(path)
	}

	return nil, fmt.Errorf("unsupported plugin type: %s", path)
}

// LoadFromConfig loads a plugin from configuration
func (l *Loader) LoadFromConfig(config *PluginConfig) (Plugin, error) {
	if config.Binary != "" {
		return l.loadExecutablePluginFromConfig(config)
	}

	return nil, fmt.Errorf("no plugin binary specified in config")
}

// ValidatePlugin validates a plugin
func (l *Loader) ValidatePlugin(plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin is nil")
	}

	if plugin.Name() == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	if plugin.Version() == "" {
		return fmt.Errorf("plugin version cannot be empty")
	}

	if plugin.Type() == "" {
		return fmt.Errorf("plugin type cannot be empty")
	}

	return nil
}

// GetPluginInfo extracts plugin information from a file
func (l *Loader) GetPluginInfo(path string) (*PluginInfo, error) {
	// Try to get info from plugin metadata file
	metadataPath := path + ".json"
	if _, err := os.Stat(metadataPath); err == nil {
		return l.loadPluginInfoFromMetadata(metadataPath)
	}

	// Try to load plugin and extract info
	plugin, err := l.LoadPlugin(path)
	if err != nil {
		return nil, err
	}

	return &PluginInfo{
		Name:        plugin.Name(),
		Version:     plugin.Version(),
		Type:        plugin.Type(),
		Description: plugin.Description(),
		Author:      plugin.Author(),
		BuildDate:   time.Now(), // Would be set during build
	}, nil
}

// loadGoPlugin loads a Go plugin (.so/.dll file)
func (l *Loader) loadGoPlugin(path string) (Plugin, error) {
	// Check if already loaded
	if p, exists := l.loadedPlugins[path]; exists {
		return l.extractPluginFromGoPlugin(p)
	}

	// Load the plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open Go plugin: %w", err)
	}

	l.loadedPlugins[path] = p
	return l.extractPluginFromGoPlugin(p)
}

// extractPluginFromGoPlugin extracts the Plugin interface from a Go plugin
func (l *Loader) extractPluginFromGoPlugin(p *plugin.Plugin) (Plugin, error) {
	// Look for the standard plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("plugin symbol not found: %w", err)
	}

	pluginInstance, ok := sym.(Plugin)
	if !ok {
		return nil, fmt.Errorf("symbol does not implement Plugin interface")
	}

	return pluginInstance, nil
}

// loadExecutablePlugin loads an executable plugin
func (l *Loader) loadExecutablePlugin(path string) (Plugin, error) {
	return NewExecutablePlugin(path)
}

// loadConfigPlugin loads a plugin from configuration file
func (l *Loader) loadConfigPlugin(path string) (Plugin, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config PluginConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return l.LoadFromConfig(&config)
}

// loadExecutablePluginFromConfig loads an executable plugin from config
func (l *Loader) loadExecutablePluginFromConfig(config *PluginConfig) (Plugin, error) {
	return NewExecutablePluginFromConfig(config)
}

// loadPluginInfoFromMetadata loads plugin info from metadata file
func (l *Loader) loadPluginInfoFromMetadata(path string) (*PluginInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var info PluginInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse metadata file: %w", err)
	}

	return &info, nil
}

// Helper functions
func isGoPlugin(path string) bool {
	ext := filepath.Ext(path)
	switch runtime.GOOS {
	case "linux", "darwin":
		return ext == ".so"
	case "windows":
		return ext == ".dll"
	default:
		return false
	}
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if file is executable
	if runtime.GOOS == "windows" {
		return filepath.Ext(path) == ".exe"
	}

	return info.Mode().Perm()&0111 != 0
}

func isConfigPlugin(path string) bool {
	return filepath.Ext(path) == ".json"
}

// ExecutablePlugin wraps an executable as a plugin
type ExecutablePlugin struct {
	name        string
	version     string
	pluginType  PluginType
	description string
	author      string
	binaryPath  string
	config      *PluginConfig
	running     bool
}

// NewExecutablePlugin creates a new executable plugin
func NewExecutablePlugin(binaryPath string) (*ExecutablePlugin, error) {
	// Extract plugin info from the executable
	info, err := extractExecutableInfo(binaryPath)
	if err != nil {
		return nil, err
	}

	return &ExecutablePlugin{
		name:        info.Name,
		version:     info.Version,
		pluginType:  info.Type,
		description: info.Description,
		author:      info.Author,
		binaryPath:  binaryPath,
	}, nil
}

// NewExecutablePluginFromConfig creates an executable plugin from config
func NewExecutablePluginFromConfig(config *PluginConfig) (*ExecutablePlugin, error) {
	plugin := &ExecutablePlugin{
		name:        config.Name,
		version:     config.Version,
		pluginType:  config.Type,
		binaryPath:  config.Binary,
		config:      config,
	}

	return plugin, nil
}

// Plugin interface implementation
func (e *ExecutablePlugin) Name() string        { return e.name }
func (e *ExecutablePlugin) Version() string     { return e.version }
func (e *ExecutablePlugin) Type() PluginType    { return e.pluginType }
func (e *ExecutablePlugin) Description() string { return e.description }
func (e *ExecutablePlugin) Author() string      { return e.author }

func (e *ExecutablePlugin) Initialize(config *PluginConfig) error {
	e.config = config
	// Send initialization message to executable
	return e.sendCommand("initialize", config)
}

func (e *ExecutablePlugin) Start() error {
	if e.running {
		return fmt.Errorf("plugin is already running")
	}

	if err := e.sendCommand("start", nil); err != nil {
		return err
	}

	e.running = true
	return nil
}

func (e *ExecutablePlugin) Stop() error {
	if !e.running {
		return fmt.Errorf("plugin is not running")
	}

	if err := e.sendCommand("stop", nil); err != nil {
		return err
	}

	e.running = false
	return nil
}

func (e *ExecutablePlugin) IsRunning() bool {
	return e.running
}

func (e *ExecutablePlugin) GetConfigSchema() *ConfigSchema {
	// This would typically be loaded from plugin metadata
	return &ConfigSchema{
		Properties: make(map[string]*PropertySchema),
		Required:   []string{},
	}
}

func (e *ExecutablePlugin) Validate(config *PluginConfig) error {
	// Basic validation
	if config.Name != e.name {
		return fmt.Errorf("config name mismatch")
	}
	return nil
}

// sendCommand sends a command to the executable plugin
func (e *ExecutablePlugin) sendCommand(command string, data interface{}) error {
	// This would implement the actual communication with the executable
	// For now, it's a placeholder
	return nil
}

// extractExecutableInfo extracts plugin information from an executable
func extractExecutableInfo(binaryPath string) (*PluginInfo, error) {
	// This would typically run the executable with --info flag
	// For now, return default info
	return &PluginInfo{
		Name:        filepath.Base(binaryPath),
		Version:     "1.0.0",
		Type:        PluginTypeCommand,
		Description: "Executable plugin",
		Author:      "Unknown",
		BuildDate:   time.Now(),
	}, nil
}

// PluginDiscovery handles plugin discovery and management
type PluginDiscovery struct {
	searchPaths []string
	loader      PluginLoader
}

// NewPluginDiscovery creates a new plugin discovery service
func NewPluginDiscovery(searchPaths []string) *PluginDiscovery {
	return &PluginDiscovery{
		searchPaths: searchPaths,
		loader:      NewPluginLoader(),
	}
}

// DiscoverAll discovers all plugins in search paths
func (pd *PluginDiscovery) DiscoverAll() ([]Plugin, error) {
	var allPlugins []Plugin
	
	for _, path := range pd.searchPaths {
		plugins, err := pd.DiscoverInPath(path)
		if err != nil {
			return nil, fmt.Errorf("failed to discover plugins in %s: %w", path, err)
		}
		allPlugins = append(allPlugins, plugins...)
	}
	
	return allPlugins, nil
}

// DiscoverInPath discovers plugins in a specific path
func (pd *PluginDiscovery) DiscoverInPath(path string) ([]Plugin, error) {
	var plugins []Plugin
	
	err := filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		// Check if this is a plugin file
		if pd.isPluginFile(filePath) {
			plugin, err := pd.loader.LoadPlugin(filePath)
			if err != nil {
				// Log warning but continue
				fmt.Printf("Warning: failed to load plugin %s: %v\n", filePath, err)
				return nil
			}
			
			plugins = append(plugins, plugin)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return plugins, nil
}

// isPluginFile checks if a file is a plugin
func (pd *PluginDiscovery) isPluginFile(path string) bool {
	ext := filepath.Ext(path)
	
	// Go plugins
	if isGoPlugin(path) {
		return true
	}
	
	// Executable plugins
	if isExecutable(path) {
		return true
	}
	
	// Config-based plugins
	if ext == ".json" {
		return pd.isPluginConfigFile(path)
	}
	
	return false
}

// isPluginConfigFile checks if a JSON file is a plugin config
func (pd *PluginDiscovery) isPluginConfigFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	
	var config PluginConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return false
	}
	
	// Check if it has required plugin fields
	return config.Name != "" && config.Type != "" && config.Version != ""
}

// AddSearchPath adds a new search path
func (pd *PluginDiscovery) AddSearchPath(path string) {
	pd.searchPaths = append(pd.searchPaths, path)
}

// GetSearchPaths returns current search paths
func (pd *PluginDiscovery) GetSearchPaths() []string {
	return pd.searchPaths
}

// PluginValidator validates plugin compatibility and security
type PluginValidator struct {
	minVersion string
	security   SecurityManager
}

// NewPluginValidator creates a new plugin validator
func NewPluginValidator(minVersion string, security SecurityManager) *PluginValidator {
	return &PluginValidator{
		minVersion: minVersion,
		security:   security,
	}
}

// Validate performs comprehensive plugin validation
func (pv *PluginValidator) Validate(plugin Plugin) error {
	// Basic validation
	if plugin.Name() == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}
	
	if plugin.Version() == "" {
		return fmt.Errorf("plugin version cannot be empty")
	}
	
	// Security validation
	if pv.security != nil {
		if err := pv.security.ValidatePlugin(plugin); err != nil {
			return fmt.Errorf("security validation failed: %w", err)
		}
	}
	
	// Version compatibility check
	if err := pv.validateVersion(plugin.Version()); err != nil {
		return fmt.Errorf("version validation failed: %w", err)
	}
	
	return nil
}

// validateVersion checks if plugin version is compatible
func (pv *PluginValidator) validateVersion(version string) error {
	// Simplified version check - in production, use proper semver comparison
	if pv.minVersion != "" && version < pv.minVersion {
		return fmt.Errorf("plugin version %s is below minimum required version %s", version, pv.minVersion)
	}
	return nil
}