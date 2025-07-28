/*
 * GitHubber - CI/CD Integration Types and Interfaces
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: CI/CD abstraction layer for various platforms
 */

package ci

import (
	"context"
	"time"
)

// CIPlatform represents the CI/CD platform type
type CIPlatform string

const (
	PlatformGitHubActions CIPlatform = "github_actions"
	PlatformGitLabCI      CIPlatform = "gitlab_ci"
	PlatformJenkins       CIPlatform = "jenkins"
	PlatformCircleCI      CIPlatform = "circle_ci"
	PlatformTravisCI      CIPlatform = "travis_ci"
	PlatformBitbucket     CIPlatform = "bitbucket_pipelines"
	PlatformAzureDevOps   CIPlatform = "azure_devops"
	PlatformCustom        CIPlatform = "custom"
)

// BuildStatus represents the status of a build
type BuildStatus string

const (
	StatusPending   BuildStatus = "pending"
	StatusRunning   BuildStatus = "running"
	StatusSuccess   BuildStatus = "success"
	StatusFailure   BuildStatus = "failure"
	StatusCanceled  BuildStatus = "canceled"
	StatusSkipped   BuildStatus = "skipped"
	StatusError     BuildStatus = "error"
	StatusUnknown   BuildStatus = "unknown"
)

// CIProvider defines the interface for CI/CD platform integrations
type CIProvider interface {
	// Platform info
	GetPlatform() CIPlatform
	GetName() string
	IsConnected() bool
	
	// Authentication
	Authenticate(ctx context.Context, config *CIConfig) error
	
	// Pipeline/Workflow operations
	ListPipelines(ctx context.Context, repoURL string, options *ListPipelineOptions) ([]*Pipeline, error)
	GetPipeline(ctx context.Context, repoURL, pipelineID string) (*Pipeline, error)
	TriggerPipeline(ctx context.Context, repoURL string, request *TriggerPipelineRequest) (*Pipeline, error)
	CancelPipeline(ctx context.Context, repoURL, pipelineID string) error
	RetryPipeline(ctx context.Context, repoURL, pipelineID string) (*Pipeline, error)
	
	// Build operations
	GetBuild(ctx context.Context, repoURL, buildID string) (*Build, error)
	CancelBuild(ctx context.Context, repoURL, buildID string) error
	GetBuildLogs(ctx context.Context, repoURL, buildID string) (*BuildLogs, error)
	
	// Artifact operations
	ListArtifacts(ctx context.Context, repoURL, pipelineID string) ([]*Artifact, error)
	DownloadArtifact(ctx context.Context, repoURL, artifactID string) ([]byte, error)
	
	// Environment operations
	ListEnvironments(ctx context.Context, repoURL string) ([]*Environment, error)
	GetEnvironment(ctx context.Context, repoURL, envName string) (*Environment, error)
	
	// Webhook operations
	CreateWebhook(ctx context.Context, repoURL string, config *WebhookConfig) (*Webhook, error)
	UpdateWebhook(ctx context.Context, repoURL string, webhook *Webhook) error
	DeleteWebhook(ctx context.Context, repoURL, webhookID string) error
	
	// Template/Configuration operations
	GetPipelineTemplate(ctx context.Context, templateName string) (*PipelineTemplate, error)
	ValidatePipelineConfig(ctx context.Context, config []byte) (*ValidationResult, error)
}

// CIConfig represents CI/CD platform configuration
type CIConfig struct {
	Platform    CIPlatform        `json:"platform"`
	BaseURL     string            `json:"base_url,omitempty"`
	Token       string            `json:"token"`
	Username    string            `json:"username,omitempty"`
	Password    string            `json:"password,omitempty"`
	APIKey      string            `json:"api_key,omitempty"`
	Settings    map[string]string `json:"settings,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
	
	// Advanced options
	RetryCount  int    `json:"retry_count,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	Insecure    bool   `json:"insecure,omitempty"`
}

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Status        BuildStatus            `json:"status"`
	URL           string                 `json:"url"`
	Branch        string                 `json:"branch"`
	Commit        *Commit                `json:"commit,omitempty"`
	Trigger       *PipelineTrigger       `json:"trigger,omitempty"`
	Environment   *Environment           `json:"environment,omitempty"`
	Variables     map[string]string      `json:"variables,omitempty"`
	
	// Timing
	CreatedAt     time.Time              `json:"created_at"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	Duration      time.Duration          `json:"duration,omitempty"`
	
	// Structure
	Stages        []*Stage               `json:"stages,omitempty"`
	Jobs          []*Job                 `json:"jobs,omitempty"`
	
	// Metadata
	Platform      CIPlatform             `json:"platform"`
	Repository    string                 `json:"repository"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Build represents a single build within a pipeline
type Build struct {
	ID            string                 `json:"id"`
	PipelineID    string                 `json:"pipeline_id"`
	JobID         string                 `json:"job_id"`
	Name          string                 `json:"name"`
	Status        BuildStatus            `json:"status"`
	URL           string                 `json:"url"`
	
	// Configuration
	Image         string                 `json:"image,omitempty"`
	Commands      []string               `json:"commands,omitempty"`
	Environment   map[string]string      `json:"environment,omitempty"`
	Services      []*Service             `json:"services,omitempty"`
	
	// Timing
	CreatedAt     time.Time              `json:"created_at"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	Duration      time.Duration          `json:"duration,omitempty"`
	
	// Results
	ExitCode      *int                   `json:"exit_code,omitempty"`
	Artifacts     []*Artifact            `json:"artifacts,omitempty"`
	TestResults   *TestResults           `json:"test_results,omitempty"`
	CoverageData  *CoverageData          `json:"coverage_data,omitempty"`
	
	// Metadata
	Platform      CIPlatform             `json:"platform"`
	Repository    string                 `json:"repository"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Stage represents a pipeline stage
type Stage struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Status      BuildStatus            `json:"status"`
	Order       int                    `json:"order"`
	Jobs        []*Job                 `json:"jobs"`
	
	// Timing
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Duration    time.Duration          `json:"duration,omitempty"`
	
	// Configuration
	Condition   string                 `json:"condition,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Job represents a pipeline job
type Job struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Status        BuildStatus            `json:"status"`
	URL           string                 `json:"url,omitempty"`
	Stage         string                 `json:"stage,omitempty"`
	
	// Configuration
	Image         string                 `json:"image,omitempty"`
	Script        []string               `json:"script,omitempty"`
	BeforeScript  []string               `json:"before_script,omitempty"`
	AfterScript   []string               `json:"after_script,omitempty"`
	Variables     map[string]string      `json:"variables,omitempty"`
	Services      []*Service             `json:"services,omitempty"`
	Cache         *CacheConfig           `json:"cache,omitempty"`
	
	// Timing
	CreatedAt     time.Time              `json:"created_at"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	Duration      time.Duration          `json:"duration,omitempty"`
	
	// Results
	ExitCode      *int                   `json:"exit_code,omitempty"`
	Artifacts     []*Artifact            `json:"artifacts,omitempty"`
	Coverage      float64                `json:"coverage,omitempty"`
	
	// Metadata
	Platform      CIPlatform             `json:"platform"`
	Runner        *Runner                `json:"runner,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Service represents a service container
type Service struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Command     []string          `json:"command,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Ports       []string          `json:"ports,omitempty"`
	Volumes     []string          `json:"volumes,omitempty"`
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Key      string   `json:"key"`
	Paths    []string `json:"paths"`
	Policy   string   `json:"policy,omitempty"`
	Fallback []string `json:"fallback,omitempty"`
}

// Runner represents a CI/CD runner
type Runner struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Status       string            `json:"status"`
	Architecture string            `json:"architecture"`
	OS           string            `json:"os"`
	Tags         []string          `json:"tags"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Artifact represents a build artifact
type Artifact struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	Size        int64     `json:"size"`
	URL         string    `json:"url,omitempty"`
	DownloadURL string    `json:"download_url,omitempty"`
	Checksum    string    `json:"checksum,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// Environment represents a deployment environment
type Environment struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	URL         string                 `json:"url,omitempty"`
	Description string                 `json:"description,omitempty"`
	Variables   map[string]string      `json:"variables,omitempty"`
	Secrets     []string               `json:"secrets,omitempty"`
	
	// Deployment info
	LastDeployment *Deployment          `json:"last_deployment,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Deployment represents a deployment
type Deployment struct {
	ID            string                 `json:"id"`
	Environment   string                 `json:"environment"`
	Status        string                 `json:"status"`
	Description   string                 `json:"description,omitempty"`
	URL           string                 `json:"url,omitempty"`
	Ref           string                 `json:"ref"`
	SHA           string                 `json:"sha"`
	
	// Timing
	CreatedAt     time.Time              `json:"created_at"`
	StartedAt     *time.Time             `json:"started_at,omitempty"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	
	// Metadata
	Creator       *User                  `json:"creator,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Commit represents a Git commit
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    *User     `json:"author,omitempty"`
	Committer *User     `json:"committer,omitempty"`
	URL       string    `json:"url,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// User represents a user
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// PipelineTrigger represents what triggered the pipeline
type PipelineTrigger struct {
	Type      string    `json:"type"`
	Source    string    `json:"source"`
	User      *User     `json:"user,omitempty"`
	Event     string    `json:"event,omitempty"`
	Ref       string    `json:"ref,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// TestResults represents test execution results
type TestResults struct {
	Total   int             `json:"total"`
	Passed  int             `json:"passed"`
	Failed  int             `json:"failed"`
	Skipped int             `json:"skipped"`
	Suites  []*TestSuite    `json:"suites,omitempty"`
}

// TestSuite represents a test suite
type TestSuite struct {
	Name      string      `json:"name"`
	Tests     int         `json:"tests"`
	Failures  int         `json:"failures"`
	Errors    int         `json:"errors"`
	Skipped   int         `json:"skipped"`
	Time      float64     `json:"time"`
	TestCases []*TestCase `json:"testcases,omitempty"`
}

// TestCase represents a test case
type TestCase struct {
	Name      string  `json:"name"`
	ClassName string  `json:"classname"`
	Time      float64 `json:"time"`
	Status    string  `json:"status"`
	Error     string  `json:"error,omitempty"`
	Failure   string  `json:"failure,omitempty"`
}

// CoverageData represents code coverage information
type CoverageData struct {
	Percentage float64               `json:"percentage"`
	Lines      *CoverageLines        `json:"lines,omitempty"`
	Branches   *CoverageBranches     `json:"branches,omitempty"`
	Functions  *CoverageFunctions    `json:"functions,omitempty"`
	Files      []*FileCoverage       `json:"files,omitempty"`
}

// CoverageLines represents line coverage
type CoverageLines struct {
	Total   int     `json:"total"`
	Covered int     `json:"covered"`
	Percent float64 `json:"percent"`
}

// CoverageBranches represents branch coverage
type CoverageBranches struct {
	Total   int     `json:"total"`
	Covered int     `json:"covered"`
	Percent float64 `json:"percent"`
}

// CoverageFunctions represents function coverage
type CoverageFunctions struct {
	Total   int     `json:"total"`
	Covered int     `json:"covered"`
	Percent float64 `json:"percent"`
}

// FileCoverage represents coverage for a single file
type FileCoverage struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Lines    int     `json:"lines"`
	Covered  int     `json:"covered"`
	Percent  float64 `json:"percent"`
}

// BuildLogs represents build log data
type BuildLogs struct {
	ID       string    `json:"id"`
	Content  string    `json:"content"`
	Size     int64     `json:"size"`
	URL      string    `json:"url,omitempty"`
	Encoding string    `json:"encoding,omitempty"`
	FetchedAt time.Time `json:"fetched_at"`
}

// Webhook represents a CI/CD webhook
type Webhook struct {
	ID       string            `json:"id"`
	URL      string            `json:"url"`
	Events   []string          `json:"events"`
	Active   bool              `json:"active"`
	Config   map[string]string `json:"config"`
	Secret   string            `json:"secret,omitempty"`
	InsecureSSL bool           `json:"insecure_ssl,omitempty"`
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	URL         string   `json:"url"`
	Events      []string `json:"events"`
	Secret      string   `json:"secret,omitempty"`
	InsecureSSL bool     `json:"insecure_ssl,omitempty"`
}

// PipelineTemplate represents a pipeline template
type PipelineTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Language    string                 `json:"language,omitempty"`
	Framework   string                 `json:"framework,omitempty"`
	Content     string                 `json:"content"`
	Variables   []TemplateVariable     `json:"variables,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TemplateVariable represents a template variable
type TemplateVariable struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Required    bool        `json:"required"`
}

// ValidationResult represents pipeline configuration validation result
type ValidationResult struct {
	Valid   bool               `json:"valid"`
	Errors  []ValidationError  `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

// ValidationError represents a validation error or warning
type ValidationError struct {
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Request types
type ListPipelineOptions struct {
	Status    BuildStatus `json:"status,omitempty"`
	Branch    string      `json:"branch,omitempty"`
	Limit     int         `json:"limit,omitempty"`
	Offset    int         `json:"offset,omitempty"`
	Sort      string      `json:"sort,omitempty"`
	Order     string      `json:"order,omitempty"`
}

type TriggerPipelineRequest struct {
	Ref         string            `json:"ref"`
	Variables   map[string]string `json:"variables,omitempty"`
	Environment string            `json:"environment,omitempty"`
	Message     string            `json:"message,omitempty"`
}

// CIManager manages multiple CI/CD providers
type CIManager interface {
	// Provider management
	RegisterProvider(name string, provider CIProvider) error
	GetProvider(name string) (CIProvider, error)
	ListProviders() map[string]CIProvider
	
	// Repository management
	RegisterRepository(repoURL string, providers []string) error
	GetRepositoryProviders(repoURL string) ([]CIProvider, error)
	
	// Unified operations
	TriggerBuilds(ctx context.Context, repoURL string, request *TriggerPipelineRequest) ([]*Pipeline, error)
	GetPipelineStatus(ctx context.Context, repoURL, pipelineID string) (*Pipeline, error)
	CancelPipelines(ctx context.Context, repoURL, pipelineID string) error
}

// CIFactory creates CI providers
type CIFactory func(config *CIConfig) (CIProvider, error)

// CIEvent represents CI/CD events for webhooks
type CIEvent struct {
	Type       string                 `json:"type"`
	Platform   CIPlatform             `json:"platform"`
	Repository string                 `json:"repository"`
	Pipeline   *Pipeline              `json:"pipeline,omitempty"`
	Build      *Build                 `json:"build,omitempty"`
	Job        *Job                   `json:"job,omitempty"`
	Deployment *Deployment            `json:"deployment,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Statistics and monitoring
type CIMetrics struct {
	PipelinesTotal     int64         `json:"pipelines_total"`
	PipelinesSuccess   int64         `json:"pipelines_success"`
	PipelinesFailure   int64         `json:"pipelines_failure"`
	AverageDuration    time.Duration `json:"average_duration"`
	SuccessRate        float64       `json:"success_rate"`
	DeploymentsTotal   int64         `json:"deployments_total"`
	ActiveEnvironments int64         `json:"active_environments"`
	CollectedAt        time.Time     `json:"collected_at"`
}

// Status monitoring
type CIStatus struct {
	Provider      CIPlatform `json:"provider"`
	Connected     bool       `json:"connected"`
	LastCheck     time.Time  `json:"last_check"`
	ResponseTime  time.Duration `json:"response_time"`
	ErrorCount    int64      `json:"error_count"`
	Version       string     `json:"version,omitempty"`
	Capabilities  []string   `json:"capabilities,omitempty"`
}