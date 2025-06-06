package git

import (
	"os"
	"os/exec"
	"testing"
)

// setupTestRepo creates a temporary Git repository for testing
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "git-tool-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to get current directory: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	if err := cmd.Run(); err != nil {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.name", "Test User")
	if err := cmd.Run(); err != nil {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git user name: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	if err := cmd.Run(); err != nil {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git user email: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.Chdir(originalDir)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// createTestFile creates a file with given content in the test repository
func createTestFile(t *testing.T, name, content string) {
	t.Helper()

	err := os.WriteFile(name, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
}

// createTestCommit creates a commit with the given message
func createTestCommit(t *testing.T, message string) {
	t.Helper()

	cmd := exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage changes: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", message)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}
}

// assertFileContent checks if a file has the expected content
func assertFileContent(t *testing.T, path, expected string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	if string(content) != expected {
		t.Errorf("File content mismatch.\nExpected: %s\nGot: %s", expected, string(content))
	}
}

// assertFileExists checks if a file exists
func assertFileExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", path)
	}
}

// assertDirExists checks if a directory exists
func assertDirExists(t *testing.T, path string) {
	t.Helper()

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Expected directory %s to exist", path)
		return
	}

	if !info.IsDir() {
		t.Errorf("Expected %s to be a directory", path)
	}
}

// assertGitBranch checks if we're on the expected git branch
func assertGitBranch(t *testing.T, expected string) {
	t.Helper()

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	branch := string(output)
	if branch != expected+"\n" {
		t.Errorf("Expected to be on branch %s, but was on %s", expected, branch)
	}
}
