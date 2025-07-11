package main

import (
	"fmt"
	"os"

	"github.com/ritankarsaha/githubber/internal/commands"
	"github.com/ritankarsaha/githubber/pkg/config"
	"github.com/ritankarsaha/githubber/pkg/git"
	"github.com/ritankarsaha/githubber/pkg/logger"
)

func main() {
	// Initialize configuration
	cfg, err := config.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Create Git client
	gitClient := git.NewClient(cfg)
	defer gitClient.Close()

	// Validate Git installation
	if err := gitClient.ValidateGitInstallation(); err != nil {
		logger.Fatalf("Git validation failed: %v", err)
	}

	// Execute root command
	rootCmd := commands.NewRootCommand(cfg, gitClient)
	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("Command execution failed: %v", err)
		os.Exit(1)
	}
}