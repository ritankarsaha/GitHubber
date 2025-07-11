package git

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type TagManager struct {
	executor CommandExecutor
}

func NewTagManager() *TagManager {
	return &TagManager{
		executor: NewCommandExecutor(),
	}
}

// Tag Creation

func (t *TagManager) Create(name, message, ref string, options TagCreateOptions) error {
	args := []string{"tag"}
	
	if options.Annotated || message != "" {
		args = append(args, "-a")
	}
	
	if message != "" {
		args = append(args, "-m", message)
	}
	
	if options.Force {
		args = append(args, "-f")
	}
	
	if options.Sign {
		args = append(args, "-s")
	}
	
	if options.LocalUser != "" {
		args = append(args, "-u", options.LocalUser)
	}
	
	args = append(args, name)
	
	if ref != "" {
		args = append(args, ref)
	}
	
	_, err := t.executor.Execute("git", args...)
	return err
}

func (t *TagManager) CreateLightweight(name, ref string, force bool) error {
	args := []string{"tag"}
	
	if force {
		args = append(args, "-f")
	}
	
	args = append(args, name)
	
	if ref != "" {
		args = append(args, ref)
	}
	
	_, err := t.executor.Execute("git", args...)
	return err
}

func (t *TagManager) CreateAnnotated(name, message, ref string, force bool) error {
	options := TagCreateOptions{
		Annotated: true,
		Force:     force,
	}
	return t.Create(name, message, ref, options)
}

// Tag Deletion

func (t *TagManager) Delete(names []string) error {
	if len(names) == 0 {
		return fmt.Errorf("no tag names provided")
	}
	
	args := append([]string{"tag", "-d"}, names...)
	_, err := t.executor.Execute("git", args...)
	return err
}

func (t *TagManager) DeleteRemote(remote string, names []string) error {
	if len(names) == 0 {
		return fmt.Errorf("no tag names provided")
	}
	
	// Format tag names for deletion
	refspecs := make([]string, len(names))
	for i, name := range names {
		refspecs[i] = fmt.Sprintf(":refs/tags/%s", name)
	}
	
	args := append([]string{"push", remote}, refspecs...)
	_, err := t.executor.Execute("git", args...)
	return err
}

// Tag Information

func (t *TagManager) List(options TagListOptions) ([]*types.TagInfo, error) {
	args := []string{"tag"}
	
	if options.List != "" {
		args = append(args, "-l", options.List)
	}
	
	if options.Sort != "" {
		args = append(args, "--sort", options.Sort)
	}
	
	if options.Merged != "" {
		args = append(args, "--merged", options.Merged)
	}
	
	if options.NoMerged != "" {
		args = append(args, "--no-merged", options.NoMerged)
	}
	
	if options.Contains != "" {
		args = append(args, "--contains", options.Contains)
	}
	
	output, err := t.executor.Execute("git", args...)
	if err != nil {
		return nil, err
	}
	
	return utils.ParseTagList(output), nil
}

func (t *TagManager) Show(name string) (*types.TagInfo, error) {
	// Get tag info
	output, err := t.executor.Execute("git", "show", name)
	if err != nil {
		return nil, err
	}
	
	tag := &types.TagInfo{
		Name: name,
	}
	
	// Parse tag information
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "tag ") {
			tag.Name = strings.TrimSpace(strings.TrimPrefix(line, "tag "))
		} else if strings.HasPrefix(line, "Tagger:") {
			// Parse tagger and date info if needed
		} else if strings.HasPrefix(line, "Date:") {
			// Parse date if needed
		} else if line == "" && i > 0 {
			// Message starts after empty line
			if i+1 < len(lines) {
				tag.Message = strings.Join(lines[i+1:], "\n")
				break
			}
		}
	}
	
	// Check if it's an annotated tag
	_, err = t.executor.Execute("git", "cat-file", "-t", name)
	if err == nil {
		output, err := t.executor.Execute("git", "cat-file", "-t", name)
		if err == nil && strings.TrimSpace(output) == "tag" {
			tag.IsAnnotated = true
		}
	}
	
	// Get tag hash
	hash, err := t.executor.Execute("git", "rev-list", "-n", "1", name)
	if err == nil {
		tag.Hash = strings.TrimSpace(hash)
	}
	
	return tag, nil
}

func (t *TagManager) Exists(name string) bool {
	_, err := t.executor.Execute("git", "rev-parse", "--verify", fmt.Sprintf("refs/tags/%s", name))
	return err == nil
}

func (t *TagManager) GetHash(name string) (string, error) {
	output, err := t.executor.Execute("git", "rev-list", "-n", "1", name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (t *TagManager) GetMessage(name string) (string, error) {
	// Check if it's an annotated tag
	output, err := t.executor.Execute("git", "cat-file", "-t", name)
	if err != nil {
		return "", err
	}
	
	if strings.TrimSpace(output) != "tag" {
		return "", fmt.Errorf("tag %s is not annotated", name)
	}
	
	output, err = t.executor.Execute("git", "tag", "-l", "-n99", name)
	if err != nil {
		return "", err
	}
	
	// Parse message from output
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		// Remove tag name from first line
		firstLine := strings.TrimSpace(lines[0])
		if strings.HasPrefix(firstLine, name) {
			firstLine = strings.TrimSpace(strings.TrimPrefix(firstLine, name))
			lines[0] = firstLine
		}
		return strings.Join(lines, "\n"), nil
	}
	
	return "", nil
}

// Tag Operations

func (t *TagManager) Push(remote, name string, force bool) error {
	args := []string{"push"}
	
	if force {
		args = append(args, "--force")
	}
	
	args = append(args, remote, fmt.Sprintf("refs/tags/%s", name))
	_, err := t.executor.Execute("git", args...)
	return err
}

func (t *TagManager) PushAll(remote string, force bool) error {
	args := []string{"push"}
	
	if force {
		args = append(args, "--force")
	}
	
	args = append(args, remote, "--tags")
	_, err := t.executor.Execute("git", args...)
	return err
}

func (t *TagManager) Fetch(remote string, force bool) error {
	args := []string{"fetch"}
	
	if force {
		args = append(args, "--force")
	}
	
	args = append(args, remote, "refs/tags/*:refs/tags/*")
	_, err := t.executor.Execute("git", args...)
	return err
}

// Tag Verification

func (t *TagManager) Verify(name string) error {
	_, err := t.executor.Execute("git", "tag", "-v", name)
	return err
}

func (t *TagManager) IsAnnotated(name string) (bool, error) {
	output, err := t.executor.Execute("git", "cat-file", "-t", name)
	if err != nil {
		return false, err
	}
	
	return strings.TrimSpace(output) == "tag", nil
}

func (t *TagManager) IsSigned(name string) (bool, error) {
	output, err := t.executor.Execute("git", "tag", "-v", name)
	if err != nil {
		return false, nil // Not signed or verification failed
	}
	
	return strings.Contains(output, "Good signature") || strings.Contains(output, "gpg:"), nil
}

// Tag Searching

func (t *TagManager) FindByPattern(pattern string) ([]*types.TagInfo, error) {
	options := TagListOptions{
		List: pattern,
	}
	return t.List(options)
}

func (t *TagManager) FindByCommit(commitHash string) ([]*types.TagInfo, error) {
	output, err := t.executor.Execute("git", "tag", "--points-at", commitHash)
	if err != nil {
		return nil, err
	}
	
	return utils.ParseTagList(output), nil
}

func (t *TagManager) FindContaining(commitHash string) ([]*types.TagInfo, error) {
	options := TagListOptions{
		Contains: commitHash,
	}
	return t.List(options)
}

// Tag Utilities

func (t *TagManager) GetLatestTag() (*types.TagInfo, error) {
	output, err := t.executor.Execute("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return nil, err
	}
	
	tagName := strings.TrimSpace(output)
	return t.Show(tagName)
}

func (t *TagManager) GetNextVersion(prefix string) (string, error) {
	// This is a simple implementation that finds the latest numeric tag
	// and increments it. In a real implementation, you might want to use
	// semantic versioning logic.
	
	pattern := prefix + "*"
	if prefix == "" {
		pattern = "*"
	}
	
	options := TagListOptions{
		List: pattern,
		Sort: "-version:refname",
	}
	
	tags, err := t.List(options)
	if err != nil {
		return "", err
	}
	
	if len(tags) == 0 {
		return prefix + "1.0.0", nil
	}
	
	// This is a simplified version - you'd want proper semantic versioning here
	return fmt.Sprintf("%s%s-next", prefix, tags[0].Name), nil
}

// Options types

type TagCreateOptions struct {
	Annotated bool
	Force     bool
	Sign      bool
	LocalUser string
}

type TagListOptions struct {
	List     string
	Sort     string
	Merged   string
	NoMerged string
	Contains string
}