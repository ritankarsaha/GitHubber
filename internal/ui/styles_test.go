package ui

import (
	"strings"
	"testing"
)

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple title",
			input:    "Test Title",
			expected: "Test Title",
		},
		{
			name:     "empty title",
			input:    "",
			expected: "",
		},
		{
			name:     "title with special characters",
			input:    "Test & Title!",
			expected: "Test & Title!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTitle(tt.input)
			// The result should contain the input text (ignoring styling codes)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatTitle(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
		})
	}
}

func TestFormatSubtitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple subtitle",
			input:    "Test Subtitle",
			expected: "Test Subtitle",
		},
		{
			name:     "empty subtitle",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSubtitle(tt.input)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatSubtitle(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
		})
	}
}

func TestFormatSuccess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "success message",
			input:    "Operation completed",
			expected: "Operation completed",
		},
		{
			name:     "empty message",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSuccess(tt.input)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatSuccess(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
			// Should contain success icon
			if !strings.Contains(result, IconSuccess) && tt.expected != "" {
				t.Errorf("FormatSuccess(%q) should contain success icon, got %q", tt.input, result)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "error message",
			input:    "Something went wrong",
			expected: "Something went wrong",
		},
		{
			name:     "empty message",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatError(tt.input)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatError(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
			// Should contain error icon
			if !strings.Contains(result, IconError) && tt.expected != "" {
				t.Errorf("FormatError(%q) should contain error icon, got %q", tt.input, result)
			}
		})
	}
}

func TestFormatWarning(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "warning message",
			input:    "This is a warning",
			expected: "This is a warning",
		},
		{
			name:     "empty message",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatWarning(tt.input)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatWarning(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
			// Should contain warning icon
			if !strings.Contains(result, IconWarning) && tt.expected != "" {
				t.Errorf("FormatWarning(%q) should contain warning icon, got %q", tt.input, result)
			}
		})
	}
}

func TestFormatInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "info message",
			input:    "Information here",
			expected: "Information here",
		},
		{
			name:     "empty message",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatInfo(tt.input)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatInfo(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
			// Should contain info icon
			if !strings.Contains(result, IconInfo) && tt.expected != "" {
				t.Errorf("FormatInfo(%q) should contain info icon, got %q", tt.input, result)
			}
		})
	}
}

func TestFormatRepoInfo(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		branch   string
		expected []string // Multiple expected substrings
	}{
		{
			name:     "github repository",
			url:      "https://github.com/user/repo.git",
			branch:   "main",
			expected: []string{"github.com/user/repo", "main"},
		},
		{
			name:     "gitlab repository",
			url:      "https://gitlab.com/user/repo.git",
			branch:   "develop",
			expected: []string{"gitlab.com/user/repo", "develop"},
		},
		{
			name:     "empty values",
			url:      "",
			branch:   "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatRepoInfo(tt.url, tt.branch)
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("FormatRepoInfo(%q, %q) should contain %q, got %q", tt.url, tt.branch, expected, result)
				}
			}
		})
	}
}

func TestFormatBox(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple content",
			content:  "test content",
			expected: "test content",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBox(tt.content)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatBox(%q) should contain %q, got %q", tt.content, tt.expected, result)
			}
		})
	}
}

func TestFormatCode(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple code",
			content:  "git status",
			expected: "git status",
		},
		{
			name:     "multi-line code",
			content:  "git add .\ngit commit",
			expected: "git add .",
		},
		{
			name:     "empty code",
			content:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCode(tt.content)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatCode(%q) should contain %q, got %q", tt.content, tt.expected, result)
			}
		})
	}
}

func TestFormatMenuItem(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		text     string
		expected []string
	}{
		{
			name:     "normal menu item",
			number:   1,
			text:     "Create Branch",
			expected: []string{"1", "Create Branch"},
		},
		{
			name:     "higher number",
			number:   10,
			text:     "Settings",
			expected: []string{"10", "Settings"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMenuItem(tt.number, tt.text)
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("FormatMenuItem(%d, %q) should contain %q, got %q", tt.number, tt.text, expected, result)
				}
			}
		})
	}
}

func TestStyles(t *testing.T) {
	// Test that basic styles are defined (they are lipgloss styles, so we just verify they exist)
	_ = BaseStyle
	_ = TitleStyle  
	_ = SubtitleStyle
	_ = MenuHeaderStyle

	// We can't easily test the actual style output, but we can verify they don't panic
	testText := "test"
	_ = TitleStyle.Render(testText)
	_ = SubtitleStyle.Render(testText)
}

func TestPromptFormatting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple prompt",
			input:    "Enter value:",
			expected: "Enter value:",
		},
		{
			name:     "empty prompt",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPrompt(tt.input)
			if !strings.Contains(result, tt.expected) && tt.expected != "" {
				t.Errorf("FormatPrompt(%q) should contain %q, got %q", tt.input, tt.expected, result)
			}
		})
	}
}