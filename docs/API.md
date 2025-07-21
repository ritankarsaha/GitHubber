# GitHubber API Documentation

**Author**: Ritankar Saha <ritankar.saha786@gmail.com>

This document provides comprehensive documentation for GitHubber's internal APIs and modules.

## Table of Contents

- [Git Operations Module](#git-operations-module)
- [GitHub API Client](#github-api-client)
- [Configuration System](#configuration-system)
- [UI Styling System](#ui-styling-system)
- [CLI Interface](#cli-interface)

## Git Operations Module

**Package**: `internal/git`

### Core Functions

#### `RunCommand(cmd string) (string, error)`
Executes a git command and returns the output.

**Parameters:**
- `cmd`: Git command to execute

**Returns:**
- `string`: Command output
- `error`: Error if command fails

**Example:**
```go
output, err := git.RunCommand("git status")
if err != nil {
    log.Fatal(err)
}
fmt.Println(output)
```

#### `GetRepositoryInfo() (*RepositoryInfo, error)`
Gets information about the current Git repository.

**Returns:**
- `*RepositoryInfo`: Repository information structure
- `error`: Error if not in a Git repository

**RepositoryInfo Structure:**
```go
type RepositoryInfo struct {
    URL           string
    CurrentBranch string
}
```

### Repository Operations

#### `Init() error`
Initializes a new Git repository in the current directory.

#### `Clone(url string) error`
Clones a repository from the specified URL.

**Parameters:**
- `url`: Repository URL (HTTPS or SSH)

### Branch Operations

#### `CreateBranch(name string) error`
Creates and switches to a new branch.

**Parameters:**
- `name`: Branch name

#### `DeleteBranch(name string) error`
Deletes a branch (force delete with -D flag).

**Parameters:**
- `name`: Branch name to delete

#### `SwitchBranch(name string) error`
Switches to an existing branch.

**Parameters:**
- `name`: Branch name to switch to

#### `ListBranches() ([]string, error)`
Lists all local branches.

**Returns:**
- `[]string`: List of branch names
- `error`: Error if operation fails

### Commit Operations

#### `Status() (string, error)`
Returns the current repository status.

#### `AddFiles(files ...string) error`
Stages files for commit. If no files specified, stages all changes.

**Parameters:**
- `files`: Optional list of files to stage

#### `Commit(message string) error`
Creates a commit with the specified message.

**Parameters:**
- `message`: Commit message

### Advanced Operations

#### `SquashCommits(baseCommit, message string) error`
Squashes commits from HEAD to the specified base commit.

**Parameters:**
- `baseCommit`: Base commit hash to squash into
- `message`: New commit message

#### `GetRecentCommits(n int) ([]Commit, error)`
Gets the most recent commits.

**Parameters:**
- `n`: Number of commits to retrieve

**Returns:**
- `[]Commit`: List of commit structures
- `error`: Error if operation fails

**Commit Structure:**
```go
type Commit struct {
    Hash    string
    Message string
}
```

## GitHub API Client

**Package**: `internal/github`

### Client Creation

#### `NewClient() (*Client, error)`
Creates a new GitHub API client using environment token or configuration.

**Returns:**
- `*Client`: GitHub API client
- `error`: Error if authentication fails

#### `NewClientWithToken(token string) *Client`
Creates a new GitHub API client with a specific token.

**Parameters:**
- `token`: GitHub personal access token

### Repository Operations

#### `(c *Client) GetRepository(owner, repo string) (*Repository, error)`
Gets detailed information about a GitHub repository.

**Parameters:**
- `owner`: Repository owner
- `repo`: Repository name

**Returns:**
- `*Repository`: Repository information
- `error`: Error if operation fails

**Repository Structure:**
```go
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
```

### Pull Request Operations

#### `(c *Client) CreatePullRequest(owner, repo, title, body, head, base string) (*PullRequest, error)`
Creates a new pull request.

**Parameters:**
- `owner`: Repository owner
- `repo`: Repository name
- `title`: PR title
- `body`: PR description
- `head`: Source branch
- `base`: Target branch

#### `(c *Client) ListPullRequests(owner, repo, state string) ([]*PullRequest, error)`
Lists pull requests for a repository.

**Parameters:**
- `owner`: Repository owner
- `repo`: Repository name
- `state`: PR state ("open", "closed", "all")

**PullRequest Structure:**
```go
type PullRequest struct {
    Number int
    Title  string
    State  string
    Author string
    URL    string
}
```

### Issue Operations

#### `(c *Client) ListIssues(owner, repo, state string) ([]*Issue, error)`
Lists issues for a repository.

**Parameters:**
- `owner`: Repository owner
- `repo`: Repository name
- `state`: Issue state ("open", "closed", "all")

**Issue Structure:**
```go
type Issue struct {
    Number int
    Title  string
    State  string
    Author string
    URL    string
}
```

### Utility Functions

#### `ParseRepoURL(url string) (owner, repo string, err error)`
Parses a GitHub repository URL to extract owner and repository name.

**Parameters:**
- `url`: Repository URL (HTTPS or SSH format)

**Returns:**
- `owner`: Repository owner
- `repo`: Repository name
- `err`: Error if URL format is invalid

**Supported URL Formats:**
- `https://github.com/owner/repo`
- `https://github.com/owner/repo.git`
- `git@github.com:owner/repo.git`

## Configuration System

**Package**: `internal/config`

### Configuration Structure

```go
type Config struct {
    GitHub GitHubConfig `json:"github"`
    UI     UIConfig     `json:"ui"`
    Git    GitConfig    `json:"git"`
}

type GitHubConfig struct {
    Token        string `json:"token,omitempty"`
    DefaultOwner string `json:"default_owner"`
    DefaultRepo  string `json:"default_repo"`
    APIBaseURL   string `json:"api_base_url,omitempty"`
}

type UIConfig struct {
    Theme       string `json:"theme"`
    ShowEmojis  bool   `json:"show_emojis"`
    PageSize    int    `json:"page_size"`
    BorderStyle string `json:"border_style"`
}

type GitConfig struct {
    DefaultBranch string `json:"default_branch"`
    AutoPush      bool   `json:"auto_push"`
    SignCommits   bool   `json:"sign_commits"`
}
```

### Configuration Functions

#### `Load() (*Config, error)`
Loads configuration from file or returns default configuration.

#### `(c *Config) Save() error`
Saves the configuration to file.

#### `(c *Config) SetGitHubToken(token string) error`
Sets and saves the GitHub token.

#### `(c *Config) GetGitHubToken() string`
Gets the GitHub token from environment or configuration.

#### `GetDefaultConfig() *Config`
Returns the default configuration.

## UI Styling System

**Package**: `internal/ui`

### Style Constants

The UI system provides pre-defined styles for consistent terminal output:

- `TitleStyle`: For main titles
- `SubtitleStyle`: For subtitles
- `MenuHeaderStyle`: For menu section headers
- `MenuItemStyle`: For menu items
- `SuccessStyle`: For success messages
- `ErrorStyle`: For error messages
- `WarningStyle`: For warning messages
- `InfoStyle`: For informational messages

### Icon Constants

```go
const (
    IconRepository = "📂"
    IconBranch     = "🌿"
    IconCommit     = "💾"
    IconRemote     = "🔄"
    IconHistory    = "📜"
    IconStash      = "📦"
    IconTag        = "🏷️"
    IconSuccess    = "✅"
    IconError      = "❌"
    IconWarning    = "⚠️"
    IconInfo       = "ℹ️"
    IconTool       = "🛠"
    IconGitHub     = "🐙"
    IconConfig     = "⚙️"
    IconExit       = "👋"
)
```

### Formatting Functions

#### `FormatTitle(text string) string`
Formats text as a main title with styling and icon.

#### `FormatMenuHeader(icon, text string) string`
Formats a menu section header with icon.

#### `FormatMenuItem(number int, text string) string`
Formats a numbered menu item.

#### `FormatSuccess(text string) string`
Formats a success message with green styling and checkmark icon.

#### `FormatError(text string) string`
Formats an error message with red styling and error icon.

#### `FormatWarning(text string) string`
Formats a warning message with yellow styling and warning icon.

#### `FormatInfo(text string) string`
Formats an informational message with blue styling and info icon.

#### `FormatPrompt(text string) string`
Formats a user input prompt.

#### `FormatRepoInfo(url, branch string) string`
Formats repository information in a styled box.

#### `FormatBox(content string) string`
Wraps content in a styled border box.

#### `FormatCode(content string) string`
Formats content as code with monospace styling.

## CLI Interface

**Package**: `internal/cli`

### Menu System

#### `StartMenu()`
Starts the main interactive menu loop. This function displays the menu options and handles user input.

#### `GetInput(prompt string) string`
Gets user input with the specified prompt.

**Parameters:**
- `prompt`: Prompt message to display

**Returns:**
- `string`: User input (trimmed)

### Menu Handlers

Each menu option has a corresponding handler function:

- `handleInit()`: Initialize repository
- `handleClone()`: Clone repository
- `handleCreateBranch()`: Create new branch
- `handleDeleteBranch()`: Delete branch
- `handleSwitchBranch()`: Switch branch
- `handleListBranches()`: List branches
- `handleStatus()`: Show repository status
- `handleAddFiles()`: Stage files
- `handleCommit()`: Create commit
- `handlePush()`: Push changes
- `handlePull()`: Pull changes
- `handleFetch()`: Fetch updates
- `handleLog()`: Show commit log
- `handleDiff()`: Show file differences
- `handleSquash()`: Squash commits
- `handleStashSave()`: Save stash
- `handleStashPop()`: Apply stash
- `handleStashList()`: List stashes
- `handleCreateTag()`: Create tag
- `handleDeleteTag()`: Delete tag
- `handleListTags()`: List tags
- `handleRepoInfo()`: Show GitHub repository info
- `handleCreatePR()`: Create pull request
- `handleListIssues()`: List GitHub issues
- `handleSettings()`: Manage settings

### Error Handling

All handlers follow a consistent error handling pattern:

1. Perform the operation
2. Check for errors
3. Display appropriate success or error message using UI formatting
4. Return gracefully

## Error Handling Patterns

### Git Operations
```go
if err := git.SomeOperation(); err != nil {
    fmt.Println(ui.FormatError(fmt.Sprintf("Operation failed: %v", err)))
    return
}
fmt.Println(ui.FormatSuccess("Operation completed successfully!"))
```

### GitHub API Operations
```go
client, err := github.NewClient()
if err != nil {
    fmt.Println(ui.FormatError(fmt.Sprintf("Failed to create GitHub client: %v", err)))
    fmt.Println(ui.FormatInfo("Please set GITHUB_TOKEN environment variable"))
    return
}
```

## Development Guidelines

### Adding New Features

1. **Follow the existing package structure**
2. **Add comprehensive error handling**
3. **Use the UI styling system for consistent output**
4. **Add tests for new functionality**
5. **Update documentation**

### Code Style

- Use meaningful variable and function names
- Add comments for exported functions
- Follow Go conventions and best practices
- Handle errors gracefully
- Use the UI formatting functions for all output

### Testing

- Write unit tests for all new functions
- Use test helpers for Git operations
- Mock external dependencies (GitHub API)
- Test error conditions

## Examples

### Adding a New Git Operation

```go
// internal/git/commands.go
func NewGitOperation(param string) error {
    _, err := RunCommand(fmt.Sprintf("git some-command %s", param))
    return err
}

// internal/cli/menu.go
func handleNewOperation() {
    param := GetInput(ui.FormatPrompt("Enter parameter: "))
    if err := git.NewGitOperation(param); err != nil {
        fmt.Println(ui.FormatError(fmt.Sprintf("Operation failed: %v", err)))
        return
    }
    fmt.Println(ui.FormatSuccess("Operation completed successfully!"))
}
```

### Adding a New GitHub API Function

```go
// internal/github/client.go
func (c *Client) NewGitHubOperation(param string) (*Result, error) {
    // Implementation using c.client
    result, _, err := c.client.SomeService.SomeMethod(c.ctx, param)
    if err != nil {
        return nil, fmt.Errorf("failed to perform operation: %w", err)
    }
    return &Result{...}, nil
}
```

This API documentation provides a comprehensive guide for understanding and extending GitHubber's functionality.