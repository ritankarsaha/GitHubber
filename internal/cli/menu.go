package cli

import (
    "fmt"
    "os"
    "github.com/ritankarsaha/git-tool/internal/git"
)

func StartMenu() {
    for {
        fmt.Println("\n📋 Available Commands:")
        fmt.Println("1. Squash Commits")
        fmt.Println("2. View Recent Commits")
        fmt.Println("3. Exit")

        choice := GetInput("Enter your choice (1-3): ")

        switch choice {
        case "1":
            handleSquash()
        case "2":
            handleViewCommits()
        case "3":
            fmt.Println("👋 Goodbye!")
            os.Exit(0)
        default:
            fmt.Println("❌ Invalid choice. Please try again.")
        }
    }
}

func handleSquash() {
    // Check if working directory is clean
    if clean, err := git.IsWorkingDirectoryClean(); err != nil || !clean {
        fmt.Println("❌ Please commit or stash your changes before squashing")
        return
    }

    // Show recent commits
    commits, err := git.GetRecentCommits(10)
    if err != nil {
        fmt.Printf("❌ Error fetching commits: %v\n", err)
        return
    }

    fmt.Println("\n📜 Recent Commits:")
    for i, commit := range commits {
        fmt.Printf("%d. %s: %s\n", i+1, commit.Hash, commit.Message)
    }

    baseCommit := GetInput("\n🎯 Enter the hash of the base commit to squash into: ")

    // Validate commit hash
    if _, err := git.RunCommand(fmt.Sprintf("git rev-parse --verify %s", baseCommit)); err != nil {
        fmt.Println("❌ Invalid commit hash")
        return
    }

    message := GetInput("✏️  Enter the new commit message: ")
    if message == "" {
        fmt.Println("❌ Commit message cannot be empty")
        return
    }

    fmt.Println("\n🔄 Squashing commits...")
    if err := git.SquashCommits(baseCommit, message); err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    fmt.Println("✅ Commits squashed successfully!")
    fmt.Println("⚠️  Note: If this branch was already pushed, you'll need to force push:")
    fmt.Printf("git push -f origin %s\n", getCurrentBranch())
}

func handleViewCommits() {
    commits, err := git.GetRecentCommits(10)
    if err != nil {
        fmt.Printf("❌ Error fetching commits: %v\n", err)
        return
    }

    fmt.Println("\n📜 Recent Commits:")
    for i, commit := range commits {
        fmt.Printf("%d. %s: %s\n", i+1, commit.Hash, commit.Message)
    }
}

func getCurrentBranch() string {
    branch, err := git.RunCommand("git rev-parse --abbrev-ref HEAD")
    if err != nil {
        return "current-branch"
    }
    return branch
}