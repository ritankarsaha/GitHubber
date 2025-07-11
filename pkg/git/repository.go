package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type Repository struct {
	executor CommandExecutor
}

func NewRepository() *Repository {
	return &Repository{
		executor: NewCommandExecutor(),
	}
}

// Basic Repository Operations

func (r *Repository) Init(path string, bare bool) error {
	args := []string{"init"}
	if bare {
		args = append(args, "--bare")
	}
	if path != "" {
		args = append(args, path)
	}
	
	_, err := r.executor.Execute("git", args...)
	return err
}

func (r *Repository) Clone(url, destination string, options CloneOptions) error {
	args := []string{"clone"}
	
	if options.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", options.Depth))
	}
	if options.Branch != "" {
		args = append(args, "--branch", options.Branch)
	}
	if options.SingleBranch {
		args = append(args, "--single-branch")
	}
	if options.Recursive {
		args = append(args, "--recursive")
	}
	if options.Mirror {
		args = append(args, "--mirror")
	}
	if options.Bare {
		args = append(args, "--bare")
	}
	
	args = append(args, url)
	if destination != "" {
		args = append(args, destination)
	}
	
	_, err := r.executor.Execute("git", args...)
	return err
}

func (r *Repository) GetInfo() (*types.RepositoryInfo, error) {
	// Check if we're in a git repository
	if !r.IsGitRepository() {
		return nil, fmt.Errorf("not in a git repository")
	}
	
	// Get repository root
	rootPath, err := r.executor.Execute("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("failed to get repository root: %w", err)
	}
	
	// Get current branch
	currentBranch, err := r.GetCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	
	// Get remote URL
	url, _ := r.executor.Execute("git", "remote", "get-url", "origin")
	remoteName := "origin"
	
	return &types.RepositoryInfo{
		URL:           strings.TrimSpace(url),
		CurrentBranch: currentBranch,
		RemoteName:    remoteName,
		IsGitRepo:     true,
		RootPath:      strings.TrimSpace(rootPath),
	}, nil
}

func (r *Repository) IsGitRepository() bool {
	_, err := r.executor.Execute("git", "rev-parse", "--git-dir")
	return err == nil
}

func (r *Repository) GetCurrentBranch() (string, error) {
	output, err := r.executor.Execute("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (r *Repository) GetRootPath() (string, error) {
	output, err := r.executor.Execute("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (r *Repository) GetGitDir() (string, error) {
	output, err := r.executor.Execute("git", "rev-parse", "--git-dir")
	if err != nil {
		return "", err
	}
	gitDir := strings.TrimSpace(output)
	if !filepath.IsAbs(gitDir) {
		rootPath, err := r.GetRootPath()
		if err != nil {
			return "", err
		}
		gitDir = filepath.Join(rootPath, gitDir)
	}
	return gitDir, nil
}

func (r *Repository) IsWorkingTreeClean() (bool, error) {
	output, err := r.executor.Execute("git", "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) == "", nil
}

func (r *Repository) HasUncommittedChanges() (bool, error) {
	clean, err := r.IsWorkingTreeClean()
	return !clean, err
}

func (r *Repository) GetHeadCommit() (*types.CommitInfo, error) {
	return r.GetCommit("HEAD")
}

func (r *Repository) GetCommit(ref string) (*types.CommitInfo, error) {
	format := "--pretty=format:%H|%h|%s|%an|%ae|%at"
	output, err := r.executor.Execute("git", "show", "-s", format, ref)
	if err != nil {
		return nil, err
	}
	
	return utils.ParseCommitInfo(strings.TrimSpace(output))
}

func (r *Repository) ValidateRef(ref string) error {
	_, err := r.executor.Execute("git", "rev-parse", "--verify", ref)
	return err
}

type CloneOptions struct {
	Depth        int
	Branch       string
	SingleBranch bool
	Recursive    bool
	Mirror       bool
	Bare         bool
}