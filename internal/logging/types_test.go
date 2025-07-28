package logging

import (
	"testing"
	"time"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
		str   string
	}{
		{"debug level", DebugLevel, "debug"},
		{"info level", InfoLevel, "info"}, 
		{"warn level", WarnLevel, "warn"},
		{"error level", ErrorLevel, "error"},
		{"fatal level", FatalLevel, "fatal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.level.String() != tt.str {
				t.Errorf("LogLevel.String() = %v, want %v", tt.level.String(), tt.str)
			}
		})
	}
}

func TestLogEntry(t *testing.T) {
	now := time.Now()
	entry := &LogEntry{
		Timestamp: now,
		Level:     InfoLevel,
		Message:   "test message",
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
		Component: "test-logger",
		Source: &LogSource{
			Function: "TestFunc",
			File:     "test.go",
			Line:     123,
		},
	}

	if entry.Level != InfoLevel {
		t.Errorf("Expected Level to be InfoLevel, got %v", entry.Level)
	}

	if entry.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %q", entry.Message)
	}

	if entry.Fields["key1"] != "value1" {
		t.Errorf("Expected Fields[key1] to be 'value1', got %v", entry.Fields["key1"])
	}

	if entry.Fields["key2"] != 42 {
		t.Errorf("Expected Fields[key2] to be 42, got %v", entry.Fields["key2"])
	}
}

func TestLogConfig(t *testing.T) {
	config := &LogConfig{
		Level:        InfoLevel,
		Format:       FormatJSON,
		EnableColors: true,
		EnableTime:   true,
		EnableCaller: false,
		Components:   make(map[string]*ComponentLogConfig),
		Outputs:      make(map[string]*LogOutputConfig),
		Hooks:        make([]LogHookConfig, 0),
		Filters:      make([]LogFilterConfig, 0),
	}

	if config.Level != InfoLevel {
		t.Errorf("Expected Level to be InfoLevel, got %v", config.Level)
	}

	if config.Format != FormatJSON {
		t.Errorf("Expected Format to be FormatJSON, got %q", config.Format)
	}

	if !config.EnableColors {
		t.Errorf("Expected EnableColors to be true, got %v", config.EnableColors)
	}

	if !config.EnableTime {
		t.Errorf("Expected EnableTime to be true, got %v", config.EnableTime)
	}
}

func TestCallerInfo(t *testing.T) {
	caller := CallerInfo{
		Function: "main.TestFunction",
		File:     "/path/to/file.go",
		Line:     42,
	}

	if caller.Function != "main.TestFunction" {
		t.Errorf("Expected Function to be 'main.TestFunction', got %q", caller.Function)
	}

	if caller.File != "/path/to/file.go" {
		t.Errorf("Expected File to be '/path/to/file.go', got %q", caller.File)
	}

	if caller.Line != 42 {
		t.Errorf("Expected Line to be 42, got %d", caller.Line)
	}
}

func TestLogFilter(t *testing.T) {
	filter := &LogFilter{
		Name:        "test-filter",
		Description: "A test filter",
		Enabled:     true,
	}

	if filter.Name != "test-filter" {
		t.Errorf("Expected Name to be 'test-filter', got %q", filter.Name)
	}

	if !filter.Enabled {
		t.Errorf("Expected Enabled to be true, got %v", filter.Enabled)
	}
}

func TestLogRotation(t *testing.T) {
	rotation := &LogRotation{
		Enabled:    true,
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 3,
		Compress:   true,
	}

	if !rotation.Enabled {
		t.Errorf("Expected Enabled to be true, got %v", rotation.Enabled)
	}

	if rotation.MaxSize != 100 {
		t.Errorf("Expected MaxSize to be 100, got %d", rotation.MaxSize)
	}

	if rotation.MaxAge != 7 {
		t.Errorf("Expected MaxAge to be 7, got %d", rotation.MaxAge)
	}
}