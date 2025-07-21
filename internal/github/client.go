/*
 * GitHubber - GitHub API Client
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: GitHub API integration and client management
 */

package github

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	ctx    context.Context
}

type Repository struct {
	Name        string
	Owner       string
	Description string
	URL         string
	Private     bool
	Language    string
	Stars       int
	Forks       int
}

type PullRequest struct {
	Number int
	Title  string
	State  string
	Author string
	URL    string
}

type Issue struct {
	Number int
	Title  string
	State  string
	Author string
	URL    string
}

// NewClient creates a new GitHub API client
func NewClient() (*Client, error) {
	ctx := context.Background()
	
	// Try to get token from environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// Try to get token from GitHub CLI if available
		token = getGitHubCLIToken()
		if token == "" {
			return nil, fmt.Errorf("GitHub token not found. Please set GITHUB_TOKEN environment variable or use 'gh auth login'")
		}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	
	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// NewClientWithToken creates a new GitHub API client with provided token
func NewClientWithToken(token string) *Client {
	ctx := context.Background()
	
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	
	return &Client{
		client: client,
		ctx:    ctx,
	}
}

// GetRepository gets repository information
func (c *Client) GetRepository(owner, repo string) (*Repository, error) {
	githubRepo, _, err := c.client.Repositories.Get(c.ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	repository := &Repository{
		Name:        githubRepo.GetName(),
		Owner:       githubRepo.GetOwner().GetLogin(),
		Description: githubRepo.GetDescription(),
		URL:         githubRepo.GetHTMLURL(),
		Private:     githubRepo.GetPrivate(),
		Language:    githubRepo.GetLanguage(),
		Stars:       githubRepo.GetStargazersCount(),
		Forks:       githubRepo.GetForksCount(),
	}

	return repository, nil
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(owner, repo, title, body, head, base string) (*PullRequest, error) {
	pr := &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	}

	createdPR, _, err := c.client.PullRequests.Create(c.ctx, owner, repo, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	pullRequest := &PullRequest{
		Number: createdPR.GetNumber(),
		Title:  createdPR.GetTitle(),
		State:  createdPR.GetState(),
		Author: createdPR.GetUser().GetLogin(),
		URL:    createdPR.GetHTMLURL(),
	}

	return pullRequest, nil
}

// ListPullRequests lists pull requests for a repository
func (c *Client) ListPullRequests(owner, repo string, state string) ([]*PullRequest, error) {
	opts := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	}

	prs, _, err := c.client.PullRequests.List(c.ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	var pullRequests []*PullRequest
	for _, pr := range prs {
		pullRequest := &PullRequest{
			Number: pr.GetNumber(),
			Title:  pr.GetTitle(),
			State:  pr.GetState(),
			Author: pr.GetUser().GetLogin(),
			URL:    pr.GetHTMLURL(),
		}
		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil
}

// ListIssues lists issues for a repository
func (c *Client) ListIssues(owner, repo string, state string) ([]*Issue, error) {
	opts := &github.IssueListByRepoOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	}

	issues, _, err := c.client.Issues.ListByRepo(c.ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}

	var issueList []*Issue
	for _, issue := range issues {
		// Skip pull requests (GitHub API treats PRs as issues)
		if issue.IsPullRequest() {
			continue
		}

		issueItem := &Issue{
			Number: issue.GetNumber(),
			Title:  issue.GetTitle(),
			State:  issue.GetState(),
			Author: issue.GetUser().GetLogin(),
			URL:    issue.GetHTMLURL(),
		}
		issueList = append(issueList, issueItem)
	}

	return issueList, nil
}

// GetUser gets the authenticated user information
func (c *Client) GetUser() (*github.User, error) {
	user, _, err := c.client.Users.Get(c.ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// TestConnection tests the GitHub API connection
func (c *Client) TestConnection() error {
	_, err := c.GetUser()
	return err
}

// Helper function to get GitHub CLI token
func getGitHubCLIToken() string {
	// This is a simplified approach - in a real implementation,
	// you might want to parse the GitHub CLI config files
	return ""
}

// ParseRepoURL parses a GitHub repository URL to extract owner and repo name
func ParseRepoURL(url string) (owner, repo string, err error) {
	// Handle different URL formats:
	// https://github.com/owner/repo
	// https://github.com/owner/repo.git
	// git@github.com:owner/repo.git
	
	url = strings.TrimSpace(url)
	
	if strings.HasPrefix(url, "git@github.com:") {
		// SSH URL format
		url = strings.TrimPrefix(url, "git@github.com:")
		url = strings.TrimSuffix(url, ".git")
		parts := strings.Split(url, "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid SSH URL format")
		}
		return parts[0], parts[1], nil
	}
	
	if strings.HasPrefix(url, "https://github.com/") {
		// HTTPS URL format
		url = strings.TrimPrefix(url, "https://github.com/")
		url = strings.TrimSuffix(url, ".git")
		parts := strings.Split(url, "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid HTTPS URL format")
		}
		return parts[0], parts[1], nil
	}
	
	return "", "", fmt.Errorf("unsupported URL format")
}