package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ritankarsaha/githubber/pkg/types"
)

// ParseCommitInfo parses a single commit line in format: hash|short_hash|message|author|email|timestamp
func ParseCommitInfo(line string) (*types.CommitInfo, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid commit format: %s", line)
	}
	
	timestamp, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %s", parts[5])
	}
	
	return &types.CommitInfo{
		Hash:        parts[0],
		ShortHash:   parts[1],
		Message:     parts[2],
		Author:      parts[3],
		AuthorEmail: parts[4],
		Date:        time.Unix(timestamp, 0),
	}, nil
}

// ParseCommitHistory parses multiple commit lines
func ParseCommitHistory(output string) []*types.CommitInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	commits := make([]*types.CommitInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		if commit, err := ParseCommitInfo(line); err == nil {
			commits = append(commits, commit)
		}
	}
	
	return commits
}

// ParseCommitList parses commit list from git rev-list output
func ParseCommitList(output string) []*types.CommitInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	commits := make([]*types.CommitInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Format: short_hash message
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			commits = append(commits, &types.CommitInfo{
				ShortHash: parts[0],
				Message:   parts[1],
			})
		}
	}
	
	return commits
}

// ParseAheadBehind parses ahead/behind count from git rev-list --count output
func ParseAheadBehind(output string) (ahead, behind int, err error) {
	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid ahead/behind format: %s", output)
	}
	
	behind, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid behind count: %s", parts[0])
	}
	
	ahead, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid ahead count: %s", parts[1])
	}
	
	return ahead, behind, nil
}

// ParseAheadBehindFromStatus parses ahead/behind from git status output
func ParseAheadBehindFromStatus(statusLine string) (ahead, behind int, err error) {
	// Format: "remote [ahead N, behind M]" or similar
	start := strings.Index(statusLine, "[")
	end := strings.Index(statusLine, "]")
	
	if start == -1 || end == -1 {
		return 0, 0, nil
	}
	
	info := statusLine[start+1 : end]
	parts := strings.Split(info, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "ahead ") {
			if count, err := strconv.Atoi(strings.TrimPrefix(part, "ahead ")); err == nil {
				ahead = count
			}
		} else if strings.HasPrefix(part, "behind ") {
			if count, err := strconv.Atoi(strings.TrimPrefix(part, "behind ")); err == nil {
				behind = count
			}
		}
	}
	
	return ahead, behind, nil
}

// ParseFileStatusList parses file status list from git diff --name-status output
func ParseFileStatusList(output string) []types.FileStatusInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	files := make([]types.FileStatusInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			files = append(files, types.FileStatusInfo{
				Status: parts[0],
				Path:   parts[1],
			})
		}
	}
	
	return files
}

// ParseTagList parses tag list from git tag output
func ParseTagList(output string) []*types.TagInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	tags := make([]*types.TagInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		tags = append(tags, &types.TagInfo{
			Name: line,
		})
	}
	
	return tags
}

// ParseRemoteList parses remote list from git remote output
func ParseRemoteList(output string) []*types.RemoteInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	remotes := make([]*types.RemoteInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			remote := &types.RemoteInfo{
				Name: parts[0],
				URL:  parts[1],
			}
			
			// Check for (fetch) or (push) suffix
			if len(parts) >= 3 {
				if strings.Contains(parts[2], "fetch") {
					remote.FetchURL = parts[1]
				} else if strings.Contains(parts[2], "push") {
					remote.PushURL = parts[1]
				}
			}
			
			remotes = append(remotes, remote)
		}
	}
	
	return remotes
}

// ParseStashList parses stash list from git stash list output
func ParseStashList(output string) []*types.StashInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	stashes := make([]*types.StashInfo, 0, len(lines))
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Format: stash@{0}: WIP on branch: message
		parts := strings.SplitN(line, ":", 3)
		if len(parts) >= 2 {
			// Extract stash index
			indexStr := strings.TrimPrefix(parts[0], "stash@{")
			indexStr = strings.TrimSuffix(indexStr, "}")
			
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				continue
			}
			
			stash := &types.StashInfo{
				Index: index,
			}
			
			if len(parts) >= 3 {
				// Extract branch and message
				branchPart := strings.TrimSpace(parts[1])
				if strings.HasPrefix(branchPart, "WIP on ") || strings.HasPrefix(branchPart, "On ") {
					branchName := strings.Fields(branchPart)
					if len(branchName) >= 3 {
						stash.Branch = branchName[len(branchName)-1]
					}
				}
				
				stash.Message = strings.TrimSpace(parts[2])
			}
			
			stashes = append(stashes, stash)
		}
	}
	
	return stashes
}

// SanitizeRef sanitizes a git reference name
func SanitizeRef(ref string) string {
	// Remove invalid characters and replace with valid ones
	ref = strings.ReplaceAll(ref, " ", "-")
	ref = strings.ReplaceAll(ref, "/", "-")
	ref = strings.ReplaceAll(ref, "\\", "-")
	ref = strings.ReplaceAll(ref, ":", "-")
	ref = strings.ReplaceAll(ref, "?", "")
	ref = strings.ReplaceAll(ref, "*", "")
	ref = strings.ReplaceAll(ref, "[", "")
	ref = strings.ReplaceAll(ref, "]", "")
	ref = strings.ReplaceAll(ref, "~", "-")
	ref = strings.ReplaceAll(ref, "^", "")
	ref = strings.ReplaceAll(ref, "..", "-")
	
	// Remove leading/trailing dots and dashes
	ref = strings.Trim(ref, ".-")
	
	return ref
}

// ValidateRef validates if a string is a valid git reference
func ValidateRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("reference cannot be empty")
	}
	
	// Check for invalid characters
	invalidChars := []string{" ", "\\", ":", "?", "*", "[", "]", "~", "^", ".."}
	for _, char := range invalidChars {
		if strings.Contains(ref, char) {
			return fmt.Errorf("reference contains invalid character: %s", char)
		}
	}
	
	// Check for leading/trailing dots
	if strings.HasPrefix(ref, ".") || strings.HasSuffix(ref, ".") {
		return fmt.Errorf("reference cannot start or end with a dot")
	}
	
	return nil
}