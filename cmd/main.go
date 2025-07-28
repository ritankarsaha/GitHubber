/*
 * GitHubber - Advanced GitHub and Git CLI Tool
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: A powerful command-line interface for Git and GitHub operations
 */

package main

import (
	"fmt"
	"github.com/ritankarsaha/git-tool/internal/cli"
	"github.com/ritankarsaha/git-tool/internal/git"
	"github.com/ritankarsaha/git-tool/internal/ui"
	"os"
)

func main() {
	// Check if Git is installed
	if _, err := git.RunCommand("git --version"); err != nil {
		fmt.Println(ui.FormatError("Git is not installed or not in PATH"))
		os.Exit(1)
	}

	// Parse command line arguments
	args := os.Args[1:] // Exclude program name

	// If arguments are provided, execute command directly
	if len(args) > 0 {
		if err := cli.ParseAndExecute(args); err != nil {
			fmt.Fprintf(os.Stderr, "%s Error: %v\n", ui.IconError, err)
			os.Exit(1)
		}
		return
	}

	// Interactive mode - display beautiful header
	fmt.Println(ui.FormatTitle("GitHubber - Advanced Git & GitHub CLI"))
	fmt.Println(ui.FormatSubtitle("Created by Ritankar Saha <ritankar.saha786@gmail.com>"))

	// Check if we're in a git repository (only for interactive mode)
	if repoInfo, err := git.GetRepositoryInfo(); err == nil {
		fmt.Println(ui.FormatRepoInfo(repoInfo.URL, repoInfo.CurrentBranch))
	} else {
		fmt.Println(ui.FormatWarning("Not in a Git repository - limited functionality available"))
	}

	// Start the interactive CLI menu
	cli.StartMenu()
}
