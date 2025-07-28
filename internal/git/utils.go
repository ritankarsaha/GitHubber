/*
 * GitHubber - Git Utilities
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Git utility functions and helpers
 */

package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type RepositoryInfo struct {
	URL           string
	CurrentBranch string
}

// RunCommand executes a git command and returns its output
func RunCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// GetRepositoryInfo retrieves current repository information
func GetRepositoryInfo() (*RepositoryInfo, error) {
	// Get remote URL
	url, err := RunCommand("git remote get-url origin")
	if err != nil {
		return nil, fmt.Errorf("failed to get repository URL: %w", err)
	}

	// Get current branch
	branch, err := RunCommand("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	return &RepositoryInfo{
		URL:           url,
		CurrentBranch: branch,
	}, nil
}

// IsWorkingDirectoryClean checks if there are any uncommitted changes
func IsWorkingDirectoryClean() (bool, error) {
	output, err := RunCommand("git status --porcelain")
	if err != nil {
		return false, err
	}
	return output == "", nil
}
