package git

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type BranchManager struct {
	executor CommandExecutor
}

func NewBranchManager() *BranchManager {
	return &BranchManager{
		executor: NewCommandExecutor(),
	}
}

// Branch Creation and Deletion

func (b *BranchManager) Create(name string, startPoint string) error {
	args := []string{"checkout", "-b", name}
	if startPoint != "" {
		args = append(args, startPoint)
	}
	
	_, err := b.executor.Execute("git", args...)
	return err
}

func (b *BranchManager) CreateFromCommit(name, commitHash string) error {
	return b.Create(name, commitHash)
}

func (b *BranchManager) Delete(name string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	
	_, err := b.executor.Execute("git", "branch", flag, name)
	return err
}

func (b *BranchManager) DeleteRemote(remote, name string) error {
	_, err := b.executor.Execute("git", "push", remote, "--delete", name)
	return err
}

// Branch Navigation

func (b *BranchManager) Switch(name string) error {
	_, err := b.executor.Execute("git", "checkout", name)
	return err
}

func (b *BranchManager) SwitchOrCreate(name string) error {
	_, err := b.executor.Execute("git", "checkout", "-B", name)
	return err
}

func (b *BranchManager) SwitchToPrevious() error {
	_, err := b.executor.Execute("git", "checkout", "-")
	return err
}

// Branch Information

func (b *BranchManager) List(options ListBranchOptions) ([]*types.BranchInfo, error) {
	args := []string{"branch"}
	
	if options.All {
		args = append(args, "-a")
	} else if options.Remote {
		args = append(args, "-r")
	}
	
	if options.Verbose {
		args = append(args, "-v")
	}
	
	if options.Merged != "" {
		args = append(args, "--merged", options.Merged)
	}
	
	if options.NoMerged != "" {
		args = append(args, "--no-merged", options.NoMerged)
	}
	
	output, err := b.executor.Execute("git", args...)
	if err != nil {
		return nil, err
	}
	
	return b.parseBranchList(output), nil
}

func (b *BranchManager) GetCurrent() (string, error) {
	output, err := b.executor.Execute("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (b *BranchManager) Exists(name string) bool {
	_, err := b.executor.Execute("git", "rev-parse", "--verify", fmt.Sprintf("refs/heads/%s", name))
	return err == nil
}

func (b *BranchManager) ExistsRemote(remote, name string) bool {
	_, err := b.executor.Execute("git", "rev-parse", "--verify", fmt.Sprintf("refs/remotes/%s/%s", remote, name))
	return err == nil
}

// Branch Tracking

func (b *BranchManager) SetUpstream(branch, upstream string) error {
	_, err := b.executor.Execute("git", "branch", "--set-upstream-to", upstream, branch)
	return err
}

func (b *BranchManager) UnsetUpstream(branch string) error {
	_, err := b.executor.Execute("git", "branch", "--unset-upstream", branch)
	return err
}

func (b *BranchManager) GetUpstream(branch string) (string, error) {
	if branch == "" {
		current, err := b.GetCurrent()
		if err != nil {
			return "", err
		}
		branch = current
	}
	
	output, err := b.executor.Execute("git", "rev-parse", "--abbrev-ref", fmt.Sprintf("%s@{upstream}", branch))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Branch Comparison

func (b *BranchManager) GetAheadBehind(branch, upstream string) (ahead, behind int, err error) {
	if upstream == "" {
		upstream, err = b.GetUpstream(branch)
		if err != nil {
			return 0, 0, err
		}
	}
	
	output, err := b.executor.Execute("git", "rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", upstream, branch))
	if err != nil {
		return 0, 0, err
	}
	
	return utils.ParseAheadBehind(strings.TrimSpace(output))
}

func (b *BranchManager) GetCommitsBetween(from, to string) ([]*types.CommitInfo, error) {
	output, err := b.executor.Execute("git", "rev-list", "--oneline", fmt.Sprintf("%s..%s", from, to))
	if err != nil {
		return nil, err
	}
	
	return utils.ParseCommitList(output), nil
}

// Branch Renaming

func (b *BranchManager) Rename(oldName, newName string) error {
	_, err := b.executor.Execute("git", "branch", "-m", oldName, newName)
	return err
}

func (b *BranchManager) RenameCurrent(newName string) error {
	_, err := b.executor.Execute("git", "branch", "-m", newName)
	return err
}

// Helper functions

func (b *BranchManager) parseBranchList(output string) []*types.BranchInfo {
	lines := strings.Split(output, "\n")
	branches := make([]*types.BranchInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		branch := &types.BranchInfo{}
		
		// Check if current branch
		if strings.HasPrefix(line, "*") {
			branch.IsCurrent = true
			line = strings.TrimSpace(line[1:])
		}
		
		// Check if remote branch
		if strings.Contains(line, "remotes/") {
			branch.IsRemote = true
			line = strings.TrimPrefix(line, "remotes/")
		}
		
		// Extract branch name (first word)
		parts := strings.Fields(line)
		if len(parts) > 0 {
			branch.Name = parts[0]
			branches = append(branches, branch)
		}
	}
	
	return branches
}

type ListBranchOptions struct {
	All      bool
	Remote   bool
	Verbose  bool
	Merged   string
	NoMerged string
}