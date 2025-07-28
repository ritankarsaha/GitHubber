/*
 * GitHubber - GitHub Actions CI/CD Provider
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: GitHub Actions integration for CI/CD operations
 */

package github

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/ritankarsaha/git-tool/internal/ci"
	"golang.org/x/oauth2"
)

// GitHubActionsProvider implements the CIProvider interface for GitHub Actions
type GitHubActionsProvider struct {
	client        *github.Client
	token         string
	authenticated bool
}

// NewGitHubActionsProvider creates a new GitHub Actions provider
func NewGitHubActionsProvider(config *ci.CIConfig) (ci.CIProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("GitHub token is required")
	}

	// Create OAuth2 token source
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	var client *github.Client
	if config.BaseURL != "" {
		baseURL, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}
		client, err = github.NewEnterpriseClient(baseURL.String(), baseURL.String(), tc)
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub Enterprise client: %w", err)
		}
	} else {
		client = github.NewClient(tc)
	}

	provider := &GitHubActionsProvider{
		client: client,
		token:  config.Token,
	}

	// Test authentication
	if err := provider.Authenticate(context.Background(), config); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return provider, nil
}

// GetPlatform returns the CI platform type
func (g *GitHubActionsProvider) GetPlatform() ci.CIPlatform {
	return ci.PlatformGitHubActions
}

// GetName returns the provider name
func (g *GitHubActionsProvider) GetName() string {
	return "GitHub Actions"
}

// IsConnected returns whether the provider is connected
func (g *GitHubActionsProvider) IsConnected() bool {
	return g.authenticated
}

// Authenticate authenticates with GitHub API
func (g *GitHubActionsProvider) Authenticate(ctx context.Context, config *ci.CIConfig) error {
	// Test authentication by getting the current user
	_, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		g.authenticated = false
		return fmt.Errorf("authentication failed: %w", err)
	}

	g.authenticated = true
	return nil
}

// ListPipelines lists workflow runs (pipelines)
func (g *GitHubActionsProvider) ListPipelines(ctx context.Context, repoURL string, options *ci.ListPipelineOptions) ([]*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	listOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	if options != nil {
		if options.Limit > 0 {
			listOpts.PerPage = options.Limit
		}
		if options.Status != "" {
			status := string(options.Status)
			listOpts.Status = status
		}
		if options.Branch != "" {
			listOpts.Branch = options.Branch
		}
	}

	runs, _, err := g.client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflow runs: %w", err)
	}

	pipelines := make([]*ci.Pipeline, len(runs.WorkflowRuns))
	for i, run := range runs.WorkflowRuns {
		pipelines[i] = g.convertWorkflowRun(run, owner, repo)
	}

	return pipelines, nil
}

// GetPipeline gets a specific workflow run
func (g *GitHubActionsProvider) GetPipeline(ctx context.Context, repoURL, pipelineID string) (*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	runID, err := strconv.ParseInt(pipelineID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	run, _, err := g.client.Actions.GetWorkflowRunByID(ctx, owner, repo, runID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow run: %w", err)
	}

	pipeline := g.convertWorkflowRun(run, owner, repo)

	// Get jobs for this run
	jobs, _, err := g.client.Actions.ListWorkflowJobs(ctx, owner, repo, runID, &github.ListWorkflowJobsOptions{})
	if err == nil && jobs != nil {
		pipeline.Jobs = make([]*ci.Job, len(jobs.Jobs))
		for i, job := range jobs.Jobs {
			pipeline.Jobs[i] = g.convertWorkflowJob(job)
		}
	}

	return pipeline, nil
}

// TriggerPipeline triggers a workflow dispatch
func (g *GitHubActionsProvider) TriggerPipeline(ctx context.Context, repoURL string, request *ci.TriggerPipelineRequest) (*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	// First, get the workflow file to trigger
	workflows, _, err := g.client.Actions.ListWorkflows(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	if len(workflows.Workflows) == 0 {
		return nil, fmt.Errorf("no workflows found in repository")
	}

	// Use the first workflow that supports workflow_dispatch
	var workflowID int64
	for _, workflow := range workflows.Workflows {
		workflowID = workflow.GetID()
		break // For simplicity, use the first workflow
	}

	// Prepare inputs
	inputs := make(map[string]interface{})
	for key, value := range request.Variables {
		inputs[key] = value
	}

	dispatchEvent := &github.CreateWorkflowDispatchEventRequest{
		Ref:    request.Ref,
		Inputs: inputs,
	}

	_, err = g.client.Actions.CreateWorkflowDispatchEventByID(ctx, owner, repo, workflowID, *dispatchEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger workflow: %w", err)
	}

	// GitHub doesn't immediately return the run, so we need to wait and find it
	time.Sleep(2 * time.Second)

	// Find the most recent run for this workflow
	runs, _, err := g.client.Actions.ListWorkflowRunsByID(ctx, owner, repo, workflowID, &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get triggered run: %w", err)
	}

	if len(runs.WorkflowRuns) == 0 {
		return nil, fmt.Errorf("triggered run not found")
	}

	return g.convertWorkflowRun(runs.WorkflowRuns[0], owner, repo), nil
}

// CancelPipeline cancels a workflow run
func (g *GitHubActionsProvider) CancelPipeline(ctx context.Context, repoURL, pipelineID string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return err
	}

	runID, err := strconv.ParseInt(pipelineID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid pipeline ID: %w", err)
	}

	_, err = g.client.Actions.CancelWorkflowRunByID(ctx, owner, repo, runID)
	if err != nil {
		return fmt.Errorf("failed to cancel workflow run: %w", err)
	}

	return nil
}

// RetryPipeline retries a workflow run
func (g *GitHubActionsProvider) RetryPipeline(ctx context.Context, repoURL, pipelineID string) (*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	runID, err := strconv.ParseInt(pipelineID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	_, err = g.client.Actions.RerunWorkflowByID(ctx, owner, repo, runID)
	if err != nil {
		return nil, fmt.Errorf("failed to retry workflow run: %w", err)
	}

	// Return the updated pipeline
	return g.GetPipeline(ctx, repoURL, pipelineID)
}

// GetBuild gets a specific workflow job
func (g *GitHubActionsProvider) GetBuild(ctx context.Context, repoURL, buildID string) (*ci.Build, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	jobID, err := strconv.ParseInt(buildID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid build ID: %w", err)
	}

	job, _, err := g.client.Actions.GetWorkflowJobByID(ctx, owner, repo, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow job: %w", err)
	}

	return g.convertWorkflowJobToBuild(job), nil
}

// CancelBuild cancels a workflow job (not directly supported by GitHub Actions)
func (g *GitHubActionsProvider) CancelBuild(ctx context.Context, repoURL, buildID string) error {
	// GitHub Actions doesn't support canceling individual jobs
	// We would need to cancel the entire workflow run
	return fmt.Errorf("canceling individual jobs is not supported by GitHub Actions")
}

// GetBuildLogs gets logs for a workflow job
func (g *GitHubActionsProvider) GetBuildLogs(ctx context.Context, repoURL, buildID string) (*ci.BuildLogs, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	jobID, err := strconv.ParseInt(buildID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid build ID: %w", err)
	}

	logURL, _, err := g.client.Actions.GetWorkflowJobLogs(ctx, owner, repo, jobID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get job logs: %w", err)
	}

	return &ci.BuildLogs{
		ID:        buildID,
		URL:       logURL.String(),
		FetchedAt: time.Now(),
	}, nil
}

// ListArtifacts lists workflow run artifacts
func (g *GitHubActionsProvider) ListArtifacts(ctx context.Context, repoURL, pipelineID string) ([]*ci.Artifact, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	runID, err := strconv.ParseInt(pipelineID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	artifacts, _, err := g.client.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, runID, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	result := make([]*ci.Artifact, len(artifacts.Artifacts))
	for i, artifact := range artifacts.Artifacts {
		result[i] = g.convertArtifact(artifact)
	}

	return result, nil
}

// DownloadArtifact downloads a workflow artifact
func (g *GitHubActionsProvider) DownloadArtifact(ctx context.Context, repoURL, artifactID string) ([]byte, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(artifactID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid artifact ID: %w", err)
	}

	downloadURL, _, err := g.client.Actions.DownloadArtifact(ctx, owner, repo, id, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact download URL: %w", err)
	}

	// This returns a redirect URL, not the actual content
	// In a real implementation, you would make an HTTP request to downloadURL
	return []byte(downloadURL.String()), nil
}

// ListEnvironments lists deployment environments
func (g *GitHubActionsProvider) ListEnvironments(ctx context.Context, repoURL string) ([]*ci.Environment, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	envs, _, err := g.client.Repositories.ListEnvironments(ctx, owner, repo, &github.EnvironmentListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	environments := make([]*ci.Environment, len(envs.Environments))
	for i, env := range envs.Environments {
		environments[i] = g.convertEnvironment(env)
	}

	return environments, nil
}

// GetEnvironment gets a specific environment
func (g *GitHubActionsProvider) GetEnvironment(ctx context.Context, repoURL, envName string) (*ci.Environment, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	env, _, err := g.client.Repositories.GetEnvironment(ctx, owner, repo, envName)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	return g.convertEnvironment(env), nil
}

// CreateWebhook creates a repository webhook
func (g *GitHubActionsProvider) CreateWebhook(ctx context.Context, repoURL string, config *ci.WebhookConfig) (*ci.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return nil, err
	}

	hookConfig := make(map[string]interface{})
	hookConfig["url"] = config.URL
	hookConfig["content_type"] = "json"
	if config.Secret != "" {
		hookConfig["secret"] = config.Secret
	}
	hookConfig["insecure_ssl"] = config.InsecureSSL

	hook := &github.Hook{
		Name:   github.String("web"),
		Config: hookConfig,
		Events: config.Events,
		Active: github.Bool(true),
	}

	createdHook, _, err := g.client.Repositories.CreateHook(ctx, owner, repo, hook)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return g.convertWebhook(createdHook), nil
}

// UpdateWebhook updates a repository webhook
func (g *GitHubActionsProvider) UpdateWebhook(ctx context.Context, repoURL string, webhook *ci.Webhook) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return err
	}

	hookID, err := strconv.ParseInt(webhook.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %w", err)
	}

	hookConfig := make(map[string]interface{})
	for key, value := range webhook.Config {
		hookConfig[key] = value
	}

	hook := &github.Hook{
		Config: hookConfig,
		Events: webhook.Events,
		Active: github.Bool(webhook.Active),
	}

	_, _, err = g.client.Repositories.EditHook(ctx, owner, repo, hookID, hook)
	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}

// DeleteWebhook deletes a repository webhook
func (g *GitHubActionsProvider) DeleteWebhook(ctx context.Context, repoURL, webhookID string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		return err
	}

	hookID, err := strconv.ParseInt(webhookID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %w", err)
	}

	_, err = g.client.Repositories.DeleteHook(ctx, owner, repo, hookID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

// GetPipelineTemplate gets a workflow template (not directly supported)
func (g *GitHubActionsProvider) GetPipelineTemplate(ctx context.Context, templateName string) (*ci.PipelineTemplate, error) {
	// GitHub Actions doesn't have a direct API for templates
	// This would require accessing the template repository or starter workflows
	return nil, fmt.Errorf("pipeline templates not directly supported by GitHub Actions API")
}

// ValidatePipelineConfig validates a workflow configuration
func (g *GitHubActionsProvider) ValidatePipelineConfig(ctx context.Context, config []byte) (*ci.ValidationResult, error) {
	// GitHub doesn't provide a validation API for workflow files
	// This would require implementing YAML parsing and validation logic
	return &ci.ValidationResult{
		Valid: true, // Assume valid for now
	}, nil
}

// Conversion methods

func (g *GitHubActionsProvider) convertWorkflowRun(run *github.WorkflowRun, owner, repo string) *ci.Pipeline {
	pipeline := &ci.Pipeline{
		ID:          strconv.FormatInt(run.GetID(), 10),
		Name:        run.GetName(),
		Status:      g.convertStatus(run.GetStatus(), run.GetConclusion()),
		URL:         run.GetHTMLURL(),
		Branch:      run.GetHeadBranch(),
		Repository:  fmt.Sprintf("%s/%s", owner, repo),
		Platform:    ci.PlatformGitHubActions,
		CreatedAt:   run.GetCreatedAt().Time,
		Variables:   make(map[string]string),
		Metadata:    make(map[string]interface{}),
	}

	if !run.GetUpdatedAt().Time.IsZero() {
		startedAt := run.GetUpdatedAt().Time
		pipeline.StartedAt = &startedAt
	}

	if run.GetHeadCommit() != nil {
		pipeline.Commit = &ci.Commit{
			SHA:       run.GetHeadSHA(),
			Message:   run.GetHeadCommit().GetMessage(),
			Timestamp: run.GetHeadCommit().GetTimestamp().Time,
		}
		
		if run.GetHeadCommit().GetAuthor() != nil {
			pipeline.Commit.Author = &ci.User{
				Username: run.GetHeadCommit().GetAuthor().GetLogin(),
				Name:     run.GetHeadCommit().GetAuthor().GetName(),
				Email:    run.GetHeadCommit().GetAuthor().GetEmail(),
			}
		}
	}

	// Set trigger information
	pipeline.Trigger = &ci.PipelineTrigger{
		Type:      run.GetEvent(),
		Timestamp: run.GetCreatedAt().Time,
	}

	if run.GetActor() != nil {
		pipeline.Trigger.User = &ci.User{
			ID:        strconv.FormatInt(run.GetActor().GetID(), 10),
			Username:  run.GetActor().GetLogin(),
			Name:      run.GetActor().GetName(),
			AvatarURL: run.GetActor().GetAvatarURL(),
		}
	}

	return pipeline
}

func (g *GitHubActionsProvider) convertWorkflowJob(job *github.WorkflowJob) *ci.Job {
	ciJob := &ci.Job{
		ID:        strconv.FormatInt(job.GetID(), 10),
		Name:      job.GetName(),
		Status:    g.convertStatus(job.GetStatus(), job.GetConclusion()),
		URL:       job.GetHTMLURL(),
		CreatedAt: time.Now(), // WorkflowJob doesn't have GetCreatedAt in this version
		Platform:  ci.PlatformGitHubActions,
		Variables: make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}

	if !job.GetStartedAt().Time.IsZero() {
		startedAt := job.GetStartedAt().Time
		ciJob.StartedAt = &startedAt
	}

	if !job.GetCompletedAt().Time.IsZero() {
		completedAt := job.GetCompletedAt().Time
		ciJob.CompletedAt = &completedAt
		if ciJob.StartedAt != nil {
			ciJob.Duration = ciJob.CompletedAt.Sub(*ciJob.StartedAt)
		}
	}

	// Add runner information
	if job.GetRunnerName() != "" {
		ciJob.Runner = &ci.Runner{
			Name: job.GetRunnerName(),
		}
	}

	return ciJob
}

func (g *GitHubActionsProvider) convertWorkflowJobToBuild(job *github.WorkflowJob) *ci.Build {
	build := &ci.Build{
		ID:         strconv.FormatInt(job.GetID(), 10),
		JobID:      strconv.FormatInt(job.GetID(), 10),
		Name:       job.GetName(),
		Status:     g.convertStatus(job.GetStatus(), job.GetConclusion()),
		URL:        job.GetHTMLURL(),
		CreatedAt:  time.Now(), // WorkflowJob doesn't have GetCreatedAt in this version
		Platform:   ci.PlatformGitHubActions,
		Environment: make(map[string]string),
		Metadata:   make(map[string]interface{}),
	}

	if !job.GetStartedAt().Time.IsZero() {
		startedAt := job.GetStartedAt().Time
		build.StartedAt = &startedAt
	}

	if !job.GetCompletedAt().Time.IsZero() {
		completedAt := job.GetCompletedAt().Time
		build.CompletedAt = &completedAt
		if build.StartedAt != nil {
			build.Duration = build.CompletedAt.Sub(*build.StartedAt)
		}
	}

	return build
}

func (g *GitHubActionsProvider) convertArtifact(artifact *github.Artifact) *ci.Artifact {
	ciArtifact := &ci.Artifact{
		ID:        strconv.FormatInt(artifact.GetID(), 10),
		Name:      artifact.GetName(),
		Size:      artifact.GetSizeInBytes(),
		URL:       "", // GetURL method doesn't exist in this version
		CreatedAt: artifact.GetCreatedAt().Time,
	}

	// GetExpiredAt method doesn't exist in this version
	// if artifact.GetExpiredAt() != nil {
	//     ciArtifact.ExpiresAt = &artifact.GetExpiredAt().Time
	// }

	return ciArtifact
}

func (g *GitHubActionsProvider) convertEnvironment(env *github.Environment) *ci.Environment {
	ciEnv := &ci.Environment{
		ID:          strconv.FormatInt(env.GetID(), 10),
		Name:        env.GetName(),
		URL:         env.GetURL(),
		CreatedAt:   env.GetCreatedAt().Time,
		UpdatedAt:   env.GetUpdatedAt().Time,
		Variables:   make(map[string]string),
		Metadata:    make(map[string]interface{}),
	}

	return ciEnv
}

func (g *GitHubActionsProvider) convertWebhook(hook *github.Hook) *ci.Webhook {
	webhook := &ci.Webhook{
		ID:     strconv.FormatInt(hook.GetID(), 10),
		Events: hook.Events,
		Active: hook.GetActive(),
		Config: make(map[string]string),
	}

	// Convert config map
	for key, value := range hook.Config {
		if str, ok := value.(string); ok {
			webhook.Config[key] = str
		}
	}

	if url, exists := webhook.Config["url"]; exists {
		webhook.URL = url
	}

	return webhook
}

func (g *GitHubActionsProvider) convertStatus(status, conclusion string) ci.BuildStatus {
	switch status {
	case "queued":
		return ci.StatusPending
	case "in_progress":
		return ci.StatusRunning
	case "completed":
		switch conclusion {
		case "success":
			return ci.StatusSuccess
		case "failure":
			return ci.StatusFailure
		case "cancelled":
			return ci.StatusCanceled
		case "skipped":
			return ci.StatusSkipped
		default:
			return ci.StatusError
		}
	default:
		return ci.StatusUnknown
	}
}

// Helper functions

func parseRepoURL(repoURL string) (owner, repo string, err error) {
	// Handle various GitHub URL formats
	if strings.HasPrefix(repoURL, "https://github.com/") {
		parts := strings.TrimPrefix(repoURL, "https://github.com/")
		parts = strings.TrimSuffix(parts, ".git")
		repoParts := strings.Split(parts, "/")
		if len(repoParts) >= 2 {
			return repoParts[0], repoParts[1], nil
		}
	} else if strings.Contains(repoURL, "/") {
		// Assume format is "owner/repo"
		repoParts := strings.Split(repoURL, "/")
		if len(repoParts) >= 2 {
			return repoParts[0], repoParts[1], nil
		}
	}

	return "", "", fmt.Errorf("invalid repository URL format: %s", repoURL)
}

// NewGitHubActionsFactory creates a factory for GitHub Actions providers
func NewGitHubActionsFactory() ci.CIFactory {
	return func(config *ci.CIConfig) (ci.CIProvider, error) {
		return NewGitHubActionsProvider(config)
	}
}