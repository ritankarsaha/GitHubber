/*
 * GitHubber - Provider Types and Interfaces
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Core interfaces and types for multi-platform Git hosting support
 */

package providers

import (
	"context"
	"time"
)

// ProviderType represents the type of Git hosting provider
type ProviderType string

const (
	ProviderGitHub    ProviderType = "github"
	ProviderGitLab    ProviderType = "gitlab"
	ProviderBitbucket ProviderType = "bitbucket"
	ProviderGitea     ProviderType = "gitea"
	ProviderCustom    ProviderType = "custom"
)

// Provider defines the interface that all Git hosting providers must implement
type Provider interface {
	// Provider metadata
	GetType() ProviderType
	GetName() string
	GetBaseURL() string
	
	// Authentication
	Authenticate(ctx context.Context, token string) error
	IsAuthenticated() bool
	
	// Repository operations
	GetRepository(ctx context.Context, owner, repo string) (*Repository, error)
	ListRepositories(ctx context.Context, options *ListOptions) ([]*Repository, error)
	CreateRepository(ctx context.Context, repo *CreateRepositoryRequest) (*Repository, error)
	UpdateRepository(ctx context.Context, owner, repo string, update *UpdateRepositoryRequest) (*Repository, error)
	DeleteRepository(ctx context.Context, owner, repo string) error
	ForkRepository(ctx context.Context, owner, repo string) (*Repository, error)
	
	// Pull Request operations
	GetPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error)
	ListPullRequests(ctx context.Context, owner, repo string, options *ListOptions) ([]*PullRequest, error)
	CreatePullRequest(ctx context.Context, owner, repo string, pr *CreatePullRequestRequest) (*PullRequest, error)
	UpdatePullRequest(ctx context.Context, owner, repo string, number int, update *UpdatePullRequestRequest) (*PullRequest, error)
	MergePullRequest(ctx context.Context, owner, repo string, number int, options *MergeOptions) error
	ClosePullRequest(ctx context.Context, owner, repo string, number int) error
	
	// Issue operations
	GetIssue(ctx context.Context, owner, repo string, number int) (*Issue, error)
	ListIssues(ctx context.Context, owner, repo string, options *ListOptions) ([]*Issue, error)
	CreateIssue(ctx context.Context, owner, repo string, issue *CreateIssueRequest) (*Issue, error)
	UpdateIssue(ctx context.Context, owner, repo string, number int, update *UpdateIssueRequest) (*Issue, error)
	CloseIssue(ctx context.Context, owner, repo string, number int) error
	
	// Branch operations
	ListBranches(ctx context.Context, owner, repo string) ([]*Branch, error)
	CreateBranch(ctx context.Context, owner, repo string, branch *CreateBranchRequest) (*Branch, error)
	DeleteBranch(ctx context.Context, owner, repo string, branch string) error
	
	// Tag operations
	ListTags(ctx context.Context, owner, repo string) ([]*Tag, error)
	CreateTag(ctx context.Context, owner, repo string, tag *CreateTagRequest) (*Tag, error)
	DeleteTag(ctx context.Context, owner, repo string, tag string) error
	
	// Release operations
	ListReleases(ctx context.Context, owner, repo string) ([]*Release, error)
	CreateRelease(ctx context.Context, owner, repo string, release *CreateReleaseRequest) (*Release, error)
	UpdateRelease(ctx context.Context, owner, repo string, id string, update *UpdateReleaseRequest) (*Release, error)
	DeleteRelease(ctx context.Context, owner, repo string, id string) error
	
	// Webhook operations
	ListWebhooks(ctx context.Context, owner, repo string) ([]*Webhook, error)
	CreateWebhook(ctx context.Context, owner, repo string, webhook *CreateWebhookRequest) (*Webhook, error)
	UpdateWebhook(ctx context.Context, owner, repo string, id string, update *UpdateWebhookRequest) (*Webhook, error)
	DeleteWebhook(ctx context.Context, owner, repo string, id string) error
	
	// CI/CD operations
	ListPipelines(ctx context.Context, owner, repo string) ([]*Pipeline, error)
	GetPipeline(ctx context.Context, owner, repo string, id string) (*Pipeline, error)
	TriggerPipeline(ctx context.Context, owner, repo string, options *TriggerPipelineOptions) (*Pipeline, error)
	CancelPipeline(ctx context.Context, owner, repo string, id string) error
	
	// User operations
	GetUser(ctx context.Context) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

// Repository represents a Git repository
type Repository struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Owner       *User     `json:"owner"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CloneURL    string    `json:"clone_url"`
	SSHURL      string    `json:"ssh_url"`
	Private     bool      `json:"private"`
	Fork        bool      `json:"fork"`
	Language    string    `json:"language"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	OpenIssues  int       `json:"open_issues"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Permissions struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

// PullRequest represents a pull/merge request
type PullRequest struct {
	ID          string    `json:"id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Author      *User     `json:"author"`
	Assignee    *User     `json:"assignee,omitempty"`
	Head        *Branch   `json:"head"`
	Base        *Branch   `json:"base"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Mergeable   bool      `json:"mergeable"`
	Labels      []string  `json:"labels"`
}

// Issue represents an issue
type Issue struct {
	ID          string    `json:"id"`
	Number      int       `json:"number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Author      *User     `json:"author"`
	Assignee    *User     `json:"assignee,omitempty"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Labels      []string  `json:"labels"`
}

// Branch represents a Git branch
type Branch struct {
	Name      string `json:"name"`
	SHA       string `json:"sha"`
	Protected bool   `json:"protected"`
}

// Tag represents a Git tag
type Tag struct {
	Name   string `json:"name"`
	SHA    string `json:"sha"`
	URL    string `json:"url"`
	Tagger *User  `json:"tagger,omitempty"`
}

// Release represents a repository release
type Release struct {
	ID          string    `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Author      *User     `json:"author"`
}

// User represents a user
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	URL       string `json:"url"`
}

// Webhook represents a repository webhook
type Webhook struct {
	ID     string            `json:"id"`
	URL    string            `json:"url"`
	Events []string          `json:"events"`
	Config map[string]string `json:"config"`
	Active bool              `json:"active"`
}

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Ref       string            `json:"ref"`
	SHA       string            `json:"sha"`
	URL       string            `json:"url"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Variables map[string]string `json:"variables,omitempty"`
}

// Request types
type ListOptions struct {
	Page     int    `json:"page,omitempty"`
	PerPage  int    `json:"per_page,omitempty"`
	State    string `json:"state,omitempty"`
	Sort     string `json:"sort,omitempty"`
	Order    string `json:"order,omitempty"`
	Since    string `json:"since,omitempty"`
	Labels   string `json:"labels,omitempty"`
	Assignee string `json:"assignee,omitempty"`
}

type CreateRepositoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init,omitempty"`
}

type UpdateRepositoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Private     *bool  `json:"private,omitempty"`
}

type CreatePullRequestRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Head        string `json:"head"`
	Base        string `json:"base"`
}

type UpdatePullRequestRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
}

type MergeOptions struct {
	CommitTitle     string `json:"commit_title,omitempty"`
	CommitMessage   string `json:"commit_message,omitempty"`
	MergeMethod     string `json:"merge_method,omitempty"`
	DeleteHeadBranch bool   `json:"delete_head_branch,omitempty"`
}

type CreateIssueRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Assignee    string   `json:"assignee,omitempty"`
}

type UpdateIssueRequest struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	State       string   `json:"state,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Assignee    string   `json:"assignee,omitempty"`
}

type CreateBranchRequest struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
}

type CreateTagRequest struct {
	Name    string `json:"name"`
	SHA     string `json:"sha"`
	Message string `json:"message,omitempty"`
}

type CreateReleaseRequest struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Draft       bool   `json:"draft,omitempty"`
	Prerelease  bool   `json:"prerelease,omitempty"`
}

type UpdateReleaseRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Draft       *bool  `json:"draft,omitempty"`
	Prerelease  *bool  `json:"prerelease,omitempty"`
}

type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Secret string   `json:"secret,omitempty"`
	Active bool     `json:"active"`
}

type UpdateWebhookRequest struct {
	URL    string   `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Active *bool    `json:"active,omitempty"`
}

type TriggerPipelineOptions struct {
	Ref       string            `json:"ref"`
	Variables map[string]string `json:"variables,omitempty"`
}

// ProviderConfig represents provider-specific configuration
type ProviderConfig struct {
	Type     ProviderType      `json:"type"`
	Name     string            `json:"name"`
	BaseURL  string            `json:"base_url"`
	Token    string            `json:"token"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Options  map[string]string `json:"options,omitempty"`
}

// ProviderRegistry manages available providers
type ProviderRegistry interface {
	Register(providerType ProviderType, factory ProviderFactory) error
	Create(config *ProviderConfig) (Provider, error)
	GetSupportedTypes() []ProviderType
}

// ProviderFactory creates provider instances
type ProviderFactory func(config *ProviderConfig) (Provider, error)

// RepositoryInfo contains basic repository information
type RepositoryInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}