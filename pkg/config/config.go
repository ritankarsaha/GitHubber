package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Git      GitConfig      `mapstructure:"git"`
	GitHub   GitHubConfig   `mapstructure:"github"`
	UI       UIConfig       `mapstructure:"ui"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	Debug        bool   `mapstructure:"debug"`
	ConfigDir    string `mapstructure:"config_dir"`
	DataDir      string `mapstructure:"data_dir"`
}

type GitConfig struct {
	DefaultRemote   string `mapstructure:"default_remote"`
	DefaultBranch   string `mapstructure:"default_branch"`
	AutoPush        bool   `mapstructure:"auto_push"`
	SignCommits     bool   `mapstructure:"sign_commits"`
	PrettyFormat    string `mapstructure:"pretty_format"`
}

type GitHubConfig struct {
	Token        string `mapstructure:"token"`
	DefaultOrg   string `mapstructure:"default_org"`
	DefaultRepo  string `mapstructure:"default_repo"`
	APIBaseURL   string `mapstructure:"api_base_url"`
}

type UIConfig struct {
	Theme           string `mapstructure:"theme"`
	ShowIcons       bool   `mapstructure:"show_icons"`
	ShowSpinner     bool   `mapstructure:"show_spinner"`
	ConfirmActions  bool   `mapstructure:"confirm_actions"`
	PageSize        int    `mapstructure:"page_size"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputFile string `mapstructure:"output_file"`
}

var defaultConfig = Config{
	App: AppConfig{
		Name:    "GitHubber",
		Version: "2.0.0",
		Debug:   false,
	},
	Git: GitConfig{
		DefaultRemote: "origin",
		DefaultBranch: "main",
		AutoPush:      false,
		SignCommits:   false,
		PrettyFormat:  "%h %s (%an, %ar)",
	},
	GitHub: GitHubConfig{
		APIBaseURL: "https://api.github.com",
	},
	UI: UIConfig{
		Theme:          "default",
		ShowIcons:      true,
		ShowSpinner:    true,
		ConfirmActions: true,
		PageSize:       10,
	},
	Logging: LoggingConfig{
		Level:  "info",
		Format: "text",
	},
}

func Init() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// Add config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.githubber")
	viper.AddConfigPath("/etc/githubber")
	
	// Set environment variable prefix
	viper.SetEnvPrefix("GITHUBBER")
	viper.AutomaticEnv()
	
	// Set defaults
	setDefaults()
	
	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			if err := createDefaultConfig(); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Set runtime directories
	if err := setupDirectories(&config); err != nil {
		return nil, fmt.Errorf("failed to setup directories: %w", err)
	}
	
	return &config, nil
}

func setDefaults() {
	viper.SetDefault("app.name", defaultConfig.App.Name)
	viper.SetDefault("app.version", defaultConfig.App.Version)
	viper.SetDefault("app.debug", defaultConfig.App.Debug)
	
	viper.SetDefault("git.default_remote", defaultConfig.Git.DefaultRemote)
	viper.SetDefault("git.default_branch", defaultConfig.Git.DefaultBranch)
	viper.SetDefault("git.auto_push", defaultConfig.Git.AutoPush)
	viper.SetDefault("git.sign_commits", defaultConfig.Git.SignCommits)
	viper.SetDefault("git.pretty_format", defaultConfig.Git.PrettyFormat)
	
	viper.SetDefault("github.api_base_url", defaultConfig.GitHub.APIBaseURL)
	
	viper.SetDefault("ui.theme", defaultConfig.UI.Theme)
	viper.SetDefault("ui.show_icons", defaultConfig.UI.ShowIcons)
	viper.SetDefault("ui.show_spinner", defaultConfig.UI.ShowSpinner)
	viper.SetDefault("ui.confirm_actions", defaultConfig.UI.ConfirmActions)
	viper.SetDefault("ui.page_size", defaultConfig.UI.PageSize)
	
	viper.SetDefault("logging.level", defaultConfig.Logging.Level)
	viper.SetDefault("logging.format", defaultConfig.Logging.Format)
}

func createDefaultConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configDir := filepath.Join(homeDir, ".githubber")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	configFile := filepath.Join(configDir, "config.yaml")
	viper.SetConfigFile(configFile)
	
	return viper.WriteConfig()
}

func setupDirectories(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configDir := filepath.Join(homeDir, ".githubber")
	dataDir := filepath.Join(homeDir, ".githubber", "data")
	
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	
	config.App.ConfigDir = configDir
	config.App.DataDir = dataDir
	
	if config.Logging.OutputFile == "" {
		config.Logging.OutputFile = filepath.Join(dataDir, "githubber.log")
	}
	
	return nil
}

func (c *Config) Save() error {
	return viper.WriteConfig()
}