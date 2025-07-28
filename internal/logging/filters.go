/*
 * GitHubber - Logging Filters Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Implementation of log filters for the logging system
 */

package logging

import (
	"fmt"
	"math/rand"
	"strings"
)

// LevelFilter filters log entries based on their level
type LevelFilter struct {
	Name       string
	Enabled    bool
	MinLevel   LogLevel
	MaxLevel   LogLevel
	AllowList  []LogLevel
	DenyList   []LogLevel
}

// NewLevelFilter creates a new level filter
func NewLevelFilter(config *LogFilterConfig) (LogFilter, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	filter := &LevelFilter{
		Name:    config.Name,
		Enabled: config.Enabled,
	}

	// Parse settings
	if config.Rules != nil && len(config.Rules) > 0 {
		for _, rule := range config.Rules {
			if rule.Field == "min_level" {
				if level, ok := rule.Value.(string); ok {
					filter.MinLevel = LogLevelFromString(level)
				}
			}
			if rule.Field == "max_level" {
				if level, ok := rule.Value.(string); ok {
					filter.MaxLevel = LogLevelFromString(level)
				}
			}
		}
	}

	return filter, nil
}

// Apply applies the level filter to a log entry
func (f *LevelFilter) Apply(entry *LogEntry) (*LogEntry, bool, error) {
	if !f.Enabled {
		return entry, true, nil
	}

	// Check deny list first
	for _, denyLevel := range f.DenyList {
		if entry.Level == denyLevel {
			return nil, false, nil // Drop entry
		}
	}

	// Check allow list if set
	if len(f.AllowList) > 0 {
		allowed := false
		for _, allowLevel := range f.AllowList {
			if entry.Level == allowLevel {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, false, nil // Drop entry
		}
	}

	// Check min/max levels
	if f.MinLevel != "" && !entry.Level.IsEnabledFor(f.MinLevel) {
		return nil, false, nil // Drop entry
	}

	return entry, true, nil
}

// IsEnabled returns whether the filter is enabled
func (f *LevelFilter) IsEnabled() bool {
	return f.Enabled
}

// SetEnabled sets the enabled state of the filter
func (f *LevelFilter) SetEnabled(enabled bool) {
	f.Enabled = enabled
}

// SamplingFilter implements sampling-based log filtering
type SamplingFilter struct {
	Name       string
	Enabled    bool
	SampleRate float64
	KeyField   string
	MaxSamples int
	samples    map[string]int
}

// NewSamplingFilter creates a new sampling filter
func NewSamplingFilter(config *LogFilterConfig) (LogFilter, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	filter := &SamplingFilter{
		Name:       config.Name,
		Enabled:    config.Enabled,
		SampleRate: 1.0, // Default to no sampling
		samples:    make(map[string]int),
	}

	// Parse settings from rules
	if config.Rules != nil && len(config.Rules) > 0 {
		for _, rule := range config.Rules {
			if rule.Field == "sample_rate" {
				if rate, ok := rule.Value.(float64); ok {
					filter.SampleRate = rate
				}
			}
			if rule.Field == "key_field" {
				if field, ok := rule.Value.(string); ok {
					filter.KeyField = field
				}
			}
			if rule.Field == "max_samples" {
				if max, ok := rule.Value.(int); ok {
					filter.MaxSamples = max
				}
			}
		}
	}

	return filter, nil
}

// Apply applies sampling to a log entry
func (f *SamplingFilter) Apply(entry *LogEntry) (*LogEntry, bool, error) {
	if !f.Enabled {
		return entry, true, nil
	}

	// Simple random sampling if no key field specified
	if f.KeyField == "" {
		if rand.Float64() > f.SampleRate {
			return nil, false, nil // Drop entry
		}
		return entry, true, nil
	}

	// Key-based sampling
	var key string
	if value, exists := entry.Fields[f.KeyField]; exists {
		key = fmt.Sprintf("%v", value)
	} else {
		key = "default"
	}

	if f.MaxSamples > 0 && f.samples[key] >= f.MaxSamples {
		return nil, false, nil // Drop entry - max samples reached
	}

	if rand.Float64() > f.SampleRate {
		return nil, false, nil // Drop entry
	}

	f.samples[key]++
	return entry, true, nil
}

// IsEnabled returns whether the filter is enabled
func (f *SamplingFilter) IsEnabled() bool {
	return f.Enabled
}

// SetEnabled sets the enabled state of the filter
func (f *SamplingFilter) SetEnabled(enabled bool) {
	f.Enabled = enabled
}

// ComponentFilter filters log entries based on component names
type ComponentFilter struct {
	Name         string
	Enabled      bool
	AllowedComponents []string
	DeniedComponents  []string
}

// Apply applies the component filter to a log entry
func (f *ComponentFilter) Apply(entry *LogEntry) (*LogEntry, bool, error) {
	if !f.Enabled {
		return entry, true, nil
	}

	component := entry.Component

	// Check denied components first
	for _, denied := range f.DeniedComponents {
		if strings.Contains(component, denied) {
			return nil, false, nil // Drop entry
		}
	}

	// Check allowed components if set
	if len(f.AllowedComponents) > 0 {
		allowed := false
		for _, allowedComp := range f.AllowedComponents {
			if strings.Contains(component, allowedComp) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, false, nil // Drop entry
		}
	}

	return entry, true, nil
}

// IsEnabled returns whether the filter is enabled
func (f *ComponentFilter) IsEnabled() bool {
	return f.Enabled
}

// SetEnabled sets the enabled state of the filter
func (f *ComponentFilter) SetEnabled(enabled bool) {
	f.Enabled = enabled
}

// FieldFilter filters log entries based on field values
type FieldFilter struct {
	Name     string
	Enabled  bool
	Rules    []FilterRule
	Action   FilterAction
}

// Apply applies the field filter to a log entry
func (f *FieldFilter) Apply(entry *LogEntry) (*LogEntry, bool, error) {
	if !f.Enabled {
		return entry, true, nil
	}

	for _, rule := range f.Rules {
		matches, err := f.evaluateRule(entry, rule)
		if err != nil {
			return entry, true, err
		}

		if matches {
			switch f.Action {
			case ActionDrop:
				return nil, false, nil
			case ActionAllow:
				return entry, true, nil
			case ActionModify:
				// Implement modification logic here
				return entry, true, nil
			default:
				return entry, true, nil
			}
		}
	}

	return entry, true, nil
}

func (f *FieldFilter) evaluateRule(entry *LogEntry, rule FilterRule) (bool, error) {
	var fieldValue interface{}
	var exists bool

	switch rule.Field {
	case "level":
		fieldValue = string(entry.Level)
		exists = true
	case "message":
		fieldValue = entry.Message
		exists = true
	case "component":
		fieldValue = entry.Component
		exists = true
	default:
		fieldValue, exists = entry.Fields[rule.Field]
	}

	if !exists {
		return false, nil
	}

	return f.compareValues(fieldValue, rule.Value, rule.Operator, rule.CaseSensitive)
}

func (f *FieldFilter) compareValues(actual, expected interface{}, operator string, caseSensitive bool) (bool, error) {
	switch operator {
	case "eq", "==":
		return f.isEqual(actual, expected, caseSensitive), nil
	case "ne", "!=":
		return !f.isEqual(actual, expected, caseSensitive), nil
	case "contains":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		if !caseSensitive {
			actualStr = strings.ToLower(actualStr)
			expectedStr = strings.ToLower(expectedStr)
		}
		return strings.Contains(actualStr, expectedStr), nil
	case "starts_with":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		if !caseSensitive {
			actualStr = strings.ToLower(actualStr)
			expectedStr = strings.ToLower(expectedStr)
		}
		return strings.HasPrefix(actualStr, expectedStr), nil
	case "ends_with":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		if !caseSensitive {
			actualStr = strings.ToLower(actualStr)
			expectedStr = strings.ToLower(expectedStr)
		}
		return strings.HasSuffix(actualStr, expectedStr), nil
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

func (f *FieldFilter) isEqual(actual, expected interface{}, caseSensitive bool) bool {
	if !caseSensitive {
		actualStr := strings.ToLower(fmt.Sprintf("%v", actual))
		expectedStr := strings.ToLower(fmt.Sprintf("%v", expected))
		return actualStr == expectedStr
	}
	return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected)
}

// IsEnabled returns whether the filter is enabled
func (f *FieldFilter) IsEnabled() bool {
	return f.Enabled
}

// SetEnabled sets the enabled state of the filter
func (f *FieldFilter) SetEnabled(enabled bool) {
	f.Enabled = enabled
}