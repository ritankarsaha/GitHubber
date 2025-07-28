/*
 * GitHubber - Configuration Manager Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Advanced configuration management with validation, migration, and templating
 */

package config

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// Manager implements ConfigManager
type Manager struct {
	mu               sync.RWMutex
	watchers         map[string]*fsnotify.Watcher
	watchCallbacks   map[string]func(*ApplicationConfig)
	migrations       map[string]*ConfigMigration
	templates        map[string]*ConfigTemplate
	profiles         map[string]*ConfigProfile
	backups          []*ConfigBackup
	validationRules  map[string][]ValidationRule
	defaultConfig    *ApplicationConfig
}

// ValidationRule represents a configuration validation rule
type ValidationRule struct {
	Field       string
	Type        string
	Required    bool
	Min         interface{}
	Max         interface{}
	Pattern     string
	Enum        []string
	Custom      func(interface{}) error
	Description string
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	manager := &Manager{
		watchers:        make(map[string]*fsnotify.Watcher),
		watchCallbacks:  make(map[string]func(*ApplicationConfig)),
		migrations:      make(map[string]*ConfigMigration),
		templates:       make(map[string]*ConfigTemplate),
		profiles:        make(map[string]*ConfigProfile),
		backups:         make([]*ConfigBackup, 0),
		validationRules: make(map[string][]ValidationRule),
	}

	manager.initializeDefaultConfig()
	manager.registerBuiltinMigrations()
	manager.registerBuiltinTemplates()
	manager.registerValidationRules()

	return manager
}

// Load loads configuration from file
func (m *Manager) Load(path string) (*ApplicationConfig, error) {
	if path == "" {
		return nil, fmt.Errorf("configuration path cannot be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	format := m.detectFormat(path)
	config, err := m.parseConfig(data, format)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply environment variables
	if err := m.LoadFromEnvironment(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment variables: %w", err)
	}

	// Validate configuration
	if err := m.Validate(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Migrate if necessary
	if config.Version != ConfigVersionV2 {
		migrated, err := m.Migrate(config, ConfigVersionV2)
		if err != nil {
			return nil, fmt.Errorf("configuration migration failed: %w", err)
		}
		config = migrated
	}

	return config, nil
}

// Save saves configuration to file
func (m *Manager) Save(config *ApplicationConfig, path string) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Update metadata
	if config.Metadata == nil {
		config.Metadata = &ConfigMetadata{}
	}
	config.Metadata.UpdatedAt = time.Now()

	// Create backup before saving
	if err := m.createBackup(config, path); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to create backup: %v\n", err)
	}

	// Validate before saving
	if err := m.Validate(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	format := m.detectFormat(path)
	data, err := m.serializeConfig(config, format)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write to temporary file first
	tempPath := path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp config file: %w", err)
	}

	// Atomic move
	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to move config file: %w", err)
	}

	return nil
}

// Validate validates the entire configuration
func (m *Manager) Validate(config *ApplicationConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	var errors []ConfigValidationError

	// Validate version
	if config.Version == "" {
		errors = append(errors, ConfigValidationError{
			Field:   "version",
			Message: "version is required",
			Code:    "required",
		})
	}

	// Validate core configuration
	if config.Core != nil {
		if errs := m.validateCore(config.Core); errs != nil {
			errors = append(errors, errs...)
		}
	} else {
		errors = append(errors, ConfigValidationError{
			Field:   "core",
			Message: "core configuration is required",
			Code:    "required",
		})
	}

	// Validate providers configuration
	if config.Providers != nil {
		if errs := m.validateProviders(config.Providers); errs != nil {
			errors = append(errors, errs...)
		}
	}

	// Validate plugins configuration
	if config.Plugins != nil {
		if errs := m.validatePlugins(config.Plugins); errs != nil {
			errors = append(errors, errs...)
		}
	}

	// Validate CI configuration
	if config.CI != nil {
		if errs := m.validateCI(config.CI); errs != nil {
			errors = append(errors, errs...)
		}
	}

	// Validate webhooks configuration
	if config.Webhooks != nil {
		if errs := m.validateWebhooks(config.Webhooks); errs != nil {
			errors = append(errors, errs...)
		}
	}

	// Validate security configuration
	if config.Security != nil {
		if errs := m.validateSecurity(config.Security); errs != nil {
			errors = append(errors, errs...)
		}
	}

	if len(errors) > 0 {
		return &ConfigValidationErrors{Errors: errors}
	}

	return nil
}

// ValidatePartial validates a specific configuration section
func (m *Manager) ValidatePartial(section string, data interface{}) error {
	rules, exists := m.validationRules[section]
	if !exists {
		return fmt.Errorf("no validation rules for section: %s", section)
	}

	var errors []ConfigValidationError

	for _, rule := range rules {
		if err := m.validateField(rule, data); err != nil {
			if validationErr, ok := err.(ConfigValidationError); ok {
				errors = append(errors, validationErr)
			} else {
				errors = append(errors, ConfigValidationError{
					Field:   rule.Field,
					Message: err.Error(),
					Code:    "validation_error",
				})
			}
		}
	}

	if len(errors) > 0 {
		return &ConfigValidationErrors{Errors: errors}
	}

	return nil
}

// Migrate migrates configuration to target version
func (m *Manager) Migrate(config *ApplicationConfig, targetVersion ConfigVersion) (*ApplicationConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	if config.Version == targetVersion {
		return config, nil
	}

	migrationPath, err := m.GetMigrationPath(config.Version, targetVersion)
	if err != nil {
		return nil, err
	}

	current := config
	for i := 0; i < len(migrationPath)-1; i++ {
		from := migrationPath[i]
		to := migrationPath[i+1]
		
		migrationKey := fmt.Sprintf("%s->%s", from, to)
		migration, exists := m.migrations[migrationKey]
		if !exists {
			return nil, fmt.Errorf("no migration found from %s to %s", from, to)
		}

		migrated, err := migration.Migrate(current)
		if err != nil {
			return nil, fmt.Errorf("migration from %s to %s failed: %w", from, to, err)
		}

		current = migrated
		current.Version = to
	}

	return current, nil
}

// GetMigrationPath returns the migration path between versions
func (m *Manager) GetMigrationPath(fromVersion, toVersion ConfigVersion) ([]ConfigVersion, error) {
	// Simple linear migration path for now
	versions := []ConfigVersion{ConfigVersionV1, ConfigVersionV2}
	
	fromIndex := -1
	toIndex := -1
	
	for i, v := range versions {
		if v == fromVersion {
			fromIndex = i
		}
		if v == toVersion {
			toIndex = i
		}
	}
	
	if fromIndex == -1 {
		return nil, fmt.Errorf("unknown source version: %s", fromVersion)
	}
	if toIndex == -1 {
		return nil, fmt.Errorf("unknown target version: %s", toVersion)
	}
	
	if fromIndex == toIndex {
		return []ConfigVersion{fromVersion}, nil
	}
	
	var path []ConfigVersion
	if fromIndex < toIndex {
		for i := fromIndex; i <= toIndex; i++ {
			path = append(path, versions[i])
		}
	} else {
		// Reverse migration not implemented
		return nil, fmt.Errorf("reverse migration from %s to %s not supported", fromVersion, toVersion)
	}
	
	return path, nil
}

// Merge merges two configurations
func (m *Manager) Merge(base, override *ApplicationConfig) (*ApplicationConfig, error) {
	if base == nil {
		return override, nil
	}
	if override == nil {
		return base, nil
	}

	// Deep copy base configuration
	result, err := m.deepCopy(base)
	if err != nil {
		return nil, fmt.Errorf("failed to copy base config: %w", err)
	}

	// Merge override into result
	if err := m.mergeConfigs(result, override); err != nil {
		return nil, fmt.Errorf("failed to merge configs: %w", err)
	}

	return result, nil
}

// ApplyTemplate applies template variables to configuration
func (m *Manager) ApplyTemplate(config *ApplicationConfig, variables map[string]string) (*ApplicationConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Serialize to JSON for template processing
	data, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Apply template variables
	content := string(data)
	for key, value := range variables {
		placeholder := fmt.Sprintf("${%s}", key)
		content = strings.ReplaceAll(content, placeholder, value)
	}

	// Parse back to configuration
	var result ApplicationConfig
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templated config: %w", err)
	}

	return &result, nil
}

// LoadFromEnvironment loads values from environment variables
func (m *Manager) LoadFromEnvironment(config *ApplicationConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Define environment variable mappings
	envMappings := map[string]func(string){
		"GITHUBBER_DEBUG": func(value string) {
			if config.Core != nil {
				config.Core.Debug = strings.ToLower(value) == "true"
			}
		},
		"GITHUBBER_ENVIRONMENT": func(value string) {
			if config.Core != nil {
				config.Core.Environment = value
			}
		},
		"GITHUBBER_LOG_LEVEL": func(value string) {
			if config.Logging != nil {
				config.Logging.Level = value
			}
		},
		"GITHUBBER_WEBHOOK_SECRET": func(value string) {
			if config.Webhooks != nil {
				config.Webhooks.Secret = value
			}
		},
	}

	// Apply environment variables
	for envVar, setter := range envMappings {
		if value := os.Getenv(envVar); value != "" {
			setter(value)
		}
	}

	return nil
}

// ExportToEnvironment exports configuration to environment variables
func (m *Manager) ExportToEnvironment(config *ApplicationConfig) map[string]string {
	env := make(map[string]string)

	if config.Core != nil {
		env["GITHUBBER_DEBUG"] = strconv.FormatBool(config.Core.Debug)
		env["GITHUBBER_ENVIRONMENT"] = config.Core.Environment
	}

	if config.Logging != nil {
		env["GITHUBBER_LOG_LEVEL"] = config.Logging.Level
	}

	if config.Webhooks != nil {
		env["GITHUBBER_WEBHOOK_SECRET"] = config.Webhooks.Secret
	}

	return env
}

// GetSchema returns the configuration schema for a version
func (m *Manager) GetSchema(version ConfigVersion) (interface{}, error) {
	switch version {
	case ConfigVersionV1, ConfigVersionV2:
		return m.generateSchema(), nil
	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
}

// GenerateExample generates an example configuration
func (m *Manager) GenerateExample(version ConfigVersion) (*ApplicationConfig, error) {
	switch version {
	case ConfigVersionV1, ConfigVersionV2:
		return m.defaultConfig, nil
	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
}

// Watch watches for configuration file changes
func (m *Manager) Watch(path string, callback func(*ApplicationConfig)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop existing watcher if any
	if watcher, exists := m.watchers[path]; exists {
		watcher.Close()
		delete(m.watchers, path)
	}

	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	if err := watcher.Add(path); err != nil {
		watcher.Close()
		return fmt.Errorf("failed to add path to watcher: %w", err)
	}

	m.watchers[path] = watcher
	m.watchCallbacks[path] = callback

	// Start watching in goroutine
	go m.watchFile(path, watcher, callback)

	return nil
}

// StopWatching stops watching a configuration file
func (m *Manager) StopWatching(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	watcher, exists := m.watchers[path]
	if !exists {
		return fmt.Errorf("no watcher found for path: %s", path)
	}

	watcher.Close()
	delete(m.watchers, path)
	delete(m.watchCallbacks, path)

	return nil
}

// Helper methods

func (m *Manager) detectFormat(path string) ConfigFormat {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return FormatYAML
	case ".toml":
		return FormatTOML
	case ".hcl":
		return FormatHCL
	default:
		return FormatJSON
	}
}

func (m *Manager) parseConfig(data []byte, format ConfigFormat) (*ApplicationConfig, error) {
	var config ApplicationConfig

	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return &config, nil
}

func (m *Manager) serializeConfig(config *ApplicationConfig, format ConfigFormat) ([]byte, error) {
	switch format {
	case FormatJSON:
		return json.MarshalIndent(config, "", "  ")
	case FormatYAML:
		return yaml.Marshal(config)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func (m *Manager) deepCopy(config *ApplicationConfig) (*ApplicationConfig, error) {
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	var copy ApplicationConfig
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, err
	}

	return &copy, nil
}

func (m *Manager) mergeConfigs(base, override *ApplicationConfig) error {
	// Use reflection to merge fields
	baseValue := reflect.ValueOf(base).Elem()
	overrideValue := reflect.ValueOf(override).Elem()

	return m.mergeValues(baseValue, overrideValue)
}

func (m *Manager) mergeValues(base, override reflect.Value) error {
	if !override.IsValid() {
		return nil
	}

	switch override.Kind() {
	case reflect.Ptr:
		if override.IsNil() {
			return nil
		}
		if base.IsNil() {
			base.Set(reflect.New(base.Type().Elem()))
		}
		return m.mergeValues(base.Elem(), override.Elem())

	case reflect.Struct:
		for i := 0; i < override.NumField(); i++ {
			if err := m.mergeValues(base.Field(i), override.Field(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		if override.IsNil() {
			return nil
		}
		if base.IsNil() {
			base.Set(reflect.MakeMap(base.Type()))
		}
		for _, key := range override.MapKeys() {
			base.SetMapIndex(key, override.MapIndex(key))
		}

	case reflect.Slice:
		if override.IsNil() {
			return nil
		}
		base.Set(override)

	default:
		if !override.IsZero() {
			base.Set(override)
		}
	}

	return nil
}

func (m *Manager) watchFile(path string, watcher *fsnotify.Watcher, callback func(*ApplicationConfig)) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Reload configuration
				if config, err := m.Load(path); err == nil {
					callback(config)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

func (m *Manager) createBackup(config *ApplicationConfig, originalPath string) error {
	backupID := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", originalPath, time.Now().UnixNano()))))
	backupPath := fmt.Sprintf("%s.backup.%s", originalPath, backupID[:8])

	if err := m.Save(config, backupPath); err != nil {
		return err
	}

	// Calculate checksum
	data, _ := os.ReadFile(backupPath)
	checksum := fmt.Sprintf("%x", md5.Sum(data))

	backup := &ConfigBackup{
		ID:          backupID,
		Path:        backupPath,
		Config:      config,
		CreatedAt:   time.Now(),
		Description: fmt.Sprintf("Backup before modifying %s", originalPath),
		Checksum:    checksum,
	}

	m.mu.Lock()
	m.backups = append(m.backups, backup)
	m.mu.Unlock()

	return nil
}

// ConfigValidationErrors represents multiple validation errors
type ConfigValidationErrors struct {
	Errors []ConfigValidationError `json:"errors"`
}

func (e *ConfigValidationErrors) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%d configuration validation errors", len(e.Errors))
}

// Validation helper methods will be implemented in separate validation files
func (m *Manager) validateField(rule ValidationRule, data interface{}) error {
	// Implementation would validate individual fields based on rules
	return nil
}

func (m *Manager) validateCore(core *CoreConfig) []ConfigValidationError {
	var errors []ConfigValidationError
	// Implementation would validate core configuration
	return errors
}

func (m *Manager) validateProviders(providers *ProvidersConfig) []ConfigValidationError {
	var errors []ConfigValidationError
	// Implementation would validate providers configuration
	return errors
}

func (m *Manager) validatePlugins(plugins *PluginsConfig) []ConfigValidationError {
	var errors []ConfigValidationError
	// Implementation would validate plugins configuration
	return errors
}

func (m *Manager) validateCI(ci *CIConfig) []ConfigValidationError {
	var errors []ConfigValidationError
	// Implementation would validate CI configuration
	return errors
}

func (m *Manager) validateWebhooks(webhooks *WebhooksConfig) []ConfigValidationError {
	var errors []ConfigValidationError
	// Implementation would validate webhooks configuration
	return errors
}

func (m *Manager) validateSecurity(security *SecurityConfig) []ConfigValidationError {
	var errors []ConfigValidationError
	// Implementation would validate security configuration
	return errors
}

func (m *Manager) initializeDefaultConfig() {
	// Initialize default configuration
	m.defaultConfig = &ApplicationConfig{
		Version: ConfigVersionV2,
		Metadata: &ConfigMetadata{
			Name:      "GitHubber Configuration",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Core: &CoreConfig{
			AppName:     "GitHubber",
			Environment: "development",
			Debug:       false,
			Timeout:     30 * time.Second,
		},
		// ... other default values
	}
}

func (m *Manager) registerBuiltinMigrations() {
	// Register built-in migrations
	m.migrations["v1->v2"] = &ConfigMigration{
		FromVersion: ConfigVersionV1,
		ToVersion:   ConfigVersionV2,
		Description: "Migrate from v1 to v2 schema",
		Migrate: func(config *ApplicationConfig) (*ApplicationConfig, error) {
			// Migration logic would go here
			return config, nil
		},
	}
}

func (m *Manager) registerBuiltinTemplates() {
	// Register built-in templates
}

func (m *Manager) registerValidationRules() {
	// Register validation rules for different sections
}

func (m *Manager) generateSchema() interface{} {
	// Generate JSON schema for configuration
	return map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type":    "object",
		// ... schema definition
	}
}