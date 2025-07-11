package git

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type CommitManager struct {
	executor CommandExecutor
}

func NewCommitManager() *CommitManager {
	return &CommitManager{
		executor: NewCommandExecutor(),
	}
}

// Basic Commit Operations

func (c *CommitManager) Commit(message string, options CommitOptions) error {
	args := []string{"commit", "-m", message}
	
	if options.All {
		args = append(args, "-a")
	}
	if options.Amend {
		args = append(args, "--amend")
	}
	if options.NoEdit && options.Amend {
		args = append(args, "--no-edit")
	}
	if options.SignOff {
		args = append(args, "--signoff")
	}
	if options.NoVerify {
		args = append(args, "--no-verify")
	}
	if options.Author != "" {
		args = append(args, "--author", options.Author)
	}
	if options.Date != "" {
		args = append(args, "--date", options.Date)
	}
	
	_, err := c.executor.Execute("git", args...)
	return err
}

func (c *CommitManager) CommitWithFiles(message string, files []string, options CommitOptions) error {
	// Stage the files first
	if len(files) > 0 {
		stager := NewStager()
		if err := stager.AddFiles(files); err != nil {
			return fmt.Errorf("failed to stage files: %w", err)
		}
	}
	
	return c.Commit(message, options)
}

func (c *CommitManager) AmendLastCommit(newMessage string, noEdit bool) error {
	options := CommitOptions{
		Amend:  true,
		NoEdit: noEdit,
	}
	
	if newMessage != "" {
		return c.Commit(newMessage, options)
	}
	
	_, err := c.executor.Execute("git", "commit", "--amend", "--no-edit")
	return err
}

// Commit History

func (c *CommitManager) GetHistory(options LogOptions) ([]*types.CommitInfo, error) {
	args := []string{"log", "--pretty=format:%H|%h|%s|%an|%ae|%at"}
	
	if options.MaxCount > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", options.MaxCount))
	}
	if options.Skip > 0 {
		args = append(args, "--skip", fmt.Sprintf("%d", options.Skip))
	}
	if options.Since != "" {
		args = append(args, "--since", options.Since)
	}
	if options.Until != "" {
		args = append(args, "--until", options.Until)
	}
	if options.Author != "" {
		args = append(args, "--author", options.Author)
	}
	if options.Grep != "" {
		args = append(args, "--grep", options.Grep)
	}
	if options.OneLine {
		args = append(args, "--oneline")
	}
	if options.Graph {
		args = append(args, "--graph")
	}
	if options.All {
		args = append(args, "--all")
	}
	if len(options.Paths) > 0 {
		args = append(args, "--")
		args = append(args, options.Paths...)
	}
	
	output, err := c.executor.Execute("git", args...)
	if err != nil {
		return nil, err
	}
	
	return utils.ParseCommitHistory(output), nil
}

func (c *CommitManager) GetCommit(ref string) (*types.CommitInfo, error) {
	output, err := c.executor.Execute("git", "show", "-s", "--pretty=format:%H|%h|%s|%an|%ae|%at", ref)
	if err != nil {
		return nil, err
	}
	
	return utils.ParseCommitInfo(strings.TrimSpace(output))
}

func (c *CommitManager) GetCommitMessage(ref string) (string, error) {
	output, err := c.executor.Execute("git", "log", "-1", "--pretty=format:%B", ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (c *CommitManager) GetCommitFiles(ref string) ([]string, error) {
	output, err := c.executor.Execute("git", "show", "--name-only", "--pretty=format:", ref)
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

// Commit Manipulation

func (c *CommitManager) Reset(ref string, mode ResetMode) error {
	args := []string{"reset"}
	
	switch mode {
	case ResetSoft:
		args = append(args, "--soft")
	case ResetMixed:
		args = append(args, "--mixed")
	case ResetHard:
		args = append(args, "--hard")
	}
	
	args = append(args, ref)
	_, err := c.executor.Execute("git", args...)
	return err
}

func (c *CommitManager) Revert(ref string, options RevertOptions) error {
	args := []string{"revert"}
	
	if options.NoCommit {
		args = append(args, "--no-commit")
	}
	if options.NoEdit {
		args = append(args, "--no-edit")
	}
	if options.SignOff {
		args = append(args, "--signoff")
	}
	if options.MainlineParent > 0 {
		args = append(args, "-m", fmt.Sprintf("%d", options.MainlineParent))
	}
	
	args = append(args, ref)
	_, err := c.executor.Execute("git", args...)
	return err
}

func (c *CommitManager) CherryPick(ref string, options CherryPickOptions) error {
	args := []string{"cherry-pick"}
	
	if options.NoCommit {
		args = append(args, "--no-commit")
	}
	if options.Edit {
		args = append(args, "--edit")
	}
	if options.SignOff {
		args = append(args, "--signoff")
	}
	if options.MainlineParent > 0 {
		args = append(args, "-m", fmt.Sprintf("%d", options.MainlineParent))
	}
	
	args = append(args, ref)
	_, err := c.executor.Execute("git", args...)
	return err
}

// Commit Searching

func (c *CommitManager) FindCommitByMessage(pattern string) ([]*types.CommitInfo, error) {
	options := LogOptions{
		Grep:     pattern,
		MaxCount: 50,
	}
	return c.GetHistory(options)
}

func (c *CommitManager) FindCommitByAuthor(author string) ([]*types.CommitInfo, error) {
	options := LogOptions{
		Author:   author,
		MaxCount: 50,
	}
	return c.GetHistory(options)
}

func (c *CommitManager) FindCommitsByFile(filePath string) ([]*types.CommitInfo, error) {
	options := LogOptions{
		Paths:    []string{filePath},
		MaxCount: 50,
	}
	return c.GetHistory(options)
}

// Commit Validation

func (c *CommitManager) ValidateCommit(ref string) error {
	_, err := c.executor.Execute("git", "rev-parse", "--verify", ref)
	return err
}

func (c *CommitManager) GetCommitHash(ref string) (string, error) {
	output, err := c.executor.Execute("git", "rev-parse", ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (c *CommitManager) GetShortHash(ref string, length int) (string, error) {
	args := []string{"rev-parse"}
	if length > 0 {
		args = append(args, fmt.Sprintf("--short=%d", length))
	} else {
		args = append(args, "--short")
	}
	args = append(args, ref)
	
	output, err := c.executor.Execute("git", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Types and Options

type CommitOptions struct {
	All       bool
	Amend     bool
	NoEdit    bool
	SignOff   bool
	NoVerify  bool
	Author    string
	Date      string
}

type LogOptions struct {
	MaxCount int
	Skip     int
	Since    string
	Until    string
	Author   string
	Grep     string
	OneLine  bool
	Graph    bool
	All      bool
	Paths    []string
}

type ResetMode int

const (
	ResetSoft ResetMode = iota
	ResetMixed
	ResetHard
)

type RevertOptions struct {
	NoCommit       bool
	NoEdit         bool
	SignOff        bool
	MainlineParent int
}

type CherryPickOptions struct {
	NoCommit       bool
	Edit           bool
	SignOff        bool
	MainlineParent int
}