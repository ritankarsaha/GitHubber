/*
 * GitHubber - CLI Arguments Parser
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Command-line argument parsing and routing
 */

package cli

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/git-tool/internal/git"
	"github.com/ritankarsaha/git-tool/internal/github"
	"github.com/ritankarsaha/git-tool/internal/ui"
)

// CommandInfo represents a CLI command
type CommandInfo struct {
	Name        string
	Usage       string
	Description string
	Handler     func(args []string) error
}

// GetCommands returns all available CLI commands
func GetCommands() map[string]CommandInfo {
	return map[string]CommandInfo{
		// Git Repository Commands
		"init": {
			Name:        "init",
			Usage:       "githubber init",
			Description: "Initialize a new Git repository",
			Handler:     handleInitCmd,
		},
		"clone": {
			Name:        "clone",
			Usage:       "githubber clone <repository-url>",
			Description: "Clone a Git repository",
			Handler:     handleCloneCmd,
		},

		// Branch Commands
		"branch": {
			Name:        "branch",
			Usage:       "githubber branch [create|delete|list|switch] [branch-name]",
			Description: "Manage Git branches",
			Handler:     handleBranchCmd,
		},
		"checkout": {
			Name:        "checkout",
			Usage:       "githubber checkout <branch-name>",
			Description: "Switch to a branch",
			Handler:     handleCheckoutCmd,
		},

		// Staging and Commit Commands
		"add": {
			Name:        "add",
			Usage:       "githubber add [files...]",
			Description: "Add files to staging area",
			Handler:     handleAddCmd,
		},
		"commit": {
			Name:        "commit",
			Usage:       "githubber commit -m <message>",
			Description: "Commit changes",
			Handler:     handleCommitCmd,
		},
		"status": {
			Name:        "status",
			Usage:       "githubber status",
			Description: "Show working tree status",
			Handler:     handleStatusCmd,
		},

		// Remote Commands
		"push": {
			Name:        "push",
			Usage:       "githubber push [remote] [branch]",
			Description: "Push changes to remote",
			Handler:     handlePushCmd,
		},
		"pull": {
			Name:        "pull",
			Usage:       "githubber pull [remote] [branch]",
			Description: "Pull changes from remote",
			Handler:     handlePullCmd,
		},
		"fetch": {
			Name:        "fetch",
			Usage:       "githubber fetch [remote]",
			Description: "Fetch changes from remote",
			Handler:     handleFetchCmd,
		},

		// History Commands
		"log": {
			Name:        "log",
			Usage:       "githubber log [-n <number>]",
			Description: "Show commit history",
			Handler:     handleLogCmd,
		},
		"diff": {
			Name:        "diff",
			Usage:       "githubber diff [file]",
			Description: "Show changes between commits",
			Handler:     handleDiffCmd,
		},

		// Advanced Git Commands
		"rebase": {
			Name:        "rebase",
			Usage:       "githubber rebase [-i] [base]",
			Description: "Reapply commits on top of another base tip",
			Handler:     handleRebaseCmd,
		},
		"cherry-pick": {
			Name:        "cherry-pick",
			Usage:       "githubber cherry-pick <commit-hash>",
			Description: "Apply changes from specific commits",
			Handler:     handleCherryPickCmd,
		},
		"reset": {
			Name:        "reset",
			Usage:       "githubber reset [--soft|--mixed|--hard] [commit]",
			Description: "Reset current HEAD to specified state",
			Handler:     handleResetCmd,
		},
		"revert": {
			Name:        "revert",
			Usage:       "githubber revert <commit-hash>",
			Description: "Revert a commit",
			Handler:     handleRevertCmd,
		},
		"merge": {
			Name:        "merge",
			Usage:       "githubber merge <branch>",
			Description: "Merge branches",
			Handler:     handleMergeCmd,
		},
		"bisect": {
			Name:        "bisect",
			Usage:       "githubber bisect [start|bad|good|reset]",
			Description: "Use binary search to find the commit that introduced a bug",
			Handler:     handleBisectCmd,
		},

		// Stash Commands
		"stash": {
			Name:        "stash",
			Usage:       "githubber stash [push|pop|list|show|drop] [options]",
			Description: "Temporarily store uncommitted changes",
			Handler:     handleStashCmd,
		},

		// Tag Commands
		"tag": {
			Name:        "tag",
			Usage:       "githubber tag [create|delete|list] [tag-name]",
			Description: "Manage Git tags",
			Handler:     handleTagCmd,
		},

		// GitHub Commands
		"github": {
			Name:        "github",
			Usage:       "githubber github [repo|pr|issue] [action] [options]",
			Description: "GitHub operations",
			Handler:     handleGitHubCmd,
		},
		"pr": {
			Name:        "pr",
			Usage:       "githubber pr [create|list|view|close|merge] [options]",
			Description: "Pull request operations",
			Handler:     handlePRCmd,
		},
		"issue": {
			Name:        "issue",
			Usage:       "githubber issue [create|list|view|close] [options]",
			Description: "Issue operations",
			Handler:     handleIssueCmd,
		},

		// Utility Commands
		"help": {
			Name:        "help",
			Usage:       "githubber help [command]",
			Description: "Show help information",
			Handler:     handleHelpCmd,
		},
		"version": {
			Name:        "version",
			Usage:       "githubber version",
			Description: "Show version information",
			Handler:     handleVersionCmd,
		},
		"completion": {
			Name:        "completion",
			Usage:       "githubber completion [bash|zsh|fish]",
			Description: "Generate shell completion scripts",
			Handler:     handleCompletionCmd,
		},
		"resolve-conflicts": {
			Name:        "resolve-conflicts",
			Usage:       "githubber resolve-conflicts",
			Description: "Interactive conflict resolution interface",
			Handler:     handleResolveConflictsCmd,
		},
	}
}

// ParseAndExecute parses command line arguments and executes the appropriate command
func ParseAndExecute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	commands := GetCommands()
	commandName := args[0]

	if command, exists := commands[commandName]; exists {
		return command.Handler(args[1:])
	}

	return fmt.Errorf("unknown command: %s", commandName)
}

// Command Handlers

func handleInitCmd(args []string) error {
	if err := git.Init(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Repository initialized successfully!"))
	return nil
}

func handleCloneCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("repository URL is required")
	}

	url := args[0]
	if err := git.Clone(url); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Repository cloned successfully!"))
	return nil
}

func handleBranchCmd(args []string) error {
	if len(args) == 0 {
		// List branches by default
		branches, err := git.ListBranches()
		if err != nil {
			return fmt.Errorf("failed to list branches: %w", err)
		}
		fmt.Println(ui.FormatInfo("Branches:"))
		for _, branch := range branches {
			fmt.Println(ui.FormatCode(branch))
		}
		return nil
	}

	action := args[0]
	switch action {
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("branch name is required")
		}
		if err := git.CreateBranch(args[1]); err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}
		fmt.Println(ui.FormatSuccess("Branch created successfully!"))
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("branch name is required")
		}
		if err := git.DeleteBranch(args[1]); err != nil {
			return fmt.Errorf("failed to delete branch: %w", err)
		}
		fmt.Println(ui.FormatSuccess("Branch deleted successfully!"))
	case "list":
		return handleBranchCmd([]string{}) // Recursive call to list
	case "switch":
		if len(args) < 2 {
			return fmt.Errorf("branch name is required")
		}
		if err := git.SwitchBranch(args[1]); err != nil {
			return fmt.Errorf("failed to switch branch: %w", err)
		}
		fmt.Println(ui.FormatSuccess("Switched to branch successfully!"))
	default:
		return fmt.Errorf("unknown branch action: %s", action)
	}
	return nil
}

func handleCheckoutCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("branch name is required")
	}

	if err := git.SwitchBranch(args[0]); err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Switched to branch successfully!"))
	return nil
}

func handleAddCmd(args []string) error {
	if len(args) == 0 {
		if err := git.AddFiles(); err != nil {
			return fmt.Errorf("failed to add files: %w", err)
		}
	} else {
		if err := git.AddFiles(args...); err != nil {
			return fmt.Errorf("failed to add files: %w", err)
		}
	}
	fmt.Println(ui.FormatSuccess("Files added successfully!"))
	return nil
}

func handleCommitCmd(args []string) error {
	var message string

	for i, arg := range args {
		if arg == "-m" && i+1 < len(args) {
			message = args[i+1]
			break
		}
	}

	if message == "" {
		return fmt.Errorf("commit message is required (use -m)")
	}

	if err := git.Commit(message); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Changes committed successfully!"))
	return nil
}

func handleStatusCmd(args []string) error {
	status, err := git.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	fmt.Printf("\n%s Git Status:\n%s\n", ui.IconRepository, status)
	return nil
}

func handlePushCmd(args []string) error {
	remote := "origin"
	branch := ""

	if len(args) > 0 {
		remote = args[0]
	}
	if len(args) > 1 {
		branch = args[1]
	}

	// Get current branch if not specified
	if branch == "" {
		repoInfo, err := git.GetRepositoryInfo()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		branch = repoInfo.CurrentBranch
	}

	if err := git.Push(remote, branch); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Changes pushed successfully!"))
	return nil
}

func handlePullCmd(args []string) error {
	remote := "origin"
	branch := ""

	if len(args) > 0 {
		remote = args[0]
	}
	if len(args) > 1 {
		branch = args[1]
	}

	// Get current branch if not specified
	if branch == "" {
		repoInfo, err := git.GetRepositoryInfo()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		branch = repoInfo.CurrentBranch
	}

	if err := git.Pull(remote, branch); err != nil {
		return fmt.Errorf("failed to pull changes: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Changes pulled successfully!"))
	return nil
}

func handleFetchCmd(args []string) error {
	remote := "origin"
	if len(args) > 0 {
		remote = args[0]
	}

	if err := git.Fetch(remote); err != nil {
		return fmt.Errorf("failed to fetch updates: %w", err)
	}
	fmt.Println(ui.FormatSuccess("Updates fetched successfully!"))
	return nil
}

func handleLogCmd(args []string) error {
	n := 10 // Default

	for i, arg := range args {
		if arg == "-n" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &n)
			break
		}
	}

	logs, err := git.Log(n)
	if err != nil {
		return fmt.Errorf("failed to get log: %w", err)
	}
	fmt.Printf("\n%s Last %d commits:\n%s\n", ui.IconHistory, n, logs)
	return nil
}

func handleDiffCmd(args []string) error {
	file := ""
	if len(args) > 0 {
		file = args[0]
	}

	diff, err := git.Diff(file)
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}
	fmt.Printf("\n%s Diff:\n%s\n", ui.IconCommit, diff)
	return nil
}

func handleStashCmd(args []string) error {
	if len(args) == 0 {
		return git.StashSave("WIP")
	}

	action := args[0]
	switch action {
	case "push":
		message := "WIP"
		if len(args) > 1 {
			message = strings.Join(args[1:], " ")
		}
		return git.StashSave(message)
	case "pop":
		return git.StashPop()
	case "list":
		list, err := git.StashList()
		if err != nil {
			return err
		}
		fmt.Printf("\n%s Stash list:\n%s\n", ui.IconStash, list)
		return nil
	default:
		return fmt.Errorf("unknown stash action: %s", action)
	}
}

func handleTagCmd(args []string) error {
	if len(args) == 0 {
		// List tags by default
		tags, err := git.ListTags()
		if err != nil {
			return fmt.Errorf("failed to list tags: %w", err)
		}
		fmt.Printf("\n%s Tags:\n%s\n", ui.IconTag, tags)
		return nil
	}

	action := args[0]
	switch action {
	case "create":
		if len(args) < 2 {
			return fmt.Errorf("tag name is required")
		}
		name := args[1]
		message := name
		if len(args) > 2 {
			message = strings.Join(args[2:], " ")
		}
		if err := git.CreateTag(name, message); err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}
		fmt.Println(ui.FormatSuccess("Tag created successfully!"))
	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("tag name is required")
		}
		if err := git.DeleteTag(args[1]); err != nil {
			return fmt.Errorf("failed to delete tag: %w", err)
		}
		fmt.Println(ui.FormatSuccess("Tag deleted successfully!"))
	case "list":
		return handleTagCmd([]string{}) // Recursive call to list
	default:
		return fmt.Errorf("unknown tag action: %s", action)
	}
	return nil
}

func handleHelpCmd(args []string) error {
	commands := GetCommands()

	if len(args) == 0 {
		fmt.Println(ui.FormatTitle("GitHubber - Advanced Git & GitHub CLI"))
		fmt.Println(ui.FormatInfo("Available Commands:"))
		fmt.Println()

		for _, cmd := range commands {
			fmt.Printf("  %-15s %s\n", cmd.Name, cmd.Description)
		}
		fmt.Println()
		fmt.Println("Use 'githubber help <command>' for more information about a specific command.")
		return nil
	}

	commandName := args[0]
	if cmd, exists := commands[commandName]; exists {
		fmt.Printf("Command: %s\n", cmd.Name)
		fmt.Printf("Usage: %s\n", cmd.Usage)
		fmt.Printf("Description: %s\n", cmd.Description)
	} else {
		return fmt.Errorf("unknown command: %s", commandName)
	}

	return nil
}

func handleVersionCmd(args []string) error {
	fmt.Println(ui.FormatTitle("GitHubber v2.0.0"))
	fmt.Println(ui.FormatInfo("Advanced Git & GitHub CLI Tool"))
	fmt.Println(ui.FormatSubtitle("Created by Ritankar Saha <ritankar.saha786@gmail.com>"))
	return nil
}

// Advanced Git Command Handlers

func handleRebaseCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("base commit is required")
	}

	interactive := false
	base := ""

	for _, arg := range args {
		if arg == "-i" {
			interactive = true
		} else if base == "" {
			base = arg
		} else if arg == "--continue" {
			return git.RebaseContinue()
		} else if arg == "--abort" {
			return git.RebaseAbort()
		} else if arg == "--skip" {
			return git.RebaseSkip()
		}
	}

	if base == "" {
		return fmt.Errorf("base commit is required")
	}

	if interactive {
		fmt.Println(ui.FormatInfo("Starting interactive rebase..."))
		if err := git.InteractiveRebase(base); err != nil {
			return fmt.Errorf("failed to start interactive rebase: %w", err)
		}
	} else {
		fmt.Println(ui.FormatInfo("Starting rebase..."))
		if _, err := git.RunCommand(fmt.Sprintf("git rebase %s", base)); err != nil {
			return fmt.Errorf("failed to rebase: %w", err)
		}
	}

	fmt.Println(ui.FormatSuccess("Rebase completed successfully!"))
	return nil
}

func handleCherryPickCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("commit hash is required")
	}

	for _, arg := range args {
		if arg == "--continue" {
			return git.CherryPickContinue()
		} else if arg == "--abort" {
			return git.CherryPickAbort()
		}
	}

	commitHash := args[0]
	if err := git.CherryPick(commitHash); err != nil {
		return fmt.Errorf("failed to cherry-pick commit: %w", err)
	}

	fmt.Println(ui.FormatSuccess("Commit cherry-picked successfully!"))
	return nil
}

func handleResetCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("commit is required")
	}

	var mode string
	var commit string

	for _, arg := range args {
		if arg == "--soft" {
			mode = "soft"
		} else if arg == "--mixed" {
			mode = "mixed"
		} else if arg == "--hard" {
			mode = "hard"
		} else if commit == "" {
			commit = arg
		}
	}

	if commit == "" {
		return fmt.Errorf("commit is required")
	}

	if mode == "" {
		mode = "mixed" // Default
	}

	var err error
	switch mode {
	case "soft":
		err = git.ResetSoft(commit)
	case "mixed":
		err = git.ResetMixed(commit)
	case "hard":
		err = git.ResetHard(commit)
	}

	if err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	fmt.Printf("%s Reset (%s) completed successfully!\n", ui.IconSuccess, mode)
	return nil
}

func handleRevertCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("commit hash is required")
	}

	commitHash := args[0]
	noCommit := false

	for _, arg := range args {
		if arg == "--no-commit" {
			noCommit = true
		}
	}

	var err error
	if noCommit {
		err = git.RevertNoCommit(commitHash)
	} else {
		err = git.Revert(commitHash)
	}

	if err != nil {
		return fmt.Errorf("failed to revert commit: %w", err)
	}

	fmt.Println(ui.FormatSuccess("Commit reverted successfully!"))
	return nil
}

func handleMergeCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("branch name is required")
	}

	branch := args[0]
	noFF := false
	squash := false

	for _, arg := range args {
		if arg == "--no-ff" {
			noFF = true
		} else if arg == "--squash" {
			squash = true
		} else if arg == "--abort" {
			return git.MergeAbort()
		} else if arg == "--continue" {
			return git.MergeContinue()
		}
	}

	var err error
	if noFF {
		err = git.MergeNoFF(branch)
	} else if squash {
		err = git.MergeSquash(branch)
	} else {
		err = git.Merge(branch)
	}

	if err != nil {
		// Check if it's a merge conflict
		if strings.Contains(err.Error(), "conflict") {
			fmt.Println(ui.FormatWarning("Merge conflicts detected!"))
			fmt.Println(ui.FormatInfo("Use 'githubber resolve-conflicts' to resolve them"))
			return nil
		}
		return fmt.Errorf("failed to merge: %w", err)
	}

	fmt.Println(ui.FormatSuccess("Merge completed successfully!"))
	return nil
}

func handleBisectCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("bisect action is required (start, bad, good, reset)")
	}

	action := args[0]
	var err error

	switch action {
	case "start":
		err = git.BisectStart()
		fmt.Println(ui.FormatInfo("Bisect started. Mark commits as 'good' or 'bad'"))
	case "bad":
		commit := ""
		if len(args) > 1 {
			commit = args[1]
		}
		err = git.BisectBad(commit)
		fmt.Println(ui.FormatInfo("Commit marked as bad"))
	case "good":
		commit := ""
		if len(args) > 1 {
			commit = args[1]
		}
		err = git.BisectGood(commit)
		fmt.Println(ui.FormatInfo("Commit marked as good"))
	case "reset":
		err = git.BisectReset()
		fmt.Println(ui.FormatSuccess("Bisect session reset"))
	case "skip":
		err = git.BisectSkip()
		fmt.Println(ui.FormatInfo("Commit skipped"))
	default:
		return fmt.Errorf("unknown bisect action: %s", action)
	}

	return err
}

// GitHub Command Handlers

func handleGitHubCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("github action is required (repo, pr, issue)")
	}

	action := args[0]
	switch action {
	case "repo":
		return handleRepoOperations(args[1:])
	case "pr":
		return handlePRCmd(args[1:])
	case "issue":
		return handleIssueCmd(args[1:])
	default:
		return fmt.Errorf("unknown github action: %s", action)
	}
}

func handleRepoOperations(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("repo action is required (info, create, fork, list)")
	}

	client, err := github.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	action := args[0]
	switch action {
	case "info":
		return showRepoInfo(client, args[1:])
	case "create":
		return createRepo(client, args[1:])
	case "fork":
		return forkRepo(client, args[1:])
	case "list":
		return listRepos(client, args[1:])
	default:
		return fmt.Errorf("unknown repo action: %s", action)
	}
}

func handlePRCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pr action is required (create, list, view, close, merge)")
	}

	client, err := github.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	action := args[0]
	switch action {
	case "create":
		return createPR(client, args[1:])
	case "list":
		return listPRs(client, args[1:])
	case "view":
		return viewPR(client, args[1:])
	case "close":
		return closePR(client, args[1:])
	case "merge":
		return mergePR(client, args[1:])
	default:
		return fmt.Errorf("unknown pr action: %s", action)
	}
}

func handleIssueCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("issue action is required (create, list, view, close)")
	}

	client, err := github.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	action := args[0]
	switch action {
	case "create":
		return createIssue(client, args[1:])
	case "list":
		return listIssues(client, args[1:])
	case "view":
		return viewIssue(client, args[1:])
	case "close":
		return closeIssue(client, args[1:])
	default:
		return fmt.Errorf("unknown issue action: %s", action)
	}
}

func handleCompletionCmd(args []string) error {
	if len(args) == 0 {
		fmt.Println(ui.FormatInfo("Available shells:"))
		for _, shell := range GetAvailableShells() {
			fmt.Printf("  - %s\n", shell)
		}
		return nil
	}

	shell := args[0]
	completion, err := GenerateCompletion(shell)
	if err != nil {
		return err
	}

	fmt.Print(completion)

	if len(args) > 1 && args[1] == "--instructions" {
		fmt.Println(ShowCompletionInstructions(shell))
	}

	return nil
}

// Helper functions for GitHub operations (simplified implementations)
func showRepoInfo(client *github.Client, args []string) error {
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	owner, repo, err := github.ParseRepoURL(repoInfo.URL)
	if err != nil {
		return fmt.Errorf("failed to parse repo URL: %w", err)
	}

	repository, err := client.GetRepository(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	fmt.Printf("Repository: %s/%s\n", repository.Owner, repository.Name)
	fmt.Printf("Description: %s\n", repository.Description)
	fmt.Printf("Language: %s\n", repository.Language)
	fmt.Printf("Stars: %d\n", repository.Stars)
	fmt.Printf("Forks: %d\n", repository.Forks)
	fmt.Printf("URL: %s\n", repository.URL)

	return nil
}

func createRepo(client *github.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("repository name is required")
	}

	name := args[0]
	description := ""
	private := false

	if len(args) > 1 {
		description = args[1]
	}
	for _, arg := range args {
		if arg == "--private" {
			private = true
		}
	}

	repo, err := client.CreateRepository(name, description, private)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	fmt.Printf("%s Repository created: %s\n", ui.IconSuccess, repo.URL)
	return nil
}

func forkRepo(client *github.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("owner and repo name are required")
	}

	owner := args[0]
	repo := args[1]

	forkedRepo, err := client.ForkRepository(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to fork repository: %w", err)
	}

	fmt.Printf("%s Repository forked: %s\n", ui.IconSuccess, forkedRepo.URL)
	return nil
}

func listRepos(client *github.Client, args []string) error {
	visibility := "all"
	if len(args) > 0 {
		visibility = args[0]
	}

	repos, err := client.ListRepositories(visibility)
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	fmt.Printf("Found %d repositories:\n", len(repos))
	for _, repo := range repos {
		fmt.Printf("  %s/%s (%s) - %d stars\n",
			repo.Owner, repo.Name, repo.Language, repo.Stars)
	}

	return nil
}

func createPR(client *github.Client, args []string) error {
	// Simplified implementation - in practice you'd parse more options
	fmt.Println(ui.FormatInfo("Creating pull request..."))
	fmt.Println(ui.FormatWarning("Use interactive mode for full PR creation functionality"))
	return nil
}

func listPRs(client *github.Client, args []string) error {
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	owner, repo, err := github.ParseRepoURL(repoInfo.URL)
	if err != nil {
		return fmt.Errorf("failed to parse repo URL: %w", err)
	}

	state := "open"
	if len(args) > 0 {
		state = args[0]
	}

	prs, err := client.ListPullRequests(owner, repo, state)
	if err != nil {
		return fmt.Errorf("failed to list pull requests: %w", err)
	}

	fmt.Printf("Found %d pull requests (%s):\n", len(prs), state)
	for _, pr := range prs {
		fmt.Printf("  #%d: %s (%s) by %s\n",
			pr.Number, pr.Title, pr.State, pr.Author)
	}

	return nil
}

func viewPR(client *github.Client, args []string) error {
	fmt.Println(ui.FormatInfo("Use interactive mode for detailed PR viewing"))
	return nil
}

func closePR(client *github.Client, args []string) error {
	fmt.Println(ui.FormatInfo("Use interactive mode for PR management"))
	return nil
}

func mergePR(client *github.Client, args []string) error {
	fmt.Println(ui.FormatInfo("Use interactive mode for PR management"))
	return nil
}

func createIssue(client *github.Client, args []string) error {
	fmt.Println(ui.FormatInfo("Use interactive mode for issue creation"))
	return nil
}

func listIssues(client *github.Client, args []string) error {
	repoInfo, err := git.GetRepositoryInfo()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	owner, repo, err := github.ParseRepoURL(repoInfo.URL)
	if err != nil {
		return fmt.Errorf("failed to parse repo URL: %w", err)
	}

	state := "open"
	if len(args) > 0 {
		state = args[0]
	}

	issues, err := client.ListIssues(owner, repo, state)
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	fmt.Printf("Found %d issues (%s):\n", len(issues), state)
	for _, issue := range issues {
		fmt.Printf("  #%d: %s (%s) by %s\n",
			issue.Number, issue.Title, issue.State, issue.Author)
	}

	return nil
}

func viewIssue(client *github.Client, args []string) error {
	fmt.Println(ui.FormatInfo("Use interactive mode for detailed issue viewing"))
	return nil
}

func closeIssue(client *github.Client, args []string) error {
	fmt.Println(ui.FormatInfo("Use interactive mode for issue management"))
	return nil
}

func handleResolveConflictsCmd(args []string) error {
	return StartConflictResolution()
}
