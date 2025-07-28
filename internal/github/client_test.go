package github

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test creating a new client
	client, err := NewClient()
	if err != nil {
		t.Skip("Skipping NewClient test - requires authentication setup")
	}

	if client == nil {
		t.Errorf("NewClient() returned nil")
	}
}

func TestNewClientWithToken(t *testing.T) {
	// Test creating a client with token
	token := "test-token"
	client := NewClientWithToken(token)
	
	if client == nil {
		t.Errorf("NewClientWithToken() returned nil")
	}
}

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		name      string
		repoURL   string
		wantOwner string
		wantRepo  string
		wantError bool
	}{
		{
			name:      "https github url",
			repoURL:   "https://github.com/user/repo.git",
			wantOwner: "user",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "https github url without .git",
			repoURL:   "https://github.com/user/repo",
			wantOwner: "user",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "ssh github url",
			repoURL:   "git@github.com:user/repo.git",
			wantOwner: "user",
			wantRepo:  "repo",
			wantError: false,
		},
		{
			name:      "invalid url",
			repoURL:   "not-a-url",
			wantOwner: "",
			wantRepo:  "",
			wantError: true,
		},
		{
			name:      "non-github url",
			repoURL:   "https://gitlab.com/user/repo.git",
			wantOwner: "",
			wantRepo:  "",
			wantError: true,
		},
		{
			name:      "empty url",
			repoURL:   "",
			wantOwner: "",
			wantRepo:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := ParseRepoURL(tt.repoURL)

			if tt.wantError && err == nil {
				t.Errorf("parseRepoURL() expected error but got none")
			}

			if !tt.wantError && err != nil {
				t.Errorf("parseRepoURL() unexpected error: %v", err)
			}

			if owner != tt.wantOwner {
				t.Errorf("parseRepoURL() owner = %q, want %q", owner, tt.wantOwner)
			}

			if repo != tt.wantRepo {
				t.Errorf("parseRepoURL() repo = %q, want %q", repo, tt.wantRepo)
			}
		})
	}
}

func TestRepositoryStruct(t *testing.T) {
	// Test Repository struct
	repo := &Repository{
		Name:        "test-repo",
		Owner:       "test-owner", 
		Description: "Test repository",
		URL:         "https://github.com/test-owner/test-repo",
		Private:     false,
		Language:    "Go",
		Stars:       100,
		Forks:       25,
	}

	if repo.Name != "test-repo" {
		t.Errorf("Expected Name to be 'test-repo', got %q", repo.Name)
	}

	if repo.Stars != 100 {
		t.Errorf("Expected Stars to be 100, got %d", repo.Stars)
	}
}

func TestPullRequestStruct(t *testing.T) {
	// Test PullRequest struct
	pr := &PullRequest{
		Number: 123,
		Title:  "Test PR",
		State:  "open",
		Author: "test-user",
		URL:    "https://github.com/test/repo/pull/123",
	}

	if pr.Number != 123 {
		t.Errorf("Expected Number to be 123, got %d", pr.Number)
	}

	if pr.Title != "Test PR" {
		t.Errorf("Expected Title to be 'Test PR', got %q", pr.Title)
	}

	if pr.State != "open" {
		t.Errorf("Expected State to be 'open', got %q", pr.State)
	}
}

func TestIssueStruct(t *testing.T) {
	// Test Issue struct
	issue := &Issue{
		Number: 456,
		Title:  "Test Issue",
		State:  "open",
		Author: "issue-author",
		URL:    "https://github.com/test/repo/issues/456",
	}

	if issue.Number != 456 {
		t.Errorf("Expected Number to be 456, got %d", issue.Number)
	}

	if issue.Title != "Test Issue" {
		t.Errorf("Expected Title to be 'Test Issue', got %q", issue.Title)
	}

	if issue.State != "open" {
		t.Errorf("Expected State to be 'open', got %q", issue.State)
	}

	if issue.Author != "issue-author" {
		t.Errorf("Expected Author to be 'issue-author', got %q", issue.Author)
	}
}