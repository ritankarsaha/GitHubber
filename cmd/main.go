package main

import (
    "fmt"
    "os"
    "github.com/ritankarsaha/git-tool/internal/cli"
    "github.com/ritankarsaha/git-tool/internal/git"
)

func main() {
    // Check if Git is installed
    if _, err := git.RunCommand("git --version"); err != nil {
        fmt.Println("Error: Git is not installed or not in PATH")
        os.Exit(1)
    }

    fmt.Println("🛠 Git CLI Tool")
    fmt.Println("---------------")

    // Check if we're in a git repository
    if repoInfo, err := git.GetRepositoryInfo(); err == nil {
        fmt.Printf("📂 Repository: %s\n", repoInfo.URL)
        fmt.Printf("🌿 Current Branch: %s\n", repoInfo.CurrentBranch)
    } else {
        fmt.Println("❌ Error: Not in a Git repository")
        os.Exit(1)
    }

    // Start the CLI menu
    cli.StartMenu()
}