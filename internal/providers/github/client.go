/*
 * GitHubber - GitHub Provider Implementation
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: GitHub API provider implementation
 */

package github

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/v66/github"
	"github.com/ritankarsaha/git-tool/internal/providers"
	"golang.org/x/oauth2"
)

// GitHubProvider implements the Provider interface for GitHub
type GitHubProvider struct {
	client    *github.Client
	ctx       context.Context
	baseURL   string
	token     string
	authenticated bool
}

// NewGitHubProvider creates a new GitHub provider
func NewGitHubProvider(config *providers.ProviderConfig) (providers.Provider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	ctx := context.Background()
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	provider := &GitHubProvider{
		ctx:     ctx,
		baseURL: baseURL,
		token:   config.Token,
	}

	if config.Token != "" {
		if err := provider.Authenticate(ctx, config.Token); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	return provider, nil
}

// GetType returns the provider type
func (g *GitHubProvider) GetType() providers.ProviderType {
	return providers.ProviderGitHub
}

// GetName returns the provider name
func (g *GitHubProvider) GetName() string {
	return "GitHub"
}

// GetBaseURL returns the base URL
func (g *GitHubProvider) GetBaseURL() string {
	return g.baseURL
}

// Authenticate authenticates with the GitHub API
func (g *GitHubProvider) Authenticate(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	g.client = github.NewClient(tc)
	g.token = token

	// Test the authentication
	_, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		g.authenticated = false
		return fmt.Errorf("authentication test failed: %w", err)
	}

	g.authenticated = true
	return nil
}

// IsAuthenticated returns whether the provider is authenticated
func (g *GitHubProvider) IsAuthenticated() bool {
	return g.authenticated
}

// GetRepository gets a repository
func (g *GitHubProvider) GetRepository(ctx context.Context, owner, repo string) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghRepo, _, err := g.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return g.convertRepository(ghRepo), nil
}

// ListRepositories lists repositories
func (g *GitHubProvider) ListRepositories(ctx context.Context, options *providers.ListOptions) ([]*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	listOpts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	if options.Sort != "" {
		listOpts.Sort = options.Sort
	}

	ghRepos, _, err := g.client.Repositories.List(ctx, "", listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	repos := make([]*providers.Repository, len(ghRepos))
	for i, ghRepo := range ghRepos {
		repos[i] = g.convertRepository(ghRepo)
	}

	return repos, nil
}

// CreateRepository creates a new repository
func (g *GitHubProvider) CreateRepository(ctx context.Context, req *providers.CreateRepositoryRequest) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghRepo := &github.Repository{
		Name:        &req.Name,
		Description: &req.Description,
		Private:     &req.Private,
		AutoInit:    &req.AutoInit,
	}

	createdRepo, _, err := g.client.Repositories.Create(ctx, "", ghRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return g.convertRepository(createdRepo), nil
}

// UpdateRepository updates a repository
func (g *GitHubProvider) UpdateRepository(ctx context.Context, owner, repo string, update *providers.UpdateRepositoryRequest) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghRepo := &github.Repository{}
	if update.Name != "" {
		ghRepo.Name = &update.Name
	}
	if update.Description != "" {
		ghRepo.Description = &update.Description
	}
	if update.Private != nil {
		ghRepo.Private = update.Private
	}

	updatedRepo, _, err := g.client.Repositories.Edit(ctx, owner, repo, ghRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}

	return g.convertRepository(updatedRepo), nil
}

// DeleteRepository deletes a repository
func (g *GitHubProvider) DeleteRepository(ctx context.Context, owner, repo string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	_, err := g.client.Repositories.Delete(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	return nil
}

// ForkRepository forks a repository
func (g *GitHubProvider) ForkRepository(ctx context.Context, owner, repo string) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	forkedRepo, _, err := g.client.Repositories.CreateFork(ctx, owner, repo, &github.RepositoryCreateForkOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fork repository: %w", err)
	}

	return g.convertRepository(forkedRepo), nil
}

// GetPullRequest gets a pull request
func (g *GitHubProvider) GetPullRequest(ctx context.Context, owner, repo string, number int) (*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghPR, _, err := g.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return g.convertPullRequest(ghPR), nil
}

// ListPullRequests lists pull requests
func (g *GitHubProvider) ListPullRequests(ctx context.Context, owner, repo string, options *providers.ListOptions) ([]*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	listOpts := &github.PullRequestListOptions{
		State: options.State,
		ListOptions: github.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	ghPRs, _, err := g.client.PullRequests.List(ctx, owner, repo, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	prs := make([]*providers.PullRequest, len(ghPRs))
	for i, ghPR := range ghPRs {
		prs[i] = g.convertPullRequest(ghPR)
	}

	return prs, nil
}

// CreatePullRequest creates a pull request
func (g *GitHubProvider) CreatePullRequest(ctx context.Context, owner, repo string, req *providers.CreatePullRequestRequest) (*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghPR := &github.NewPullRequest{
		Title: &req.Title,
		Body:  &req.Description,
		Head:  &req.Head,
		Base:  &req.Base,
	}

	createdPR, _, err := g.client.PullRequests.Create(ctx, owner, repo, ghPR)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return g.convertPullRequest(createdPR), nil
}

// UpdatePullRequest updates a pull request
func (g *GitHubProvider) UpdatePullRequest(ctx context.Context, owner, repo string, number int, update *providers.UpdatePullRequestRequest) (*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghPR := &github.PullRequest{}
	if update.Title != "" {
		ghPR.Title = &update.Title
	}
	if update.Description != "" {
		ghPR.Body = &update.Description
	}
	if update.State != "" {
		ghPR.State = &update.State
	}

	updatedPR, _, err := g.client.PullRequests.Edit(ctx, owner, repo, number, ghPR)
	if err != nil {
		return nil, fmt.Errorf("failed to update pull request: %w", err)
	}

	return g.convertPullRequest(updatedPR), nil
}

// MergePullRequest merges a pull request
func (g *GitHubProvider) MergePullRequest(ctx context.Context, owner, repo string, number int, options *providers.MergeOptions) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	mergeOpts := &github.PullRequestOptions{
		CommitTitle: options.CommitTitle,
		MergeMethod: options.MergeMethod,
	}

	_, _, err := g.client.PullRequests.Merge(ctx, owner, repo, number, options.CommitMessage, mergeOpts)
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}

	return nil
}

// ClosePullRequest closes a pull request
func (g *GitHubProvider) ClosePullRequest(ctx context.Context, owner, repo string, number int) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	state := "closed"
	ghPR := &github.PullRequest{State: &state}

	_, _, err := g.client.PullRequests.Edit(ctx, owner, repo, number, ghPR)
	if err != nil {
		return fmt.Errorf("failed to close pull request: %w", err)
	}

	return nil
}

// GetIssue gets an issue
func (g *GitHubProvider) GetIssue(ctx context.Context, owner, repo string, number int) (*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghIssue, _, err := g.client.Issues.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return g.convertIssue(ghIssue), nil
}

// ListIssues lists issues
func (g *GitHubProvider) ListIssues(ctx context.Context, owner, repo string, options *providers.ListOptions) ([]*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	listOpts := &github.IssueListByRepoOptions{
		State: options.State,
		ListOptions: github.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	ghIssues, _, err := g.client.Issues.ListByRepo(ctx, owner, repo, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}

	issues := make([]*providers.Issue, 0, len(ghIssues))
	for _, ghIssue := range ghIssues {
		// Skip pull requests (GitHub API treats PRs as issues)
		if ghIssue.IsPullRequest() {
			continue
		}
		issues = append(issues, g.convertIssue(ghIssue))
	}

	return issues, nil
}

// CreateIssue creates an issue
func (g *GitHubProvider) CreateIssue(ctx context.Context, owner, repo string, req *providers.CreateIssueRequest) (*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghIssue := &github.IssueRequest{
		Title:  &req.Title,
		Body:   &req.Description,
		Labels: &req.Labels,
	}

	if req.Assignee != "" {
		ghIssue.Assignee = &req.Assignee
	}

	createdIssue, _, err := g.client.Issues.Create(ctx, owner, repo, ghIssue)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return g.convertIssue(createdIssue), nil
}

// UpdateIssue updates an issue
func (g *GitHubProvider) UpdateIssue(ctx context.Context, owner, repo string, number int, update *providers.UpdateIssueRequest) (*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghIssue := &github.IssueRequest{}
	if update.Title != "" {
		ghIssue.Title = &update.Title
	}
	if update.Description != "" {
		ghIssue.Body = &update.Description
	}
	if update.State != "" {
		ghIssue.State = &update.State
	}
	if len(update.Labels) > 0 {
		ghIssue.Labels = &update.Labels
	}
	if update.Assignee != "" {
		ghIssue.Assignee = &update.Assignee
	}

	updatedIssue, _, err := g.client.Issues.Edit(ctx, owner, repo, number, ghIssue)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	return g.convertIssue(updatedIssue), nil
}

// CloseIssue closes an issue
func (g *GitHubProvider) CloseIssue(ctx context.Context, owner, repo string, number int) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	state := "closed"
	ghIssue := &github.IssueRequest{State: &state}

	_, _, err := g.client.Issues.Edit(ctx, owner, repo, number, ghIssue)
	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	return nil
}

// ListBranches lists branches
func (g *GitHubProvider) ListBranches(ctx context.Context, owner, repo string) ([]*providers.Branch, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghBranches, _, err := g.client.Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	branches := make([]*providers.Branch, len(ghBranches))
	for i, ghBranch := range ghBranches {
		branches[i] = g.convertBranch(ghBranch)
	}

	return branches, nil
}

// CreateBranch creates a branch
func (g *GitHubProvider) CreateBranch(ctx context.Context, owner, repo string, req *providers.CreateBranchRequest) (*providers.Branch, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ref := fmt.Sprintf("refs/heads/%s", req.Name)
	ghRef := &github.Reference{
		Ref: &ref,
		Object: &github.GitObject{
			SHA: &req.SHA,
		},
	}

	createdRef, _, err := g.client.Git.CreateRef(ctx, owner, repo, ghRef)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return &providers.Branch{
		Name: req.Name,
		SHA:  createdRef.GetObject().GetSHA(),
	}, nil
}

// DeleteBranch deletes a branch
func (g *GitHubProvider) DeleteBranch(ctx context.Context, owner, repo string, branch string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	ref := fmt.Sprintf("heads/%s", branch)
	_, err := g.client.Git.DeleteRef(ctx, owner, repo, ref)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

// ListTags lists tags
func (g *GitHubProvider) ListTags(ctx context.Context, owner, repo string) ([]*providers.Tag, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghTags, _, err := g.client.Repositories.ListTags(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tags := make([]*providers.Tag, len(ghTags))
	for i, ghTag := range ghTags {
		tags[i] = g.convertTag(ghTag)
	}

	return tags, nil
}

// CreateTag creates a tag
func (g *GitHubProvider) CreateTag(ctx context.Context, owner, repo string, req *providers.CreateTagRequest) (*providers.Tag, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	// Create tag object if message is provided
	if req.Message != "" {
		tagObject := &github.Tag{
			Tag:     &req.Name,
			Message: &req.Message,
			Object: &github.GitObject{
				SHA: &req.SHA,
			},
		}

		createdTag, _, err := g.client.Git.CreateTag(ctx, owner, repo, tagObject)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag object: %w", err)
		}

		// Create reference
		ref := fmt.Sprintf("refs/tags/%s", req.Name)
		ghRef := &github.Reference{
			Ref: &ref,
			Object: &github.GitObject{
				SHA: createdTag.SHA,
			},
		}

		_, _, err = g.client.Git.CreateRef(ctx, owner, repo, ghRef)
		if err != nil {
			return nil, fmt.Errorf("failed to create tag reference: %w", err)
		}

		return &providers.Tag{
			Name: req.Name,
			SHA:  *createdTag.SHA,
		}, nil
	}

	// Create lightweight tag
	ref := fmt.Sprintf("refs/tags/%s", req.Name)
	ghRef := &github.Reference{
		Ref: &ref,
		Object: &github.GitObject{
			SHA: &req.SHA,
		},
	}

	createdRef, _, err := g.client.Git.CreateRef(ctx, owner, repo, ghRef)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &providers.Tag{
		Name: req.Name,
		SHA:  createdRef.GetObject().GetSHA(),
	}, nil
}

// DeleteTag deletes a tag
func (g *GitHubProvider) DeleteTag(ctx context.Context, owner, repo string, tag string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	ref := fmt.Sprintf("tags/%s", tag)
	_, err := g.client.Git.DeleteRef(ctx, owner, repo, ref)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// ListReleases lists releases
func (g *GitHubProvider) ListReleases(ctx context.Context, owner, repo string) ([]*providers.Release, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghReleases, _, err := g.client.Repositories.ListReleases(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	releases := make([]*providers.Release, len(ghReleases))
	for i, ghRelease := range ghReleases {
		releases[i] = g.convertRelease(ghRelease)
	}

	return releases, nil
}

// CreateRelease creates a release
func (g *GitHubProvider) CreateRelease(ctx context.Context, owner, repo string, req *providers.CreateReleaseRequest) (*providers.Release, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghRelease := &github.RepositoryRelease{
		TagName:    &req.TagName,
		Name:       &req.Name,
		Body:       &req.Description,
		Draft:      &req.Draft,
		Prerelease: &req.Prerelease,
	}

	createdRelease, _, err := g.client.Repositories.CreateRelease(ctx, owner, repo, ghRelease)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	return g.convertRelease(createdRelease), nil
}

// UpdateRelease updates a release
func (g *GitHubProvider) UpdateRelease(ctx context.Context, owner, repo string, id string, update *providers.UpdateReleaseRequest) (*providers.Release, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	releaseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid release ID: %w", err)
	}

	ghRelease := &github.RepositoryRelease{}
	if update.Name != "" {
		ghRelease.Name = &update.Name
	}
	if update.Description != "" {
		ghRelease.Body = &update.Description
	}
	if update.Draft != nil {
		ghRelease.Draft = update.Draft
	}
	if update.Prerelease != nil {
		ghRelease.Prerelease = update.Prerelease
	}

	updatedRelease, _, err := g.client.Repositories.EditRelease(ctx, owner, repo, releaseID, ghRelease)
	if err != nil {
		return nil, fmt.Errorf("failed to update release: %w", err)
	}

	return g.convertRelease(updatedRelease), nil
}

// DeleteRelease deletes a release
func (g *GitHubProvider) DeleteRelease(ctx context.Context, owner, repo string, id string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	releaseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid release ID: %w", err)
	}

	_, err = g.client.Repositories.DeleteRelease(ctx, owner, repo, releaseID)
	if err != nil {
		return fmt.Errorf("failed to delete release: %w", err)
	}

	return nil
}

// ListWebhooks lists webhooks
func (g *GitHubProvider) ListWebhooks(ctx context.Context, owner, repo string) ([]*providers.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghHooks, _, err := g.client.Repositories.ListHooks(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	webhooks := make([]*providers.Webhook, len(ghHooks))
	for i, ghHook := range ghHooks {
		webhooks[i] = g.convertWebhook(ghHook)
	}

	return webhooks, nil
}

// CreateWebhook creates a webhook
func (g *GitHubProvider) CreateWebhook(ctx context.Context, owner, repo string, req *providers.CreateWebhookRequest) (*providers.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	config := map[string]interface{}{
		"url":          req.URL,
		"content_type": "json",
	}

	if req.Secret != "" {
		config["secret"] = req.Secret
	}

	ghHook := &github.Hook{
		Name:   github.String("web"),
		Config: config,
		Events: req.Events,
		Active: &req.Active,
	}

	createdHook, _, err := g.client.Repositories.CreateHook(ctx, owner, repo, ghHook)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return g.convertWebhook(createdHook), nil
}

// UpdateWebhook updates a webhook
func (g *GitHubProvider) UpdateWebhook(ctx context.Context, owner, repo string, id string, update *providers.UpdateWebhookRequest) (*providers.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	hookID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook ID: %w", err)
	}

	ghHook := &github.Hook{}
	if update.URL != "" {
		config := map[string]interface{}{
			"url":          update.URL,
			"content_type": "json",
		}
		ghHook.Config = config
	}
	if len(update.Events) > 0 {
		ghHook.Events = update.Events
	}
	if update.Active != nil {
		ghHook.Active = update.Active
	}

	updatedHook, _, err := g.client.Repositories.EditHook(ctx, owner, repo, hookID, ghHook)
	if err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}

	return g.convertWebhook(updatedHook), nil
}

// DeleteWebhook deletes a webhook
func (g *GitHubProvider) DeleteWebhook(ctx context.Context, owner, repo string, id string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	hookID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %w", err)
	}

	_, err = g.client.Repositories.DeleteHook(ctx, owner, repo, hookID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

// ListPipelines lists GitHub Actions workflows (pipelines)
func (g *GitHubProvider) ListPipelines(ctx context.Context, owner, repo string) ([]*providers.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	workflows, _, err := g.client.Actions.ListWorkflows(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	pipelines := make([]*providers.Pipeline, len(workflows.Workflows))
	for i, workflow := range workflows.Workflows {
		pipelines[i] = &providers.Pipeline{
			ID:     strconv.FormatInt(workflow.GetID(), 10),
			Status: workflow.GetState(),
			URL:    workflow.GetHTMLURL(),
		}
	}

	return pipelines, nil
}

// GetPipeline gets a workflow run (pipeline)
func (g *GitHubProvider) GetPipeline(ctx context.Context, owner, repo string, id string) (*providers.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	runID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid run ID: %w", err)
	}

	run, _, err := g.client.Actions.GetWorkflowRunByID(ctx, owner, repo, runID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow run: %w", err)
	}

	return &providers.Pipeline{
		ID:        strconv.FormatInt(run.GetID(), 10),
		Status:    run.GetStatus(),
		Ref:       run.GetHeadBranch(),
		SHA:       run.GetHeadSHA(),
		URL:       run.GetHTMLURL(),
		CreatedAt: run.GetCreatedAt().Time,
		UpdatedAt: run.GetUpdatedAt().Time,
	}, nil
}

// TriggerPipeline triggers a workflow
func (g *GitHubProvider) TriggerPipeline(ctx context.Context, owner, repo string, options *providers.TriggerPipelineOptions) (*providers.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	event := github.CreateWorkflowDispatchEventRequest{
		Ref:    options.Ref,
		Inputs: options.Variables,
	}

	_, err := g.client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, "main.yml", event)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger workflow: %w", err)
	}

	// Return a placeholder pipeline as GitHub doesn't immediately return the run
	return &providers.Pipeline{
		Status: "queued",
		Ref:    options.Ref,
	}, nil
}

// CancelPipeline cancels a workflow run
func (g *GitHubProvider) CancelPipeline(ctx context.Context, owner, repo string, id string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	runID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid run ID: %w", err)
	}

	_, err = g.client.Actions.CancelWorkflowRunByID(ctx, owner, repo, runID)
	if err != nil {
		return fmt.Errorf("failed to cancel workflow run: %w", err)
	}

	return nil
}

// GetUser gets the authenticated user
func (g *GitHubProvider) GetUser(ctx context.Context) (*providers.User, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghUser, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return g.convertUser(ghUser), nil
}

// GetUserByUsername gets a user by username
func (g *GitHubProvider) GetUserByUsername(ctx context.Context, username string) (*providers.User, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	ghUser, _, err := g.client.Users.Get(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return g.convertUser(ghUser), nil
}

// Conversion methods
func (g *GitHubProvider) convertRepository(ghRepo *github.Repository) *providers.Repository {
	repo := &providers.Repository{
		ID:          strconv.FormatInt(ghRepo.GetID(), 10),
		Name:        ghRepo.GetName(),
		FullName:    ghRepo.GetFullName(),
		Description: ghRepo.GetDescription(),
		URL:         ghRepo.GetHTMLURL(),
		CloneURL:    ghRepo.GetCloneURL(),
		SSHURL:      ghRepo.GetSSHURL(),
		Private:     ghRepo.GetPrivate(),
		Fork:        ghRepo.GetFork(),
		Language:    ghRepo.GetLanguage(),
		Stars:       ghRepo.GetStargazersCount(),
		Forks:       ghRepo.GetForksCount(),
		OpenIssues:  ghRepo.GetOpenIssuesCount(),
		CreatedAt:   ghRepo.GetCreatedAt().Time,
		UpdatedAt:   ghRepo.GetUpdatedAt().Time,
	}

	if ghRepo.Owner != nil {
		repo.Owner = g.convertUser(ghRepo.Owner)
	}

	if ghRepo.Permissions != nil {
		repo.Permissions.Admin = ghRepo.Permissions["admin"]
		repo.Permissions.Push = ghRepo.Permissions["push"]
		repo.Permissions.Pull = ghRepo.Permissions["pull"]
	}

	return repo
}

func (g *GitHubProvider) convertPullRequest(ghPR *github.PullRequest) *providers.PullRequest {
	pr := &providers.PullRequest{
		ID:          strconv.Itoa(ghPR.GetNumber()),
		Number:      ghPR.GetNumber(),
		Title:       ghPR.GetTitle(),
		Description: ghPR.GetBody(),
		State:       ghPR.GetState(),
		URL:         ghPR.GetHTMLURL(),
		CreatedAt:   ghPR.GetCreatedAt().Time,
		UpdatedAt:   ghPR.GetUpdatedAt().Time,
		Mergeable:   ghPR.GetMergeable(),
	}

	if ghPR.User != nil {
		pr.Author = g.convertUser(ghPR.User)
	}

	if ghPR.Assignee != nil {
		pr.Assignee = g.convertUser(ghPR.Assignee)
	}

	if ghPR.Head != nil {
		pr.Head = &providers.Branch{
			Name: ghPR.Head.GetRef(),
			SHA:  ghPR.Head.GetSHA(),
		}
	}

	if ghPR.Base != nil {
		pr.Base = &providers.Branch{
			Name: ghPR.Base.GetRef(),
			SHA:  ghPR.Base.GetSHA(),
		}
	}

	for _, label := range ghPR.Labels {
		pr.Labels = append(pr.Labels, label.GetName())
	}

	return pr
}

func (g *GitHubProvider) convertIssue(ghIssue *github.Issue) *providers.Issue {
	issue := &providers.Issue{
		ID:          strconv.Itoa(ghIssue.GetNumber()),
		Number:      ghIssue.GetNumber(),
		Title:       ghIssue.GetTitle(),
		Description: ghIssue.GetBody(),
		State:       ghIssue.GetState(),
		URL:         ghIssue.GetHTMLURL(),
		CreatedAt:   ghIssue.GetCreatedAt().Time,
		UpdatedAt:   ghIssue.GetUpdatedAt().Time,
	}

	if ghIssue.User != nil {
		issue.Author = g.convertUser(ghIssue.User)
	}

	if ghIssue.Assignee != nil {
		issue.Assignee = g.convertUser(ghIssue.Assignee)
	}

	for _, label := range ghIssue.Labels {
		issue.Labels = append(issue.Labels, label.GetName())
	}

	return issue
}

func (g *GitHubProvider) convertUser(ghUser *github.User) *providers.User {
	user := &providers.User{
		ID:        strconv.FormatInt(ghUser.GetID(), 10),
		Username:  ghUser.GetLogin(),
		Email:     ghUser.GetEmail(),
		Name:      ghUser.GetName(),
		AvatarURL: ghUser.GetAvatarURL(),
		URL:       ghUser.GetHTMLURL(),
	}

	return user
}

func (g *GitHubProvider) convertBranch(ghBranch *github.Branch) *providers.Branch {
	branch := &providers.Branch{
		Name:      ghBranch.GetName(),
		Protected: ghBranch.GetProtected(),
	}

	if ghBranch.Commit != nil {
		branch.SHA = ghBranch.Commit.GetSHA()
	}

	return branch
}

func (g *GitHubProvider) convertTag(ghTag *github.RepositoryTag) *providers.Tag {
	tag := &providers.Tag{
		Name: ghTag.GetName(),
	}

	if ghTag.Commit != nil {
		tag.SHA = ghTag.Commit.GetSHA()
	}

	return tag
}

func (g *GitHubProvider) convertRelease(ghRelease *github.RepositoryRelease) *providers.Release {
	release := &providers.Release{
		ID:          strconv.FormatInt(ghRelease.GetID(), 10),
		TagName:     ghRelease.GetTagName(),
		Name:        ghRelease.GetName(),
		Description: ghRelease.GetBody(),
		URL:         ghRelease.GetHTMLURL(),
		Draft:       ghRelease.GetDraft(),
		Prerelease:  ghRelease.GetPrerelease(),
		CreatedAt:   ghRelease.GetCreatedAt().Time,
		PublishedAt: ghRelease.GetPublishedAt().Time,
	}

	if ghRelease.Author != nil {
		release.Author = g.convertUser(ghRelease.Author)
	}

	return release
}

func (g *GitHubProvider) convertWebhook(ghHook *github.Hook) *providers.Webhook {
	webhook := &providers.Webhook{
		ID:     strconv.FormatInt(ghHook.GetID(), 10),
		Events: ghHook.Events,
		Active: ghHook.GetActive(),
		Config: make(map[string]string),
	}

	if ghHook.Config != nil {
		if url, ok := ghHook.Config["url"].(string); ok {
			webhook.URL = url
		}
		for k, v := range ghHook.Config {
			if str, ok := v.(string); ok {
				webhook.Config[k] = str
			}
		}
	}

	return webhook
}

// Factory function for provider registration
func NewGitHubProviderFactory() providers.ProviderFactory {
	return func(config *providers.ProviderConfig) (providers.Provider, error) {
		return NewGitHubProvider(config)
	}
}