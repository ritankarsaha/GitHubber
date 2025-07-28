/*
 * GitHubber - Logging Manager Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Advanced logging system with structured logging and multiple outputs
 */

package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Manager implements LogManager
type Manager struct {
	mu           sync.RWMutex
	loggers      map[string]Logger
	config       *LogConfig
	outputs      map[string]LogOutput
	hooks        map[string]LogHook
	filters      map[string]LogFilter
	factories    *Factories
	metrics      *LogMetrics
	metricsEnabled bool
	started      bool
	startTime    time.Time
}

// Factories holds all registered factories
type Factories struct {
	Outputs map[string]LogOutputFactory
	Hooks   map[string]LogHookFactory
	Filters map[string]LogFilterFactory
}

// NewManager creates a new logging manager
func NewManager() *Manager {
	return &Manager{
		loggers:  make(map[string]Logger),
		outputs:  make(map[string]LogOutput),
		hooks:    make(map[string]LogHook),
		filters:  make(map[string]LogFilter),
		factories: &Factories{
			Outputs: make(map[string]LogOutputFactory),
			Hooks:   make(map[string]LogHookFactory),
			Filters: make(map[string]LogFilterFactory),
		},
		metrics: &LogMetrics{
			EntriesByLevel: make(map[LogLevel]int64),
			OutputMetrics:  make(map[string]*OutputMetrics),
			LastActivity:   time.Now(),
		},
		metricsEnabled: true,
	}
}

// DefaultManager provides a default logging manager instance
var DefaultManager = NewManager()

// GetLogger returns a logger by name
func (m *Manager) GetLogger(name string) Logger {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if logger, exists := m.loggers[name]; exists {
		return logger
	}

	// Create default logger if not found
	config := m.getDefaultLogConfig()
	logger, _ := m.CreateLogger(name, config)
	return logger
}

// CreateLogger creates a new logger with the given configuration
func (m *Manager) CreateLogger(name string, config *LogConfig) (Logger, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config == nil {
		config = m.getDefaultLogConfig()
	}

	// Create zap logger based on configuration
	zapConfig := m.buildZapConfig(config)
	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	logger := &ZapLogger{
		logger:    zapLogger.Named(name),
		name:      name,
		config:    config,
		manager:   m,
		fields:    make(map[string]interface{}),
		component: name,
	}

	m.loggers[name] = logger
	return logger, nil
}

// RemoveLogger removes a logger
func (m *Manager) RemoveLogger(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if logger, exists := m.loggers[name]; exists {
		if err := logger.Close(); err != nil {
			return fmt.Errorf("failed to close logger: %w", err)
		}
		delete(m.loggers, name)
	}

	return nil
}

// ListLoggers returns the names of all loggers
func (m *Manager) ListLoggers() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.loggers))
	for name := range m.loggers {
		names = append(names, name)
	}
	return names
}

// LoadConfig loads configuration from file
func (m *Manager) LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config LogConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return m.UpdateConfig(&config)
}

// SaveConfig saves configuration to file
func (m *Manager) SaveConfig(path string) error {
	m.mu.RLock()
	config := m.config
	m.mu.RUnlock()

	if config == nil {
		return fmt.Errorf("no configuration to save")
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// UpdateConfig updates the logging configuration
func (m *Manager) UpdateConfig(config *LogConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate configuration
	if err := m.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create outputs
	for name, outputConfig := range config.Outputs {
		if output, err := m.createOutputLocked(outputConfig); err != nil {
			return fmt.Errorf("failed to create output %s: %w", name, err)
		} else {
			if existing, exists := m.outputs[name]; exists {
				existing.Close()
			}
			m.outputs[name] = output
		}
	}

	// Create hooks
	for _, hookConfig := range config.Hooks {
		if hook, err := m.createHookLocked(&hookConfig); err != nil {
			return fmt.Errorf("failed to create hook %s: %w", hookConfig.Name, err)
		} else {
			m.hooks[hookConfig.Name] = hook
		}
	}

	// Create filters
	for _, filterConfig := range config.Filters {
		if filter, err := m.createFilterLocked(&filterConfig); err != nil {
			return fmt.Errorf("failed to create filter %s: %w", filterConfig.Name, err)
		} else {
			m.filters[filterConfig.Name] = filter
		}
	}

	m.config = config
	return nil
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *LogConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// RegisterOutput registers an output factory
func (m *Manager) RegisterOutput(name string, factory LogOutputFactory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.factories.Outputs[name] = factory
	return nil
}

// CreateOutput creates a log output
func (m *Manager) CreateOutput(config *LogOutputConfig) (LogOutput, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createOutputLocked(config)
}

// RegisterHook registers a hook factory
func (m *Manager) RegisterHook(name string, factory LogHookFactory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.factories.Hooks[name] = factory
	return nil
}

// CreateHook creates a log hook
func (m *Manager) CreateHook(config *LogHookConfig) (LogHook, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createHookLocked(config)
}

// RegisterFilter registers a filter factory
func (m *Manager) RegisterFilter(name string, factory LogFilterFactory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.factories.Filters[name] = factory
	return nil
}

// CreateFilter creates a log filter
func (m *Manager) CreateFilter(config *LogFilterConfig) (LogFilter, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createFilterLocked(config)
}

// GetMetrics returns logging metrics
func (m *Manager) GetMetrics() *LogMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Update uptime
	if m.started {
		m.metrics.Uptime = time.Since(m.startTime)
	}

	// Deep copy to avoid concurrent modifications
	metrics := &LogMetrics{
		EntriesTotal:   m.metrics.EntriesTotal,
		EntriesByLevel: make(map[LogLevel]int64),
		ErrorsTotal:    m.metrics.ErrorsTotal,
		DroppedTotal:   m.metrics.DroppedTotal,
		OutputMetrics:  make(map[string]*OutputMetrics),
		SampleRate:     m.metrics.SampleRate,
		LastActivity:   m.metrics.LastActivity,
		Uptime:         m.metrics.Uptime,
	}

	for level, count := range m.metrics.EntriesByLevel {
		metrics.EntriesByLevel[level] = count
	}

	for name, outputMetrics := range m.metrics.OutputMetrics {
		metrics.OutputMetrics[name] = &OutputMetrics{
			Name:         outputMetrics.Name,
			EntriesTotal: outputMetrics.EntriesTotal,
			ErrorsTotal:  outputMetrics.ErrorsTotal,
			BytesTotal:   outputMetrics.BytesTotal,
			LastWrite:    outputMetrics.LastWrite,
			IsHealthy:    outputMetrics.IsHealthy,
			LatencyP50:   outputMetrics.LatencyP50,
			LatencyP95:   outputMetrics.LatencyP95,
			LatencyP99:   outputMetrics.LatencyP99,
		}
	}

	return metrics
}

// EnableMetrics enables or disables metrics collection
func (m *Manager) EnableMetrics(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metricsEnabled = enabled
}

// Start starts the logging manager
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return fmt.Errorf("logging manager already started")
	}

	// Initialize built-in outputs, hooks, and filters
	m.registerBuiltinFactories()

	// Apply default configuration if none is set
	if m.config == nil {
		m.config = m.getDefaultLogConfig()
	}

	m.started = true
	m.startTime = time.Now()

	return nil
}

// Stop stops the logging manager
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.started {
		return nil
	}

	// Flush all loggers
	var errors []string
	for name, logger := range m.loggers {
		if err := logger.Flush(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", name, err))
		}
		if err := logger.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", name, err))
		}
	}

	// Close all outputs
	for name, output := range m.outputs {
		if err := output.Flush(); err != nil {
			errors = append(errors, fmt.Sprintf("output %s: %v", name, err))
		}
		if err := output.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("output %s: %v", name, err))
		}
	}

	m.started = false

	if len(errors) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errors)
	}

	return nil
}

// Flush flushes all loggers and outputs
func (m *Manager) Flush() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errors []string

	// Flush all loggers
	for name, logger := range m.loggers {
		if err := logger.Flush(); err != nil {
			errors = append(errors, fmt.Sprintf("logger %s: %v", name, err))
		}
	}

	// Flush all outputs
	for name, output := range m.outputs {
		if err := output.Flush(); err != nil {
			errors = append(errors, fmt.Sprintf("output %s: %v", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("flush errors: %v", errors)
	}

	return nil
}

// Helper methods

func (m *Manager) createOutputLocked(config *LogOutputConfig) (LogOutput, error) {
	factory, exists := m.factories.Outputs[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown output type: %s", config.Type)
	}

	return factory(config)
}

func (m *Manager) createHookLocked(config *LogHookConfig) (LogHook, error) {
	factory, exists := m.factories.Hooks[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown hook type: %s", config.Type)
	}

	return factory(config)
}

func (m *Manager) createFilterLocked(config *LogFilterConfig) (LogFilter, error) {
	factory, exists := m.factories.Filters[config.Type]
	if !exists {
		return nil, fmt.Errorf("unknown filter type: %s", config.Type)
	}

	return factory(config)
}

func (m *Manager) validateConfig(config *LogConfig) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate log level
	if config.Level == "" {
		config.Level = LevelInfo
	}

	// Validate outputs
	for name, outputConfig := range config.Outputs {
		if outputConfig.Name == "" {
			outputConfig.Name = name
		}
		if outputConfig.Type == "" {
			return fmt.Errorf("output %s: type is required", name)
		}
	}

	return nil
}

func (m *Manager) buildZapConfig(config *LogConfig) zap.Config {
	zapConfig := zap.NewProductionConfig()

	// Set log level
	switch config.Level {
	case LevelTrace:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel) // Zap doesn't have trace
	case LevelDebug:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case LevelInfo:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LevelWarn:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LevelError:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case LevelFatal:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	case LevelPanic:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	}

	// Set encoding
	switch config.Format {
	case FormatJSON:
		zapConfig.Encoding = "json"
	case FormatConsole:
		zapConfig.Encoding = "console"
	default:
		zapConfig.Encoding = "json"
	}

	// Configure encoder
	zapConfig.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if config.EnableCaller {
		zapConfig.Development = true
	}

	return zapConfig
}

func (m *Manager) getDefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:          LevelInfo,
		Format:         FormatJSON,
		EnableColors:   false,
		EnableCaller:   true,
		EnableTime:     true,
		TimeFormat:     time.RFC3339,
		ComponentField: "component",
		SampleRate:     1.0,
		Outputs: map[string]*LogOutputConfig{
			"console": {
				Name:    "console",
				Type:    "console",
				Level:   LevelInfo,
				Format:  FormatConsole,
				Enabled: true,
			},
		},
		Components: make(map[string]*ComponentLogConfig),
		Hooks:      make([]LogHookConfig, 0),
		Filters:    make([]LogFilterConfig, 0),
	}
}

func (m *Manager) registerBuiltinFactories() {
	// Register console output
	m.factories.Outputs["console"] = func(config *LogOutputConfig) (LogOutput, error) {
		return NewConsoleOutput(config)
	}

	// Register file output
	m.factories.Outputs["file"] = func(config *LogOutputConfig) (LogOutput, error) {
		return NewFileOutput(config)
	}

	// Register syslog output
	m.factories.Outputs["syslog"] = func(config *LogOutputConfig) (LogOutput, error) {
		return NewSyslogOutput(config)
	}

	// Register level filter
	m.factories.Filters["level"] = func(config *LogFilterConfig) (LogFilter, error) {
		return NewLevelFilter(config)
	}

	// Register sampling filter
	m.factories.Filters["sampling"] = func(config *LogFilterConfig) (LogFilter, error) {
		return NewSamplingFilter(config)
	}
}

func (m *Manager) recordMetrics(entry *LogEntry) {
	if !m.metricsEnabled {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.EntriesTotal++
	m.metrics.EntriesByLevel[entry.Level]++
	m.metrics.LastActivity = time.Now()
}

// ZapLogger wraps a zap logger to implement our Logger interface
type ZapLogger struct {
	logger    *zap.Logger
	name      string
	config    *LogConfig
	manager   *Manager
	fields    map[string]interface{}
	component string
	requestID string
	userID    string
}

// Implement Logger interface methods for ZapLogger
func (l *ZapLogger) Trace(msg string, fields ...Field) {
	l.log(LevelTrace, msg, fields...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal, msg, fields...)
}

func (l *ZapLogger) Panic(msg string, fields ...Field) {
	l.log(LevelPanic, msg, fields...)
}

func (l *ZapLogger) log(level LogLevel, msg string, fields ...Field) {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]interface{}),
		Component: l.component,
		RequestID: l.requestID,
		UserID:    l.userID,
	}

	// Add fields
	for key, value := range l.fields {
		entry.Fields[key] = value
	}

	for _, field := range fields {
		entry.Fields[field.Key] = field.Value
	}

	// Add source information if enabled
	if l.config.EnableCaller {
		if pc, file, line, ok := runtime.Caller(2); ok {
			entry.Source = &LogSource{
				File:     file,
				Line:     line,
				Function: runtime.FuncForPC(pc).Name(),
			}
		}
	}

	// Record metrics
	l.manager.recordMetrics(entry)

	// Convert to zap fields
	zapFields := make([]zap.Field, 0, len(entry.Fields)+3)
	
	if entry.Component != "" {
		zapFields = append(zapFields, zap.String("component", entry.Component))
	}
	if entry.RequestID != "" {
		zapFields = append(zapFields, zap.String("request_id", entry.RequestID))
	}
	if entry.UserID != "" {
		zapFields = append(zapFields, zap.String("user_id", entry.UserID))
	}

	for key, value := range entry.Fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	// Log with zap
	switch level {
	case LevelTrace, LevelDebug:
		l.logger.Debug(msg, zapFields...)
	case LevelInfo:
		l.logger.Info(msg, zapFields...)
	case LevelWarn:
		l.logger.Warn(msg, zapFields...)
	case LevelError:
		l.logger.Error(msg, zapFields...)
	case LevelFatal:
		l.logger.Fatal(msg, zapFields...)
	case LevelPanic:
		l.logger.Panic(msg, zapFields...)
	}
}

func (l *ZapLogger) TraceContext(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Trace(msg, fields...)
}

func (l *ZapLogger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Debug(msg, fields...)
}

func (l *ZapLogger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Info(msg, fields...)
}

func (l *ZapLogger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Warn(msg, fields...)
}

func (l *ZapLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Error(msg, fields...)
}

func (l *ZapLogger) Tracef(format string, args ...interface{}) {
	l.Trace(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.Panic(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) WithFields(fields ...Field) Logger {
	newFields := make(map[string]interface{})
	for key, value := range l.fields {
		newFields[key] = value
	}
	
	for _, field := range fields {
		newFields[field.Key] = field.Value
	}

	return &ZapLogger{
		logger:    l.logger,
		name:      l.name,
		config:    l.config,
		manager:   l.manager,
		fields:    newFields,
		component: l.component,
		requestID: l.requestID,
		userID:    l.userID,
	}
}

func (l *ZapLogger) WithComponent(component string) Logger {
	return &ZapLogger{
		logger:    l.logger,
		name:      l.name,
		config:    l.config,
		manager:   l.manager,
		fields:    l.fields,
		component: component,
		requestID: l.requestID,
		userID:    l.userID,
	}
}

func (l *ZapLogger) WithRequestID(requestID string) Logger {
	return &ZapLogger{
		logger:    l.logger,
		name:      l.name,
		config:    l.config,
		manager:   l.manager,
		fields:    l.fields,
		component: l.component,
		requestID: requestID,
		userID:    l.userID,
	}
}

func (l *ZapLogger) WithUserID(userID string) Logger {
	return &ZapLogger{
		logger:    l.logger,
		name:      l.name,
		config:    l.config,
		manager:   l.manager,
		fields:    l.fields,
		component: l.component,
		requestID: l.requestID,
		userID:    userID,
	}
}

func (l *ZapLogger) SetLevel(level LogLevel) {
	// This would require reconfiguring the zap logger
	// For now, we'll just update the config
	l.config.Level = level
}

func (l *ZapLogger) GetLevel() LogLevel {
	return l.config.Level
}

func (l *ZapLogger) AddOutput(output LogOutput) error {
	// This would require reconfiguring the zap logger
	return fmt.Errorf("dynamic output addition not supported yet")
}

func (l *ZapLogger) RemoveOutput(name string) error {
	// This would require reconfiguring the zap logger
	return fmt.Errorf("dynamic output removal not supported yet")
}

func (l *ZapLogger) Flush() error {
	return l.logger.Sync()
}

func (l *ZapLogger) Close() error {
	return l.logger.Sync()
}

func (l *ZapLogger) withContext(ctx context.Context) Logger {
	logger := &ZapLogger{
		logger:    l.logger,
		name:      l.name,
		config:    l.config,
		manager:   l.manager,
		fields:    l.fields,
		component: l.component,
		requestID: l.requestID,
		userID:    l.userID,
	}

	// Extract values from context
	if requestID := ctx.Value(ContextKeyRequestID); requestID != nil {
		if id, ok := requestID.(string); ok {
			logger.requestID = id
		}
	}

	if userID := ctx.Value(ContextKeyUserID); userID != nil {
		if id, ok := userID.(string); ok {
			logger.userID = id
		}
	}

	if component := ctx.Value(ContextKeyComponent); component != nil {
		if comp, ok := component.(string); ok {
			logger.component = comp
		}
	}

	return logger
}

// Global convenience functions
func GetLogger(name string) Logger {
	return DefaultManager.GetLogger(name)
}

func Start() error {
	return DefaultManager.Start()
}

func Stop() error {
	return DefaultManager.Stop()
}

func Flush() error {
	return DefaultManager.Flush()
}