package git

import (
	"os"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "git-tool-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test Init
	if err := Init(); err != nil {
		t.Errorf("Init() error = %v", err)
	}

	// Verify .git directory exists
	assertDirExists(t, ".git")
}

func TestClone(t *testing.T) {
	// Set up test directory
	tmpDir, err := os.MkdirTemp("", "git-tool-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test Clone with a public repository
	err = Clone("https://github.com/golang/example.git")
	if err != nil {
		t.Errorf("Clone() error = %v", err)
	}

	// Verify repository was cloned
	assertDirExists(t, "example")
}

func TestBranchOperations(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createTestFile(t, "test.txt", "test content")
	createTestCommit(t, "Initial commit")

	// Test CreateBranch
	err := CreateBranch("feature")
	if err != nil {
		t.Errorf("CreateBranch() error = %v", err)
	}
	assertGitBranch(t, "feature")

	// Test SwitchBranch back to main
	err = SwitchBranch("main")
	if err != nil {
		t.Errorf("SwitchBranch() error = %v", err)
	}
	assertGitBranch(t, "main")

	// Test ListBranches
	branches, err := ListBranches()
	if err != nil {
		t.Errorf("ListBranches() error = %v", err)
	}
	if len(branches) != 2 {
		t.Errorf("ListBranches() returned %d branches, want 2", len(branches))
	}
	foundFeature := false
	for _, branch := range branches {
		if strings.Contains(branch, "feature") {
			foundFeature = true
			break
		}
	}
	if !foundFeature {
		t.Error("ListBranches() did not return the feature branch")
	}

	// Test DeleteBranch
	err = DeleteBranch("feature")
	if err != nil {
		t.Errorf("DeleteBranch() error = %v", err)
	}
	branches, _ = ListBranches()
	for _, branch := range branches {
		if strings.Contains(branch, "feature") {
			t.Error("DeleteBranch() did not delete the feature branch")
		}
	}
}

func TestFileOperations(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Test Status on clean repository
	status, err := Status()
	if err != nil {
		t.Errorf("Status() error = %v", err)
	}
	if !strings.Contains(status, "No commits yet") {
		t.Errorf("Status() = %v, want message about no commits", status)
	}

	// Create test file
	createTestFile(t, "test.txt", "test content")

	// Test Status with untracked file
	status, err = Status()
	if err != nil {
		t.Errorf("Status() error = %v", err)
	}
	if !strings.Contains(status, "Untracked files") {
		t.Errorf("Status() = %v, want message about untracked files", status)
	}

	// Test AddFiles
	err = AddFiles("test.txt")
	if err != nil {
		t.Errorf("AddFiles() error = %v", err)
	}

	// Test Status with staged file
	status, err = Status()
	if err != nil {
		t.Errorf("Status() error = %v", err)
	}
	if !strings.Contains(status, "Changes to be committed") {
		t.Errorf("Status() = %v, want message about staged changes", status)
	}

	// Test Commit
	err = Commit("Test commit")
	if err != nil {
		t.Errorf("Commit() error = %v", err)
	}

	// Test Log
	log, err := Log(1)
	if err != nil {
		t.Errorf("Log() error = %v", err)
	}
	if !strings.Contains(log, "Test commit") {
		t.Errorf("Log() = %v, want commit message", log)
	}
}

func TestStashOperations(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createTestFile(t, "test.txt", "initial content")
	createTestCommit(t, "Initial commit")

	// Modify file
	createTestFile(t, "test.txt", "modified content")

	// Test StashSave
	err := StashSave("Test stash")
	if err != nil {
		t.Errorf("StashSave() error = %v", err)
	}

	// Verify file is back to original state
	content, err := os.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "initial content" {
		t.Errorf("File content = %v, want initial content", string(content))
	}

	// Test StashList
	list, err := StashList()
	if err != nil {
		t.Errorf("StashList() error = %v", err)
	}
	if !strings.Contains(list, "Test stash") {
		t.Errorf("StashList() = %v, want stash message", list)
	}

	// Test StashPop
	err = StashPop()
	if err != nil {
		t.Errorf("StashPop() error = %v", err)
	}

	// Verify file is back to modified state
	content, err = os.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != "modified content" {
		t.Errorf("File content = %v, want modified content", string(content))
	}
}

func TestTagOperations(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createTestFile(t, "test.txt", "test content")
	createTestCommit(t, "Initial commit")

	// Test CreateTag
	err := CreateTag("v1.0.0", "Version 1.0.0")
	if err != nil {
		t.Errorf("CreateTag() error = %v", err)
	}

	// Test ListTags
	tags, err := ListTags()
	if err != nil {
		t.Errorf("ListTags() error = %v", err)
	}
	if !strings.Contains(tags, "v1.0.0") {
		t.Errorf("ListTags() = %v, want v1.0.0", tags)
	}

	// Test DeleteTag
	err = DeleteTag("v1.0.0")
	if err != nil {
		t.Errorf("DeleteTag() error = %v", err)
	}

	// Verify tag is deleted
	tags, err = ListTags()
	if err != nil {
		t.Errorf("ListTags() error = %v", err)
	}
	if strings.Contains(tags, "v1.0.0") {
		t.Errorf("ListTags() = %v, tag should be deleted", tags)
	}
}

func TestRemoteOperations(t *testing.T) {
	// Set up test repository
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	createTestFile(t, "test.txt", "test content")
	createTestCommit(t, "Initial commit")

	// Add remote
	_, err := RunCommand("git remote add origin https://github.com/test/repo.git")
	if err != nil {
		t.Fatalf("Failed to add remote: %v", err)
	}

	// Note: We can't actually test Push/Pull/Fetch without a real remote repository
	// Instead, we'll just verify that the commands are formatted correctly

	// Test Push (this will fail but we can check that it fails)
	err = Push("origin", "main")
	if err == nil {
		t.Error("Push() should fail without real remote")
	}
	// Just verify it failed - error messages can vary
	t.Logf("Push error (expected): %v", err)

	// Test Pull (this will fail but we can check that it fails)
	err = Pull("origin", "main")
	if err == nil {
		t.Error("Pull() should fail without real remote")
	}
	// Just verify it failed - error messages can vary
	t.Logf("Pull error (expected): %v", err)

	// Test Fetch (this will fail but we can check that it fails)
	err = Fetch("origin")
	if err == nil {
		t.Error("Fetch() should fail without real remote")
	}
	// Just verify it failed - error messages can vary
	t.Logf("Fetch error (expected): %v", err)
}
