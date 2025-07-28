/*
 * GitHubber - Log Output Implementations
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Various log output implementations (console, file, syslog, etc.)
 */

package logging

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// ConsoleOutput writes log entries to console
type ConsoleOutput struct {
	name      string
	config    *LogOutputConfig
	formatter LogFormatter
	enabled   bool
	mu        sync.RWMutex
}

// NewConsoleOutput creates a new console output
func NewConsoleOutput(config *LogOutputConfig) (LogOutput, error) {
	formatter, err := NewFormatter(config.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to create formatter: %w", err)
	}

	return &ConsoleOutput{
		name:      config.Name,
		config:    config,
		formatter: formatter,
		enabled:   config.Enabled,
	}, nil
}

func (c *ConsoleOutput) Write(entry *LogEntry) error {
	if !c.IsEnabled() {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check level
	if !c.config.Level.IsEnabledFor(entry.Level) {
		return nil
	}

	formatted, err := c.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to format entry: %w", err)
	}

	_, err = fmt.Print(formatted)
	return err
}

func (c *ConsoleOutput) GetName() string {
	return c.name
}

func (c *ConsoleOutput) IsEnabled() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.enabled
}

func (c *ConsoleOutput) SetEnabled(enabled bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enabled = enabled
}

func (c *ConsoleOutput) Flush() error {
	return nil // Console output doesn't need flushing
}

func (c *ConsoleOutput) Close() error {
	return nil // Console output doesn't need closing
}

// FileOutput writes log entries to file with rotation
type FileOutput struct {
	name      string
	config    *LogOutputConfig
	formatter LogFormatter
	writer    *lumberjack.Logger
	enabled   bool
	mu        sync.RWMutex
}

// NewFileOutput creates a new file output
func NewFileOutput(config *LogOutputConfig) (LogOutput, error) {
	formatter, err := NewFormatter(config.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to create formatter: %w", err)
	}

	// Parse file-specific settings
	var fileConfig FileOutputConfig
	if config.Settings != nil {
		settingsBytes, _ := json.Marshal(config.Settings)
		json.Unmarshal(settingsBytes, &fileConfig)
	}

	// Set defaults
	if fileConfig.Path == "" {
		fileConfig.Path = "logs/app.log"
	}
	if fileConfig.MaxSize == 0 {
		fileConfig.MaxSize = 100 // 100MB
	}
	if fileConfig.MaxAge == 0 {
		fileConfig.MaxAge = 7 // 7 days
	}
	if fileConfig.MaxBackups == 0 {
		fileConfig.MaxBackups = 3
	}

	writer := &lumberjack.Logger{
		Filename:   fileConfig.Path,
		MaxSize:    fileConfig.MaxSize,
		MaxAge:     fileConfig.MaxAge,
		MaxBackups: fileConfig.MaxBackups,
		Compress:   fileConfig.Compress,
		LocalTime:  fileConfig.LocalTime,
	}

	return &FileOutput{
		name:      config.Name,
		config:    config,
		formatter: formatter,
		writer:    writer,
		enabled:   config.Enabled,
	}, nil
}

func (f *FileOutput) Write(entry *LogEntry) error {
	if !f.IsEnabled() {
		return nil
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Check level
	if !f.config.Level.IsEnabledFor(entry.Level) {
		return nil
	}

	formatted, err := f.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to format entry: %w", err)
	}

	_, err = f.writer.Write([]byte(formatted))
	return err
}

func (f *FileOutput) GetName() string {
	return f.name
}

func (f *FileOutput) IsEnabled() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.enabled
}

func (f *FileOutput) SetEnabled(enabled bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.enabled = enabled
}

func (f *FileOutput) Flush() error {
	// lumberjack doesn't have a flush method, so we don't need to do anything
	return nil
}

func (f *FileOutput) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.writer.Close()
}

// SyslogOutput writes log entries to syslog
type SyslogOutput struct {
	name      string
	config    *LogOutputConfig
	formatter LogFormatter
	writer    *syslog.Writer
	enabled   bool
	mu        sync.RWMutex
}

// NewSyslogOutput creates a new syslog output
func NewSyslogOutput(config *LogOutputConfig) (LogOutput, error) {
	formatter, err := NewFormatter(config.Format)
	if err != nil {
		return nil, fmt.Errorf("failed to create formatter: %w", err)
	}

	// Parse syslog-specific settings
	var syslogConfig SyslogOutputConfig
	if config.Settings != nil {
		settingsBytes, _ := json.Marshal(config.Settings)
		json.Unmarshal(settingsBytes, &syslogConfig)
	}

	// Set defaults
	if syslogConfig.Network == "" {
		syslogConfig.Network = ""
	}
	if syslogConfig.Address == "" {
		syslogConfig.Address = ""
	}
	if syslogConfig.Tag == "" {
		syslogConfig.Tag = "githubber"
	}

	// Parse priority
	priority := syslog.LOG_INFO | syslog.LOG_LOCAL0
	if syslogConfig.Priority != "" {
		// Parse priority string (implementation would convert string to syslog.Priority)
	}

	writer, err := syslog.Dial(syslogConfig.Network, syslogConfig.Address, priority, syslogConfig.Tag)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to syslog: %w", err)
	}

	return &SyslogOutput{
		name:      config.Name,
		config:    config,
		formatter: formatter,
		writer:    writer,
		enabled:   config.Enabled,
	}, nil
}

func (s *SyslogOutput) Write(entry *LogEntry) error {
	if !s.IsEnabled() {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check level
	if !s.config.Level.IsEnabledFor(entry.Level) {
		return nil
	}

	formatted, err := s.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to format entry: %w", err)
	}

	// Write to appropriate syslog level
	switch entry.Level {
	case LevelTrace, LevelDebug:
		return s.writer.Debug(formatted)
	case LevelInfo:
		return s.writer.Info(formatted)
	case LevelWarn:
		return s.writer.Warning(formatted)
	case LevelError:
		return s.writer.Err(formatted)
	case LevelFatal:
		return s.writer.Crit(formatted)
	case LevelPanic:
		return s.writer.Emerg(formatted)
	default:
		return s.writer.Info(formatted)
	}
}

func (s *SyslogOutput) GetName() string {
	return s.name
}

func (s *SyslogOutput) IsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.enabled
}

func (s *SyslogOutput) SetEnabled(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.enabled = enabled
}

func (s *SyslogOutput) Flush() error {
	// Syslog doesn't need explicit flushing
	return nil
}

func (s *SyslogOutput) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.writer.Close()
}

// LogFormatter interface for formatting log entries
type LogFormatter interface {
	Format(entry *LogEntry) (string, error)
}

// JSONFormatter formats log entries as JSON
type JSONFormatter struct{}

func (j *JSONFormatter) Format(entry *LogEntry) (string, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// TextFormatter formats log entries as plain text
type TextFormatter struct {
	EnableColors bool
	TimeFormat   string
}

func (t *TextFormatter) Format(entry *LogEntry) (string, error) {
	timeStr := entry.Timestamp.Format(t.TimeFormat)
	if t.TimeFormat == "" {
		timeStr = entry.Timestamp.Format(time.RFC3339)
	}

	levelStr := string(entry.Level)
	if t.EnableColors {
		levelStr = t.colorizeLevel(entry.Level)
	}

	var fieldsStr string
	if len(entry.Fields) > 0 {
		fieldsBytes, _ := json.Marshal(entry.Fields)
		fieldsStr = " " + string(fieldsBytes)
	}

	component := ""
	if entry.Component != "" {
		component = fmt.Sprintf("[%s] ", entry.Component)
	}

	return fmt.Sprintf("%s [%s] %s%s%s\n", 
		timeStr, levelStr, component, entry.Message, fieldsStr), nil
}

func (t *TextFormatter) colorizeLevel(level LogLevel) string {
	switch level {
	case LevelTrace:
		return "\033[37mTRACE\033[0m" // White
	case LevelDebug:
		return "\033[36mDEBUG\033[0m" // Cyan
	case LevelInfo:
		return "\033[32mINFO\033[0m"  // Green
	case LevelWarn:
		return "\033[33mWARN\033[0m"  // Yellow
	case LevelError:
		return "\033[31mERROR\033[0m" // Red
	case LevelFatal:
		return "\033[35mFATAL\033[0m" // Magenta
	case LevelPanic:
		return "\033[41mPANIC\033[0m" // Red background
	default:
		return string(level)
	}
}

// ConsoleFormatter formats log entries for console display
type ConsoleFormatter struct {
	EnableColors bool
	TimeFormat   string
}

func (c *ConsoleFormatter) Format(entry *LogEntry) (string, error) {
	timeStr := entry.Timestamp.Format(c.TimeFormat)
	if c.TimeFormat == "" {
		timeStr = entry.Timestamp.Format("15:04:05")
	}

	levelStr := string(entry.Level)
	if c.EnableColors {
		levelStr = c.colorizeLevel(entry.Level)
	}

	component := ""
	if entry.Component != "" {
		component = fmt.Sprintf("[%s] ", entry.Component)
	}

	message := entry.Message

	// Add important fields inline
	var inlineFields []string
	if entry.RequestID != "" {
		inlineFields = append(inlineFields, fmt.Sprintf("req=%s", entry.RequestID))
	}
	if entry.UserID != "" {
		inlineFields = append(inlineFields, fmt.Sprintf("user=%s", entry.UserID))
	}

	if len(inlineFields) > 0 {
		message += fmt.Sprintf(" (%s)", strings.Join(inlineFields, " "))
	}

	// Add remaining fields as key=value pairs
	var extraFields []string
	for key, value := range entry.Fields {
		if key != "request_id" && key != "user_id" {
			extraFields = append(extraFields, fmt.Sprintf("%s=%v", key, value))
		}
	}

	if len(extraFields) > 0 {
		message += fmt.Sprintf(" %s", strings.Join(extraFields, " "))
	}

	return fmt.Sprintf("%s %s %s%s\n", 
		timeStr, levelStr, component, message), nil
}

func (c *ConsoleFormatter) colorizeLevel(level LogLevel) string {
	switch level {
	case LevelTrace:
		return "\033[37mTRC\033[0m" // White
	case LevelDebug:
		return "\033[36mDBG\033[0m" // Cyan
	case LevelInfo:
		return "\033[32mINF\033[0m" // Green
	case LevelWarn:
		return "\033[33mWRN\033[0m" // Yellow
	case LevelError:
		return "\033[31mERR\033[0m" // Red
	case LevelFatal:
		return "\033[35mFTL\033[0m" // Magenta
	case LevelPanic:
		return "\033[41mPNC\033[0m" // Red background
	default:
		return strings.ToUpper(string(level))[:3]
	}
}

// NewFormatter creates a formatter based on the format type
func NewFormatter(format LogFormat) (LogFormatter, error) {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatText:
		return &TextFormatter{
			EnableColors: false,
			TimeFormat:   time.RFC3339,
		}, nil
	case FormatConsole:
		return &ConsoleFormatter{
			EnableColors: true,
			TimeFormat:   "15:04:05",
		}, nil
	default:
		return &JSONFormatter{}, nil
	}
}

// BufferedOutput wraps an output with buffering
type BufferedOutput struct {
	output    LogOutput
	buffer    []*LogEntry
	batchSize int
	ticker    *time.Ticker
	done      chan struct{}
	mu        sync.Mutex
}

// NewBufferedOutput creates a new buffered output
func NewBufferedOutput(output LogOutput, batchSize int, flushInterval time.Duration) *BufferedOutput {
	bo := &BufferedOutput{
		output:    output,
		buffer:    make([]*LogEntry, 0, batchSize),
		batchSize: batchSize,
		ticker:    time.NewTicker(flushInterval),
		done:      make(chan struct{}),
	}

	// Start background flushing
	go bo.flushLoop()

	return bo
}

func (b *BufferedOutput) Write(entry *LogEntry) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buffer = append(b.buffer, entry)
	
	if len(b.buffer) >= b.batchSize {
		return b.flushLocked()
	}

	return nil
}

func (b *BufferedOutput) GetName() string {
	return b.output.GetName()
}

func (b *BufferedOutput) IsEnabled() bool {
	return b.output.IsEnabled()
}

func (b *BufferedOutput) SetEnabled(enabled bool) {
	b.output.SetEnabled(enabled)
}

func (b *BufferedOutput) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.flushLocked()
}

func (b *BufferedOutput) Close() error {
	close(b.done)
	b.ticker.Stop()
	
	if err := b.Flush(); err != nil {
		return err
	}
	
	return b.output.Close()
}

func (b *BufferedOutput) flushLoop() {
	for {
		select {
		case <-b.ticker.C:
			b.Flush()
		case <-b.done:
			return
		}
	}
}

func (b *BufferedOutput) flushLocked() error {
	if len(b.buffer) == 0 {
		return nil
	}

	var err error
	for _, entry := range b.buffer {
		if writeErr := b.output.Write(entry); writeErr != nil {
			err = writeErr // Keep last error
		}
	}

	b.buffer = b.buffer[:0] // Clear buffer
	return err
}

// MultiOutput writes to multiple outputs
type MultiOutput struct {
	name    string
	outputs []LogOutput
	enabled bool
	mu      sync.RWMutex
}

// NewMultiOutput creates a new multi-output
func NewMultiOutput(name string, outputs ...LogOutput) *MultiOutput {
	return &MultiOutput{
		name:    name,
		outputs: outputs,
		enabled: true,
	}
}

func (m *MultiOutput) Write(entry *LogEntry) error {
	if !m.IsEnabled() {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for _, output := range m.outputs {
		if err := output.Write(entry); err != nil {
			lastErr = err // Keep last error
		}
	}

	return lastErr
}

func (m *MultiOutput) GetName() string {
	return m.name
}

func (m *MultiOutput) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

func (m *MultiOutput) SetEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

func (m *MultiOutput) Flush() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for _, output := range m.outputs {
		if err := output.Flush(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (m *MultiOutput) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for _, output := range m.outputs {
		if err := output.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (m *MultiOutput) AddOutput(output LogOutput) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.outputs = append(m.outputs, output)
}

func (m *MultiOutput) RemoveOutput(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, output := range m.outputs {
		if output.GetName() == name {
			m.outputs = append(m.outputs[:i], m.outputs[i+1:]...)
			break
		}
	}
}

