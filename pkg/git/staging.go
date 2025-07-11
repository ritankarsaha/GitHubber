package git

import (
	"fmt"
	"strings"

	"github.com/ritankarsaha/githubber/pkg/types"
	"github.com/ritankarsaha/githubber/pkg/utils"
)

type Stager struct {
	executor CommandExecutor
}

func NewStager() *Stager {
	return &Stager{
		executor: NewCommandExecutor(),
	}
}

// Basic Staging Operations

func (s *Stager) AddFiles(files []string) error {
	if len(files) == 0 {
		return fmt.Errorf("no files specified")
	}
	
	args := append([]string{"add"}, files...)
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *Stager) AddAll() error {
	_, err := s.executor.Execute("git", "add", ".")
	return err
}

func (s *Stager) AddByPattern(pattern string) error {
	_, err := s.executor.Execute("git", "add", pattern)
	return err
}

func (s *Stager) AddInteractive() error {
	_, err := s.executor.Execute("git", "add", "-i")
	return err
}

func (s *Stager) AddPatch(file string) error {
	args := []string{"add", "-p"}
	if file != "" {
		args = append(args, file)
	}
	_, err := s.executor.Execute("git", args...)
	return err
}

// Unstaging Operations

func (s *Stager) Reset(files []string) error {
	args := []string{"reset", "HEAD"}
	if len(files) > 0 {
		args = append(args, "--")
		args = append(args, files...)
	}
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *Stager) ResetAll() error {
	_, err := s.executor.Execute("git", "reset", "HEAD")
	return err
}

func (s *Stager) RestoreStaged(files []string) error {
	args := []string{"restore", "--staged"}
	if len(files) > 0 {
		args = append(args, files...)
	} else {
		args = append(args, ".")
	}
	_, err := s.executor.Execute("git", args...)
	return err
}

// Working Directory Operations

func (s *Stager) RestoreFiles(files []string) error {
	args := []string{"restore"}
	if len(files) > 0 {
		args = append(args, files...)
	} else {
		args = append(args, ".")
	}
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *Stager) CheckoutFiles(files []string) error {
	args := []string{"checkout", "HEAD", "--"}
	args = append(args, files...)
	_, err := s.executor.Execute("git", args...)
	return err
}

func (s *Stager) CleanUntracked(options CleanOptions) error {
	args := []string{"clean"}
	
	if options.DryRun {
		args = append(args, "-n")
	} else {
		args = append(args, "-f")
	}
	
	if options.Directories {
		args = append(args, "-d")
	}
	
	if options.IgnoredFiles {
		args = append(args, "-x")
	}
	
	if options.Interactive {
		args = append(args, "-i")
	}
	
	_, err := s.executor.Execute("git", args...)
	return err
}

// Status Operations

func (s *Stager) GetStatus() (*types.StatusInfo, error) {
	output, err := s.executor.Execute("git", "status", "--porcelain=v1", "-b")
	if err != nil {
		return nil, err
	}
	
	return s.parseStatus(output)
}

func (s *Stager) GetDetailedStatus() (string, error) {
	return s.executor.Execute("git", "status")
}

func (s *Stager) IsClean() (bool, error) {
	output, err := s.executor.Execute("git", "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) == "", nil
}

func (s *Stager) HasStagedChanges() (bool, error) {
	output, err := s.executor.Execute("git", "diff", "--cached", "--name-only")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) != "", nil
}

func (s *Stager) HasUnstagedChanges() (bool, error) {
	output, err := s.executor.Execute("git", "diff", "--name-only")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) != "", nil
}

func (s *Stager) GetStagedFiles() ([]types.FileStatusInfo, error) {
	output, err := s.executor.Execute("git", "diff", "--cached", "--name-status")
	if err != nil {
		return nil, err
	}
	
	return utils.ParseFileStatusList(output), nil
}

func (s *Stager) GetModifiedFiles() ([]types.FileStatusInfo, error) {
	output, err := s.executor.Execute("git", "diff", "--name-status")
	if err != nil {
		return nil, err
	}
	
	return utils.ParseFileStatusList(output), nil
}

func (s *Stager) GetUntrackedFiles() ([]string, error) {
	output, err := s.executor.Execute("git", "ls-files", "--others", "--exclude-standard")
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

// Diff Operations

func (s *Stager) DiffStaged(file string) (string, error) {
	args := []string{"diff", "--cached"}
	if file != "" {
		args = append(args, file)
	}
	return s.executor.Execute("git", args...)
}

func (s *Stager) DiffWorking(file string) (string, error) {
	args := []string{"diff"}
	if file != "" {
		args = append(args, file)
	}
	return s.executor.Execute("git", args...)
}

func (s *Stager) DiffStat(staged bool) (string, error) {
	args := []string{"diff", "--stat"}
	if staged {
		args = append(args, "--cached")
	}
	return s.executor.Execute("git", args...)
}

// Helper Methods

func (s *Stager) parseStatus(output string) (*types.StatusInfo, error) {
	lines := strings.Split(output, "\n")
	status := &types.StatusInfo{
		StagedFiles:     []types.FileStatusInfo{},
		ModifiedFiles:   []types.FileStatusInfo{},
		UntrackedFiles:  []string{},
		DeletedFiles:    []types.FileStatusInfo{},
		RenamedFiles:    []types.RenamedFileInfo{},
		ConflictedFiles: []types.FileStatusInfo{},
	}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Parse branch info
		if strings.HasPrefix(line, "##") {
			s.parseBranchInfo(line, status)
			continue
		}
		
		// Parse file status
		if len(line) >= 3 {
			s.parseFileStatus(line, status)
		}
	}
	
	status.IsClean = len(status.StagedFiles) == 0 &&
		len(status.ModifiedFiles) == 0 &&
		len(status.UntrackedFiles) == 0 &&
		len(status.DeletedFiles) == 0 &&
		len(status.ConflictedFiles) == 0
	
	return status, nil
}

func (s *Stager) parseBranchInfo(line string, status *types.StatusInfo) {
	// Format: ## branch...remote [ahead N, behind M]
	line = strings.TrimPrefix(line, "## ")
	
	if strings.Contains(line, "...") {
		parts := strings.Split(line, "...")
		status.Branch = parts[0]
		
		if len(parts) > 1 {
			remotePart := parts[1]
			// Parse ahead/behind info
			if strings.Contains(remotePart, "[") {
				if ahead, behind, err := utils.ParseAheadBehindFromStatus(remotePart); err == nil {
					status.Ahead = ahead
					status.Behind = behind
				}
			}
		}
	} else {
		status.Branch = line
	}
}

func (s *Stager) parseFileStatus(line string, status *types.StatusInfo) {
	if len(line) < 3 {
		return
	}
	
	indexStatus := line[0]
	workingStatus := line[1]
	fileName := line[3:]
	
	// Handle renames
	if strings.Contains(fileName, " -> ") {
		parts := strings.Split(fileName, " -> ")
		if len(parts) == 2 {
			renamed := types.RenamedFileInfo{
				OldPath: parts[0],
				NewPath: parts[1],
				Status:  string(indexStatus),
			}
			status.RenamedFiles = append(status.RenamedFiles, renamed)
			return
		}
	}
	
	fileInfo := types.FileStatusInfo{
		Path:   fileName,
		Status: string(indexStatus) + string(workingStatus),
	}
	
	// Categorize file based on status
	switch {
	case indexStatus != ' ' && indexStatus != '?':
		status.StagedFiles = append(status.StagedFiles, fileInfo)
	case workingStatus == 'M':
		status.ModifiedFiles = append(status.ModifiedFiles, fileInfo)
	case workingStatus == 'D':
		status.DeletedFiles = append(status.DeletedFiles, fileInfo)
	case indexStatus == '?' && workingStatus == '?':
		status.UntrackedFiles = append(status.UntrackedFiles, fileName)
	case indexStatus == 'U' || workingStatus == 'U':
		status.ConflictedFiles = append(status.ConflictedFiles, fileInfo)
	}
}

type CleanOptions struct {
	DryRun       bool
	Directories  bool
	IgnoredFiles bool
	Interactive  bool
}