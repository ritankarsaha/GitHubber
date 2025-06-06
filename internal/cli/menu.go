package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ritankarsaha/git-tool/internal/git"
)

func StartMenu() {
	for {
		fmt.Println("\n📋 Git Tool Menu:")
		fmt.Println("\n🔧 Repository Operations:")
		fmt.Println("1. Initialize Repository")
		fmt.Println("2. Clone Repository")

		fmt.Println("\n🌿 Branch Operations:")
		fmt.Println("3. Create Branch")
		fmt.Println("4. Delete Branch")
		fmt.Println("5. Switch Branch")
		fmt.Println("6. List Branches")

		fmt.Println("\n💾 Changes and Staging:")
		fmt.Println("7. View Status")
		fmt.Println("8. Add Files")
		fmt.Println("9. Commit Changes")

		fmt.Println("\n🔄 Remote Operations:")
		fmt.Println("10. Push Changes")
		fmt.Println("11. Pull Changes")
		fmt.Println("12. Fetch Updates")

		fmt.Println("\n📜 History and Diff:")
		fmt.Println("13. View Log")
		fmt.Println("14. View Diff")
		fmt.Println("15. Squash Commits")

		fmt.Println("\n📦 Stash Operations:")
		fmt.Println("16. Stash Save")
		fmt.Println("17. Stash Pop")
		fmt.Println("18. List Stashes")

		fmt.Println("\n🏷️  Tag Operations:")
		fmt.Println("19. Create Tag")
		fmt.Println("20. Delete Tag")
		fmt.Println("21. List Tags")

		fmt.Println("\n❌ Exit:")
		fmt.Println("22. Exit")

		choice := GetInput("\nEnter your choice (1-22): ")

		switch choice {
		case "1":
			handleInit()
		case "2":
			handleClone()
		case "3":
			handleCreateBranch()
		case "4":
			handleDeleteBranch()
		case "5":
			handleSwitchBranch()
		case "6":
			handleListBranches()
		case "7":
			handleStatus()
		case "8":
			handleAddFiles()
		case "9":
			handleCommit()
		case "10":
			handlePush()
		case "11":
			handlePull()
		case "12":
			handleFetch()
		case "13":
			handleLog()
		case "14":
			handleDiff()
		case "15":
			handleSquash()
		case "16":
			handleStashSave()
		case "17":
			handleStashPop()
		case "18":
			handleStashList()
		case "19":
			handleCreateTag()
		case "20":
			handleDeleteTag()
		case "21":
			handleListTags()
		case "22":
			fmt.Println("👋 Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}
	}
}

func handleInit() {
	if err := git.Init(); err != nil {
		fmt.Printf("❌ Error initializing repository: %v\n", err)
		return
	}
	fmt.Println("✅ Repository initialized successfully!")
}

func handleClone() {
	url := GetInput("Enter repository URL: ")
	if err := git.Clone(url); err != nil {
		fmt.Printf("❌ Error cloning repository: %v\n", err)
		return
	}
	fmt.Println("✅ Repository cloned successfully!")
}

func handleCreateBranch() {
	name := GetInput("Enter branch name: ")
	if err := git.CreateBranch(name); err != nil {
		fmt.Printf("❌ Error creating branch: %v\n", err)
		return
	}
	fmt.Println("✅ Branch created successfully!")
}

func handleDeleteBranch() {
	name := GetInput("Enter branch name to delete: ")
	if err := git.DeleteBranch(name); err != nil {
		fmt.Printf("❌ Error deleting branch: %v\n", err)
		return
	}
	fmt.Println("✅ Branch deleted successfully!")
}

func handleSwitchBranch() {
	name := GetInput("Enter branch name to switch to: ")
	if err := git.SwitchBranch(name); err != nil {
		fmt.Printf("❌ Error switching branch: %v\n", err)
		return
	}
	fmt.Println("✅ Switched to branch successfully!")
}

func handleListBranches() {
	branches, err := git.ListBranches()
	if err != nil {
		fmt.Printf("❌ Error listing branches: %v\n", err)
		return
	}
	fmt.Println("\n🌿 Branches:")
	for _, branch := range branches {
		fmt.Println(branch)
	}
}

func handleStatus() {
	status, err := git.Status()
	if err != nil {
		fmt.Printf("❌ Error getting status: %v\n", err)
		return
	}
	fmt.Printf("\n📊 Git Status:\n%s\n", status)
}

func handleAddFiles() {
	files := GetInput("Enter files to add (space-separated, or press enter for all): ")
	var err error
	if files == "" {
		err = git.AddFiles()
	} else {
		err = git.AddFiles(strings.Fields(files)...)
	}
	if err != nil {
		fmt.Printf("❌ Error adding files: %v\n", err)
		return
	}
	fmt.Println("✅ Files added successfully!")
}

func handleCommit() {
	message := GetInput("Enter commit message: ")
	if err := git.Commit(message); err != nil {
		fmt.Printf("❌ Error committing changes: %v\n", err)
		return
	}
	fmt.Println("✅ Changes committed successfully!")
}

func handlePush() {
	remote := GetInput("Enter remote name (default: origin): ")
	if remote == "" {
		remote = "origin"
	}
	branch := GetInput("Enter branch name: ")
	if err := git.Push(remote, branch); err != nil {
		fmt.Printf("❌ Error pushing changes: %v\n", err)
		return
	}
	fmt.Println("✅ Changes pushed successfully!")
}

func handlePull() {
	remote := GetInput("Enter remote name (default: origin): ")
	if remote == "" {
		remote = "origin"
	}
	branch := GetInput("Enter branch name: ")
	if err := git.Pull(remote, branch); err != nil {
		fmt.Printf("❌ Error pulling changes: %v\n", err)
		return
	}
	fmt.Println("✅ Changes pulled successfully!")
}

func handleFetch() {
	remote := GetInput("Enter remote name (default: origin): ")
	if remote == "" {
		remote = "origin"
	}
	if err := git.Fetch(remote); err != nil {
		fmt.Printf("❌ Error fetching updates: %v\n", err)
		return
	}
	fmt.Println("✅ Updates fetched successfully!")
}

func handleLog() {
	n := 10 // Default to last 10 commits
	logs, err := git.Log(n)
	if err != nil {
		fmt.Printf("❌ Error viewing log: %v\n", err)
		return
	}
	fmt.Printf("\n📜 Last %d commits:\n%s\n", n, logs)
}

func handleDiff() {
	file := GetInput("Enter file to diff (press enter for all files): ")
	diff, err := git.Diff(file)
	if err != nil {
		fmt.Printf("❌ Error viewing diff: %v\n", err)
		return
	}
	fmt.Printf("\n📝 Diff:\n%s\n", diff)
}

func handleStashSave() {
	message := GetInput("Enter stash message: ")
	if err := git.StashSave(message); err != nil {
		fmt.Printf("❌ Error stashing changes: %v\n", err)
		return
	}
	fmt.Println("✅ Changes stashed successfully!")
}

func handleStashPop() {
	if err := git.StashPop(); err != nil {
		fmt.Printf("❌ Error popping stash: %v\n", err)
		return
	}
	fmt.Println("✅ Stash applied successfully!")
}

func handleStashList() {
	list, err := git.StashList()
	if err != nil {
		fmt.Printf("❌ Error listing stashes: %v\n", err)
		return
	}
	fmt.Printf("\n📦 Stash list:\n%s\n", list)
}

func handleCreateTag() {
	name := GetInput("Enter tag name: ")
	message := GetInput("Enter tag message: ")
	if err := git.CreateTag(name, message); err != nil {
		fmt.Printf("❌ Error creating tag: %v\n", err)
		return
	}
	fmt.Println("✅ Tag created successfully!")
}

func handleDeleteTag() {
	name := GetInput("Enter tag name to delete: ")
	if err := git.DeleteTag(name); err != nil {
		fmt.Printf("❌ Error deleting tag: %v\n", err)
		return
	}
	fmt.Println("✅ Tag deleted successfully!")
}

func handleListTags() {
	tags, err := git.ListTags()
	if err != nil {
		fmt.Printf("❌ Error listing tags: %v\n", err)
		return
	}
	fmt.Printf("\n🏷️  Tags:\n%s\n", tags)
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

func getCurrentBranch() string {
	branch, err := git.RunCommand("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "current-branch"
	}
	return branch
}
