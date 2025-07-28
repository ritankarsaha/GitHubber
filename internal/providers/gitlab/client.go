/*
 * GitHubber - GitLab Provider Implementation  
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: GitLab API provider implementation
 */

package gitlab

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ritankarsaha/git-tool/internal/providers"
	"github.com/xanzy/go-gitlab"
)

// GitLabProvider implements the Provider interface for GitLab
type GitLabProvider struct {
	client        *gitlab.Client
	baseURL       string
	token         string
	authenticated bool
}

// NewGitLabProvider creates a new GitLab provider
func NewGitLabProvider(config *providers.ProviderConfig) (providers.Provider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	var client *gitlab.Client
	var err error

	if config.Token != "" {
		client, err = gitlab.NewClient(config.Token, gitlab.WithBaseURL(baseURL))
	} else {
		return nil, fmt.Errorf("GitLab token is required")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	provider := &GitLabProvider{
		client:  client,
		baseURL: baseURL,
		token:   config.Token,
	}

	// Test authentication
	if err := provider.Authenticate(context.Background(), config.Token); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return provider, nil
}

// GetType returns the provider type
func (g *GitLabProvider) GetType() providers.ProviderType {
	return providers.ProviderGitLab
}

// GetName returns the provider name
func (g *GitLabProvider) GetName() string {
	return "GitLab"
}

// GetBaseURL returns the base URL
func (g *GitLabProvider) GetBaseURL() string {
	return g.baseURL
}

// Authenticate authenticates with the GitLab API
func (g *GitLabProvider) Authenticate(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Test the authentication by getting current user
	_, _, err := g.client.Users.CurrentUser()
	if err != nil {
		g.authenticated = false
		return fmt.Errorf("authentication test failed: %w", err)
	}

	g.authenticated = true
	g.token = token
	return nil
}

// IsAuthenticated returns whether the provider is authenticated
func (g *GitLabProvider) IsAuthenticated() bool {
	return g.authenticated
}

// GetRepository gets a repository
func (g *GitLabProvider) GetRepository(ctx context.Context, owner, repo string) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glProject, _, err := g.client.Projects.GetProject(projectPath, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return g.convertProject(glProject), nil
}

// ListRepositories lists repositories
func (g *GitLabProvider) ListRepositories(ctx context.Context, options *providers.ListOptions) ([]*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	listOpts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	if options.Sort != "" {
		sort := gitlab.SortOptions(options.Sort)
		listOpts.Sort = &sort
	}

	glProjects, _, err := g.client.Projects.ListProjects(listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	repos := make([]*providers.Repository, len(glProjects))
	for i, glProject := range glProjects {
		repos[i] = g.convertProject(glProject)
	}

	return repos, nil
}

// CreateRepository creates a new repository
func (g *GitLabProvider) CreateRepository(ctx context.Context, req *providers.CreateRepositoryRequest) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	visibility := gitlab.PublicVisibility
	if req.Private {
		visibility = gitlab.PrivateVisibility
	}

	createOpts := &gitlab.CreateProjectOptions{
		Name:                 &req.Name,
		Description:          &req.Description,
		Visibility:           &visibility,
		InitializeWithReadme: &req.AutoInit,
	}

	glProject, _, err := g.client.Projects.CreateProject(createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return g.convertProject(glProject), nil
}

// UpdateRepository updates a repository
func (g *GitLabProvider) UpdateRepository(ctx context.Context, owner, repo string, update *providers.UpdateRepositoryRequest) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	
	editOpts := &gitlab.EditProjectOptions{}
	if update.Name != "" {
		editOpts.Name = &update.Name
	}
	if update.Description != "" {
		editOpts.Description = &update.Description
	}
	if update.Private != nil {
		if *update.Private {
			visibility := gitlab.PrivateVisibility
			editOpts.Visibility = &visibility
		} else {
			visibility := gitlab.PublicVisibility
			editOpts.Visibility = &visibility
		}
	}

	glProject, _, err := g.client.Projects.EditProject(projectPath, editOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}

	return g.convertProject(glProject), nil
}

// DeleteRepository deletes a repository
func (g *GitLabProvider) DeleteRepository(ctx context.Context, owner, repo string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	_, err := g.client.Projects.DeleteProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	return nil
}

// ForkRepository forks a repository
func (g *GitLabProvider) ForkRepository(ctx context.Context, owner, repo string) (*providers.Repository, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glProject, _, err := g.client.Projects.ForkProject(projectPath, &gitlab.ForkProjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fork repository: %w", err)
	}

	return g.convertProject(glProject), nil
}

// GetPullRequest gets a merge request (GitLab's equivalent of PR)
func (g *GitLabProvider) GetPullRequest(ctx context.Context, owner, repo string, number int) (*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glMR, _, err := g.client.MergeRequests.GetMergeRequest(projectPath, number, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request: %w", err)
	}

	return g.convertMergeRequest(glMR), nil
}

// ListPullRequests lists merge requests
func (g *GitLabProvider) ListPullRequests(ctx context.Context, owner, repo string, options *providers.ListOptions) ([]*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	listOpts := &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	if options.State != "" {
		listOpts.State = &options.State
	}

	glMRs, _, err := g.client.MergeRequests.ListProjectMergeRequests(projectPath, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list merge requests: %w", err)
	}

	prs := make([]*providers.PullRequest, len(glMRs))
	for i, glMR := range glMRs {
		prs[i] = g.convertMergeRequest(glMR)
	}

	return prs, nil
}

// CreatePullRequest creates a merge request
func (g *GitLabProvider) CreatePullRequest(ctx context.Context, owner, repo string, req *providers.CreatePullRequestRequest) (*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.CreateMergeRequestOptions{
		Title:        &req.Title,
		Description:  &req.Description,
		SourceBranch: &req.Head,
		TargetBranch: &req.Base,
	}

	glMR, _, err := g.client.MergeRequests.CreateMergeRequest(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create merge request: %w", err)
	}

	return g.convertMergeRequest(glMR), nil
}

// UpdatePullRequest updates a merge request
func (g *GitLabProvider) UpdatePullRequest(ctx context.Context, owner, repo string, number int, update *providers.UpdatePullRequestRequest) (*providers.PullRequest, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	updateOpts := &gitlab.UpdateMergeRequestOptions{}

	if update.Title != "" {
		updateOpts.Title = &update.Title
	}
	if update.Description != "" {
		updateOpts.Description = &update.Description
	}
	if update.State != "" {
		updateOpts.StateEvent = &update.State
	}

	glMR, _, err := g.client.MergeRequests.UpdateMergeRequest(projectPath, number, updateOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to update merge request: %w", err)
	}

	return g.convertMergeRequest(glMR), nil
}

// MergePullRequest merges a merge request
func (g *GitLabProvider) MergePullRequest(ctx context.Context, owner, repo string, number int, options *providers.MergeOptions) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	mergeOpts := &gitlab.AcceptMergeRequestOptions{}

	if options.CommitMessage != "" {
		mergeOpts.MergeCommitMessage = &options.CommitMessage
	}
	if options.DeleteHeadBranch {
		mergeOpts.ShouldRemoveSourceBranch = &options.DeleteHeadBranch
	}

	_, _, err := g.client.MergeRequests.AcceptMergeRequest(projectPath, number, mergeOpts)
	if err != nil {
		return fmt.Errorf("failed to merge request: %w", err)
	}

	return nil
}

// ClosePullRequest closes a merge request
func (g *GitLabProvider) ClosePullRequest(ctx context.Context, owner, repo string, number int) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	state := "close"
	updateOpts := &gitlab.UpdateMergeRequestOptions{
		StateEvent: &state,
	}

	_, _, err := g.client.MergeRequests.UpdateMergeRequest(projectPath, number, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to close merge request: %w", err)
	}

	return nil
}

// GetIssue gets an issue
func (g *GitLabProvider) GetIssue(ctx context.Context, owner, repo string, number int) (*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glIssue, _, err := g.client.Issues.GetIssue(projectPath, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	return g.convertIssue(glIssue), nil
}

// ListIssues lists issues
func (g *GitLabProvider) ListIssues(ctx context.Context, owner, repo string, options *providers.ListOptions) ([]*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	listOpts := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    options.Page,
			PerPage: options.PerPage,
		},
	}

	if options.State != "" {
		listOpts.State = &options.State
	}

	glIssues, _, err := g.client.Issues.ListProjectIssues(projectPath, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}

	issues := make([]*providers.Issue, len(glIssues))
	for i, glIssue := range glIssues {
		issues[i] = g.convertIssue(glIssue)
	}

	return issues, nil
}

// CreateIssue creates an issue
func (g *GitLabProvider) CreateIssue(ctx context.Context, owner, repo string, req *providers.CreateIssueRequest) (*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.CreateIssueOptions{
		Title:       &req.Title,
		Description: &req.Description,
		Labels:      &gitlab.LabelOptions{},
	}

	if len(req.Labels) > 0 {
		labels := gitlab.Labels(req.Labels)
		createOpts.Labels = &labels
	}

	glIssue, _, err := g.client.Issues.CreateIssue(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return g.convertIssue(glIssue), nil
}

// UpdateIssue updates an issue
func (g *GitLabProvider) UpdateIssue(ctx context.Context, owner, repo string, number int, update *providers.UpdateIssueRequest) (*providers.Issue, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	updateOpts := &gitlab.UpdateIssueOptions{}

	if update.Title != "" {
		updateOpts.Title = &update.Title
	}
	if update.Description != "" {
		updateOpts.Description = &update.Description
	}
	if update.State != "" {
		updateOpts.StateEvent = &update.State
	}
	if len(update.Labels) > 0 {
		labels := gitlab.Labels(update.Labels)
		updateOpts.Labels = &labels
	}

	glIssue, _, err := g.client.Issues.UpdateIssue(projectPath, number, updateOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue: %w", err)
	}

	return g.convertIssue(glIssue), nil
}

// CloseIssue closes an issue
func (g *GitLabProvider) CloseIssue(ctx context.Context, owner, repo string, number int) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	state := "close"
	updateOpts := &gitlab.UpdateIssueOptions{
		StateEvent: &state,
	}

	_, _, err := g.client.Issues.UpdateIssue(projectPath, number, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	return nil
}

// ListBranches lists branches
func (g *GitLabProvider) ListBranches(ctx context.Context, owner, repo string) ([]*providers.Branch, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glBranches, _, err := g.client.Branches.ListBranches(projectPath, &gitlab.ListBranchesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	branches := make([]*providers.Branch, len(glBranches))
	for i, glBranch := range glBranches {
		branches[i] = g.convertBranch(glBranch)
	}

	return branches, nil
}

// CreateBranch creates a branch
func (g *GitLabProvider) CreateBranch(ctx context.Context, owner, repo string, req *providers.CreateBranchRequest) (*providers.Branch, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.CreateBranchOptions{
		Branch: &req.Name,
		Ref:    &req.SHA,
	}

	glBranch, _, err := g.client.Branches.CreateBranch(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	return g.convertBranch(glBranch), nil
}

// DeleteBranch deletes a branch
func (g *GitLabProvider) DeleteBranch(ctx context.Context, owner, repo string, branch string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	_, err := g.client.Branches.DeleteBranch(projectPath, branch)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

// ListTags lists tags
func (g *GitLabProvider) ListTags(ctx context.Context, owner, repo string) ([]*providers.Tag, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glTags, _, err := g.client.Tags.ListTags(projectPath, &gitlab.ListTagsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	tags := make([]*providers.Tag, len(glTags))
	for i, glTag := range glTags {
		tags[i] = g.convertTag(glTag)
	}

	return tags, nil
}

// CreateTag creates a tag
func (g *GitLabProvider) CreateTag(ctx context.Context, owner, repo string, req *providers.CreateTagRequest) (*providers.Tag, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.CreateTagOptions{
		TagName: &req.Name,
		Ref:     &req.SHA,
	}

	if req.Message != "" {
		createOpts.Message = &req.Message
	}

	glTag, _, err := g.client.Tags.CreateTag(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return g.convertTag(glTag), nil
}

// DeleteTag deletes a tag
func (g *GitLabProvider) DeleteTag(ctx context.Context, owner, repo string, tag string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	_, err := g.client.Tags.DeleteTag(projectPath, tag)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

// ListReleases lists releases
func (g *GitLabProvider) ListReleases(ctx context.Context, owner, repo string) ([]*providers.Release, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glReleases, _, err := g.client.Releases.ListReleases(projectPath, &gitlab.ListReleasesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	releases := make([]*providers.Release, len(glReleases))
	for i, glRelease := range glReleases {
		releases[i] = g.convertRelease(glRelease)
	}

	return releases, nil
}

// CreateRelease creates a release
func (g *GitLabProvider) CreateRelease(ctx context.Context, owner, repo string, req *providers.CreateReleaseRequest) (*providers.Release, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.CreateReleaseOptions{
		TagName:     &req.TagName,
		Name:        &req.Name,
		Description: &req.Description,
	}

	glRelease, _, err := g.client.Releases.CreateRelease(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	return g.convertRelease(glRelease), nil
}

// UpdateRelease updates a release
func (g *GitLabProvider) UpdateRelease(ctx context.Context, owner, repo string, id string, update *providers.UpdateReleaseRequest) (*providers.Release, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	updateOpts := &gitlab.UpdateReleaseOptions{}

	if update.Name != "" {
		updateOpts.Name = &update.Name
	}
	if update.Description != "" {
		updateOpts.Description = &update.Description
	}

	glRelease, _, err := g.client.Releases.UpdateRelease(projectPath, id, updateOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to update release: %w", err)
	}

	return g.convertRelease(glRelease), nil
}

// DeleteRelease deletes a release
func (g *GitLabProvider) DeleteRelease(ctx context.Context, owner, repo string, id string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	_, err := g.client.Releases.DeleteRelease(projectPath, id)
	if err != nil {
		return fmt.Errorf("failed to delete release: %w", err)
	}

	return nil
}

// ListWebhooks lists project hooks (webhooks)
func (g *GitLabProvider) ListWebhooks(ctx context.Context, owner, repo string) ([]*providers.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glHooks, _, err := g.client.Projects.ListProjectHooks(projectPath, &gitlab.ListProjectHooksOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	webhooks := make([]*providers.Webhook, len(glHooks))
	for i, glHook := range glHooks {
		webhooks[i] = g.convertWebhook(glHook)
	}

	return webhooks, nil
}

// CreateWebhook creates a project hook
func (g *GitLabProvider) CreateWebhook(ctx context.Context, owner, repo string, req *providers.CreateWebhookRequest) (*providers.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.AddProjectHookOptions{
		URL:                   &req.URL,
		EnableSSLVerification: gitlab.Bool(true),
	}

	if req.Secret != "" {
		createOpts.Token = &req.Secret
	}

	// Map events to GitLab hook options
	for _, event := range req.Events {
		switch event {
		case "push":
			createOpts.PushEvents = &req.Active
		case "issues":
			createOpts.IssuesEvents = &req.Active
		case "merge_request":
			createOpts.MergeRequestsEvents = &req.Active
		case "tag_push":
			createOpts.TagPushEvents = &req.Active
		}
	}

	glHook, _, err := g.client.Projects.AddProjectHook(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return g.convertWebhook(glHook), nil
}

// UpdateWebhook updates a project hook
func (g *GitLabProvider) UpdateWebhook(ctx context.Context, owner, repo string, id string, update *providers.UpdateWebhookRequest) (*providers.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	hookID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook ID: %w", err)
	}

	updateOpts := &gitlab.EditProjectHookOptions{}
	if update.URL != "" {
		updateOpts.URL = &update.URL
	}

	if len(update.Events) > 0 {
		for _, event := range update.Events {
			switch event {
			case "push":
				updateOpts.PushEvents = update.Active
			case "issues":
				updateOpts.IssuesEvents = update.Active
			case "merge_request":
				updateOpts.MergeRequestsEvents = update.Active
			case "tag_push":
				updateOpts.TagPushEvents = update.Active
			}
		}
	}

	glHook, _, err := g.client.Projects.EditProjectHook(projectPath, hookID, updateOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}

	return g.convertWebhook(glHook), nil
}

// DeleteWebhook deletes a project hook
func (g *GitLabProvider) DeleteWebhook(ctx context.Context, owner, repo string, id string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	hookID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %w", err)
	}

	_, err = g.client.Projects.DeleteProjectHook(projectPath, hookID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

// ListPipelines lists pipelines
func (g *GitLabProvider) ListPipelines(ctx context.Context, owner, repo string) ([]*providers.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	glPipelines, _, err := g.client.Pipelines.ListProjectPipelines(projectPath, &gitlab.ListProjectPipelinesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipelines: %w", err)
	}

	pipelines := make([]*providers.Pipeline, len(glPipelines))
	for i, glPipeline := range glPipelines {
		pipelines[i] = g.convertPipeline(glPipeline)
	}

	return pipelines, nil
}

// GetPipeline gets a pipeline
func (g *GitLabProvider) GetPipeline(ctx context.Context, owner, repo string, id string) (*providers.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	pipelineID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	glPipeline, _, err := g.client.Pipelines.GetPipeline(projectPath, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	return g.convertPipeline(glPipeline), nil
}

// TriggerPipeline triggers a pipeline
func (g *GitLabProvider) TriggerPipeline(ctx context.Context, owner, repo string, options *providers.TriggerPipelineOptions) (*providers.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	createOpts := &gitlab.CreatePipelineOptions{
		Ref: &options.Ref,
	}

	if len(options.Variables) > 0 {
		variables := make([]*gitlab.PipelineVariable, 0, len(options.Variables))
		for k, v := range options.Variables {
			variables = append(variables, &gitlab.PipelineVariable{
				Key:   k,
				Value: v,
			})
		}
		createOpts.Variables = &variables
	}

	glPipeline, _, err := g.client.Pipelines.CreatePipeline(projectPath, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger pipeline: %w", err)
	}

	return g.convertPipeline(glPipeline), nil
}

// CancelPipeline cancels a pipeline
func (g *GitLabProvider) CancelPipeline(ctx context.Context, owner, repo string, id string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectPath := fmt.Sprintf("%s/%s", owner, repo)
	pipelineID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid pipeline ID: %w", err)
	}

	_, _, err = g.client.Pipelines.CancelPipelineBuild(projectPath, pipelineID)
	if err != nil {
		return fmt.Errorf("failed to cancel pipeline: %w", err)
	}

	return nil
}

// GetUser gets the authenticated user
func (g *GitLabProvider) GetUser(ctx context.Context) (*providers.User, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	glUser, _, err := g.client.Users.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return g.convertUser(glUser), nil
}

// GetUserByUsername gets a user by username
func (g *GitLabProvider) GetUserByUsername(ctx context.Context, username string) (*providers.User, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	users, _, err := g.client.Users.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return g.convertUser(users[0]), nil
}

// Conversion methods
func (g *GitLabProvider) convertProject(glProject *gitlab.Project) *providers.Repository {
	repo := &providers.Repository{
		ID:          strconv.Itoa(glProject.ID),
		Name:        glProject.Name,
		FullName:    glProject.PathWithNamespace,
		Description: glProject.Description,
		URL:         glProject.WebURL,
		CloneURL:    glProject.HTTPURLToRepo,
		SSHURL:      glProject.SSHURLToRepo,
		Private:     glProject.Visibility != gitlab.PublicVisibility,
		Fork:        glProject.ForkedFromProject != nil,
		Language:    "", // GitLab doesn't provide primary language in project info
		Stars:       glProject.StarCount,
		Forks:       glProject.ForksCount,
		OpenIssues:  glProject.OpenIssuesCount,
		CreatedAt:   *glProject.CreatedAt,
		UpdatedAt:   *glProject.LastActivityAt,
	}

	if glProject.Owner != nil {
		repo.Owner = g.convertUser(glProject.Owner)
	}

	// GitLab permissions are handled differently
	repo.Permissions.Admin = glProject.Permissions != nil && glProject.Permissions.ProjectAccess != nil &&
		glProject.Permissions.ProjectAccess.AccessLevel >= gitlab.MaintainerPermissions
	repo.Permissions.Push = glProject.Permissions != nil && glProject.Permissions.ProjectAccess != nil &&
		glProject.Permissions.ProjectAccess.AccessLevel >= gitlab.DeveloperPermissions
	repo.Permissions.Pull = glProject.Permissions != nil && glProject.Permissions.ProjectAccess != nil &&
		glProject.Permissions.ProjectAccess.AccessLevel >= gitlab.GuestPermissions

	return repo
}

func (g *GitLabProvider) convertMergeRequest(glMR *gitlab.MergeRequest) *providers.PullRequest {
	pr := &providers.PullRequest{
		ID:          strconv.Itoa(glMR.IID),
		Number:      glMR.IID,
		Title:       glMR.Title,
		Description: glMR.Description,
		State:       glMR.State,
		URL:         glMR.WebURL,
		CreatedAt:   *glMR.CreatedAt,
		UpdatedAt:   *glMR.UpdatedAt,
		Labels:      glMR.Labels,
	}

	if glMR.Author != nil {
		pr.Author = g.convertUser(glMR.Author)
	}

	if glMR.Assignee != nil {
		pr.Assignee = g.convertUser(glMR.Assignee)
	}

	pr.Head = &providers.Branch{
		Name: glMR.SourceBranch,
		SHA:  glMR.SHA,
	}

	pr.Base = &providers.Branch{
		Name: glMR.TargetBranch,
	}

	return pr
}

func (g *GitLabProvider) convertIssue(glIssue *gitlab.Issue) *providers.Issue {
	issue := &providers.Issue{
		ID:          strconv.Itoa(glIssue.IID),
		Number:      glIssue.IID,
		Title:       glIssue.Title,
		Description: glIssue.Description,
		State:       glIssue.State,
		URL:         glIssue.WebURL,
		CreatedAt:   *glIssue.CreatedAt,
		UpdatedAt:   *glIssue.UpdatedAt,
		Labels:      glIssue.Labels,
	}

	if glIssue.Author != nil {
		issue.Author = g.convertUser(glIssue.Author)
	}

	if glIssue.Assignee != nil {
		issue.Assignee = g.convertUser(glIssue.Assignee)
	}

	return issue
}

func (g *GitLabProvider) convertUser(glUser *gitlab.User) *providers.User {
	user := &providers.User{
		ID:        strconv.Itoa(glUser.ID),
		Username:  glUser.Username,
		Email:     glUser.Email,
		Name:      glUser.Name,
		AvatarURL: glUser.AvatarURL,
		URL:       glUser.WebURL,
	}

	return user
}

func (g *GitLabProvider) convertBranch(glBranch *gitlab.Branch) *providers.Branch {
	branch := &providers.Branch{
		Name:      glBranch.Name,
		Protected: glBranch.Protected,
	}

	if glBranch.Commit != nil {
		branch.SHA = glBranch.Commit.ID
	}

	return branch
}

func (g *GitLabProvider) convertTag(glTag *gitlab.Tag) *providers.Tag {
	tag := &providers.Tag{
		Name: glTag.Name,
	}

	if glTag.Commit != nil {
		tag.SHA = glTag.Commit.ID
	}

	return tag
}

func (g *GitLabProvider) convertRelease(glRelease *gitlab.Release) *providers.Release {
	release := &providers.Release{
		ID:          glRelease.TagName, // GitLab uses tag name as release ID
		TagName:     glRelease.TagName,
		Name:        glRelease.Name,
		Description: glRelease.Description,
		CreatedAt:   *glRelease.CreatedAt,
		PublishedAt: *glRelease.ReleasedAt,
	}

	if glRelease.Author != nil {
		release.Author = g.convertUser(glRelease.Author)
	}

	return release
}

func (g *GitLabProvider) convertWebhook(glHook *gitlab.ProjectHook) *providers.Webhook {
	webhook := &providers.Webhook{
		ID:     strconv.Itoa(glHook.ID),
		URL:    glHook.URL,
		Events: make([]string, 0),
		Config: make(map[string]string),
	}

	// Map GitLab hook events
	if glHook.PushEvents {
		webhook.Events = append(webhook.Events, "push")
	}
	if glHook.IssuesEvents {
		webhook.Events = append(webhook.Events, "issues")
	}
	if glHook.MergeRequestsEvents {
		webhook.Events = append(webhook.Events, "merge_request")
	}
	if glHook.TagPushEvents {
		webhook.Events = append(webhook.Events, "tag_push")
	}

	webhook.Config["url"] = glHook.URL
	webhook.Config["enable_ssl_verification"] = strconv.FormatBool(glHook.EnableSSLVerification)

	return webhook
}

func (g *GitLabProvider) convertPipeline(glPipeline *gitlab.PipelineInfo) *providers.Pipeline {
	pipeline := &providers.Pipeline{
		ID:        strconv.Itoa(glPipeline.ID),
		Status:    glPipeline.Status,
		Ref:       glPipeline.Ref,
		SHA:       glPipeline.SHA,
		URL:       glPipeline.WebURL,
		CreatedAt: *glPipeline.CreatedAt,
		UpdatedAt: *glPipeline.UpdatedAt,
	}

	return pipeline
}

// Factory function for provider registration
func NewGitLabProviderFactory() providers.ProviderFactory {
	return func(config *providers.ProviderConfig) (providers.Provider, error) {
		return NewGitLabProvider(config)
	}
}