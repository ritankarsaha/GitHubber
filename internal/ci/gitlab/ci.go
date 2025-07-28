/*
 * GitHubber - GitLab CI/CD Provider
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: GitLab CI/CD integration for pipeline operations
 */

package gitlab

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/ritankarsaha/git-tool/internal/ci"
	"github.com/xanzy/go-gitlab"
)

// GitLabCIProvider implements the CIProvider interface for GitLab CI
type GitLabCIProvider struct {
	client        *gitlab.Client
	baseURL       string
	token         string
	authenticated bool
}

// NewGitLabCIProvider creates a new GitLab CI provider
func NewGitLabCIProvider(config *ci.CIConfig) (ci.CIProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("GitLab token is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	client, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	provider := &GitLabCIProvider{
		client:  client,
		baseURL: baseURL,
		token:   config.Token,
	}

	// Test authentication
	if err := provider.Authenticate(context.Background(), config); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return provider, nil
}

// GetPlatform returns the CI platform type
func (g *GitLabCIProvider) GetPlatform() ci.CIPlatform {
	return ci.PlatformGitLabCI
}

// GetName returns the provider name
func (g *GitLabCIProvider) GetName() string {
	return "GitLab CI"
}

// IsConnected returns whether the provider is connected
func (g *GitLabCIProvider) IsConnected() bool {
	return g.authenticated
}

// Authenticate authenticates with GitLab API
func (g *GitLabCIProvider) Authenticate(ctx context.Context, config *ci.CIConfig) error {
	// Test authentication by getting current user
	_, _, err := g.client.Users.CurrentUser()
	if err != nil {
		g.authenticated = false
		return fmt.Errorf("authentication failed: %w", err)
	}

	g.authenticated = true
	return nil
}

// ListPipelines lists GitLab pipelines
func (g *GitLabCIProvider) ListPipelines(ctx context.Context, repoURL string, options *ci.ListPipelineOptions) ([]*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	listOpts := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
		},
	}

	if options != nil {
		if options.Limit > 0 {
			listOpts.PerPage = options.Limit
		}
		if options.Status != "" {
			status := gitlab.BuildStateValue(string(options.Status))
			listOpts.Status = &status
		}
		if options.Branch != "" {
			listOpts.Ref = &options.Branch
		}
	}

	pipelines, _, err := g.client.Pipelines.ListProjectPipelines(projectID, listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pipelines: %w", err)
	}

	result := make([]*ci.Pipeline, len(pipelines))
	for i, pipeline := range pipelines {
		result[i] = g.convertPipeline(pipeline, projectID)
	}

	return result, nil
}

// GetPipeline gets a specific GitLab pipeline
func (g *GitLabCIProvider) GetPipeline(ctx context.Context, repoURL, pipelineID string) (*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(pipelineID)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	pipeline, _, err := g.client.Pipelines.GetPipeline(projectID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	result := g.convertPipelineInfo(pipeline, projectID)

	// Get jobs for this pipeline
	jobs, _, err := g.client.Jobs.ListPipelineJobs(projectID, id, &gitlab.ListJobsOptions{})
	if err == nil && jobs != nil {
		result.Jobs = make([]*ci.Job, len(jobs))
		for i, job := range jobs {
			result.Jobs[i] = g.convertJob(job)
		}
	}

	return result, nil
}

// TriggerPipeline triggers a GitLab pipeline
func (g *GitLabCIProvider) TriggerPipeline(ctx context.Context, repoURL string, request *ci.TriggerPipelineRequest) (*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	createOpts := &gitlab.CreatePipelineOptions{
		Ref: &request.Ref,
	}

	if len(request.Variables) > 0 {
		variables := make([]*gitlab.PipelineVariableOptions, 0, len(request.Variables))
		for key, value := range request.Variables {
			variables = append(variables, &gitlab.PipelineVariableOptions{
				Key:   &key,
				Value: &value,
			})
		}
		createOpts.Variables = &variables
	}

	pipeline, _, err := g.client.Pipelines.CreatePipeline(projectID, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger pipeline: %w", err)
	}

	return g.convertPipelineInfo(pipeline, projectID), nil
}

// CancelPipeline cancels a GitLab pipeline
func (g *GitLabCIProvider) CancelPipeline(ctx context.Context, repoURL, pipelineID string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(pipelineID)
	if err != nil {
		return fmt.Errorf("invalid pipeline ID: %w", err)
	}

	_, _, err = g.client.Pipelines.CancelPipelineBuild(projectID, id)
	if err != nil {
		return fmt.Errorf("failed to cancel pipeline: %w", err)
	}

	return nil
}

// RetryPipeline retries a GitLab pipeline
func (g *GitLabCIProvider) RetryPipeline(ctx context.Context, repoURL, pipelineID string) (*ci.Pipeline, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(pipelineID)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	pipeline, _, err := g.client.Pipelines.RetryPipelineBuild(projectID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retry pipeline: %w", err)
	}

	return g.convertPipelineInfo(pipeline, projectID), nil
}

// GetBuild gets a specific GitLab job
func (g *GitLabCIProvider) GetBuild(ctx context.Context, repoURL, buildID string) (*ci.Build, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(buildID)
	if err != nil {
		return nil, fmt.Errorf("invalid build ID: %w", err)
	}

	job, _, err := g.client.Jobs.GetJob(projectID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return g.convertJobToBuild(job, projectID), nil
}

// CancelBuild cancels a GitLab job
func (g *GitLabCIProvider) CancelBuild(ctx context.Context, repoURL, buildID string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(buildID)
	if err != nil {
		return fmt.Errorf("invalid build ID: %w", err)
	}

	_, _, err = g.client.Jobs.CancelJob(projectID, id)
	if err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}

	return nil
}

// GetBuildLogs gets logs for a GitLab job
func (g *GitLabCIProvider) GetBuildLogs(ctx context.Context, repoURL, buildID string) (*ci.BuildLogs, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(buildID)
	if err != nil {
		return nil, fmt.Errorf("invalid build ID: %w", err)
	}

	logs, _, err := g.client.Jobs.GetTraceFile(projectID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job logs: %w", err)
	}

	// Read the logs from the reader
	logBytes, err := io.ReadAll(logs)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	return &ci.BuildLogs{
		ID:        buildID,
		Content:   string(logBytes),
		Size:      int64(len(logBytes)),
		FetchedAt: time.Now(),
	}, nil
}

// ListArtifacts lists GitLab job artifacts
func (g *GitLabCIProvider) ListArtifacts(ctx context.Context, repoURL, pipelineID string) ([]*ci.Artifact, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(pipelineID)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	// Get jobs for the pipeline first
	jobs, _, err := g.client.Jobs.ListPipelineJobs(projectID, id, &gitlab.ListJobsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipeline jobs: %w", err)
	}

	var artifacts []*ci.Artifact
	for _, job := range jobs {
		if len(job.Artifacts) > 0 {
			for _, artifact := range job.Artifacts {
				artifacts = append(artifacts, &ci.Artifact{
					ID:        strconv.Itoa(job.ID),
					Name:      artifact.Filename,
					Path:      artifact.Filename,
					Type:      "file",
					Size:      int64(artifact.Size),
					CreatedAt: *job.CreatedAt,
				})
			}
		}
	}

	return artifacts, nil
}

// DownloadArtifact downloads a GitLab job artifact
func (g *GitLabCIProvider) DownloadArtifact(ctx context.Context, repoURL, artifactID string) ([]byte, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(artifactID)
	if err != nil {
		return nil, fmt.Errorf("invalid artifact ID: %w", err)
	}

	artifact, _, err := g.client.Jobs.GetJobArtifacts(projectID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to download artifact: %w", err)
	}

	// Read artifact bytes
	artifactBytes, err := io.ReadAll(artifact)
	if err != nil {
		return nil, fmt.Errorf("failed to read artifact: %w", err)
	}

	return artifactBytes, nil
}

// ListEnvironments lists GitLab environments
func (g *GitLabCIProvider) ListEnvironments(ctx context.Context, repoURL string) ([]*ci.Environment, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	envs, _, err := g.client.Environments.ListEnvironments(projectID, &gitlab.ListEnvironmentsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	environments := make([]*ci.Environment, len(envs))
	for i, env := range envs {
		environments[i] = g.convertEnvironment(env)
	}

	return environments, nil
}

// GetEnvironment gets a specific GitLab environment
func (g *GitLabCIProvider) GetEnvironment(ctx context.Context, repoURL, envName string) (*ci.Environment, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	// GitLab API doesn't have a direct get environment by name
	// We need to list and filter
	envs, _, err := g.client.Environments.ListEnvironments(projectID, &gitlab.ListEnvironmentsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	for _, env := range envs {
		if env.Name == envName {
			return g.convertEnvironment(env), nil
		}
	}

	return nil, fmt.Errorf("environment %s not found", envName)
}

// CreateWebhook creates a GitLab project hook
func (g *GitLabCIProvider) CreateWebhook(ctx context.Context, repoURL string, config *ci.WebhookConfig) (*ci.Webhook, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return nil, err
	}

	hookOpts := &gitlab.AddProjectHookOptions{
		URL:                   &config.URL,
		EnableSSLVerification: gitlab.Bool(!config.InsecureSSL),
	}

	if config.Secret != "" {
		hookOpts.Token = &config.Secret
	}

	// Map events to GitLab hook options
	for _, event := range config.Events {
		switch event {
		case "push":
			hookOpts.PushEvents = gitlab.Bool(true)
		case "issues":
			hookOpts.IssuesEvents = gitlab.Bool(true)
		case "merge_requests":
			hookOpts.MergeRequestsEvents = gitlab.Bool(true)
		case "pipeline":
			hookOpts.PipelineEvents = gitlab.Bool(true)
		case "job":
			hookOpts.JobEvents = gitlab.Bool(true)
		}
	}

	hook, _, err := g.client.Projects.AddProjectHook(projectID, hookOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return g.convertWebhook(hook), nil
}

// UpdateWebhook updates a GitLab project hook
func (g *GitLabCIProvider) UpdateWebhook(ctx context.Context, repoURL string, webhook *ci.Webhook) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return err
	}

	hookID, err := strconv.Atoi(webhook.ID)
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %w", err)
	}

	hookOpts := &gitlab.EditProjectHookOptions{
		URL: &webhook.URL,
	}

	// Map config to hook options
	if insecureSSL, exists := webhook.Config["insecure_ssl"]; exists {
		if insecure := insecureSSL == "true"; !insecure {
			hookOpts.EnableSSLVerification = gitlab.Bool(true)
		} else {
			hookOpts.EnableSSLVerification = gitlab.Bool(false)
		}
	}

	// Map events
	for _, event := range webhook.Events {
		switch event {
		case "push":
			hookOpts.PushEvents = gitlab.Bool(webhook.Active)
		case "issues":
			hookOpts.IssuesEvents = gitlab.Bool(webhook.Active)
		case "merge_requests":
			hookOpts.MergeRequestsEvents = gitlab.Bool(webhook.Active)
		case "pipeline":
			hookOpts.PipelineEvents = gitlab.Bool(webhook.Active)
		case "job":
			hookOpts.JobEvents = gitlab.Bool(webhook.Active)
		}
	}

	_, _, err = g.client.Projects.EditProjectHook(projectID, hookID, hookOpts)
	if err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}

// DeleteWebhook deletes a GitLab project hook
func (g *GitLabCIProvider) DeleteWebhook(ctx context.Context, repoURL, webhookID string) error {
	if !g.authenticated {
		return fmt.Errorf("not authenticated")
	}

	projectID, err := g.getProjectID(repoURL)
	if err != nil {
		return err
	}

	hookID, err := strconv.Atoi(webhookID)
	if err != nil {
		return fmt.Errorf("invalid webhook ID: %w", err)
	}

	_, err = g.client.Projects.DeleteProjectHook(projectID, hookID)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

// GetPipelineTemplate gets a GitLab CI template
func (g *GitLabCIProvider) GetPipelineTemplate(ctx context.Context, templateName string) (*ci.PipelineTemplate, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	template, _, err := g.client.CIYMLTemplate.GetTemplate(templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &ci.PipelineTemplate{
		ID:          templateName,
		Name:        template.Name,
		Content:     template.Content,
		Description: "GitLab CI template",
		Metadata:    make(map[string]interface{}),
	}, nil
}

// ValidatePipelineConfig validates a GitLab CI configuration
func (g *GitLabCIProvider) ValidatePipelineConfig(ctx context.Context, config []byte) (*ci.ValidationResult, error) {
	if !g.authenticated {
		return nil, fmt.Errorf("not authenticated")
	}

	// GitLab provides a CI lint API
	configStr := string(config)
	lintOpts := &gitlab.LintOptions{
		Content: configStr,
	}
	lintResult, _, err := g.client.Validate.Lint(lintOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	result := &ci.ValidationResult{
		Valid:    lintResult.Status == "valid",
		Errors:   make([]ci.ValidationError, 0),
		Warnings: make([]ci.ValidationError, 0),
	}

	if !result.Valid {
		for _, errMsg := range lintResult.Errors {
			result.Errors = append(result.Errors, ci.ValidationError{
				Message: errMsg,
				Type:    "error",
			})
		}
	}

	return result, nil
}

// Conversion methods

func (g *GitLabCIProvider) convertPipeline(pipeline *gitlab.PipelineInfo, projectID interface{}) *ci.Pipeline {
	return &ci.Pipeline{
		ID:         strconv.Itoa(pipeline.ID),
		Status:     g.convertStatus(pipeline.Status),
		URL:        pipeline.WebURL,
		Branch:     pipeline.Ref,
		Platform:   ci.PlatformGitLabCI,
		CreatedAt:  *pipeline.CreatedAt,
		Repository: fmt.Sprintf("%v", projectID),
		Variables:  make(map[string]string),
		Metadata:   make(map[string]interface{}),
	}
}

func (g *GitLabCIProvider) convertPipelineInfo(pipeline *gitlab.Pipeline, projectID interface{}) *ci.Pipeline {
	result := &ci.Pipeline{
		ID:         strconv.Itoa(pipeline.ID),
		Status:     g.convertStatus(pipeline.Status),
		URL:        pipeline.WebURL,
		Branch:     pipeline.Ref,
		Platform:   ci.PlatformGitLabCI,
		CreatedAt:  *pipeline.CreatedAt,
		Repository: fmt.Sprintf("%v", projectID),
		Variables:  make(map[string]string),
		Metadata:   make(map[string]interface{}),
	}

	if pipeline.UpdatedAt != nil {
		result.StartedAt = pipeline.UpdatedAt
	}

	if pipeline.StartedAt != nil {
		result.StartedAt = pipeline.StartedAt
	}

	if pipeline.FinishedAt != nil {
		result.CompletedAt = pipeline.FinishedAt
		if result.StartedAt != nil {
			result.Duration = pipeline.FinishedAt.Sub(*result.StartedAt)
		}
	}

	// Set trigger information
	result.Trigger = &ci.PipelineTrigger{
		Type:      "unknown", // GitLab API doesn't always provide trigger source
		Timestamp: *pipeline.CreatedAt,
	}

	if pipeline.User != nil {
		result.Trigger.User = &ci.User{
			ID:        strconv.Itoa(pipeline.User.ID),
			Username:  pipeline.User.Username,
			Name:      pipeline.User.Name,
			Email:     "", // BasicUser doesn't have Email field
			AvatarURL: pipeline.User.AvatarURL,
		}
	}

	return result
}

func (g *GitLabCIProvider) convertJob(job *gitlab.Job) *ci.Job {
	result := &ci.Job{
		ID:        strconv.Itoa(job.ID),
		Name:      job.Name,
		Status:    g.convertStatus(job.Status),
		URL:       job.WebURL,
		Stage:     job.Stage,
		CreatedAt: *job.CreatedAt,
		Platform:  ci.PlatformGitLabCI,
		Variables: make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}

	if job.StartedAt != nil {
		result.StartedAt = job.StartedAt
	}

	if job.FinishedAt != nil {
		result.CompletedAt = job.FinishedAt
		if result.StartedAt != nil {
			result.Duration = job.FinishedAt.Sub(*result.StartedAt)
		}
	}

	if job.Runner.ID != 0 {  // Check if runner is assigned
		result.Runner = &ci.Runner{
			ID:          strconv.Itoa(job.Runner.ID),
			Name:        job.Runner.Name,
			Description: job.Runner.Description,
			Status:      "", // Runner struct doesn't have Status field
		}
	}

	return result
}

func (g *GitLabCIProvider) convertJobToBuild(job *gitlab.Job, projectID interface{}) *ci.Build {
	result := &ci.Build{
		ID:          strconv.Itoa(job.ID),
		JobID:       strconv.Itoa(job.ID),
		Name:        job.Name,
		Status:      g.convertStatus(job.Status),
		URL:         job.WebURL,
		CreatedAt:   *job.CreatedAt,
		Platform:    ci.PlatformGitLabCI,
		Repository:  fmt.Sprintf("%v", projectID),
		Environment: make(map[string]string),
		Metadata:    make(map[string]interface{}),
	}

	if job.StartedAt != nil {
		result.StartedAt = job.StartedAt
	}

	if job.FinishedAt != nil {
		result.CompletedAt = job.FinishedAt
		if result.StartedAt != nil {
			result.Duration = job.FinishedAt.Sub(*result.StartedAt)
		}
	}

	// Convert artifacts
	if len(job.Artifacts) > 0 {
		result.Artifacts = make([]*ci.Artifact, len(job.Artifacts))
		for i, artifact := range job.Artifacts {
			result.Artifacts[i] = &ci.Artifact{
				ID:        strconv.Itoa(job.ID),
				Name:      artifact.Filename,
				Path:      artifact.Filename,
				Type:      "file",
				Size:      int64(artifact.Size),
				CreatedAt: *job.CreatedAt,
			}
		}
	}

	return result
}

func (g *GitLabCIProvider) convertEnvironment(env *gitlab.Environment) *ci.Environment {
	return &ci.Environment{
		ID:        strconv.Itoa(env.ID),
		Name:      env.Name,
		Status:    env.State,
		URL:       env.ExternalURL,
		CreatedAt: *env.CreatedAt,
		UpdatedAt: *env.UpdatedAt,
		Variables: make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}
}

func (g *GitLabCIProvider) convertWebhook(hook *gitlab.ProjectHook) *ci.Webhook {
	webhook := &ci.Webhook{
		ID:     strconv.Itoa(hook.ID),
		URL:    hook.URL,
		Events: make([]string, 0),
		Active: true, // GitLab hooks are active by default
		Config: make(map[string]string),
	}

	// Map GitLab hook events
	if hook.PushEvents {
		webhook.Events = append(webhook.Events, "push")
	}
	if hook.IssuesEvents {
		webhook.Events = append(webhook.Events, "issues")
	}
	if hook.MergeRequestsEvents {
		webhook.Events = append(webhook.Events, "merge_requests")
	}
	if hook.PipelineEvents {
		webhook.Events = append(webhook.Events, "pipeline")
	}
	if hook.JobEvents {
		webhook.Events = append(webhook.Events, "job")
	}

	webhook.Config["url"] = hook.URL
	webhook.Config["enable_ssl_verification"] = strconv.FormatBool(hook.EnableSSLVerification)

	return webhook
}

func (g *GitLabCIProvider) convertStatus(status string) ci.BuildStatus {
	switch status {
	case "created", "pending":
		return ci.StatusPending
	case "running":
		return ci.StatusRunning
	case "success":
		return ci.StatusSuccess
	case "failed":
		return ci.StatusFailure
	case "canceled":
		return ci.StatusCanceled
	case "skipped":
		return ci.StatusSkipped
	default:
		return ci.StatusUnknown
	}
}

// Helper methods

func (g *GitLabCIProvider) getProjectID(repoURL string) (interface{}, error) {
	// Handle various GitLab URL formats
	if strings.HasPrefix(repoURL, "https://gitlab.com/") {
		parts := strings.TrimPrefix(repoURL, "https://gitlab.com/")
		parts = strings.TrimSuffix(parts, ".git")
		return parts, nil
	} else if strings.Contains(repoURL, "/") {
		// Assume format is "group/project"
		return repoURL, nil
	}

	return nil, fmt.Errorf("invalid repository URL format: %s", repoURL)
}

// NewGitLabCIFactory creates a factory for GitLab CI providers
func NewGitLabCIFactory() ci.CIFactory {
	return func(config *ci.CIConfig) (ci.CIProvider, error) {
		return NewGitLabCIProvider(config)
	}
}