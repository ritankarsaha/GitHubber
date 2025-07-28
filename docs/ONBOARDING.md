# 🚀 New Contributor Onboarding Guide

Welcome to GitHubber! This guide will help you get up and running as a contributor to this project. Whether you're new to Go, Git, or open source in general, this document will walk you through everything you need to know.

## 📋 Table of Contents

- [Welcome](#welcome)
- [Project Overview](#project-overview)
- [Quick Setup Guide](#quick-setup-guide)
- [Understanding the Codebase](#understanding-the-codebase)
- [Your First Contribution](#your-first-contribution)
- [Development Environment](#development-environment)
- [Testing and Quality Assurance](#testing-and-quality-assurance)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

## 🎉 Welcome

Thank you for your interest in contributing to GitHubber! We're excited to have you join our community. This project aims to make Git and GitHub operations more accessible and enjoyable through a beautiful command-line interface.

### What You'll Learn

By contributing to GitHubber, you'll gain experience with:

- **Go Programming**: Modern Go development practices
- **CLI Development**: Building user-friendly command-line tools
- **Git Operations**: Advanced Git workflows and automation
- **GitHub API Integration**: Working with REST APIs
- **Open Source Collaboration**: Contributing to open source projects
- **Testing**: Writing comprehensive tests for CLI applications

## 🔍 Project Overview

### What is GitHubber?

GitHubber is an advanced command-line interface that enhances your Git and GitHub workflow with:

- **Beautiful Terminal UI**: Colorful, interactive menus
- **Git Operations**: Repository management, branching, committing
- **GitHub Integration**: PR creation, issue management, repository stats
- **Advanced Features**: Interactive commit squashing, stash management

### Key Technologies

- **Language**: Go 1.23+
- **UI Libraries**: Charm Bracelet's Lipgloss for styling
- **APIs**: GitHub REST API, GitLab API
- **Build System**: Make-based build system
- **Testing**: Go's built-in testing framework

### Project Philosophy

- **User Experience First**: Every feature should be intuitive and helpful
- **Beautiful Output**: Terminal output should be colorful and well-formatted
- **Comprehensive Testing**: All functionality should be thoroughly tested
- **Clean Code**: Code should be readable, maintainable, and well-documented

## ⚡ Quick Setup Guide

### 1. Environment Prerequisites

```bash
# Check if you have the required tools
go version    # Should be 1.23 or higher
git --version # Any modern version
make --version # Usually pre-installed

# Install additional tools (optional but recommended)
gh --version  # GitHub CLI for easier GitHub operations
```

### 2. One-Minute Setup

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/GitHubber.git
cd GitHubber

# 2. Setup development environment
make deps     # Download dependencies
make build    # Build the application
make test     # Ensure everything works

# 3. Run the application
./build/githubber
```

### 3. Verify Your Setup

```bash
# Run the application in a test Git repository
cd /path/to/any/git/repository
/path/to/GitHubber/build/githubber

# You should see the beautiful GitHubber interface!
```

## 🏗 Understanding the Codebase

### High-Level Architecture

```
GitHubber follows a modular architecture:

┌─────────────────┐
│   cmd/main.go   │  ← Application entry point
└─────────────────┘
         │
    ┌────▼────┐
    │   CLI   │  ← Command-line interface and menus
    └────┬────┘
         │
    ┌────▼────┐
    │   Git   │  ← Git operations (status, commit, etc.)
    └────┬────┘
         │
    ┌────▼────┐
    │ GitHub  │  ← GitHub API integration
    └────┬────┘
         │
    ┌────▼────┐
    │   UI    │  ← Terminal styling and formatting
    └─────────┘
```

### Directory Structure Deep Dive

```
GitHubber/
├── cmd/
│   └── main.go                 # Entry point - starts the application
│
├── internal/                   # Private application code
│   ├── cli/                    # Command-line interface
│   │   ├── input.go           # User input handling
│   │   ├── menu.go            # Menu system and navigation
│   │   ├── args.go            # Command-line argument parsing
│   │   ├── completion.go      # Shell completion
│   │   └── conflict.go        # Merge conflict resolution
│   │
│   ├── git/                    # Git operations
│   │   ├── commands.go        # Core Git commands (status, commit, etc.)
│   │   ├── squash.go          # Interactive commit squashing
│   │   ├── utils.go           # Git utilities and helpers
│   │   └── *_test.go          # Test files
│   │
│   ├── github/                 # GitHub API integration
│   │   └── client.go          # GitHub API client and operations
│   │
│   ├── config/                 # Configuration management
│   │   ├── config.go          # Configuration loading and saving
│   │   ├── manager.go         # Configuration management
│   │   └── types.go           # Configuration data structures
│   │
│   ├── ui/                     # Terminal UI components
│   │   └── styles.go          # Color schemes and formatting
│   │
│   ├── logging/                # Logging infrastructure
│   ├── plugins/                # Plugin system (extensibility)
│   ├── providers/              # External service providers
│   ├── ci/                     # CI/CD integration
│   └── webhooks/               # Webhook handling
│
├── docs/                       # Documentation
├── scripts/                    # Build and utility scripts
├── tests/                      # Integration tests and fixtures
└── examples/                   # Usage examples
```

### Key Files to Understand

#### 1. `cmd/main.go` - Application Entry Point

```go
// This file:
// - Checks if Git is installed
// - Parses command-line arguments
// - Shows the main menu interface
// - Handles both interactive and direct command modes
```

#### 2. `internal/cli/menu.go` - Menu System

```go
// This file contains:
// - Main menu loop and navigation
// - All menu options (Repository, Branch, Changes, etc.)
// - User interaction handling
// - Menu state management
```

#### 3. `internal/git/commands.go` - Git Operations

```go
// This file implements:
// - All Git commands (status, add, commit, push, etc.)
// - Repository information gathering
// - Branch management
// - Error handling for Git operations
```

#### 4. `internal/ui/styles.go` - Terminal Styling

```go
// This file defines:
// - Color schemes and themes
// - Text formatting functions
// - Icons and emojis
// - Styling constants
```

### Code Patterns and Conventions

#### Error Handling Pattern
```go
// Standard error handling in GitHubber
func SomeGitOperation() error {
    output, err := git.RunCommand("git status")
    if err != nil {
        return fmt.Errorf("failed to get git status: %w", err)
    }
    
    // Process output...
    return nil
}
```

#### Menu Item Pattern
```go
// Menu items follow this structure
type MenuItem struct {
    Label       string
    Description string
    Handler     func() error
    Condition   func() bool  // When to show this item
}
```

#### Styling Pattern
```go
// UI components use consistent styling
fmt.Println(ui.FormatSuccess("Operation completed successfully!"))
fmt.Println(ui.FormatError("Something went wrong"))
fmt.Println(ui.FormatWarning("This is a warning"))
```

## 🎯 Your First Contribution

### Step 1: Find a Good First Issue

Look for issues labeled with:
- `good first issue` - Perfect for newcomers
- `help wanted` - Community help needed
- `bug` - Bug fixes are great starting points
- `documentation` - Improve docs and examples

### Step 2: Create Your Development Branch

```bash
# Update your local repository
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/your-feature-name

# For example:
git checkout -b feature/add-git-stash-menu
git checkout -b fix/handle-empty-repository
git checkout -b docs/improve-setup-guide
```

### Step 3: Start Small

Great first contributions include:

#### Option A: Add a New Menu Item
```go
// In internal/cli/menu.go, add a new menu option
{
    Label:       "🏷️  Manage Tags",
    Description: "Create, list, and delete Git tags",
    Handler:     handleTags,
    Condition:   git.IsGitRepository,
}
```

#### Option B: Improve Error Messages
```go
// In internal/git/commands.go, enhance error handling
if err != nil {
    return fmt.Errorf("failed to create branch '%s': %w\nHint: Make sure the branch name doesn't already exist", branchName, err)
}
```

#### Option C: Add a New Git Command
```go
// Add a new function to internal/git/commands.go
func GetTags() ([]string, error) {
    output, err := RunCommand("git tag -l")
    if err != nil {
        return nil, fmt.Errorf("failed to list tags: %w", err)
    }
    
    tags := strings.Split(strings.TrimSpace(output), "\n")
    return tags, nil
}
```

### Step 4: Test Your Changes

```bash
# Run tests
make test

# Test manually
make build
./build/githubber

# Test in different scenarios:
# - In a Git repository
# - Outside a Git repository  
# - With various Git states (clean, dirty, etc.)
```

### Step 5: Submit Your Pull Request

```bash
# Commit your changes
git add .
git commit -m "feat: add git tag management menu"

# Push to your fork
git push origin feature/add-git-stash-menu

# Create a pull request on GitHub
gh pr create --title "Add Git tag management menu" --body "Adds a new menu for creating, listing, and deleting Git tags"
```

## 💻 Development Environment

### IDE Setup

#### VS Code (Recommended)
```json
// .vscode/settings.json
{
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "editor.formatOnSave": true,
    "files.autoSave": "afterDelay"
}
```

#### Useful VS Code Extensions
- **Go** (by Google) - Official Go support
- **Go Test Explorer** - Visual test runner
- **GitLens** - Advanced Git features
- **Error Lens** - Inline error display

### Development Tools

#### Install Development Tools
```bash
# Install Go tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install air for hot reloading (used by make dev)
go install github.com/cosmtrek/air@latest

# Install GitHub CLI (optional)
brew install gh  # macOS
sudo apt install gh  # Ubuntu
```

#### Hot Reloading Setup
```bash
# Start development mode
make dev

# This will:
# 1. Watch for file changes
# 2. Automatically rebuild the application
# 3. Restart the program when changes are detected
```

### Debugging

#### Using Go's Built-in Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug ./cmd/main.go
```

#### Debug with VS Code
1. Set breakpoints in your code
2. Press F5 or use "Run and Debug"
3. Choose "Launch Package" configuration

#### Debug Menu Items
```go
// Add debug prints to understand flow
func handleBranchOperations() error {
    fmt.Printf("DEBUG: Entering branch operations menu\n")
    
    // Your code here...
    
    fmt.Printf("DEBUG: Branch operations completed\n")
    return nil
}
```

## 🧪 Testing and Quality Assurance

### Running Tests

```bash
# Run all tests
make test

# Run tests with detailed output
go test -v ./...

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test -v ./internal/git/

# Run a specific test
go test -run TestGetRepositoryInfo ./internal/git/
```

### Writing Tests

#### Unit Test Example
```go
// internal/git/commands_test.go
func TestGetCurrentBranch(t *testing.T) {
    tests := []struct {
        name          string
        gitOutput     string
        gitError      error
        expectedBranch string
        expectError   bool
    }{
        {
            name:          "main branch",
            gitOutput:     "main",
            gitError:      nil,
            expectedBranch: "main",
            expectError:   false,
        },
        {
            name:          "feature branch",
            gitOutput:     "feature/user-auth",
            gitError:      nil,
            expectedBranch: "feature/user-auth",
            expectError:   false,
        },
        {
            name:          "git command fails",
            gitOutput:     "",
            gitError:      errors.New("not a git repository"),
            expectedBranch: "",
            expectError:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock the git command
            originalRunCommand := runCommand
            runCommand = func(cmd string) (string, error) {
                return tt.gitOutput, tt.gitError
            }
            defer func() { runCommand = originalRunCommand }()

            branch, err := GetCurrentBranch()
            
            if tt.expectError && err == nil {
                t.Errorf("Expected error, but got none")
            }
            
            if !tt.expectError && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            
            if branch != tt.expectedBranch {
                t.Errorf("Expected branch %q, got %q", tt.expectedBranch, branch)
            }
        })
    }
}
```

### Testing Guidelines

1. **Test Happy Path and Error Cases**: Always test both success and failure scenarios
2. **Use Table-Driven Tests**: For multiple similar test cases
3. **Mock External Dependencies**: Mock Git commands, API calls, file system operations
4. **Test Edge Cases**: Empty repositories, no internet connection, invalid inputs
5. **Keep Tests Fast**: Unit tests should run quickly

### Manual Testing Scenarios

Before submitting a PR, test these scenarios:

#### Repository States
- [ ] Fresh Git repository (no commits)
- [ ] Repository with staged changes
- [ ] Repository with unstaged changes
- [ ] Repository with no changes (clean)
- [ ] Repository with merge conflicts
- [ ] Non-Git directory

#### Git Operations
- [ ] Creating and switching branches
- [ ] Committing changes
- [ ] Pushing and pulling
- [ ] Viewing logs and diffs
- [ ] Stash operations
- [ ] Tag operations

#### GitHub Features
- [ ] Repository information display
- [ ] Creating pull requests
- [ ] Listing issues
- [ ] Authentication (valid and invalid tokens)

## 🔧 Common Tasks

### Adding a New Git Command

1. **Add the function to `internal/git/commands.go`**:
```go
// GetRemotes returns a list of Git remotes
func GetRemotes() ([]Remote, error) {
    output, err := RunCommand("git remote -v")
    if err != nil {
        return nil, fmt.Errorf("failed to get remotes: %w", err)
    }
    
    // Parse the output and return remotes
    return parseRemotes(output), nil
}
```

2. **Add tests in `internal/git/commands_test.go`**:
```go
func TestGetRemotes(t *testing.T) {
    // Test implementation...
}
```

3. **Add menu item in `internal/cli/menu.go`**:
```go
{
    Label:       "🌐 View Remotes",
    Description: "Show configured Git remotes",
    Handler:     handleRemotes,
    Condition:   git.IsGitRepository,
}
```

4. **Implement the handler**:
```go
func handleRemotes() error {
    remotes, err := git.GetRemotes()
    if err != nil {
        return err
    }
    
    for _, remote := range remotes {
        fmt.Printf("%s %s (%s)\n", 
            ui.IconRemote, 
            remote.Name, 
            remote.URL)
    }
    
    return nil
}
```

### Adding a New GitHub Feature

1. **Extend the GitHub client in `internal/github/client.go`**
2. **Add the API call method**
3. **Add menu integration in `internal/cli/menu.go`**
4. **Write tests for the new functionality**

### Improving UI and Styling

1. **Add new styles in `internal/ui/styles.go`**:
```go
// FormatRemote formats remote information
func FormatRemote(name, url string) string {
    return fmt.Sprintf("%s %s → %s", 
        IconRemote, 
        ColorBold.Render(name), 
        ColorURL.Render(url))
}
```

2. **Use the new styles in your handlers**

### Adding Configuration Options

1. **Add fields to config struct in `internal/config/types.go`**
2. **Update config loading/saving in `internal/config/config.go`**
3. **Add UI for configuration in settings menu**

## 🔧 Troubleshooting

### Common Issues

#### "Command not found" errors
```bash
# Make sure Go bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Or check where Go installs binaries
go env GOPATH
```

#### Tests failing
```bash
# Clean and rebuild
make clean
make deps
make build
make test
```

#### Git operations failing
```bash
# Check if you're in a Git repository
git status

# Verify Git is installed and working
git --version
```

#### Import paths not working
```bash
# Make sure you're using the correct module path
go mod tidy
```

### Debug Mode

Enable debug output:
```bash
# Set debug environment variable
export GITHUBBER_DEBUG=1
./build/githubber
```

### Getting Help

If you're stuck:

1. **Check existing issues** on GitHub
2. **Search the codebase** for similar patterns
3. **Ask questions** in GitHub Discussions
4. **Create an issue** with the `question` label
5. **Reach out** to maintainers

## 🎯 Next Steps

### After Your First Contribution

1. **Explore Advanced Features**:
   - Work on GitHub API integration
   - Add new Git operations
   - Improve the plugin system

2. **Help Other Contributors**:
   - Review pull requests
   - Answer questions in issues
   - Improve documentation

3. **Suggest New Features**:
   - Share your ideas in GitHub Discussions
   - Create feature request issues
   - Propose architectural improvements

### Becoming a Regular Contributor

- **Join our community** discussions
- **Take on larger features** and improvements
- **Help with project maintenance** and code reviews
- **Mentor new contributors** joining the project

### Learning Opportunities

Contributing to GitHubber helps you learn:

- **Advanced Go patterns** and best practices
- **CLI application architecture** and design
- **Git internals** and advanced Git operations
- **API integration** and error handling
- **Open source collaboration** and project management

## 📚 Additional Resources

### Go Language Resources
- [Official Go Tutorial](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

### Git and GitHub
- [Pro Git Book](https://git-scm.com/book)
- [GitHub Docs](https://docs.github.com/)
- [Git Workflows](https://www.atlassian.com/git/tutorials/comparing-workflows)

### CLI Development
- [Charm Bracelet Libraries](https://charm.sh/)
- [CLI Guidelines](https://clig.dev/)

### Open Source
- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [First Contributions](https://github.com/firstcontributions/first-contributions)

---

**Welcome to the GitHubber community! We're excited to see what you'll build! 🚀**

*Need help? Don't hesitate to ask questions in issues or reach out to the maintainers. We're here to help you succeed!*