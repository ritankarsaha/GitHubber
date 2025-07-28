package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainWithoutArgs(t *testing.T) {
	// Create a test Git repository
	tmpDir, err := os.MkdirTemp("", "githubber-main-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure git for testing
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user name: %v", err)
	}

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Build the application
	buildCmd := exec.Command("go", "build", "-o", "githubber-test", filepath.Join(originalDir, "cmd", "main.go"))
	buildCmd.Dir = tmpDir
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}

	// Test running the application (it will exit immediately in test mode)
	// We can't easily test the interactive menu, but we can test that it starts
	cmd = exec.Command("./githubber-test")
	cmd.Dir = tmpDir
	
	// Run with timeout to prevent hanging
	err = cmd.Start()
	if err != nil {
		t.Errorf("Failed to start application: %v", err)
	}

	// Kill the process immediately since we can't interact with the menu
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func TestMainWithArgs(t *testing.T) {
	// Create a test Git repository
	tmpDir, err := os.MkdirTemp("", "githubber-main-args-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Change to test directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Build the application
	buildCmd := exec.Command("go", "build", "-o", "githubber-test", filepath.Join(originalDir, "cmd", "main.go"))
	buildCmd.Dir = tmpDir
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build application: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "help command",
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "version command",
			args:        []string{"--version"},
			expectError: false,
		},
		{
			name:        "invalid command",
			args:        []string{"--invalid-flag"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./githubber-test", tt.args...)
			cmd.Dir = tmpDir
			err := cmd.Run()
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error for args %v, but got none", tt.args)
			}
			
			// Note: We can't easily test the actual CLI functionality in unit tests
			// as it requires user interaction. Integration tests would be better for that.
		})
	}
}

func TestGitVersionCheck(t *testing.T) {
	// Test that Git is available in the system
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	if err != nil {
		t.Skip("Git is not installed or not in PATH - skipping test")
	}
}