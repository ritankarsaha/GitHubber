/*
 * GitHubber - Configuration Management
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Configuration file management and user preferences
 */

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	GitHub GitHubConfig `json:"github"`
	UI     UIConfig     `json:"ui"`
	Git    GitConfig    `json:"git"`
}

type GitHubConfig struct {
	Token        string `json:"token,omitempty"`        // GitHub personal access token
	DefaultOwner string `json:"default_owner"`          // Default repository owner
	DefaultRepo  string `json:"default_repo"`           // Default repository name
	APIBaseURL   string `json:"api_base_url,omitempty"` // For GitHub Enterprise
}

type UIConfig struct {
	Theme       string `json:"theme"`        // Color theme (dark, light, auto)
	ShowEmojis  bool   `json:"show_emojis"`  // Whether to show emojis in output
	PageSize    int    `json:"page_size"`    // Number of items to show per page
	BorderStyle string `json:"border_style"` // Border style (rounded, square, double)
}

type GitConfig struct {
	DefaultBranch string `json:"default_branch"` // Default branch name for new repos
	AutoPush      bool   `json:"auto_push"`      // Automatically push commits
	SignCommits   bool   `json:"sign_commits"`   // Sign commits with GPG
}

const (
	configFileName = "githubber.json"
	configDirName  = ".githubber"
)

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, configDirName)
	return filepath.Join(configDir, configFileName), nil
}

// Load loads the configuration from file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Fill in any missing fields with defaults
	fillDefaults(&config)

	return &config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
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
	}
}

// fillDefaults fills in any missing configuration fields with defaults
func fillDefaults(config *Config) {
	defaults := GetDefaultConfig()

	if config.GitHub.APIBaseURL == "" {
		config.GitHub.APIBaseURL = defaults.GitHub.APIBaseURL
	}
	if config.UI.Theme == "" {
		config.UI.Theme = defaults.UI.Theme
	}
	if config.UI.PageSize == 0 {
		config.UI.PageSize = defaults.UI.PageSize
	}
	if config.UI.BorderStyle == "" {
		config.UI.BorderStyle = defaults.UI.BorderStyle
	}
	if config.Git.DefaultBranch == "" {
		config.Git.DefaultBranch = defaults.Git.DefaultBranch
	}
}

// SetGitHubToken sets the GitHub token in the configuration
func (c *Config) SetGitHubToken(token string) error {
	c.GitHub.Token = token
	return c.Save()
}

// GetGitHubToken gets the GitHub token from config or environment
func (c *Config) GetGitHubToken() string {
	// First check environment variable
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}
	
	// Then check config file
	return c.GitHub.Token
}

// IsConfigured checks if the basic configuration is set up
func (c *Config) IsConfigured() bool {
	return c.GetGitHubToken() != ""
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.UI.PageSize <= 0 {
		return fmt.Errorf("page_size must be greater than 0")
	}
	
	validThemes := []string{"dark", "light", "auto"}
	validTheme := false
	for _, theme := range validThemes {
		if c.UI.Theme == theme {
			validTheme = true
			break
		}
	}
	if !validTheme {
		return fmt.Errorf("invalid theme: %s (valid themes: %v)", c.UI.Theme, validThemes)
	}
	
	return nil
}