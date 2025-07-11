package git

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type RemoteManager struct {
	executor CommandExecutor
}

func NewRemoteManager() *RemoteManager {
	return &RemoteManager{
		executor: NewCommandExecutor(),
	}
}

// Remote Management

func (r *RemoteManager) Add(name, url string) error {
	_, err := r.executor.Execute("git", "remote", "add", name, url)
	return err
}

func (r *RemoteManager) Remove(name string) error {
	_, err := r.executor.Execute("git", "remote", "remove", name)
	return err
}

func (r *RemoteManager) Rename(oldName, newName string) error {
	_, err := r.executor.Execute("git", "remote", "rename", oldName, newName)
	return err
}

func (r *RemoteManager) SetURL(name, url string, push bool) error {
	args := []string{"remote", "set-url"}
	if push {
		args = append(args, "--push")
	}
	args = append(args, name, url)
	
	_, err := r.executor.Execute("git", args...)
	return err
}

func (r *RemoteManager) GetURL(name string, push bool) (string, error) {
	args := []string{"remote", "get-url"}
	if push {
		args = append(args, "--push")
	}
	args = append(args, name)
	
	output, err := r.executor.Execute("git", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Remote Information

func (r *RemoteManager) List(verbose bool) ([]*types.RemoteInfo, error) {
	args := []string{"remote"}
	if verbose {
		args = append(args, "-v")
	}
	
	output, err := r.executor.Execute("git", args...)
	if err != nil {
		return nil, err
	}
	
	if verbose {
		return utils.ParseRemoteList(output), nil
	}
	
	// Simple list without URLs
	lines := strings.Split(strings.TrimSpace(output), "\n")
	remotes := make([]*types.RemoteInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			remotes = append(remotes, &types.RemoteInfo{
				Name: line,
			})
		}
	}
	
	return remotes, nil
}

func (r *RemoteManager) Show(name string) (*types.RemoteInfo, error) {
	output, err := r.executor.Execute("git", "remote", "show", name)
	if err != nil {
		return nil, err
	}
	
	return r.parseRemoteShow(output, name), nil
}

func (r *RemoteManager) Exists(name string) bool {
	_, err := r.executor.Execute("git", "remote", "get-url", name)
	return err == nil
}

// Fetch and Pull Operations

func (r *RemoteManager) Fetch(remote string, options FetchOptions) error {
	args := []string{"fetch"}
	
	if options.All {
		args = append(args, "--all")
	} else if remote != "" {
		args = append(args, remote)
	}
	
	if options.Prune {
		args = append(args, "--prune")
	}
	
	if options.Tags {
		args = append(args, "--tags")
	}
	
	if options.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", options.Depth))
	}
	
	if len(options.Refspecs) > 0 {
		args = append(args, options.Refspecs...)
	}
	
	_, err := r.executor.Execute("git", args...)
	return err
}

func (r *RemoteManager) Pull(remote, branch string, options PullOptions) error {
	args := []string{"pull"}
	
	if options.Rebase {
		args = append(args, "--rebase")
	}
	
	if options.NoCommit {
		args = append(args, "--no-commit")
	}
	
	if options.Squash {
		args = append(args, "--squash")
	}
	
	if options.FastForwardOnly {
		args = append(args, "--ff-only")
	}
	
	if remote != "" {
		args = append(args, remote)
		if branch != "" {
			args = append(args, branch)
		}
	}
	
	_, err := r.executor.Execute("git", args...)
	return err
}

// Push Operations

func (r *RemoteManager) Push(remote, branch string, options PushOptions) error {
	args := []string{"push"}
	
	if options.Force {
		if options.ForceWithLease {
			args = append(args, "--force-with-lease")
		} else {
			args = append(args, "--force")
		}
	}
	
	if options.SetUpstream {
		args = append(args, "--set-upstream")
	}
	
	if options.All {
		args = append(args, "--all")
	}
	
	if options.Tags {
		args = append(args, "--tags")
	}
	
	if options.Delete {
		args = append(args, "--delete")
	}
	
	if remote != "" {
		args = append(args, remote)
		if branch != "" {
			args = append(args, branch)
		}
	}
	
	_, err := r.executor.Execute("git", args...)
	return err
}

func (r *RemoteManager) PushTags(remote string, force bool) error {
	args := []string{"push"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, remote, "--tags")
	
	_, err := r.executor.Execute("git", args...)
	return err
}

// Remote Tracking

func (r *RemoteManager) GetTrackingBranch(branch string) (string, error) {
	if branch == "" {
		// Get current branch
		current, err := r.executor.Execute("git", "rev-parse", "--abbrev-ref", "HEAD")
		if err != nil {
			return "", err
		}
		branch = strings.TrimSpace(current)
	}
	
	output, err := r.executor.Execute("git", "rev-parse", "--abbrev-ref", fmt.Sprintf("%s@{upstream}", branch))
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(output), nil
}

func (r *RemoteManager) SetTrackingBranch(localBranch, remoteBranch string) error {
	_, err := r.executor.Execute("git", "branch", "--set-upstream-to", remoteBranch, localBranch)
	return err
}

func (r *RemoteManager) UnsetTrackingBranch(branch string) error {
	_, err := r.executor.Execute("git", "branch", "--unset-upstream", branch)
	return err
}

// Remote Branch Operations

func (r *RemoteManager) GetRemoteBranches(remote string) ([]string, error) {
	args := []string{"branch", "-r"}
	if remote != "" {
		args = append(args, "--list", fmt.Sprintf("%s/*", remote))
	}
	
	output, err := r.executor.Execute("git", args...)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(output, "\n")
	branches := make([]string, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "->") {
			// Remove "origin/" prefix if present
			if remote != "" && strings.HasPrefix(line, remote+"/") {
				line = strings.TrimPrefix(line, remote+"/")
			}
			branches = append(branches, line)
		}
	}
	
	return branches, nil
}

func (r *RemoteManager) CreateTrackingBranch(localName, remoteBranch string) error {
	_, err := r.executor.Execute("git", "checkout", "-b", localName, remoteBranch)
	return err
}

// Cleanup Operations

func (r *RemoteManager) Prune(remote string) error {
	_, err := r.executor.Execute("git", "remote", "prune", remote)
	return err
}

func (r *RemoteManager) Update(remote string) error {
	if remote == "" {
		_, err := r.executor.Execute("git", "remote", "update")
		return err
	}
	
	_, err := r.executor.Execute("git", "remote", "update", remote)
	return err
}

// Helper Methods

func (r *RemoteManager) parseRemoteShow(output, name string) *types.RemoteInfo {
	remote := &types.RemoteInfo{
		Name: name,
	}
	
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Fetch URL:") {
			remote.FetchURL = strings.TrimSpace(strings.TrimPrefix(line, "Fetch URL:"))
		} else if strings.HasPrefix(line, "Push  URL:") {
			remote.PushURL = strings.TrimSpace(strings.TrimPrefix(line, "Push  URL:"))
		}
	}
	
	// Set URL to fetch URL if not set
	if remote.URL == "" {
		remote.URL = remote.FetchURL
	}
	
	return remote
}

// Options types

type FetchOptions struct {
	All      bool
	Prune    bool
	Tags     bool
	Depth    int
	Refspecs []string
}

type PullOptions struct {
	Rebase          bool
	NoCommit        bool
	Squash          bool
	FastForwardOnly bool
}

type PushOptions struct {
	Force           bool
	ForceWithLease  bool
	SetUpstream     bool
	All             bool
	Tags            bool
	Delete          bool
}