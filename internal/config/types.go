/*
 * GitHubber - Advanced Configuration Management Types
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Comprehensive configuration management system
 */

package config

import (
	"fmt"
	"time"

	"github.com/ritankarsaha/git-tool/internal/providers"
	"github.com/ritankarsaha/git-tool/internal/plugins"
	"github.com/ritankarsaha/git-tool/internal/ci"
)

// ConfigVersion represents the configuration schema version
type ConfigVersion string

const (
	ConfigVersionV1 ConfigVersion = "v1"
	ConfigVersionV2 ConfigVersion = "v2"
)

// ConfigFormat represents the configuration file format
type ConfigFormat string

const (
	FormatJSON ConfigFormat = "json"
	FormatYAML ConfigFormat = "yaml"
	FormatTOML ConfigFormat = "toml"
	FormatHCL  ConfigFormat = "hcl"
)

// ApplicationConfig represents the main application configuration
type ApplicationConfig struct {
	Version     ConfigVersion    `json:"version" yaml:"version" toml:"version"`
	Metadata    *ConfigMetadata  `json:"metadata,omitempty" yaml:"metadata,omitempty" toml:"metadata,omitempty"`
	Core        *CoreConfig      `json:"core" yaml:"core" toml:"core"`
	Providers   *ProvidersConfig `json:"providers" yaml:"providers" toml:"providers"`
	Plugins     *PluginsConfig   `json:"plugins" yaml:"plugins" toml:"plugins"`
	CI          *CIConfig        `json:"ci" yaml:"ci" toml:"ci"`
	Webhooks    *WebhooksConfig  `json:"webhooks" yaml:"webhooks" toml:"webhooks"`
	Security    *SecurityConfig  `json:"security" yaml:"security" toml:"security"`
	Logging     *LoggingConfig   `json:"logging" yaml:"logging" toml:"logging"`
	Monitoring  *MonitoringConfig `json:"monitoring" yaml:"monitoring" toml:"monitoring"`
	Features    *FeaturesConfig  `json:"features" yaml:"features" toml:"features"`
	Extensions  map[string]interface{} `json:"extensions,omitempty" yaml:"extensions,omitempty" toml:"extensions,omitempty"`
}

// ConfigMetadata contains configuration metadata
type ConfigMetadata struct {
	Name        string    `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Author      string    `json:"author,omitempty" yaml:"author,omitempty" toml:"author,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty" toml:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty" toml:"updated_at,omitempty"`
	Tags        []string  `json:"tags,omitempty" yaml:"tags,omitempty" toml:"tags,omitempty"`
}

// CoreConfig contains core application settings
type CoreConfig struct {
	AppName         string        `json:"app_name" yaml:"app_name" toml:"app_name"`
	Environment     string        `json:"environment" yaml:"environment" toml:"environment"`
	Debug           bool          `json:"debug" yaml:"debug" toml:"debug"`
	Timeout         time.Duration `json:"timeout" yaml:"timeout" toml:"timeout"`
	MaxRetries      int           `json:"max_retries" yaml:"max_retries" toml:"max_retries"`
	RetryDelay      time.Duration `json:"retry_delay" yaml:"retry_delay" toml:"retry_delay"`
	CacheDir        string        `json:"cache_dir" yaml:"cache_dir" toml:"cache_dir"`
	ConfigDir       string        `json:"config_dir" yaml:"config_dir" toml:"config_dir"`
	DataDir         string        `json:"data_dir" yaml:"data_dir" toml:"data_dir"`
	TempDir         string        `json:"temp_dir" yaml:"temp_dir" toml:"temp_dir"`
	
	// Performance settings
	Concurrency     int           `json:"concurrency" yaml:"concurrency" toml:"concurrency"`
	BatchSize       int           `json:"batch_size" yaml:"batch_size" toml:"batch_size"`
	RequestRate     int           `json:"request_rate" yaml:"request_rate" toml:"request_rate"`
	
	// UI settings
	Theme           string        `json:"theme" yaml:"theme" toml:"theme"`
	ColorScheme     string        `json:"color_scheme" yaml:"color_scheme" toml:"color_scheme"`
	DateFormat      string        `json:"date_format" yaml:"date_format" toml:"date_format"`
	TimeFormat      string        `json:"time_format" yaml:"time_format" toml:"time_format"`
	Timezone        string        `json:"timezone" yaml:"timezone" toml:"timezone"`
}

// ProvidersConfig contains provider configurations
type ProvidersConfig struct {
	Default   string                              `json:"default" yaml:"default" toml:"default"`
	Providers map[string]*providers.ProviderConfig `json:"providers" yaml:"providers" toml:"providers"`
	
	// Global provider settings
	ConnectTimeout time.Duration `json:"connect_timeout" yaml:"connect_timeout" toml:"connect_timeout"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout" toml:"request_timeout"`
	MaxRetries     int          `json:"max_retries" yaml:"max_retries" toml:"max_retries"`
	RateLimit      int          `json:"rate_limit" yaml:"rate_limit" toml:"rate_limit"`
	
	// Authentication
	AuthCache      bool          `json:"auth_cache" yaml:"auth_cache" toml:"auth_cache"`
	AuthExpiry     time.Duration `json:"auth_expiry" yaml:"auth_expiry" toml:"auth_expiry"`
}

// PluginsConfig contains plugin configurations
type PluginsConfig struct {
	Enabled     bool                               `json:"enabled" yaml:"enabled" toml:"enabled"`
	SearchPaths []string                           `json:"search_paths" yaml:"search_paths" toml:"search_paths"`
	Plugins     map[string]*plugins.PluginConfig   `json:"plugins" yaml:"plugins" toml:"plugins"`
	
	// Plugin system settings
	AutoLoad        bool          `json:"auto_load" yaml:"auto_load" toml:"auto_load"`
	AutoUpdate      bool          `json:"auto_update" yaml:"auto_update" toml:"auto_update"`
	Sandboxing      bool          `json:"sandboxing" yaml:"sandboxing" toml:"sandboxing"`
	MaxMemory       int64         `json:"max_memory" yaml:"max_memory" toml:"max_memory"`
	MaxCPU          float64       `json:"max_cpu" yaml:"max_cpu" toml:"max_cpu"`
	ExecutionTimeout time.Duration `json:"execution_timeout" yaml:"execution_timeout" toml:"execution_timeout"`
	
	// Security settings
	AllowUnsigned   bool          `json:"allow_unsigned" yaml:"allow_unsigned" toml:"allow_unsigned"`
	TrustedSources  []string      `json:"trusted_sources" yaml:"trusted_sources" toml:"trusted_sources"`
	BlockedPlugins  []string      `json:"blocked_plugins" yaml:"blocked_plugins" toml:"blocked_plugins"`
}

// CIConfig contains CI/CD configurations
type CIConfig struct {
	Enabled   bool                        `json:"enabled" yaml:"enabled" toml:"enabled"`
	Default   string                      `json:"default" yaml:"default" toml:"default"`
	Providers map[string]*ci.CIConfig     `json:"providers" yaml:"providers" toml:"providers"`
	
	// Global CI settings
	AutoTrigger     bool          `json:"auto_trigger" yaml:"auto_trigger" toml:"auto_trigger"`
	ParallelBuilds  int           `json:"parallel_builds" yaml:"parallel_builds" toml:"parallel_builds"`
	BuildTimeout    time.Duration `json:"build_timeout" yaml:"build_timeout" toml:"build_timeout"`
	ArtifactRetention time.Duration `json:"artifact_retention" yaml:"artifact_retention" toml:"artifact_retention"`
	
	// Notifications
	NotifyOnSuccess bool          `json:"notify_on_success" yaml:"notify_on_success" toml:"notify_on_success"`
	NotifyOnFailure bool          `json:"notify_on_failure" yaml:"notify_on_failure" toml:"notify_on_failure"`
	NotificationChannels []string `json:"notification_channels" yaml:"notification_channels" toml:"notification_channels"`
}

// WebhooksConfig contains webhook configurations
type WebhooksConfig struct {
	Enabled     bool                    `json:"enabled" yaml:"enabled" toml:"enabled"`
	Port        int                     `json:"port" yaml:"port" toml:"port"`
	Path        string                  `json:"path" yaml:"path" toml:"path"`
	Secret      string                  `json:"secret" yaml:"secret" toml:"secret"`
	TLS         *TLSConfig              `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
	
	// Processing settings
	QueueSize   int                     `json:"queue_size" yaml:"queue_size" toml:"queue_size"`
	Workers     int                     `json:"workers" yaml:"workers" toml:"workers"`
	Timeout     time.Duration           `json:"timeout" yaml:"timeout" toml:"timeout"`
	
	// Security settings
	IPWhitelist []string                `json:"ip_whitelist" yaml:"ip_whitelist" toml:"ip_whitelist"`
	UserAgent   string                  `json:"user_agent" yaml:"user_agent" toml:"user_agent"`
	MaxPayloadSize int64                `json:"max_payload_size" yaml:"max_payload_size" toml:"max_payload_size"`
	
	// Provider-specific webhook settings
	Providers   map[string]*WebhookProviderConfig `json:"providers" yaml:"providers" toml:"providers"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled   bool   `json:"enabled" yaml:"enabled" toml:"enabled"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file"`
	CAFile    string `json:"ca_file,omitempty" yaml:"ca_file,omitempty" toml:"ca_file,omitempty"`
	Insecure  bool   `json:"insecure" yaml:"insecure" toml:"insecure"`
}

// WebhookProviderConfig contains provider-specific webhook settings
type WebhookProviderConfig struct {
	Secret    string            `json:"secret" yaml:"secret" toml:"secret"`
	Events    []string          `json:"events" yaml:"events" toml:"events"`
	Headers   map[string]string `json:"headers" yaml:"headers" toml:"headers"`
	Enabled   bool              `json:"enabled" yaml:"enabled" toml:"enabled"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	// Authentication
	EnableAuth       bool          `json:"enable_auth" yaml:"enable_auth" toml:"enable_auth"`
	SessionTimeout   time.Duration `json:"session_timeout" yaml:"session_timeout" toml:"session_timeout"`
	MaxSessions      int           `json:"max_sessions" yaml:"max_sessions" toml:"max_sessions"`
	
	// API Security
	APIKeys          map[string]*APIKeyConfig `json:"api_keys" yaml:"api_keys" toml:"api_keys"`
	RateLimit        *RateLimitConfig         `json:"rate_limit" yaml:"rate_limit" toml:"rate_limit"`
	IPFiltering      *IPFilterConfig          `json:"ip_filtering" yaml:"ip_filtering" toml:"ip_filtering"`
	
	// Encryption
	EncryptionKey    string        `json:"encryption_key" yaml:"encryption_key" toml:"encryption_key"`
	EncryptSecrets   bool          `json:"encrypt_secrets" yaml:"encrypt_secrets" toml:"encrypt_secrets"`
	
	// Audit
	AuditLog         bool          `json:"audit_log" yaml:"audit_log" toml:"audit_log"`
	AuditLogPath     string        `json:"audit_log_path" yaml:"audit_log_path" toml:"audit_log_path"`
	RetentionDays    int           `json:"retention_days" yaml:"retention_days" toml:"retention_days"`
}

// APIKeyConfig contains API key configuration
type APIKeyConfig struct {
	Name        string    `json:"name" yaml:"name" toml:"name"`
	Permissions []string  `json:"permissions" yaml:"permissions" toml:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" yaml:"expires_at,omitempty" toml:"expires_at,omitempty"`
	LastUsed    *time.Time `json:"last_used,omitempty" yaml:"last_used,omitempty" toml:"last_used,omitempty"`
	Enabled     bool      `json:"enabled" yaml:"enabled" toml:"enabled"`
}

// RateLimitConfig contains rate limiting configuration
type RateLimitConfig struct {
	Enabled     bool          `json:"enabled" yaml:"enabled" toml:"enabled"`
	Requests    int           `json:"requests" yaml:"requests" toml:"requests"`
	Window      time.Duration `json:"window" yaml:"window" toml:"window"`
	BurstSize   int           `json:"burst_size" yaml:"burst_size" toml:"burst_size"`
	
	// Per-endpoint limits
	Endpoints   map[string]*EndpointLimit `json:"endpoints" yaml:"endpoints" toml:"endpoints"`
}

// EndpointLimit contains endpoint-specific rate limits
type EndpointLimit struct {
	Requests  int           `json:"requests" yaml:"requests" toml:"requests"`
	Window    time.Duration `json:"window" yaml:"window" toml:"window"`
	BurstSize int           `json:"burst_size" yaml:"burst_size" toml:"burst_size"`
}

// IPFilterConfig contains IP filtering configuration
type IPFilterConfig struct {
	Enabled   bool     `json:"enabled" yaml:"enabled" toml:"enabled"`
	Mode      string   `json:"mode" yaml:"mode" toml:"mode"` // "whitelist" or "blacklist"
	Whitelist []string `json:"whitelist" yaml:"whitelist" toml:"whitelist"`
	Blacklist []string `json:"blacklist" yaml:"blacklist" toml:"blacklist"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level       string        `json:"level" yaml:"level" toml:"level"`
	Format      string        `json:"format" yaml:"format" toml:"format"`
	Output      string        `json:"output" yaml:"output" toml:"output"`
	File        *FileLogConfig `json:"file,omitempty" yaml:"file,omitempty" toml:"file,omitempty"`
	Syslog      *SyslogConfig `json:"syslog,omitempty" yaml:"syslog,omitempty" toml:"syslog,omitempty"`
	
	// Advanced settings
	EnableColors    bool          `json:"enable_colors" yaml:"enable_colors" toml:"enable_colors"`
	EnableTimestamp bool          `json:"enable_timestamp" yaml:"enable_timestamp" toml:"enable_timestamp"`
	EnableCaller    bool          `json:"enable_caller" yaml:"enable_caller" toml:"enable_caller"`
	SampleRate      float64       `json:"sample_rate" yaml:"sample_rate" toml:"sample_rate"`
	
	// Component-specific logging
	Components      map[string]*ComponentLogConfig `json:"components" yaml:"components" toml:"components"`
}

// FileLogConfig contains file logging configuration
type FileLogConfig struct {
	Path       string `json:"path" yaml:"path" toml:"path"`
	MaxSize    int    `json:"max_size" yaml:"max_size" toml:"max_size"`       // MB
	MaxAge     int    `json:"max_age" yaml:"max_age" toml:"max_age"`          // days
	MaxBackups int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`
	Compress   bool   `json:"compress" yaml:"compress" toml:"compress"`
}

// SyslogConfig contains syslog configuration
type SyslogConfig struct {
	Network  string `json:"network" yaml:"network" toml:"network"`
	Address  string `json:"address" yaml:"address" toml:"address"`
	Priority string `json:"priority" yaml:"priority" toml:"priority"`
	Tag      string `json:"tag" yaml:"tag" toml:"tag"`
}

// ComponentLogConfig contains component-specific logging configuration
type ComponentLogConfig struct {
	Level    string `json:"level" yaml:"level" toml:"level"`
	Enabled  bool   `json:"enabled" yaml:"enabled" toml:"enabled"`
	SampleRate float64 `json:"sample_rate" yaml:"sample_rate" toml:"sample_rate"`
}

// MonitoringConfig contains monitoring and metrics configuration
type MonitoringConfig struct {
	Enabled     bool              `json:"enabled" yaml:"enabled" toml:"enabled"`
	MetricsPort int               `json:"metrics_port" yaml:"metrics_port" toml:"metrics_port"`
	MetricsPath string            `json:"metrics_path" yaml:"metrics_path" toml:"metrics_path"`
	
	// Metrics collection
	CollectInterval time.Duration     `json:"collect_interval" yaml:"collect_interval" toml:"collect_interval"`
	RetentionPeriod time.Duration     `json:"retention_period" yaml:"retention_period" toml:"retention_period"`
	
	// Health checks
	HealthChecks    map[string]*HealthCheckConfig `json:"health_checks" yaml:"health_checks" toml:"health_checks"`
	
	// Alerting
	Alerting        *AlertingConfig   `json:"alerting" yaml:"alerting" toml:"alerting"`
	
	// Tracing
	Tracing         *TracingConfig    `json:"tracing" yaml:"tracing" toml:"tracing"`
}

// HealthCheckConfig contains health check configuration
type HealthCheckConfig struct {
	Enabled     bool          `json:"enabled" yaml:"enabled" toml:"enabled"`
	Interval    time.Duration `json:"interval" yaml:"interval" toml:"interval"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout" toml:"timeout"`
	Threshold   int           `json:"threshold" yaml:"threshold" toml:"threshold"`
	Endpoint    string        `json:"endpoint,omitempty" yaml:"endpoint,omitempty" toml:"endpoint,omitempty"`
}

// AlertingConfig contains alerting configuration
type AlertingConfig struct {
	Enabled     bool                      `json:"enabled" yaml:"enabled" toml:"enabled"`
	Rules       map[string]*AlertRule     `json:"rules" yaml:"rules" toml:"rules"`
	Channels    map[string]*AlertChannel  `json:"channels" yaml:"channels" toml:"channels"`
}

// AlertRule contains alert rule configuration
type AlertRule struct {
	Metric      string        `json:"metric" yaml:"metric" toml:"metric"`
	Condition   string        `json:"condition" yaml:"condition" toml:"condition"`
	Threshold   float64       `json:"threshold" yaml:"threshold" toml:"threshold"`
	Duration    time.Duration `json:"duration" yaml:"duration" toml:"duration"`
	Severity    string        `json:"severity" yaml:"severity" toml:"severity"`
	Message     string        `json:"message" yaml:"message" toml:"message"`
	Channels    []string      `json:"channels" yaml:"channels" toml:"channels"`
}

// AlertChannel contains alert channel configuration
type AlertChannel struct {
	Type     string            `json:"type" yaml:"type" toml:"type"`
	URL      string            `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	Token    string            `json:"token,omitempty" yaml:"token,omitempty" toml:"token,omitempty"`
	Settings map[string]string `json:"settings,omitempty" yaml:"settings,omitempty" toml:"settings,omitempty"`
}

// TracingConfig contains distributed tracing configuration
type TracingConfig struct {
	Enabled     bool    `json:"enabled" yaml:"enabled" toml:"enabled"`
	Endpoint    string  `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
	ServiceName string  `json:"service_name" yaml:"service_name" toml:"service_name"`
	SampleRate  float64 `json:"sample_rate" yaml:"sample_rate" toml:"sample_rate"`
}

// FeaturesConfig contains feature flag configuration
type FeaturesConfig struct {
	Flags map[string]*FeatureFlag `json:"flags" yaml:"flags" toml:"flags"`
}

// FeatureFlag contains feature flag configuration
type FeatureFlag struct {
	Enabled     bool              `json:"enabled" yaml:"enabled" toml:"enabled"`
	Description string            `json:"description" yaml:"description" toml:"description"`
	Rollout     *RolloutConfig    `json:"rollout,omitempty" yaml:"rollout,omitempty" toml:"rollout,omitempty"`
	Conditions  map[string]string `json:"conditions,omitempty" yaml:"conditions,omitempty" toml:"conditions,omitempty"`
}

// RolloutConfig contains feature rollout configuration
type RolloutConfig struct {
	Percentage int      `json:"percentage" yaml:"percentage" toml:"percentage"`
	UserGroups []string `json:"user_groups" yaml:"user_groups" toml:"user_groups"`
	StartDate  *time.Time `json:"start_date,omitempty" yaml:"start_date,omitempty" toml:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty" yaml:"end_date,omitempty" toml:"end_date,omitempty"`
}

// ConfigManager interface defines configuration management operations
type ConfigManager interface {
	// Loading and saving
	Load(path string) (*ApplicationConfig, error)
	Save(config *ApplicationConfig, path string) error
	
	// Validation
	Validate(config *ApplicationConfig) error
	ValidatePartial(section string, data interface{}) error
	
	// Migration
	Migrate(config *ApplicationConfig, targetVersion ConfigVersion) (*ApplicationConfig, error)
	GetMigrationPath(fromVersion, toVersion ConfigVersion) ([]ConfigVersion, error)
	
	// Merging and templating
	Merge(base, override *ApplicationConfig) (*ApplicationConfig, error)
	ApplyTemplate(config *ApplicationConfig, variables map[string]string) (*ApplicationConfig, error)
	
	// Environment handling
	LoadFromEnvironment(config *ApplicationConfig) error
	ExportToEnvironment(config *ApplicationConfig) map[string]string
	
	// Schema operations
	GetSchema(version ConfigVersion) (interface{}, error)
	GenerateExample(version ConfigVersion) (*ApplicationConfig, error)
	
	// Watch and reload
	Watch(path string, callback func(*ApplicationConfig)) error
	StopWatching(path string) error
}

// ConfigValidationError represents a configuration validation error
type ConfigValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
	Code    string `json:"code"`
}

func (e ConfigValidationError) Error() string {
	return fmt.Sprintf("config validation error in field '%s': %s", e.Field, e.Message)
}

// ConfigMigration represents a configuration migration
type ConfigMigration struct {
	FromVersion ConfigVersion
	ToVersion   ConfigVersion
	Description string
	Migrate     func(*ApplicationConfig) (*ApplicationConfig, error)
}

// ConfigTemplate represents a configuration template
type ConfigTemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Variables   []TemplateVariable     `json:"variables"`
	Template    *ApplicationConfig     `json:"template"`
}

// TemplateVariable represents a template variable
type TemplateVariable struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	Default      interface{} `json:"default,omitempty"`
	Required     bool        `json:"required"`
	Validation   string      `json:"validation,omitempty"`
}

// ConfigProfile represents a configuration profile
type ConfigProfile struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Environment string             `json:"environment"`
	Config      *ApplicationConfig `json:"config"`
	Active      bool               `json:"active"`
}

// ConfigBackup represents a configuration backup
type ConfigBackup struct {
	ID          string             `json:"id"`
	Path        string             `json:"path"`
	Config      *ApplicationConfig `json:"config"`
	CreatedAt   time.Time          `json:"created_at"`
	Description string             `json:"description"`
	Checksum    string             `json:"checksum"`
}

