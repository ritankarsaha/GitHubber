package git

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type StashManager struct {
	executor CommandExecutor
}

func NewStashManager() *StashManager {
	return &StashManager{
		executor: NewCommandExecutor(),
	}
}

// Basic Stash Operations

func (s *StashManager) Save(message string, options StashSaveOptions) error {
	args := []string{"stash", "push"}
	
	if message != "" {
		args = append(args, "-m", message)
	}
	
	if options.KeepIndex {
		args = append(args, "--keep-index")
	}
	
	if options.IncludeUntracked {
		args = append(args, "--include-untracked")
	}
	
	if options.All {
		args = append(args, "--all")
	}
	
	if len(options.Pathspec) > 0 {
		args = append(args, "--")
		args = append(args, options.Pathspec...)
	}
	
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *StashManager) Pop(stashRef string) error {
	args := []string{"stash", "pop"}
	if stashRef != "" {
		args = append(args, stashRef)
	}
	
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *StashManager) Apply(stashRef string, index bool) error {
	args := []string{"stash", "apply"}
	
	if index {
		args = append(args, "--index")
	}
	
	if stashRef != "" {
		args = append(args, stashRef)
	}
	
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *StashManager) Drop(stashRef string) error {
	args := []string{"stash", "drop"}
	if stashRef != "" {
		args = append(args, stashRef)
	}
	
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *StashManager) Clear() error {
	_, err := s.executor.Execute("git", "stash", "clear")
	return err
}

// Stash Information

func (s *StashManager) List() ([]*types.StashInfo, error) {
	output, err := s.executor.Execute("git", "stash", "list")
	if err != nil {
		return nil, err
	}
	
	return utils.ParseStashList(output), nil
}

func (s *StashManager) Show(stashRef string, options StashShowOptions) (string, error) {
	args := []string{"stash", "show"}
	
	if options.Patch {
		args = append(args, "-p")
	}
	
	if options.Stat {
		args = append(args, "--stat")
	}
	
	if options.NameOnly {
		args = append(args, "--name-only")
	}
	
	if stashRef != "" {
		args = append(args, stashRef)
	}
	
	return s.executor.Execute("git", args...)
}

func (s *StashManager) Exists(stashRef string) bool {
	_, err := s.executor.Execute("git", "rev-parse", "--verify", stashRef)
	return err == nil
}

func (s *StashManager) Count() (int, error) {
	stashes, err := s.List()
	if err != nil {
		return 0, err
	}
	return len(stashes), nil
}

// Advanced Stash Operations

func (s *StashManager) Branch(branchName, stashRef string) error {
	args := []string{"stash", "branch", branchName}
	if stashRef != "" {
		args = append(args, stashRef)
	}
	
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *StashManager) Store(ref, message string) error {
	args := []string{"stash", "store"}
	if message != "" {
		args = append(args, "-m", message)
	}
	args = append(args, ref)
	
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *StashManager) Create(message string) (string, error) {
	args := []string{"stash", "create"}
	if message != "" {
		args = append(args, message)
	}
	
	output, err := s.executor.Execute("git", args...)
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(output), nil
}

// Stash Utilities

func (s *StashManager) GetStashHash(stashRef string) (string, error) {
	if stashRef == "" {
		stashRef = "stash@{0}"
	}
	
	output, err := s.executor.Execute("git", "rev-parse", stashRef)
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(output), nil
}

func (s *StashManager) GetStashMessage(stashRef string) (string, error) {
	if stashRef == "" {
		stashRef = "stash@{0}"
	}
	
	output, err := s.executor.Execute("git", "log", "-1", "--pretty=format:%s", stashRef)
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(output), nil
}

func (s *StashManager) GetStashFiles(stashRef string) ([]string, error) {
	if stashRef == "" {
		stashRef = "stash@{0}"
	}
	
	output, err := s.executor.Execute("git", "stash", "show", "--name-only", stashRef)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(output), "\n")
	files := make([]string, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	
	return files, nil
}

// Stash Validation

func (s *StashManager) ValidateStashRef(stashRef string) error {
	if stashRef == "" {
		return fmt.Errorf("stash reference cannot be empty")
	}
	
	if !strings.HasPrefix(stashRef, "stash@{") || !strings.HasSuffix(stashRef, "}") {
		return fmt.Errorf("invalid stash reference format: %s", stashRef)
	}
	
	return nil
}

func (s *StashManager) GetLatestStash() (*types.StashInfo, error) {
	stashes, err := s.List()
	if err != nil {
		return nil, err
	}
	
	if len(stashes) == 0 {
		return nil, fmt.Errorf("no stashes found")
	}
	
	return stashes[0], nil
}

// Options types

type StashSaveOptions struct {
	KeepIndex        bool
	IncludeUntracked bool
	All              bool
	Pathspec         []string
}

type StashShowOptions struct {
	Patch    bool
	Stat     bool
	NameOnly bool
}