/*
 * GitHubber - Logging Types and Interfaces
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Comprehensive logging system with structured logging and multiple outputs
 */

package logging

import (
	"context"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelTrace LogLevel = "trace"
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
	LevelPanic LogLevel = "panic"
	
	// Aliases for backward compatibility
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// LogFormat represents the output format of log entries
type LogFormat string

const (
	FormatJSON LogFormat = "json"
	FormatText LogFormat = "text"
	FormatConsole LogFormat = "console"
)

// Logger interface defines the logging contract
type Logger interface {
	// Basic logging methods
	Trace(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)

	// Context-aware logging
	TraceContext(ctx context.Context, msg string, fields ...Field)
	DebugContext(ctx context.Context, msg string, fields ...Field)
	InfoContext(ctx context.Context, msg string, fields ...Field)
	WarnContext(ctx context.Context, msg string, fields ...Field)
	ErrorContext(ctx context.Context, msg string, fields ...Field)

	// Formatted logging
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	// Logger configuration
	WithFields(fields ...Field) Logger
	WithComponent(component string) Logger
	WithRequestID(requestID string) Logger
	WithUserID(userID string) Logger
	SetLevel(level LogLevel)
	GetLevel() LogLevel
	
	// Output control
	AddOutput(output LogOutput) error
	RemoveOutput(name string) error
	
	// Lifecycle
	Flush() error
	Close() error
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
	Type  FieldType
}

// FieldType represents the type of a log field
type FieldType int

const (
	StringType FieldType = iota
	IntType
	Int64Type
	Float64Type
	BoolType
	TimeType
	DurationType
	ErrorType
	ObjectType
	ArrayType
)

// LogEntry represents a complete log entry
type LogEntry struct {
	Timestamp   time.Time          `json:"timestamp"`
	Level       LogLevel           `json:"level"`
	Message     string             `json:"message"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
	Component   string             `json:"component,omitempty"`
	RequestID   string             `json:"request_id,omitempty"`
	UserID      string             `json:"user_id,omitempty"`
	TraceID     string             `json:"trace_id,omitempty"`
	SpanID      string             `json:"span_id,omitempty"`
	Source      *LogSource         `json:"source,omitempty"`
	Error       *LogError          `json:"error,omitempty"`
	Duration    *time.Duration     `json:"duration,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// LogSource represents the source location of a log entry
type LogSource struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// CallerInfo is an alias for LogSource for backward compatibility
type CallerInfo struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// LogError represents error information in a log entry
type LogError struct {
	Message    string            `json:"message"`
	Type       string            `json:"type"`
	StackTrace string            `json:"stack_trace,omitempty"`
	Cause      string            `json:"cause,omitempty"`
	Code       string            `json:"code,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// LogOutput represents a log output destination
type LogOutput interface {
	Write(entry *LogEntry) error
	GetName() string
	IsEnabled() bool
	SetEnabled(enabled bool)
	Flush() error
	Close() error
}

// LogOutputConfig represents configuration for a log output
type LogOutputConfig struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Level    LogLevel          `json:"level"`
	Format   LogFormat         `json:"format"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
}

// FileOutputConfig represents file output configuration
type FileOutputConfig struct {
	Path       string `json:"path"`
	MaxSize    int    `json:"max_size"`    // MB
	MaxAge     int    `json:"max_age"`     // days
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
	LocalTime  bool   `json:"local_time"`
}

// SyslogOutputConfig represents syslog output configuration
type SyslogOutputConfig struct {
	Network   string `json:"network"`
	Address   string `json:"address"`
	Priority  string `json:"priority"`
	Tag       string `json:"tag"`
	Facility  string `json:"facility"`
	Hostname  string `json:"hostname"`
}

// ElasticsearchOutputConfig represents Elasticsearch output configuration
type ElasticsearchOutputConfig struct {
	URLs      []string `json:"urls"`
	Index     string   `json:"index"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	BatchSize int      `json:"batch_size"`
	Timeout   time.Duration `json:"timeout"`
}

// LogConfig represents the overall logging configuration
type LogConfig struct {
	Level          LogLevel                    `json:"level"`
	Format         LogFormat                   `json:"format"`
	EnableColors   bool                        `json:"enable_colors"`
	EnableCaller   bool                        `json:"enable_caller"`
	EnableTime     bool                        `json:"enable_time"`
	TimeFormat     string                      `json:"time_format"`
	ComponentField string                      `json:"component_field"`
	SampleRate     float64                     `json:"sample_rate"`
	
	// Output configurations
	Outputs        map[string]*LogOutputConfig `json:"outputs"`
	
	// Component-specific configurations
	Components     map[string]*ComponentLogConfig `json:"components"`
	
	// Hooks and filters
	Hooks          []LogHookConfig             `json:"hooks"`
	Filters        []LogFilterConfig           `json:"filters"`
}

// ComponentLogConfig represents component-specific logging configuration
type ComponentLogConfig struct {
	Level      LogLevel `json:"level"`
	Enabled    bool     `json:"enabled"`
	SampleRate float64  `json:"sample_rate"`
	Fields     map[string]interface{} `json:"fields"`
}

// LogHookConfig represents a log hook configuration
type LogHookConfig struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Levels   []LogLevel        `json:"levels"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
}

// LogFilterConfig represents a log filter configuration
type LogFilterConfig struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Rules    []FilterRule      `json:"rules"`
	Action   FilterAction      `json:"action"`
	Enabled  bool              `json:"enabled"`
}

// FilterRule represents a log filter rule
type FilterRule struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
	CaseSensitive bool    `json:"case_sensitive"`
}

// FilterAction represents the action to take when a filter matches
type FilterAction string

const (
	ActionDrop     FilterAction = "drop"
	ActionAllow    FilterAction = "allow"
	ActionModify   FilterAction = "modify"
	ActionRedirect FilterAction = "redirect"
)

// LogHook represents a log hook for processing log entries
type LogHook interface {
	Fire(entry *LogEntry) error
	GetLevels() []LogLevel
	IsEnabled() bool
	SetEnabled(enabled bool)
}

// LogFilter represents a log filter for processing log entries
type LogFilter interface {
	Apply(entry *LogEntry) (*LogEntry, bool, error)
	IsEnabled() bool
	SetEnabled(enabled bool)
}

// LogManager manages multiple loggers and their configuration
type LogManager interface {
	// Logger management
	GetLogger(name string) Logger
	CreateLogger(name string, config *LogConfig) (Logger, error)
	RemoveLogger(name string) error
	ListLoggers() []string
	
	// Configuration management
	LoadConfig(path string) error
	SaveConfig(path string) error
	UpdateConfig(config *LogConfig) error
	GetConfig() *LogConfig
	
	// Output management
	RegisterOutput(name string, factory LogOutputFactory) error
	CreateOutput(config *LogOutputConfig) (LogOutput, error)
	
	// Hook management
	RegisterHook(name string, factory LogHookFactory) error
	CreateHook(config *LogHookConfig) (LogHook, error)
	
	// Filter management
	RegisterFilter(name string, factory LogFilterFactory) error
	CreateFilter(config *LogFilterConfig) (LogFilter, error)
	
	// Metrics and monitoring
	GetMetrics() *LogMetrics
	EnableMetrics(enabled bool)
	
	// Lifecycle
	Start() error
	Stop() error
	Flush() error
}

// LogOutputFactory creates log output instances
type LogOutputFactory func(config *LogOutputConfig) (LogOutput, error)

// LogHookFactory creates log hook instances
type LogHookFactory func(config *LogHookConfig) (LogHook, error)

// LogFilterFactory creates log filter instances
type LogFilterFactory func(config *LogFilterConfig) (LogFilter, error)

// LogMetrics represents logging system metrics
type LogMetrics struct {
	EntriesTotal    int64             `json:"entries_total"`
	EntriesByLevel  map[LogLevel]int64 `json:"entries_by_level"`
	ErrorsTotal     int64             `json:"errors_total"`
	DroppedTotal    int64             `json:"dropped_total"`
	OutputMetrics   map[string]*OutputMetrics `json:"output_metrics"`
	SampleRate      float64           `json:"sample_rate"`
	LastActivity    time.Time         `json:"last_activity"`
	Uptime          time.Duration     `json:"uptime"`
}

// OutputMetrics represents metrics for a specific log output
type OutputMetrics struct {
	Name         string        `json:"name"`
	EntriesTotal int64         `json:"entries_total"`
	ErrorsTotal  int64         `json:"errors_total"`
	BytesTotal   int64         `json:"bytes_total"`
	LastWrite    time.Time     `json:"last_write"`
	IsHealthy    bool          `json:"is_healthy"`
	LatencyP50   time.Duration `json:"latency_p50"`
	LatencyP95   time.Duration `json:"latency_p95"`
	LatencyP99   time.Duration `json:"latency_p99"`
}

// Audit represents audit logging functionality
type AuditLogger interface {
	LogAccess(userID, resource, action string, success bool, details map[string]interface{})
	LogAuthentication(userID, method string, success bool, ip string)
	LogConfigChange(userID, component string, oldValue, newValue interface{})
	LogSecurityEvent(eventType, description string, severity string, details map[string]interface{})
	LogDataAccess(userID, resource string, operation string, recordCount int)
	LogError(component string, error error, context map[string]interface{})
}

// ContextKeys for logging context
type ContextKey string

const (
	ContextKeyRequestID ContextKey = "request_id"
	ContextKeyUserID    ContextKey = "user_id"
	ContextKeyTraceID   ContextKey = "trace_id"
	ContextKeySpanID    ContextKey = "span_id"
	ContextKeyComponent ContextKey = "component"
)

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value, Type: StringType}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value, Type: IntType}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value, Type: Int64Type}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value, Type: Float64Type}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value, Type: BoolType}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value, Type: TimeType}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value, Type: DurationType}
}

func Error(key string, err error) Field {
	return Field{Key: key, Value: err, Type: ErrorType}
}

func Object(key string, value interface{}) Field {
	return Field{Key: key, Value: value, Type: ObjectType}
}

func Array(key string, value interface{}) Field {
	return Field{Key: key, Value: value, Type: ArrayType}
}

// LogLevelFromString converts string to LogLevel
func LogLevelFromString(s string) LogLevel {
	switch s {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	case "panic":
		return LevelPanic
	default:
		return LevelInfo
	}
}

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	return string(l)
}

// IsEnabledFor checks if the current level is enabled for the given level
func (l LogLevel) IsEnabledFor(target LogLevel) bool {
	levels := []LogLevel{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic}
	
	currentIndex := -1
	targetIndex := -1
	
	for i, level := range levels {
		if level == l {
			currentIndex = i
		}
		if level == target {
			targetIndex = i
		}
	}
	
	return currentIndex != -1 && targetIndex != -1 && targetIndex >= currentIndex
}