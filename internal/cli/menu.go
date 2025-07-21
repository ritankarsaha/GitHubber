/*
 * GitHubber - CLI Menu System
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Interactive menu-driven interface for Git and GitHub operations
 */

package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/ritankarsaha/git-tool/internal/config"
	"github.com/ritankarsaha/git-tool/internal/git"
	"github.com/ritankarsaha/git-tool/internal/github"
	"github.com/ritankarsaha/git-tool/internal/ui"
)

func StartMenu() {
	for {
		// Repository Operations
		fmt.Println(ui.FormatMenuHeader(ui.IconRepository, "Repository Operations"))
		fmt.Println(ui.FormatMenuItem(1, "Initialize Repository"))
		fmt.Println(ui.FormatMenuItem(2, "Clone Repository"))

		// Branch Operations
		fmt.Println(ui.FormatMenuHeader(ui.IconBranch, "Branch Operations"))
		fmt.Println(ui.FormatMenuItem(3, "Create Branch"))
		fmt.Println(ui.FormatMenuItem(4, "Delete Branch"))
		fmt.Println(ui.FormatMenuItem(5, "Switch Branch"))
		fmt.Println(ui.FormatMenuItem(6, "List Branches"))

		// Changes and Staging
		fmt.Println(ui.FormatMenuHeader(ui.IconCommit, "Changes and Staging"))
		fmt.Println(ui.FormatMenuItem(7, "View Status"))
		fmt.Println(ui.FormatMenuItem(8, "Add Files"))
		fmt.Println(ui.FormatMenuItem(9, "Commit Changes"))

		// Remote Operations
		fmt.Println(ui.FormatMenuHeader(ui.IconRemote, "Remote Operations"))
		fmt.Println(ui.FormatMenuItem(10, "Push Changes"))
		fmt.Println(ui.FormatMenuItem(11, "Pull Changes"))
		fmt.Println(ui.FormatMenuItem(12, "Fetch Updates"))

		// History and Diff
		fmt.Println(ui.FormatMenuHeader(ui.IconHistory, "History and Diff"))
		fmt.Println(ui.FormatMenuItem(13, "View Log"))
		fmt.Println(ui.FormatMenuItem(14, "View Diff"))
		fmt.Println(ui.FormatMenuItem(15, "Squash Commits"))

		// Stash Operations
		fmt.Println(ui.FormatMenuHeader(ui.IconStash, "Stash Operations"))
		fmt.Println(ui.FormatMenuItem(16, "Stash Save"))
		fmt.Println(ui.FormatMenuItem(17, "Stash Pop"))
		fmt.Println(ui.FormatMenuItem(18, "List Stashes"))

		// Tag Operations
		fmt.Println(ui.FormatMenuHeader(ui.IconTag, "Tag Operations"))
		fmt.Println(ui.FormatMenuItem(19, "Create Tag"))
		fmt.Println(ui.FormatMenuItem(20, "Delete Tag"))
		fmt.Println(ui.FormatMenuItem(21, "List Tags"))

		// GitHub Operations (New section)
		fmt.Println(ui.FormatMenuHeader(ui.IconGitHub, "GitHub Operations"))
		fmt.Println(ui.FormatMenuItem(22, "View Repository Info"))
		fmt.Println(ui.FormatMenuItem(23, "Create Pull Request"))
		fmt.Println(ui.FormatMenuItem(24, "List Issues"))

		// Configuration and Exit
		fmt.Println(ui.FormatMenuHeader(ui.IconConfig, "Configuration"))
		fmt.Println(ui.FormatMenuItem(25, "Settings"))
		fmt.Println(ui.FormatMenuHeader(ui.IconExit, "Exit"))
		fmt.Println(ui.FormatMenuItem(26, "Exit"))

		choice := GetInput(ui.FormatPrompt("Enter your choice (1-26): "))

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
			handleRepoInfo()
		case "23":
			handleCreatePR()
		case "24":
			handleListIssues()
		case "25":
			handleSettings()
		case "26":
			fmt.Println(ui.FormatSuccess("Goodbye! Thank you for using GitHubber!"))
			os.Exit(0)
		default:
			fmt.Println(ui.FormatError("Invalid choice. Please try again."))
		}
	}
}

func handleInit() {
	if err := git.Init(); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Error initializing repository: %v", err)))
		return
	}
	fmt.Println(ui.FormatSuccess("Repository initialized successfully!"))
}

func handleClone() {
	url := GetInput(ui.FormatPrompt("Enter repository URL: "))
	if err := git.Clone(url); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Error cloning repository: %v", err)))
		return
	}
	fmt.Println(ui.FormatSuccess("Repository cloned successfully!"))
}

func handleCreateBranch() {
	name := GetInput(ui.FormatPrompt("Enter branch name: "))
	if err := git.CreateBranch(name); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Error creating branch: %v", err)))
		return
	}
	fmt.Println(ui.FormatSuccess("Branch created successfully!"))
}

func handleDeleteBranch() {
	name := GetInput(ui.FormatPrompt("Enter branch name to delete: "))
	if err := git.DeleteBranch(name); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Error deleting branch: %v", err)))
		return
	}
	fmt.Println(ui.FormatSuccess("Branch deleted successfully!"))
}

func handleSwitchBranch() {
	name := GetInput(ui.FormatPrompt("Enter branch name to switch to: "))
	if err := git.SwitchBranch(name); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Error switching branch: %v", err)))
		return
	}
	fmt.Println(ui.FormatSuccess("Switched to branch successfully!"))
}

func handleListBranches() {
	branches, err := git.ListBranches()
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Error listing branches: %v", err)))
		return
	}
	fmt.Println(ui.FormatInfo("Branches:"))
	for _, branch := range branches {
		fmt.Println(ui.FormatCode(branch))
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

// GitHub Operations (New handlers)
func handleRepoInfo() {
	client, err := github.NewClient()
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create GitHub client: %v", err)))
		fmt.Println(ui.FormatInfo("Please set GITHUB_TOKEN environment variable or configure authentication"))
		return
	}

	// Get current repository info from git
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		fmt.Println(ui.FormatError("Not in a Git repository or no remote origin found"))
		return
	}

	owner, repo, err := github.ParseRepoURL(repoInfo.URL)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to parse repository URL: %v", err)))
		return
	}

	repository, err := client.GetRepository(owner, repo)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to get repository information: %v", err)))
		return
	}

	fmt.Println(ui.FormatInfo("GitHub Repository Information"))
	fmt.Println(ui.FormatBox(fmt.Sprintf(
		"Name: %s\nOwner: %s\nDescription: %s\nURL: %s\nPrivate: %t\nLanguage: %s\nStars: %d\nForks: %d",
		repository.Name, repository.Owner, repository.Description,
		repository.URL, repository.Private, repository.Language,
		repository.Stars, repository.Forks,
	)))
}

func handleCreatePR() {
	client, err := github.NewClient()
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create GitHub client: %v", err)))
		return
	}

	// Get current repository info
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		fmt.Println(ui.FormatError("Not in a Git repository"))
		return
	}

	owner, repo, err := github.ParseRepoURL(repoInfo.URL)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to parse repository URL: %v", err)))
		return
	}

	// Get current branch
	currentBranch := getCurrentBranch()
	
	fmt.Println(ui.FormatInfo("Create Pull Request"))
	title := GetInput(ui.FormatPrompt("Enter PR title: "))
	body := GetInput(ui.FormatPrompt("Enter PR description: "))
	base := GetInput(ui.FormatPrompt(fmt.Sprintf("Enter base branch (default: main): ")))
	if base == "" {
		base = "main"
	}

	pr, err := client.CreatePullRequest(owner, repo, title, body, currentBranch, base)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create pull request: %v", err)))
		return
	}

	fmt.Println(ui.FormatSuccess(fmt.Sprintf("Pull request created successfully!")))
	fmt.Println(ui.FormatInfo(fmt.Sprintf("PR #%d: %s", pr.Number, pr.Title)))
	fmt.Println(ui.FormatInfo(fmt.Sprintf("URL: %s", pr.URL)))
}

func handleListIssues() {
	client, err := github.NewClient()
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create GitHub client: %v", err)))
		return
	}

	// Get current repository info
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		fmt.Println(ui.FormatError("Not in a Git repository"))
		return
	}

	owner, repo, err := github.ParseRepoURL(repoInfo.URL)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to parse repository URL: %v", err)))
		return
	}

	state := GetInput(ui.FormatPrompt("Enter issue state (open/closed/all, default: open): "))
	if state == "" {
		state = "open"
	}

	issues, err := client.ListIssues(owner, repo, state)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to list issues: %v", err)))
		return
	}

	if len(issues) == 0 {
		fmt.Println(ui.FormatInfo("No issues found"))
		return
	}

	fmt.Println(ui.FormatInfo(fmt.Sprintf("GitHub Issues (%s)", state)))
	for _, issue := range issues {
		fmt.Printf("%s #%d: %s (%s) by %s\n",
			ui.IconInfo, issue.Number, issue.Title, issue.State, issue.Author)
	}
}

func handleSettings() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to load configuration: %v", err)))
		return
	}

	fmt.Println(ui.FormatInfo("GitHubber Settings"))
	fmt.Println("1. View current settings")
	fmt.Println("2. Set GitHub token")
	fmt.Println("3. Set default repository")
	fmt.Println("4. UI preferences")
	fmt.Println("5. Back to main menu")

	choice := GetInput(ui.FormatPrompt("Enter your choice (1-5): "))

	switch choice {
	case "1":
		showCurrentSettings(cfg)
	case "2":
		setGitHubToken(cfg)
	case "3":
		setDefaultRepo(cfg)
	case "4":
		setUIPreferences(cfg)
	case "5":
		return
	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func showCurrentSettings(cfg *config.Config) {
	fmt.Println(ui.FormatInfo("Current Settings"))
	hasToken := "No"
	if cfg.GetGitHubToken() != "" {
		hasToken = "Yes"
	}
	
	settings := fmt.Sprintf(
		"GitHub Token: %s\nDefault Owner: %s\nDefault Repo: %s\nTheme: %s\nShow Emojis: %t\nPage Size: %d",
		hasToken, cfg.GitHub.DefaultOwner, cfg.GitHub.DefaultRepo,
		cfg.UI.Theme, cfg.UI.ShowEmojis, cfg.UI.PageSize,
	)
	fmt.Println(ui.FormatBox(settings))
}

func setGitHubToken(cfg *config.Config) {
	token := GetInput(ui.FormatPrompt("Enter GitHub personal access token: "))
	if token == "" {
		fmt.Println(ui.FormatWarning("Token not set"))
		return
	}
	
	if err := cfg.SetGitHubToken(token); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to save token: %v", err)))
		return
	}
	
	fmt.Println(ui.FormatSuccess("GitHub token saved successfully"))
}

func setDefaultRepo(cfg *config.Config) {
	owner := GetInput(ui.FormatPrompt("Enter default repository owner: "))
	repo := GetInput(ui.FormatPrompt("Enter default repository name: "))
	
	cfg.GitHub.DefaultOwner = owner
	cfg.GitHub.DefaultRepo = repo
	
	if err := cfg.Save(); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to save configuration: %v", err)))
		return
	}
	
	fmt.Println(ui.FormatSuccess("Default repository saved successfully"))
}

func setUIPreferences(cfg *config.Config) {
	theme := GetInput(ui.FormatPrompt("Enter theme (dark/light/auto): "))
	if theme != "" {
		cfg.UI.Theme = theme
	}
	
	if err := cfg.Save(); err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to save configuration: %v", err)))
		return
	}
	
	fmt.Println(ui.FormatSuccess("UI preferences saved successfully"))
}
