package git

import (
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	tests := []struct {
		name    string
		command string
		want    string
		wantErr bool
	}{
		{
			name:    "valid command",
			command: "git status",
			want:    "On branch main",
			wantErr: false,
		},
		{
			name:    "invalid command",
			command: "git invalid-command",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunCommand(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !strings.Contains(got, tt.want) {
				t.Errorf("RunCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRepositoryInfo(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a test file and commit
	createTestFile(t, "test.txt", "test content")
	createTestCommit(t, "Initial commit")

	// Add remote
	_, err := RunCommand("git remote add origin https://github.com/test/repo.git")
	if err != nil {
		t.Fatalf("Failed to add remote: %v", err)
	}

	// Test GetRepositoryInfo
	info, err := GetRepositoryInfo()
	if err != nil {
		t.Fatalf("GetRepositoryInfo() error = %v", err)
	}

	// Check URL
	if info.URL != "https://github.com/test/repo.git" {
		t.Errorf("GetRepositoryInfo() URL = %v, want %v", info.URL, "https://github.com/test/repo.git")
	}

	// Check branch
	if info.CurrentBranch != "main" {
		t.Errorf("GetRepositoryInfo() CurrentBranch = %v, want %v", info.CurrentBranch, "main")
	}
}

func TestIsWorkingDirectoryClean(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Test with clean directory
	clean, err := IsWorkingDirectoryClean()
	if err != nil {
		t.Fatalf("IsWorkingDirectoryClean() error = %v", err)
	}
	if !clean {
		t.Error("IsWorkingDirectoryClean() = false, want true for clean directory")
	}

	// Create an untracked file
	createTestFile(t, "untracked.txt", "untracked content")

	// Test with untracked file
	clean, err = IsWorkingDirectoryClean()
	if err != nil {
		t.Fatalf("IsWorkingDirectoryClean() error = %v", err)
	}
	if clean {
		t.Error("IsWorkingDirectoryClean() = true, want false for directory with untracked files")
	}

	// Stage the file
	_, err = RunCommand("git add untracked.txt")
	if err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	// Test with staged file
	clean, err = IsWorkingDirectoryClean()
	if err != nil {
		t.Fatalf("IsWorkingDirectoryClean() error = %v", err)
	}
	if clean {
		t.Error("IsWorkingDirectoryClean() = true, want false for directory with staged files")
	}

	// Commit the file
	createTestCommit(t, "Add test file")

	// Test with committed file
	clean, err = IsWorkingDirectoryClean()
	if err != nil {
		t.Fatalf("IsWorkingDirectoryClean() error = %v", err)
	}
	if !clean {
		t.Error("IsWorkingDirectoryClean() = false, want true for clean directory after commit")
	}
}
