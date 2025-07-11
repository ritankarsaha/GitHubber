package git

import (
	"fmt"
	"strings"
)

type RebaseManager struct {
	executor CommandExecutor
}

func NewRebaseManager() *RebaseManager {
	return &RebaseManager{
		executor: NewCommandExecutor(),
	}
}

// Basic Rebase Operations

func (r *RebaseManager) Rebase(upstream string, options RebaseOptions) error {
	args := []string{"rebase"}
	
	if options.Interactive {
		args = append(args, "-i")
	}
	
	if options.Preserve {
		args = append(args, "-p")
	}
	
	if options.Strategy != "" {
		args = append(args, "-s", options.Strategy)
	}
	
	if len(options.StrategyOptions) > 0 {
		for _, opt := range options.StrategyOptions {
			args = append(args, "-X", opt)
		}
	}
	
	if options.Onto != "" {
		args = append(args, "--onto", options.Onto)
	}
	
	if options.Root {
		args = append(args, "--root")
	}
	
	if options.AutoSquash {
		args = append(args, "--autosquash")
	}
	
	if options.AutoStash {
		args = append(args, "--autostash")
	}
	
	if upstream != "" {
		args = append(args, upstream)
	}
	
	if options.Branch != "" {
		args = append(args, options.Branch)
	}
	
	_, err := r.executor.Execute("git", args...)
	return err
}

func (r *RebaseManager) InteractiveRebase(upstream string, options InteractiveOptions) error {
	rebaseOpts := RebaseOptions{
		Interactive: true,
		AutoSquash:  options.AutoSquash,
		AutoStash:   options.AutoStash,
	}
	
	if options.Editor != "" {
		// Set custom editor for this rebase
		// This would typically be done via environment variables
	}
	
	return r.Rebase(upstream, rebaseOpts)
}

// Rebase Control

func (r *RebaseManager) Continue() error {
	_, err := r.executor.Execute("git", "rebase", "--continue")
	return err
}

func (r *RebaseManager) Skip() error {
	_, err := r.executor.Execute("git", "rebase", "--skip")
	return err
}

func (r *RebaseManager) Abort() error {
	_, err := r.executor.Execute("git", "rebase", "--abort")
	return err
}

func (r *RebaseManager) Quit() error {
	_, err := r.executor.Execute("git", "rebase", "--quit")
	return err
}

// Rebase Status

func (r *RebaseManager) IsInProgress() (bool, error) {
	// Check for rebase directory
	_, err := r.executor.Execute("test", "-d", ".git/rebase-merge")
	if err == nil {
		return true, nil
	}
	
	_, err = r.executor.Execute("test", "-d", ".git/rebase-apply")
	return err == nil, nil
}

func (r *RebaseManager) GetStatus() (*RebaseStatus, error) {
	inProgress, err := r.IsInProgress()
	if err != nil {
		return nil, err
	}
	
	status := &RebaseStatus{
		InProgress: inProgress,
	}
	
	if !inProgress {
		return status, nil
	}
	
	// Get current step and total steps
	headName, err := r.executor.Execute("cat", ".git/rebase-merge/head-name")
	if err == nil {
		status.Branch = strings.TrimSpace(headName)
		status.Branch = strings.TrimPrefix(status.Branch, "refs/heads/")
	}
	
	onto, err := r.executor.Execute("cat", ".git/rebase-merge/onto")
	if err == nil {
		status.Onto = strings.TrimSpace(onto)
	}
	
	msgnum, err := r.executor.Execute("cat", ".git/rebase-merge/msgnum")
	if err == nil {
		status.CurrentStep = strings.TrimSpace(msgnum)
	}
	
	end, err := r.executor.Execute("cat", ".git/rebase-merge/end")
	if err == nil {
		status.TotalSteps = strings.TrimSpace(end)
	}
	
	return status, nil
}

// Advanced Rebase Operations

func (r *RebaseManager) SquashCommits(from, to string, message string) error {
	// This implements a squash operation using interactive rebase
	// In practice, this would require more sophisticated handling
	
	if from == "" || to == "" {
		return fmt.Errorf("both from and to commits must be specified")
	}
	
	// Use interactive rebase with auto-squash
	options := RebaseOptions{
		Interactive: true,
		AutoSquash:  true,
	}
	
	return r.Rebase(from+"^", options)
}

func (r *RebaseManager) EditCommit(commitHash string) error {
	// Start interactive rebase to edit a specific commit
	options := RebaseOptions{
		Interactive: true,
	}
	
	return r.Rebase(commitHash+"^", options)
}

func (r *RebaseManager) ReorderCommits(commitHashes []string) error {
	if len(commitHashes) < 2 {
		return fmt.Errorf("need at least 2 commits to reorder")
	}
	
	// Start interactive rebase from the earliest commit
	earliestCommit := commitHashes[len(commitHashes)-1]
	
	options := RebaseOptions{
		Interactive: true,
	}
	
	return r.Rebase(earliestCommit+"^", options)
}

// Rebase Utilities

func (r *RebaseManager) GetRebaseCommits(upstream string) ([]string, error) {
	// Get list of commits that would be rebased
	output, err := r.executor.Execute("git", "rev-list", "--reverse", upstream+"..")
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(output), "\n")
	commits := make([]string, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			commits = append(commits, line)
		}
	}
	
	return commits, nil
}

func (r *RebaseManager) ValidateRebase(upstream, branch string) error {
	// Check if the rebase would be valid
	if upstream == "" {
		return fmt.Errorf("upstream commit must be specified")
	}
	
	// Verify upstream exists
	_, err := r.executor.Execute("git", "rev-parse", "--verify", upstream)
	if err != nil {
		return fmt.Errorf("upstream commit %s does not exist: %w", upstream, err)
	}
	
	// Verify branch exists if specified
	if branch != "" {
		_, err := r.executor.Execute("git", "rev-parse", "--verify", branch)
		if err != nil {
			return fmt.Errorf("branch %s does not exist: %w", branch, err)
		}
	}
	
	return nil
}

func (r *RebaseManager) GetConflicts() ([]string, error) {
	inProgress, err := r.IsInProgress()
	if err != nil {
		return nil, err
	}
	
	if !inProgress {
		return []string{}, nil
	}
	
	// Get conflicted files
	output, err := r.executor.Execute("git", "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(output), "\n")
	conflicts := make([]string, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			conflicts = append(conflicts, line)
		}
	}
	
	return conflicts, nil
}

// Types and Options

type RebaseOptions struct {
	Interactive     bool
	Preserve        bool
	Strategy        string
	StrategyOptions []string
	Onto            string
	Root            bool
	AutoSquash      bool
	AutoStash       bool
	Branch          string
}

type InteractiveOptions struct {
	AutoSquash bool
	AutoStash  bool
	Editor     string
}

type RebaseStatus struct {
	InProgress  bool
	Branch      string
	Onto        string
	CurrentStep string
	TotalSteps  string
}