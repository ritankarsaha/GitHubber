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

// Enhanced GitHub Operations - CRUD for Issues and Pull Requests

// UpdatePullRequest updates an existing pull request
func (c *Client) UpdatePullRequest(owner, repo string, number int, title, body string) (*PullRequest, error) {
	pr := &github.PullRequest{
		Title: &title,
		Body:  &body,
	}

	updatedPR, _, err := c.client.PullRequests.Edit(c.ctx, owner, repo, number, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to update pull request: %w", err)
	}

	return &PullRequest{
		Number: updatedPR.GetNumber(),
		Title:  updatedPR.GetTitle(),
		State:  updatedPR.GetState(),
		Author: updatedPR.GetUser().GetLogin(),
		URL:    updatedPR.GetHTMLURL(),
	}, nil
}

// ClosePullRequest closes a pull request
func (c *Client) ClosePullRequest(owner, repo string, number int) error {
	state := "closed"
	pr := &github.PullRequest{
		State: &state,
	}

	_, _, err := c.client.PullRequests.Edit(c.ctx, owner, repo, number, pr)
	if err != nil {
		return fmt.Errorf("failed to close pull request: %w", err)
	}
	return nil
}

// MergePullRequest merges a pull request
func (c *Client) MergePullRequest(owner, repo string, number int, commitMessage, mergeMethod string) error {
	options := &github.PullRequestOptions{
		CommitTitle: commitMessage,
		MergeMethod: mergeMethod,
	}

	_, _, err := c.client.PullRequests.Merge(c.ctx, owner, repo, number, commitMessage, options)
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}
	return nil
}

// GetPullRequest gets a specific pull request
func (c *Client) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(c.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &PullRequest{
		Number: pr.GetNumber(),
		Title:  pr.GetTitle(),
		State:  pr.GetState(),
		Author: pr.GetUser().GetLogin(),
		URL:    pr.GetHTMLURL(),
	}, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(owner, repo, title, body string, labels []string) (*Issue, error) {
	issue := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}

	createdIssue, _, err := c.client.Issues.Create(c.ctx, owner, repo, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return &Issue{
		Number: createdIssue.GetNumber(),
		Title:  createdIssue.GetTitle(),
		State:  createdIssue.GetState(),
		Author: createdIssue.GetUser().GetLogin(),
		URL:    createdIssue.GetHTMLURL(),
	}, nil
}

// UpdateIssue updates an existing issue
func (c *Client) UpdateIssue(owner, repo string, number int, title, body string, labels []string) (*Issue, error) {
	issue := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}

	updatedIssue, _, err := c.client.Issues.Edit(c.ctx, owner, repo, number, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	return &Issue{
		Number: updatedIssue.GetNumber(),
		Title:  updatedIssue.GetTitle(),
		State:  updatedIssue.GetState(),
		Author: updatedIssue.GetUser().GetLogin(),
		URL:    updatedIssue.GetHTMLURL(),
	}, nil
}

// CloseIssue closes an issue
func (c *Client) CloseIssue(owner, repo string, number int) error {
	state := "closed"
	issue := &github.IssueRequest{
		State: &state,
	}

	_, _, err := c.client.Issues.Edit(c.ctx, owner, repo, number, issue)
	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}
	return nil
}

// ReopenIssue reopens a closed issue
func (c *Client) ReopenIssue(owner, repo string, number int) error {
	state := "open"
	issue := &github.IssueRequest{
		State: &state,
	}

	_, _, err := c.client.Issues.Edit(c.ctx, owner, repo, number, issue)
	if err != nil {
		return fmt.Errorf("failed to reopen issue: %w", err)
	}
	return nil
}

// GetIssue gets a specific issue
func (c *Client) GetIssue(owner, repo string, number int) (*Issue, error) {
	issue, _, err := c.client.Issues.Get(c.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return &Issue{
		Number: issue.GetNumber(),
		Title:  issue.GetTitle(),
		State:  issue.GetState(),
		Author: issue.GetUser().GetLogin(),
		URL:    issue.GetHTMLURL(),
	}, nil
}

// CommentOnIssue adds a comment to an issue
func (c *Client) CommentOnIssue(owner, repo string, number int, body string) error {
	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := c.client.Issues.CreateComment(c.ctx, owner, repo, number, comment)
	if err != nil {
		return fmt.Errorf("failed to comment on issue: %w", err)
	}
	return nil
}

// CommentOnPullRequest adds a comment to a pull request
func (c *Client) CommentOnPullRequest(owner, repo string, number int, body string) error {
	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := c.client.Issues.CreateComment(c.ctx, owner, repo, number, comment)
	if err != nil {
		return fmt.Errorf("failed to comment on pull request: %w", err)
	}
	return nil
}

// ListRepositories lists repositories for the authenticated user
func (c *Client) ListRepositories(visibility string) ([]*Repository, error) {
	opts := &github.RepositoryListOptions{
		Visibility: visibility,
		Sort:       "updated",
		Direction:  "desc",
		ListOptions: github.ListOptions{
			PerPage: 30,
		},
	}

	repos, _, err := c.client.Repositories.List(c.ctx, "", opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	var repositories []*Repository
	for _, repo := range repos {
		repository := &Repository{
			Name:        repo.GetName(),
			Owner:       repo.GetOwner().GetLogin(),
			Description: repo.GetDescription(),
			URL:         repo.GetHTMLURL(),
			Private:     repo.GetPrivate(),
			Language:    repo.GetLanguage(),
			Stars:       repo.GetStargazersCount(),
			Forks:       repo.GetForksCount(),
		}
		repositories = append(repositories, repository)
	}

	return repositories, nil
}

// CreateRepository creates a new repository
func (c *Client) CreateRepository(name, description string, private bool) (*Repository, error) {
	repo := &github.Repository{
		Name:        &name,
		Description: &description,
		Private:     &private,
	}

	createdRepo, _, err := c.client.Repositories.Create(c.ctx, "", repo)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return &Repository{
		Name:        createdRepo.GetName(),
		Owner:       createdRepo.GetOwner().GetLogin(),
		Description: createdRepo.GetDescription(),
		URL:         createdRepo.GetHTMLURL(),
		Private:     createdRepo.GetPrivate(),
		Language:    createdRepo.GetLanguage(),
		Stars:       createdRepo.GetStargazersCount(),
		Forks:       createdRepo.GetForksCount(),
	}, nil
}

// DeleteRepository deletes a repository
func (c *Client) DeleteRepository(owner, repo string) error {
	_, err := c.client.Repositories.Delete(c.ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}
	return nil
}

// ForkRepository forks a repository
func (c *Client) ForkRepository(owner, repo string) (*Repository, error) {
	forkedRepo, _, err := c.client.Repositories.CreateFork(c.ctx, owner, repo, &github.RepositoryCreateForkOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fork repository: %w", err)
	}

	return &Repository{
		Name:        forkedRepo.GetName(),
		Owner:       forkedRepo.GetOwner().GetLogin(),
		Description: forkedRepo.GetDescription(),
		URL:         forkedRepo.GetHTMLURL(),
		Private:     forkedRepo.GetPrivate(),
		Language:    forkedRepo.GetLanguage(),
		Stars:       forkedRepo.GetStargazersCount(),
		Forks:       forkedRepo.GetForksCount(),
	}, nil
}

// ListLabels lists all labels in a repository
func (c *Client) ListLabels(owner, repo string) ([]*github.Label, error) {
	labels, _, err := c.client.Issues.ListLabels(c.ctx, owner, repo, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list labels: %w", err)
	}
	return labels, nil
}

// CreateLabel creates a new label
func (c *Client) CreateLabel(owner, repo, name, color, description string) error {
	label := &github.Label{
		Name:        &name,
		Color:       &color,
		Description: &description,
	}

	_, _, err := c.client.Issues.CreateLabel(c.ctx, owner, repo, label)
	if err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}
	return nil
}

// SearchRepositories searches for repositories
func (c *Client) SearchRepositories(query string, limit int) ([]*Repository, error) {
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	result, _, err := c.client.Search.Repositories(c.ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search repositories: %w", err)
	}

	var repositories []*Repository
	for _, repo := range result.Repositories {
		repository := &Repository{
			Name:        repo.GetName(),
			Owner:       repo.GetOwner().GetLogin(),
			Description: repo.GetDescription(),
			URL:         repo.GetHTMLURL(),
			Private:     repo.GetPrivate(),
			Language:    repo.GetLanguage(),
			Stars:       repo.GetStargazersCount(),
			Forks:       repo.GetForksCount(),
		}
		repositories = append(repositories, repository)
	}

	return repositories, nil
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
