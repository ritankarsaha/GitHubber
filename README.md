# 🚀 GitHubber - Advanced Git & GitHub CLI Tool

<div align="center">

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/ritankarsaha/GitHubber)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue)](https://golang.org/)
[![GitHub release](https://img.shields.io/badge/version-v2.0.0-green)](https://github.com/ritankarsaha/GitHubber/releases)

*A powerful, beautiful, and feature-rich command-line interface for Git and GitHub operations*

**Created by [Ritankar Saha](mailto:ritankar.saha786@gmail.com)**

</div>

## ✨ Features

GitHubber is a comprehensive CLI tool that supercharges your Git and GitHub workflow with:

### 🛠 **Core Git Operations**
- **Repository Management**: Initialize, clone, and manage repositories
- **Branch Operations**: Create, delete, switch, and list branches
- **Commit Management**: Stage files, create commits, and view history
- **Remote Operations**: Push, pull, and fetch from remote repositories
- **Advanced Features**: Interactive commit squashing, stash management, tag operations

### 🐙 **GitHub Integration**
- **Repository Information**: View detailed GitHub repository stats
- **Pull Request Management**: Create and manage pull requests directly from CLI
- **Issue Tracking**: List and view GitHub issues
- **Authentication**: Secure token-based authentication

### 🎨 **Beautiful Terminal UI**
- **Colored Output**: Syntax-highlighted, colorful terminal interface
- **Interactive Menus**: Easy-to-navigate menu-driven interface
- **Professional Styling**: Clean, modern design with emojis and icons
- **Customizable Themes**: Dark, light, and auto themes

### ⚙️ **Configuration & Settings**
- **User Preferences**: Customizable UI themes and preferences
- **Token Management**: Secure GitHub token storage
- **Default Settings**: Set default repositories and workflows

## 🚀 Quick Start

### Prerequisites
- Go 1.23 or higher
- Git installed and configured
- GitHub account (for GitHub features)

### Installation

#### Option 1: Using Makefile (Recommended)
```bash
# Clone the repository
git clone https://github.com/ritankarsaha/GitHubber.git
cd GitHubber

# Build the application
make build

# Install globally (requires sudo)
make install

# Verify installation
githubber --help
```

#### Option 2: Manual Build
```bash
# Clone the repository
git clone https://github.com/ritankarsaha/GitHubber.git
cd GitHubber

# Download dependencies
go mod download

# Build the application
go build -o githubber ./cmd/main.go

# Install globally (optional)
sudo mv githubber /usr/local/bin/
```

#### Option 3: Direct Go Install
```bash
go install github.com/ritankarsaha/git-tool/cmd/main.go@latest
```

#### Option 4: Using Install Script
```bash
# Download and run install script
curl -fsSL https://raw.githubusercontent.com/ritankarsaha/GitHubber/main/scripts/install.sh | bash
```

### Setup GitHub Authentication

To use GitHub features, you need to authenticate:

#### Method 1: Environment Variable
```bash
export GITHUB_TOKEN="your_github_personal_access_token"
```

#### Method 2: Configuration Menu
1. Run `githubber`
2. Select "Settings" from the menu
3. Choose "Set GitHub token"
4. Enter your personal access token

**How to create a GitHub Personal Access Token:**
1. Go to GitHub → Settings → Developer settings → Personal access tokens
2. Click "Generate new token (classic)"
3. Select scopes: `repo`, `read:user`, `read:org`
4. Copy the generated token

## 📖 Usage

### Basic Usage
```bash
# Navigate to any Git repository
cd /path/to/your/repository

# Launch GitHubber
githubber
```

### Menu Overview

#### 📂 Repository Operations
- **Initialize Repository**: Create a new Git repository
- **Clone Repository**: Clone a repository from URL

#### 🌿 Branch Operations
- **Create Branch**: Create and switch to a new branch
- **Delete Branch**: Delete local branches
- **Switch Branch**: Switch between existing branches
- **List Branches**: View all available branches

#### 💾 Changes and Staging
- **View Status**: Check repository status
- **Add Files**: Stage files for commit
- **Commit Changes**: Create commits with messages

#### 🔄 Remote Operations
- **Push Changes**: Push commits to remote repository
- **Pull Changes**: Pull updates from remote
- **Fetch Updates**: Fetch without merging

#### 📜 History and Diff
- **View Log**: Display commit history
- **View Diff**: Show file differences
- **Squash Commits**: Interactive commit squashing

#### 📦 Stash Operations
- **Stash Save**: Save current changes to stash
- **Stash Pop**: Apply stashed changes
- **List Stashes**: View all stashes

#### 🏷️ Tag Operations
- **Create Tag**: Create annotated tags
- **Delete Tag**: Remove tags
- **List Tags**: View all tags

#### 🐙 GitHub Operations
- **View Repository Info**: Display GitHub repository statistics
- **Create Pull Request**: Create PRs directly from CLI
- **List Issues**: View repository issues

#### ⚙️ Settings
- **View Settings**: Display current configuration
- **GitHub Authentication**: Manage GitHub tokens
- **UI Preferences**: Customize themes and display options

## 🏗 Project Structure

```
GitHubber/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── cli/                    # CLI interface components
│   │   ├── input.go           # User input handling
│   │   └── menu.go            # Menu system and handlers
│   ├── git/                   # Git operations
│   │   ├── commands.go        # Core Git commands
│   │   ├── squash.go          # Commit squashing functionality
│   │   ├── utils.go           # Git utilities
│   │   └── test_helpers.go    # Testing utilities
│   ├── github/                # GitHub API integration
│   │   └── client.go          # GitHub API client
│   ├── config/                # Configuration management
│   │   └── config.go          # Settings and preferences
│   └── ui/                    # Terminal UI styling
│       └── styles.go          # Styling and themes
├── tests/                     # Test files and fixtures
├── docs/                      # Documentation
├── scripts/                   # Utility scripts
├── examples/                  # Usage examples
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── README.md                  # This file
```

## 🎨 Customization

### Themes
GitHubber supports multiple color themes:
- **Dark Theme** (default): Optimized for dark terminals
- **Light Theme**: For light terminal backgrounds
- **Auto Theme**: Automatically detects terminal theme

### Configuration File
Configuration is stored in `~/.githubber/githubber.json`:

```json
{
  "github": {
    "default_owner": "your-username",
    "default_repo": "your-repo",
    "api_base_url": "https://api.github.com"
  },
  "ui": {
    "theme": "dark",
    "show_emojis": true,
    "page_size": 20,
    "border_style": "rounded"
  },
  "git": {
    "default_branch": "main",
    "auto_push": false,
    "sign_commits": false
  }
}
```

## 🔧 Advanced Features

### Interactive Commit Squashing
GitHubber provides an intuitive interface for squashing commits:
1. Ensures working directory is clean
2. Shows recent commits with hashes and messages
3. Prompts for base commit and new commit message
4. Handles the rebase automatically

### GitHub Integration
- Automatically detects repository from Git remote
- Parses both HTTPS and SSH repository URLs
- Provides detailed repository statistics
- Creates pull requests with current branch

### Error Handling
- Comprehensive error messages with helpful suggestions
- Graceful handling of authentication issues
- Validation of Git repository state before operations

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific package tests
go test ./internal/git/
```

## 🤝 Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on how to contribute to this project.

### Quick Start for Contributors
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run tests: `make test`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development Commands
```bash
make build         # Build the application
make test          # Run all tests
make test-coverage # Run tests with coverage
make lint          # Lint the code
make fmt           # Format the code
make dev           # Run in development mode
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Ritankar Saha**
- Email: [ritankar.saha786@gmail.com](mailto:ritankar.saha786@gmail.com)
- GitHub: [@ritankarsaha](https://github.com/ritankarsaha)

## 🙏 Acknowledgments

- [Charm Bracelet](https://charm.sh/) for the amazing Lipgloss library
- [Google](https://github.com/google/go-github) for the go-github library
- The Go community for excellent tooling and libraries

## 📊 Project Status

GitHubber is actively maintained and under continuous development. Current version: **v2.0.0**

### Recent Updates
- ✅ Complete UI overhaul with beautiful styling
- ✅ GitHub API integration
- ✅ Configuration system
- ✅ Enhanced error handling
- 🔄 GitHub Actions integration (coming soon)
- 🔄 Plugin system (planned)

---

<div align="center">

**⭐ If you find GitHubber useful, please give it a star on GitHub! ⭐**

[Report Bug](https://github.com/ritankarsaha/GitHubber/issues) · [Request Feature](https://github.com/ritankarsaha/GitHubber/issues) · [Documentation](https://github.com/ritankarsaha/GitHubber/wiki)

</div>