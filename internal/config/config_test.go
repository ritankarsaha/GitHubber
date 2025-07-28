package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "githubber-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .githubber directory
	configDir := filepath.Join(tmpDir, ".githubber")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, "githubber.json")

	tests := []struct {
		name           string
		setupConfig    func() error
		expectedConfig *Config
		expectError    bool
	}{
		{
			name: "load existing config",
			setupConfig: func() error {
				config := &Config{
					GitHub: GitHubConfig{
						Token:        "test-token",
						DefaultOwner: "testuser",
						DefaultRepo:  "testrepo",
						APIBaseURL:   "https://api.github.com",
					},
					UI: UIConfig{
						Theme:       "dark",
						ShowEmojis:  true,
						PageSize:    20,
						BorderStyle: "rounded",
					},
					Git: GitConfig{
						DefaultBranch: "main",
						AutoPush:      false,
						SignCommits:   false,
					},
				}
				
				data, err := json.MarshalIndent(config, "", "  ")
				if err != nil {
					return err
				}
				
				return os.WriteFile(configFile, data, 0644)
			},
			expectedConfig: &Config{
				GitHub: GitHubConfig{
					Token:        "test-token",
					DefaultOwner: "testuser",
					DefaultRepo:  "testrepo",
					APIBaseURL:   "https://api.github.com",
				},
				UI: UIConfig{
					Theme:       "dark",
					ShowEmojis:  true,
					PageSize:    20,
					BorderStyle: "rounded",
				},
				Git: GitConfig{
					DefaultBranch: "main",
					AutoPush:      false,
					SignCommits:   false,
				},
			},
			expectError: false,
		},
		{
			name: "no config file - create default",
			setupConfig: func() error {
				// Don't create any file
				return nil
			},
			expectedConfig: &Config{
				GitHub: GitHubConfig{
					APIBaseURL: "https://api.github.com",
				},
				UI: UIConfig{
					Theme:       "dark",
					ShowEmojis:  true,
					PageSize:    20,
					BorderStyle: "rounded",
				},
				Git: GitConfig{
					DefaultBranch: "main",
					AutoPush:      false,
					SignCommits:   false,
				},
			},
			expectError: false,
		},
		{
			name: "invalid json config",
			setupConfig: func() error {
				return os.WriteFile(configFile, []byte("{invalid json}"), 0644)
			},
			expectedConfig: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up temporary home directory
			originalHome := os.Getenv("HOME")
			os.Setenv("HOME", tmpDir)
			defer os.Setenv("HOME", originalHome)

			// Clean up any existing config file
			os.Remove(configFile)

			// Setup test config
			if err := tt.setupConfig(); err != nil {
				t.Fatalf("Failed to setup config: %v", err)
			}

			// Test Load function
			config, err := Load()

			if tt.expectError && err == nil {
				t.Errorf("Expected error, but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.expectError && !reflect.DeepEqual(config, tt.expectedConfig) {
				t.Errorf("LoadConfig() = %+v, want %+v", config, tt.expectedConfig)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "githubber-config-save-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up temporary home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	config := &Config{
		GitHub: GitHubConfig{
			Token:        "test-token",
			DefaultOwner: "testuser",
			DefaultRepo:  "testrepo",
			APIBaseURL:   "https://api.github.com",
		},
		UI: UIConfig{
			Theme:       "light",
			ShowEmojis:  false,
			PageSize:    30,
			BorderStyle: "square",
		},
		Git: GitConfig{
			DefaultBranch: "develop",
			AutoPush:      true,
			SignCommits:   true,
		},
	}

	// Test Save method
	if err := config.Save(); err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Verify file was created
	configFile := filepath.Join(tmpDir, ".githubber", "githubber.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
		return
	}

	// Load and verify the saved config
	loadedConfig, err := Load()
	if err != nil {
		t.Errorf("Failed to load saved config: %v", err)
		return
	}

	if !reflect.DeepEqual(loadedConfig, config) {
		t.Errorf("Saved and loaded config don't match. Got %+v, want %+v", loadedConfig, config)
	}
}

func TestGetGitHubToken(t *testing.T) {
	tests := []struct {
		name        string
		envToken    string
		configToken string
		expected    string
	}{
		{
			name:        "environment token takes precedence",
			envToken:    "env-token",
			configToken: "config-token",
			expected:    "env-token",
		},
		{
			name:        "config token when no env",
			envToken:    "",
			configToken: "config-token",
			expected:    "config-token",
		},
		{
			name:        "no token available",
			envToken:    "",
			configToken: "",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			originalToken := os.Getenv("GITHUB_TOKEN")
			if tt.envToken != "" {
				os.Setenv("GITHUB_TOKEN", tt.envToken)
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}
			defer func() {
				if originalToken != "" {
					os.Setenv("GITHUB_TOKEN", originalToken)
				} else {
					os.Unsetenv("GITHUB_TOKEN")
				}
			}()

			// Create config with token
			config := &Config{
				GitHub: GitHubConfig{
					Token: tt.configToken,
				},
			}

			result := config.GetGitHubToken()
			if result != tt.expected {
				t.Errorf("GetGitHubToken() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	// Verify default values
	if config.GitHub.APIBaseURL != "https://api.github.com" {
		t.Errorf("Default GitHub API base URL = %q, want %q", config.GitHub.APIBaseURL, "https://api.github.com")
	}

	if config.UI.Theme != "dark" {
		t.Errorf("Default UI theme = %q, want %q", config.UI.Theme, "dark")
	}

	if config.UI.ShowEmojis != true {
		t.Errorf("Default UI ShowEmojis = %v, want %v", config.UI.ShowEmojis, true)
	}

	if config.UI.PageSize != 20 {
		t.Errorf("Default UI PageSize = %d, want %d", config.UI.PageSize, 20)
	}

	if config.UI.BorderStyle != "rounded" {
		t.Errorf("Default UI BorderStyle = %q, want %q", config.UI.BorderStyle, "rounded")
	}

	if config.Git.DefaultBranch != "main" {
		t.Errorf("Default Git DefaultBranch = %q, want %q", config.Git.DefaultBranch, "main")
	}

	if config.Git.AutoPush != false {
		t.Errorf("Default Git AutoPush = %v, want %v", config.Git.AutoPush, false)
	}

	if config.Git.SignCommits != false {
		t.Errorf("Default Git SignCommits = %v, want %v", config.Git.SignCommits, false)
	}
}

func TestConfigPaths(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "githubber-config-path-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set up temporary home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Test that config directory is created when saving
	config := GetDefaultConfig()
	if err := config.Save(); err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// Verify directory structure
	configDir := filepath.Join(tmpDir, ".githubber")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created")
	}

	configFile := filepath.Join(configDir, "githubber.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("Config file was not created")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  GetDefaultConfig(),
			wantErr: false,
		},
		{
			name: "invalid theme",
			config: &Config{
				UI: UIConfig{
					Theme: "invalid-theme",
				},
			},
			wantErr: false, // Currently no validation, but structure is ready
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Currently there's no ValidateConfig function,
			// but this test structure is ready for when validation is added
			_ = tt.config
			if tt.wantErr {
				// When validation is added, test for errors here
			}
		})
	}
}