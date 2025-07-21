/*
 * GitHubber - Advanced GitHub and Git CLI Tool
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: A powerful command-line interface for Git and GitHub operations
 */

package main

import (
    "fmt"
    "os"
    "github.com/ritankarsaha/git-tool/internal/cli"
    "github.com/ritankarsaha/git-tool/internal/git"
    "github.com/ritankarsaha/git-tool/internal/ui"
)

func main() {
    // Check if Git is installed
    if _, err := git.RunCommand("git --version"); err != nil {
        fmt.Println(ui.FormatError("Git is not installed or not in PATH"))
        os.Exit(1)
    }

    // Display beautiful header
    fmt.Println(ui.FormatTitle("GitHubber - Advanced Git & GitHub CLI"))
    fmt.Println(ui.FormatSubtitle("Created by Ritankar Saha <ritankar.saha786@gmail.com>"))

    // Check if we're in a git repository
    if repoInfo, err := git.GetRepositoryInfo(); err == nil {
        fmt.Println(ui.FormatRepoInfo(repoInfo.URL, repoInfo.CurrentBranch))
    } else {
        fmt.Println(ui.FormatError("Not in a Git repository"))
        os.Exit(1)
    }

    // Start the CLI menu
    cli.StartMenu()
}