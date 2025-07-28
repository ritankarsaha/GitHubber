/*
 * GitHubber - Git Squash Operations
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Advanced commit squashing functionality
 */

package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CommitInfo struct {
	Hash    string
	Message string
}

// GetRecentCommits returns the last n commits
func GetRecentCommits(n int) ([]CommitInfo, error) {
	output, err := RunCommand(fmt.Sprintf("git log -%d --oneline", n))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")
	commits := make([]CommitInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			commits = append(commits, CommitInfo{
				Hash:    parts[0],
				Message: parts[1],
			})
		}
	}

	return commits, nil
}

// SquashCommits performs the squash operation
func SquashCommits(baseCommit, message string) error {
	// Verify working directory is clean
	if clean, err := IsWorkingDirectoryClean(); err != nil || !clean {
		return fmt.Errorf("working directory must be clean before squashing")
	}

	// Create temporary directory for scripts
	tmpDir, err := os.MkdirTemp("", "git-squash-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create editor script
	editorScript := filepath.Join(tmpDir, "editor.sh")
	// Use a more portable approach that works on both macOS and Linux
	editorContent := `#!/bin/sh
# Portable sed command that works on both macOS and Linux
if [ "$(uname)" = "Darwin" ]; then
    sed -i '' -e '2,$s/pick/squash/' "$1"
else
    sed -i -e '2,$s/pick/squash/' "$1"
fi
`
	if err := os.WriteFile(editorScript, []byte(editorContent), 0755); err != nil {
		return fmt.Errorf("failed to create editor script: %w", err)
	}

	// Set up the environment for the rebase
	os.Setenv("GIT_SEQUENCE_EDITOR", editorScript)
	os.Setenv("GIT_EDITOR", "true")

	// Start the interactive rebase
	if _, err := RunCommand(fmt.Sprintf("git rebase -i %s~1", baseCommit)); err != nil {
		// Attempt to abort the rebase if it fails
		RunCommand("git rebase --abort")
		return fmt.Errorf("rebase failed: %w", err)
	}

	// Set the final commit message
	if _, err := RunCommand(fmt.Sprintf("git commit --amend -m \"%s\"", message)); err != nil {
		return fmt.Errorf("failed to set commit message: %w", err)
	}

	return nil
}
