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

		// Advanced Git Operations
		fmt.Println(ui.FormatMenuHeader("🔧", "Advanced Git Operations"))
		fmt.Println(ui.FormatMenuItem(22, "Interactive Rebase"))
		fmt.Println(ui.FormatMenuItem(23, "Cherry Pick"))
		fmt.Println(ui.FormatMenuItem(24, "Reset (Soft/Mixed/Hard)"))
		fmt.Println(ui.FormatMenuItem(25, "Revert Commit"))
		fmt.Println(ui.FormatMenuItem(26, "Merge Branch"))
		fmt.Println(ui.FormatMenuItem(27, "Bisect"))
		fmt.Println(ui.FormatMenuItem(28, "Resolve Conflicts"))

		// GitHub Operations
		fmt.Println(ui.FormatMenuHeader(ui.IconGitHub, "GitHub Operations"))
		fmt.Println(ui.FormatMenuItem(29, "View Repository Info"))
		fmt.Println(ui.FormatMenuItem(30, "Create Pull Request"))
		fmt.Println(ui.FormatMenuItem(31, "List Pull Requests"))
		fmt.Println(ui.FormatMenuItem(32, "List Issues"))
		fmt.Println(ui.FormatMenuItem(33, "Create Issue"))
		fmt.Println(ui.FormatMenuItem(34, "Repository Management"))

		// Remote Management
		fmt.Println(ui.FormatMenuHeader("🔗", "Remote Management"))
		fmt.Println(ui.FormatMenuItem(35, "Add Remote"))
		fmt.Println(ui.FormatMenuItem(36, "Remove Remote"))
		fmt.Println(ui.FormatMenuItem(37, "Rename Remote"))
		fmt.Println(ui.FormatMenuItem(38, "List All Remotes"))
		fmt.Println(ui.FormatMenuItem(39, "Set Remote URL"))
		fmt.Println(ui.FormatMenuItem(40, "Sync with All Remotes"))

		// Advanced History & Analysis
		fmt.Println(ui.FormatMenuHeader("📊", "Advanced History & Analysis"))
		fmt.Println(ui.FormatMenuItem(41, "Interactive Log Viewer"))
		fmt.Println(ui.FormatMenuItem(42, "File History"))
		fmt.Println(ui.FormatMenuItem(43, "Blame/Annotate File"))
		fmt.Println(ui.FormatMenuItem(44, "Show Commit Details"))
		fmt.Println(ui.FormatMenuItem(45, "Compare Branches"))
		fmt.Println(ui.FormatMenuItem(46, "Find Commits by Author"))
		fmt.Println(ui.FormatMenuItem(47, "Find Commits by Message"))

		// Patch & Bundle Operations
		fmt.Println(ui.FormatMenuHeader("🧩", "Patch & Bundle Operations"))
		fmt.Println(ui.FormatMenuItem(48, "Create Patch File"))
		fmt.Println(ui.FormatMenuItem(49, "Apply Patch"))
		fmt.Println(ui.FormatMenuItem(50, "Create Bundle"))
		fmt.Println(ui.FormatMenuItem(51, "Verify/List Bundle"))
		fmt.Println(ui.FormatMenuItem(52, "Format Patch for Email"))

		// Worktree Management
		fmt.Println(ui.FormatMenuHeader("🌳", "Worktree Management"))
		fmt.Println(ui.FormatMenuItem(53, "List Worktrees"))
		fmt.Println(ui.FormatMenuItem(54, "Add Worktree"))
		fmt.Println(ui.FormatMenuItem(55, "Remove Worktree"))
		fmt.Println(ui.FormatMenuItem(56, "Move Worktree"))

		// Repository Maintenance
		fmt.Println(ui.FormatMenuHeader("📋", "Repository Maintenance"))
		fmt.Println(ui.FormatMenuItem(57, "Garbage Collection"))
		fmt.Println(ui.FormatMenuItem(58, "Verify Repository Integrity"))
		fmt.Println(ui.FormatMenuItem(59, "Optimize Repository"))
		fmt.Println(ui.FormatMenuItem(60, "Repository Statistics"))
		fmt.Println(ui.FormatMenuItem(61, "Reflog Management"))
		fmt.Println(ui.FormatMenuItem(62, "Prune Remote Branches"))

		// Smart Git Operations
		fmt.Println(ui.FormatMenuHeader("🎯", "Smart Git Operations"))
		fmt.Println(ui.FormatMenuItem(63, "Interactive Add (Patch Mode)"))
		fmt.Println(ui.FormatMenuItem(64, "Partial File Commits"))
		fmt.Println(ui.FormatMenuItem(65, "Commit Amend Helper"))
		fmt.Println(ui.FormatMenuItem(66, "Branch Comparison Tool"))
		fmt.Println(ui.FormatMenuItem(67, "Conflict Prevention Check"))

		// Configuration and Exit
		fmt.Println(ui.FormatMenuHeader(ui.IconConfig, "Configuration"))
		fmt.Println(ui.FormatMenuItem(68, "Settings"))
		fmt.Println(ui.FormatMenuHeader(ui.IconExit, "Exit"))
		fmt.Println(ui.FormatMenuItem(69, "Exit"))

		choice := GetInput(ui.FormatPrompt("Enter your choice (1-69): "))

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
			handleInteractiveRebase()
		case "23":
			handleCherryPickMenu()
		case "24":
			handleResetMenu()
		case "25":
			handleRevertMenu()
		case "26":
			handleMergeMenu()
		case "27":
			handleBisectMenu()
		case "28":
			handleResolveConflictsMenu()
		case "29":
			handleRepoInfo()
		case "30":
			handleCreatePR()
		case "31":
			handleListPRs()
		case "32":
			handleListIssues()
		case "33":
			handleCreateIssueMenu()
		case "34":
			handleRepoManagement()
		case "35":
			handleAddRemote()
		case "36":
			handleRemoveRemote()
		case "37":
			handleRenameRemote()
		case "38":
			handleListRemotes()
		case "39":
			handleSetRemoteURL()
		case "40":
			handleSyncAllRemotes()
		case "41":
			handleInteractiveLogViewer()
		case "42":
			handleFileHistory()
		case "43":
			handleBlameFile()
		case "44":
			handleShowCommitDetails()
		case "45":
			handleCompareBranches()
		case "46":
			handleFindCommitsByAuthor()
		case "47":
			handleFindCommitsByMessage()
		case "48":
			handleCreatePatch()
		case "49":
			handleApplyPatch()
		case "50":
			handleCreateBundle()
		case "51":
			handleVerifyBundle()
		case "52":
			handleFormatPatchEmail()
		case "53":
			handleListWorktrees()
		case "54":
			handleAddWorktree()
		case "55":
			handleRemoveWorktree()
		case "56":
			handleMoveWorktree()
		case "57":
			handleGarbageCollection()
		case "58":
			handleVerifyRepository()
		case "59":
			handleOptimizeRepository()
		case "60":
			handleRepositoryStatistics()
		case "61":
			handleReflogManagement()
		case "62":
			handlePruneRemoteBranches()
		case "63":
			handleInteractiveAdd()
		case "64":
			handlePartialCommit()
		case "65":
			handleCommitAmend()
		case "66":
			handleBranchComparisonTool()
		case "67":
			handleConflictPreventionCheck()
		case "68":
			handleSettings()
		case "69":
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

// Advanced Git Operations Handlers

func handleInteractiveRebase() {
	base := GetInput(ui.FormatPrompt("Enter base commit for rebase: "))
	if base == "" {
		fmt.Println(ui.FormatError("Base commit is required"))
		return
	}

	fmt.Println(ui.FormatInfo("Starting interactive rebase..."))
	if err := git.InteractiveRebase(base); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println(ui.FormatSuccess("Interactive rebase started!"))
}

func handleCherryPickMenu() {
	fmt.Println(ui.FormatInfo("Cherry Pick Operations"))
	fmt.Println("1. Cherry pick commit")
	fmt.Println("2. Cherry pick range")
	fmt.Println("3. Continue cherry pick")
	fmt.Println("4. Abort cherry pick")

	choice := GetInput(ui.FormatPrompt("Choose option (1-4): "))

	switch choice {
	case "1":
		commit := GetInput(ui.FormatPrompt("Enter commit hash: "))
		if err := git.CherryPick(commit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Cherry pick completed!"))
	case "2":
		start := GetInput(ui.FormatPrompt("Enter start commit: "))
		end := GetInput(ui.FormatPrompt("Enter end commit: "))
		if err := git.CherryPickRange(start, end); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Cherry pick range completed!"))
	case "3":
		if err := git.CherryPickContinue(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Cherry pick continued!"))
	case "4":
		if err := git.CherryPickAbort(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Cherry pick aborted!"))
	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handleResetMenu() {
	fmt.Println(ui.FormatInfo("Reset Operations"))
	fmt.Println("1. Soft reset (keep changes staged)")
	fmt.Println("2. Mixed reset (unstage changes)")
	fmt.Println("3. Hard reset (discard all changes)")
	fmt.Println("4. Reset specific file")

	choice := GetInput(ui.FormatPrompt("Choose option (1-4): "))

	switch choice {
	case "1":
		commit := GetInput(ui.FormatPrompt("Enter commit to reset to: "))
		if err := git.ResetSoft(commit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Soft reset completed!"))
	case "2":
		commit := GetInput(ui.FormatPrompt("Enter commit to reset to: "))
		if err := git.ResetMixed(commit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Mixed reset completed!"))
	case "3":
		commit := GetInput(ui.FormatPrompt("Enter commit to reset to: "))
		confirm := GetInput(ui.FormatPrompt("⚠️  This will discard all changes. Continue? (yes/no): "))
		if strings.ToLower(confirm) != "yes" {
			fmt.Println(ui.FormatInfo("Reset cancelled"))
			return
		}
		if err := git.ResetHard(commit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Hard reset completed!"))
	case "4":
		file := GetInput(ui.FormatPrompt("Enter file to reset: "))
		if err := git.ResetFile(file); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("File reset completed!"))
	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handleRevertMenu() {
	fmt.Println(ui.FormatInfo("Revert Operations"))
	fmt.Println("1. Revert commit (create new commit)")
	fmt.Println("2. Revert without committing")

	choice := GetInput(ui.FormatPrompt("Choose option (1-2): "))
	commit := GetInput(ui.FormatPrompt("Enter commit hash to revert: "))

	switch choice {
	case "1":
		if err := git.Revert(commit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Commit reverted!"))
	case "2":
		if err := git.RevertNoCommit(commit); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Changes reverted (not committed)!"))
	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handleMergeMenu() {
	fmt.Println(ui.FormatInfo("Merge Operations"))
	fmt.Println("1. Regular merge")
	fmt.Println("2. No fast-forward merge")
	fmt.Println("3. Squash merge")
	fmt.Println("4. Abort merge")
	fmt.Println("5. Continue merge")

	choice := GetInput(ui.FormatPrompt("Choose option (1-5): "))

	switch choice {
	case "1":
		branch := GetInput(ui.FormatPrompt("Enter branch to merge: "))
		if err := git.Merge(branch); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Merge completed!"))
	case "2":
		branch := GetInput(ui.FormatPrompt("Enter branch to merge: "))
		if err := git.MergeNoFF(branch); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("No-FF merge completed!"))
	case "3":
		branch := GetInput(ui.FormatPrompt("Enter branch to merge: "))
		if err := git.MergeSquash(branch); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Squash merge completed!"))
	case "4":
		if err := git.MergeAbort(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Merge aborted!"))
	case "5":
		if err := git.MergeContinue(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Merge continued!"))
	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handleBisectMenu() {
	fmt.Println(ui.FormatInfo("Git Bisect Operations"))
	fmt.Println("1. Start bisect")
	fmt.Println("2. Mark current commit as bad")
	fmt.Println("3. Mark current commit as good")
	fmt.Println("4. Skip current commit")
	fmt.Println("5. Reset bisect")

	choice := GetInput(ui.FormatPrompt("Choose option (1-5): "))

	switch choice {
	case "1":
		if err := git.BisectStart(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Bisect started! Mark commits as good or bad."))
	case "2":
		if err := git.BisectBad(""); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatInfo("Current commit marked as bad"))
	case "3":
		if err := git.BisectGood(""); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatInfo("Current commit marked as good"))
	case "4":
		if err := git.BisectSkip(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatInfo("Current commit skipped"))
	case "5":
		if err := git.BisectReset(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println(ui.FormatSuccess("Bisect reset!"))
	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handleResolveConflictsMenu() {
	fmt.Println(ui.FormatInfo("Starting conflict resolution..."))
	if err := StartConflictResolution(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
}

// Enhanced GitHub Operations Handlers

func handleListPRs() {
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

	state := GetInput(ui.FormatPrompt("Enter PR state (open/closed/all, default: open): "))
	if state == "" {
		state = "open"
	}

	prs, err := client.ListPullRequests(owner, repo, state)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to list pull requests: %v", err)))
		return
	}

	if len(prs) == 0 {
		fmt.Println(ui.FormatInfo("No pull requests found"))
		return
	}

	fmt.Println(ui.FormatInfo(fmt.Sprintf("Pull Requests (%s)", state)))
	for _, pr := range prs {
		fmt.Printf("%s #%d: %s (%s) by %s\n",
			ui.IconInfo, pr.Number, pr.Title, pr.State, pr.Author)
		fmt.Printf("   URL: %s\n", pr.URL)
	}
}

func handleCreateIssueMenu() {
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

	fmt.Println(ui.FormatInfo("Create New Issue"))
	title := GetInput(ui.FormatPrompt("Enter issue title: "))
	if title == "" {
		fmt.Println(ui.FormatError("Title is required"))
		return
	}

	body := GetInput(ui.FormatPrompt("Enter issue description: "))
	labels := []string{}
	labelsInput := GetInput(ui.FormatPrompt("Enter labels (comma-separated, optional): "))
	if labelsInput != "" {
		labels = strings.Split(labelsInput, ",")
		for i := range labels {
			labels[i] = strings.TrimSpace(labels[i])
		}
	}

	issue, err := client.CreateIssue(owner, repo, title, body, labels)
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create issue: %v", err)))
		return
	}

	fmt.Println(ui.FormatSuccess("Issue created successfully!"))
	fmt.Printf("Issue #%d: %s\n", issue.Number, issue.Title)
	fmt.Printf("URL: %s\n", issue.URL)
}

func handleRepoManagement() {
	fmt.Println(ui.FormatInfo("Repository Management"))
	fmt.Println("1. List repositories")
	fmt.Println("2. Create repository")
	fmt.Println("3. Fork repository")
	fmt.Println("4. Repository search")

	client, err := github.NewClient()
	if err != nil {
		fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create GitHub client: %v", err)))
		return
	}

	choice := GetInput(ui.FormatPrompt("Choose option (1-4): "))

	switch choice {
	case "1":
		visibility := GetInput(ui.FormatPrompt("Enter visibility (all/public/private, default: all): "))
		if visibility == "" {
			visibility = "all"
		}

		repos, err := client.ListRepositories(visibility)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("Found %d repositories:\n", len(repos))
		for _, repo := range repos {
			fmt.Printf("  %s/%s (%s) - %d stars, %d forks\n",
				repo.Owner, repo.Name, repo.Language, repo.Stars, repo.Forks)
		}

	case "2":
		name := GetInput(ui.FormatPrompt("Enter repository name: "))
		if name == "" {
			fmt.Println(ui.FormatError("Repository name is required"))
			return
		}

		description := GetInput(ui.FormatPrompt("Enter description (optional): "))
		privateInput := GetInput(ui.FormatPrompt("Make private? (y/n, default: n): "))
		private := strings.ToLower(privateInput) == "y" || strings.ToLower(privateInput) == "yes"

		repo, err := client.CreateRepository(name, description, private)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Println(ui.FormatSuccess("Repository created successfully!"))
		fmt.Printf("Repository: %s/%s\n", repo.Owner, repo.Name)
		fmt.Printf("URL: %s\n", repo.URL)

	case "3":
		owner := GetInput(ui.FormatPrompt("Enter owner/organization: "))
		repoName := GetInput(ui.FormatPrompt("Enter repository name: "))

		if owner == "" || repoName == "" {
			fmt.Println(ui.FormatError("Owner and repository name are required"))
			return
		}

		forkedRepo, err := client.ForkRepository(owner, repoName)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Println(ui.FormatSuccess("Repository forked successfully!"))
		fmt.Printf("Forked: %s/%s\n", forkedRepo.Owner, forkedRepo.Name)
		fmt.Printf("URL: %s\n", forkedRepo.URL)

	case "4":
		query := GetInput(ui.FormatPrompt("Enter search query: "))
		if query == "" {
			fmt.Println(ui.FormatError("Search query is required"))
			return
		}

		limit := 10
		repos, err := client.SearchRepositories(query, limit)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("Found %d repositories:\n", len(repos))
		for _, repo := range repos {
			fmt.Printf("  %s/%s (%s) - %d stars\n",
				repo.Owner, repo.Name, repo.Language, repo.Stars)
			fmt.Printf("    %s\n", repo.Description)
		}

	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

// Remote Management Handlers

func handleAddRemote() {
	name := GetInput(ui.FormatPrompt("Enter remote name: "))
	if name == "" {
		fmt.Println(ui.FormatError("Remote name is required"))
		return
	}

	url := GetInput(ui.FormatPrompt("Enter remote URL: "))
	if url == "" {
		fmt.Println(ui.FormatError("Remote URL is required"))
		return
	}

	if err := git.AddRemote(name, url); err != nil {
		fmt.Printf("❌ Error adding remote: %v\n", err)
		return
	}
	fmt.Printf("✅ Remote '%s' added successfully!\n", name)
}

func handleRemoveRemote() {
	remotes, err := git.ListRemotes()
	if err != nil {
		fmt.Printf("❌ Error listing remotes: %v\n", err)
		return
	}

	fmt.Printf("Current remotes:\n%s\n", remotes)
	name := GetInput(ui.FormatPrompt("Enter remote name to remove: "))
	if name == "" {
		fmt.Println(ui.FormatError("Remote name is required"))
		return
	}

	if err := git.RemoveRemote(name); err != nil {
		fmt.Printf("❌ Error removing remote: %v\n", err)
		return
	}
	fmt.Printf("✅ Remote '%s' removed successfully!\n", name)
}

func handleRenameRemote() {
	remotes, err := git.ListRemotes()
	if err != nil {
		fmt.Printf("❌ Error listing remotes: %v\n", err)
		return
	}

	fmt.Printf("Current remotes:\n%s\n", remotes)
	oldName := GetInput(ui.FormatPrompt("Enter current remote name: "))
	newName := GetInput(ui.FormatPrompt("Enter new remote name: "))

	if oldName == "" || newName == "" {
		fmt.Println(ui.FormatError("Both old and new names are required"))
		return
	}

	if err := git.RenameRemote(oldName, newName); err != nil {
		fmt.Printf("❌ Error renaming remote: %v\n", err)
		return
	}
	fmt.Printf("✅ Remote renamed from '%s' to '%s'!\n", oldName, newName)
}

func handleListRemotes() {
	remotes, err := git.ListRemotes()
	if err != nil {
		fmt.Printf("❌ Error listing remotes: %v\n", err)
		return
	}

	if remotes == "" {
		fmt.Println("📭 No remotes configured")
		return
	}

	fmt.Printf("🔗 Configured remotes:\n%s\n", remotes)
}

func handleSetRemoteURL() {
	remotes, err := git.ListRemotes()
	if err != nil {
		fmt.Printf("❌ Error listing remotes: %v\n", err)
		return
	}

	fmt.Printf("Current remotes:\n%s\n", remotes)
	name := GetInput(ui.FormatPrompt("Enter remote name: "))
	url := GetInput(ui.FormatPrompt("Enter new URL: "))

	if name == "" || url == "" {
		fmt.Println(ui.FormatError("Both remote name and URL are required"))
		return
	}

	if err := git.SetRemoteURL(name, url); err != nil {
		fmt.Printf("❌ Error setting remote URL: %v\n", err)
		return
	}
	fmt.Printf("✅ URL for remote '%s' updated successfully!\n", name)
}

func handleSyncAllRemotes() {
	fmt.Println(ui.FormatInfo("Syncing with all remotes..."))
	if err := git.SyncWithAllRemotes(); err != nil {
		fmt.Printf("❌ Error syncing remotes: %v\n", err)
		return
	}
	fmt.Println("✅ All remotes synced successfully!")
}

// Advanced History & Analysis Handlers

func handleInteractiveLogViewer() {
	logs, err := git.InteractiveLog()
	if err != nil {
		fmt.Printf("❌ Error getting log: %v\n", err)
		return
	}

	fmt.Printf("📊 Interactive Git Log:\n%s\n", logs)
}

func handleFileHistory() {
	file := GetInput(ui.FormatPrompt("Enter file path: "))
	if file == "" {
		fmt.Println(ui.FormatError("File path is required"))
		return
	}

	history, err := git.FileHistory(file)
	if err != nil {
		fmt.Printf("❌ Error getting file history: %v\n", err)
		return
	}

	fmt.Printf("📜 History for %s:\n%s\n", file, history)
}

func handleBlameFile() {
	file := GetInput(ui.FormatPrompt("Enter file path: "))
	if file == "" {
		fmt.Println(ui.FormatError("File path is required"))
		return
	}

	blame, err := git.BlameFile(file)
	if err != nil {
		fmt.Printf("❌ Error getting blame: %v\n", err)
		return
	}

	fmt.Printf("🔍 Blame for %s:\n%s\n", file, blame)
}

func handleShowCommitDetails() {
	commit := GetInput(ui.FormatPrompt("Enter commit hash: "))
	if commit == "" {
		fmt.Println(ui.FormatError("Commit hash is required"))
		return
	}

	details, err := git.ShowCommitDetails(commit)
	if err != nil {
		fmt.Printf("❌ Error getting commit details: %v\n", err)
		return
	}

	fmt.Printf("📋 Commit Details:\n%s\n", details)
}

func handleCompareBranches() {
	branch1 := GetInput(ui.FormatPrompt("Enter first branch: "))
	branch2 := GetInput(ui.FormatPrompt("Enter second branch: "))

	if branch1 == "" || branch2 == "" {
		fmt.Println(ui.FormatError("Both branch names are required"))
		return
	}

	comparison, err := git.CompareBranches(branch1, branch2)
	if err != nil {
		fmt.Printf("❌ Error comparing branches: %v\n", err)
		return
	}

	fmt.Printf("📊 Branch Comparison (%s vs %s):\n%s\n", branch1, branch2, comparison)

	// Also show numerical comparison
	numComparison, _ := git.BranchComparison(branch1, branch2)
	fmt.Printf("\n%s\n", numComparison)
}

func handleFindCommitsByAuthor() {
	author := GetInput(ui.FormatPrompt("Enter author name/email: "))
	if author == "" {
		fmt.Println(ui.FormatError("Author name is required"))
		return
	}

	commits, err := git.FindCommitsByAuthor(author)
	if err != nil {
		fmt.Printf("❌ Error finding commits: %v\n", err)
		return
	}

	fmt.Printf("👤 Commits by %s:\n%s\n", author, commits)
}

func handleFindCommitsByMessage() {
	message := GetInput(ui.FormatPrompt("Enter search term: "))
	if message == "" {
		fmt.Println(ui.FormatError("Search term is required"))
		return
	}

	commits, err := git.FindCommitsByMessage(message)
	if err != nil {
		fmt.Printf("❌ Error finding commits: %v\n", err)
		return
	}

	fmt.Printf("🔍 Commits containing '%s':\n%s\n", message, commits)
}

// Patch & Bundle Operations Handlers

func handleCreatePatch() {
	outputFile := GetInput(ui.FormatPrompt("Enter output file name (e.g., changes.patch): "))
	if outputFile == "" {
		outputFile = "changes.patch"
	}

	if err := git.CreatePatch(outputFile); err != nil {
		fmt.Printf("❌ Error creating patch: %v\n", err)
		return
	}

	fmt.Printf("✅ Patch created: %s\n", outputFile)
}

func handleApplyPatch() {
	patchFile := GetInput(ui.FormatPrompt("Enter patch file path: "))
	if patchFile == "" {
		fmt.Println(ui.FormatError("Patch file path is required"))
		return
	}

	if err := git.ApplyPatch(patchFile); err != nil {
		fmt.Printf("❌ Error applying patch: %v\n", err)
		return
	}

	fmt.Printf("✅ Patch applied: %s\n", patchFile)
}

func handleCreateBundle() {
	bundleFile := GetInput(ui.FormatPrompt("Enter bundle file name: "))
	refSpec := GetInput(ui.FormatPrompt("Enter ref specification (e.g., HEAD, main, --all): "))

	if bundleFile == "" || refSpec == "" {
		fmt.Println(ui.FormatError("Both bundle file and ref specification are required"))
		return
	}

	if err := git.CreateBundle(bundleFile, refSpec); err != nil {
		fmt.Printf("❌ Error creating bundle: %v\n", err)
		return
	}

	fmt.Printf("✅ Bundle created: %s\n", bundleFile)
}

func handleVerifyBundle() {
	bundleFile := GetInput(ui.FormatPrompt("Enter bundle file path: "))
	if bundleFile == "" {
		fmt.Println(ui.FormatError("Bundle file path is required"))
		return
	}

	verification, err := git.VerifyBundle(bundleFile)
	if err != nil {
		fmt.Printf("❌ Error verifying bundle: %v\n", err)
		return
	}

	fmt.Printf("🔍 Bundle verification:\n%s\n", verification)

	refs, err := git.ListBundleRefs(bundleFile)
	if err == nil {
		fmt.Printf("📋 Bundle refs:\n%s\n", refs)
	}
}

func handleFormatPatchEmail() {
	since := GetInput(ui.FormatPrompt("Enter since reference (e.g., origin/main, HEAD~3): "))
	if since == "" {
		fmt.Println(ui.FormatError("Since reference is required"))
		return
	}

	patches, err := git.FormatPatchForEmail(since)
	if err != nil {
		fmt.Printf("❌ Error formatting patches: %v\n", err)
		return
	}

	fmt.Printf("📧 Email patches created:\n%s\n", patches)
}

// Worktree Management Handlers

func handleListWorktrees() {
	worktrees, err := git.ListWorktrees()
	if err != nil {
		fmt.Printf("❌ Error listing worktrees: %v\n", err)
		return
	}

	fmt.Printf("🌳 Worktrees:\n%s\n", worktrees)
}

func handleAddWorktree() {
	path := GetInput(ui.FormatPrompt("Enter worktree path: "))
	if path == "" {
		fmt.Println(ui.FormatError("Worktree path is required"))
		return
	}

	branch := GetInput(ui.FormatPrompt("Enter branch (optional, leave empty for new branch): "))

	if err := git.AddWorktree(path, branch); err != nil {
		fmt.Printf("❌ Error adding worktree: %v\n", err)
		return
	}

	fmt.Printf("✅ Worktree added: %s\n", path)
}

func handleRemoveWorktree() {
	worktrees, _ := git.ListWorktrees()
	fmt.Printf("Current worktrees:\n%s\n", worktrees)

	path := GetInput(ui.FormatPrompt("Enter worktree path to remove: "))
	if path == "" {
		fmt.Println(ui.FormatError("Worktree path is required"))
		return
	}

	if err := git.RemoveWorktree(path); err != nil {
		fmt.Printf("❌ Error removing worktree: %v\n", err)
		return
	}

	fmt.Printf("✅ Worktree removed: %s\n", path)
}

func handleMoveWorktree() {
	worktrees, _ := git.ListWorktrees()
	fmt.Printf("Current worktrees:\n%s\n", worktrees)

	oldPath := GetInput(ui.FormatPrompt("Enter current worktree path: "))
	newPath := GetInput(ui.FormatPrompt("Enter new worktree path: "))

	if oldPath == "" || newPath == "" {
		fmt.Println(ui.FormatError("Both old and new paths are required"))
		return
	}

	if err := git.MoveWorktree(oldPath, newPath); err != nil {
		fmt.Printf("❌ Error moving worktree: %v\n", err)
		return
	}

	fmt.Printf("✅ Worktree moved from %s to %s\n", oldPath, newPath)
}

// Repository Maintenance Handlers

func handleGarbageCollection() {
	fmt.Println(ui.FormatInfo("Running garbage collection (this may take a while)..."))
	if err := git.GarbageCollect(); err != nil {
		fmt.Printf("❌ Error running garbage collection: %v\n", err)
		return
	}
	fmt.Println("✅ Garbage collection completed!")
}

func handleVerifyRepository() {
	fmt.Println(ui.FormatInfo("Verifying repository integrity..."))
	result, err := git.VerifyRepository()
	if err != nil {
		fmt.Printf("❌ Error verifying repository: %v\n", err)
		return
	}

	if result == "" {
		fmt.Println("✅ Repository integrity verified - no issues found!")
	} else {
		fmt.Printf("🔍 Repository verification results:\n%s\n", result)
	}
}

func handleOptimizeRepository() {
	fmt.Println(ui.FormatInfo("Optimizing repository (this may take several minutes)..."))
	if err := git.OptimizeRepository(); err != nil {
		fmt.Printf("❌ Error optimizing repository: %v\n", err)
		return
	}
	fmt.Println("✅ Repository optimized successfully!")
}

func handleRepositoryStatistics() {
	fmt.Println(ui.FormatInfo("Gathering repository statistics..."))
	stats, err := git.RepositoryStatistics()
	if err != nil {
		fmt.Printf("❌ Error getting statistics: %v\n", err)
		return
	}

	fmt.Printf("📊 %s\n", stats)
}

func handleReflogManagement() {
	fmt.Println(ui.FormatInfo("Reflog Management"))
	fmt.Println("1. Show reflog")
	fmt.Println("2. Show reflog for specific ref")
	fmt.Println("3. Expire reflog entries")

	choice := GetInput(ui.FormatPrompt("Choose option (1-3): "))

	switch choice {
	case "1":
		reflog, err := git.Reflog()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("📋 Reflog:\n%s\n", reflog)

	case "2":
		ref := GetInput(ui.FormatPrompt("Enter ref (e.g., HEAD, main): "))
		reflog, err := git.ReflogShow(ref)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Printf("📋 Reflog for %s:\n%s\n", ref, reflog)

	case "3":
		confirm := GetInput(ui.FormatPrompt("⚠️  This will expire old reflog entries. Continue? (yes/no): "))
		if strings.ToLower(confirm) != "yes" {
			fmt.Println("Cancelled")
			return
		}
		if err := git.ReflogExpire(); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Reflog entries expired!")

	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handlePruneRemoteBranches() {
	remote := GetInput(ui.FormatPrompt("Enter remote name (default: origin): "))

	if err := git.PruneRemoteBranches(remote); err != nil {
		fmt.Printf("❌ Error pruning remote branches: %v\n", err)
		return
	}

	remoteName := remote
	if remoteName == "" {
		remoteName = "origin"
	}
	fmt.Printf("✅ Pruned stale branches from '%s'!\n", remoteName)
}

// Smart Git Operations Handlers

func handleInteractiveAdd() {
	fmt.Println(ui.FormatInfo("Starting interactive add (patch mode)..."))
	if err := git.InteractiveAdd(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Interactive add completed!")
}

func handlePartialCommit() {
	message := GetInput(ui.FormatPrompt("Enter commit message: "))
	if message == "" {
		fmt.Println(ui.FormatError("Commit message is required"))
		return
	}

	fmt.Println(ui.FormatInfo("Starting partial commit with interactive add..."))
	if err := git.PartialCommit(message); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Println("✅ Partial commit completed!")
}

func handleCommitAmend() {
	fmt.Println(ui.FormatInfo("Amend Last Commit"))
	fmt.Println("1. Amend without changing message")
	fmt.Println("2. Amend with new message")

	choice := GetInput(ui.FormatPrompt("Choose option (1-2): "))

	switch choice {
	case "1":
		if err := git.AmendLastCommit(""); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Last commit amended!")

	case "2":
		message := GetInput(ui.FormatPrompt("Enter new commit message: "))
		if message == "" {
			fmt.Println(ui.FormatError("Commit message is required"))
			return
		}
		if err := git.AmendLastCommit(message); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}
		fmt.Println("✅ Last commit amended with new message!")

	default:
		fmt.Println(ui.FormatError("Invalid choice"))
	}
}

func handleBranchComparisonTool() {
	branch1 := GetInput(ui.FormatPrompt("Enter first branch: "))
	branch2 := GetInput(ui.FormatPrompt("Enter second branch: "))

	if branch1 == "" || branch2 == "" {
		fmt.Println(ui.FormatError("Both branch names are required"))
		return
	}

	comparison, err := git.BranchComparison(branch1, branch2)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("🔄 %s\n", comparison)
}

func handleConflictPreventionCheck() {
	branch := GetInput(ui.FormatPrompt("Enter branch to check merge conflicts with: "))
	if branch == "" {
		fmt.Println(ui.FormatError("Branch name is required"))
		return
	}

	result, err := git.ConflictPreventionCheck(branch)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("%s\n", result)
}
