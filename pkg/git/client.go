package git

import (
	"github.com/ritankarsaha/githubber/pkg/config"
	"github.com/ritankarsaha/githubber/pkg/logger"
)

// Client is the main Git client that provides access to all Git operations
type Client struct {
	Repository *Repository
	Branches   *BranchManager
	Commits    *CommitManager
	Staging    *Stager
	Remotes    *RemoteManager
	Stash      *StashManager
	Tags       *TagManager
	Rebase     *RebaseManager
	
	config   *config.Config
	executor CommandExecutor
}

// NewClient creates a new Git client with the provided configuration
func NewClient(cfg *config.Config) *Client {
	executor := NewCommandExecutor()
	
	client := &Client{
		Repository: NewRepository(),
		Branches:   NewBranchManager(),
		Commits:    NewCommitManager(),
		Staging:    NewStager(),
		Remotes:    NewRemoteManager(),
		Stash:      NewStashManager(),
		Tags:       NewTagManager(),
		Rebase:     NewRebaseManager(),
		
		config:   cfg,
		executor: executor,
	}
	
	// Set the same executor for all managers to ensure consistency
	client.Repository.executor = executor
	client.Branches.executor = executor
	client.Commits.executor = executor
	client.Staging.executor = executor
	client.Remotes.executor = executor
	client.Stash.executor = executor
	client.Tags.executor = executor
	client.Rebase.executor = executor
	
	return client
}

// NewTestClient creates a client for testing with a mock executor
func NewTestClient() (*Client, *TestCommandExecutor) {
	testExecutor := NewTestCommandExecutor()
	
	client := &Client{
		Repository: NewRepository(),
		Branches:   NewBranchManager(),
		Commits:    NewCommitManager(),
		Staging:    NewStager(),
		Remotes:    NewRemoteManager(),
		Stash:      NewStashManager(),
		Tags:       NewTagManager(),
		Rebase:     NewRebaseManager(),
		
		executor: testExecutor,
	}
	
	// Set the test executor for all managers
	client.Repository.executor = testExecutor
	client.Branches.executor = testExecutor
	client.Commits.executor = testExecutor
	client.Staging.executor = testExecutor
	client.Remotes.executor = testExecutor
	client.Stash.executor = testExecutor
	client.Tags.executor = testExecutor
	client.Rebase.executor = testExecutor
	
	return client, testExecutor
}

// SetWorkingDirectory sets the working directory for all Git operations
func (c *Client) SetWorkingDirectory(dir string) {
	c.executor.SetWorkingDirectory(dir)
}

// GetWorkingDirectory returns the current working directory
func (c *Client) GetWorkingDirectory() string {
	return c.executor.GetWorkingDirectory()
}

// IsGitRepository checks if the current directory is a Git repository
func (c *Client) IsGitRepository() bool {
	return c.Repository.IsGitRepository()
}

// Version returns the Git version
func (c *Client) Version() (string, error) {
	output, err := c.executor.Execute("git", "--version")
	if err != nil {
		return "", err
	}
	return output, nil
}

// Config returns the client configuration
func (c *Client) Config() *config.Config {
	return c.config
}

// ValidateGitInstallation checks if Git is properly installed and accessible
func (c *Client) ValidateGitInstallation() error {
	_, err := c.Version()
	if err != nil {
		logger.Error("Git is not installed or not accessible")
		return err
	}
	
	logger.Debug("Git installation validated successfully")
	return nil
}

// Close cleans up any resources used by the client
func (c *Client) Close() error {
	// Currently no cleanup needed, but this provides a hook for future use
	return nil
}