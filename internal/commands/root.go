package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ritankarsaha/githubber/pkg/config"
	"github.com/ritankarsaha/githubber/pkg/git"
	"github.com/ritankarsaha/githubber/pkg/logger"
)

// NewRootCommand creates the root command for the application
func NewRootCommand(cfg *config.Config, gitClient *git.Client) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "githubber",
		Short: "A comprehensive Git CLI tool with advanced functionality",
		Long: `GitHubber is a powerful Git CLI tool that provides comprehensive
Git functionality with an intuitive interface. It supports everything from
basic repository operations to advanced workflows including rebasing,
stashing, tagging, and remote management.`,
		Version: cfg.App.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set debug mode if flag is provided
			if cfg.App.Debug {
				logger.GetLogger().SetLevel(logger.GetLogger().GetLevel())
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&cfg.App.Debug, "debug", cfg.App.Debug, "Enable debug mode")
	rootCmd.PersistentFlags().StringVar(&cfg.Logging.Level, "log-level", cfg.Logging.Level, "Set log level (debug, info, warn, error)")

	// Add basic subcommands
	rootCmd.AddCommand(NewInitCommand(cfg, gitClient))
	rootCmd.AddCommand(NewCloneCommand(cfg, gitClient))
	rootCmd.AddCommand(NewStatusCommand(cfg, gitClient))

	// Set completion
	rootCmd.CompletionOptions.DisableDefaultCmd = false

	return rootCmd
}

// Execute runs the root command
func Execute(cfg *config.Config, gitClient *git.Client) error {
	rootCmd := NewRootCommand(cfg, gitClient)
	return rootCmd.Execute()
}

// Version information
func printVersion(cfg *config.Config) {
	fmt.Printf("GitHubber v%s\n", cfg.App.Version)
	fmt.Printf("A comprehensive Git CLI tool\n")
}