# GitHubber 🚀

A comprehensive, production-ready Git CLI tool written in Go that provides advanced Git functionality with an intuitive interface. GitHubber covers everything from basic repository operations to complex workflows including interactive rebasing, advanced stashing, sophisticated tagging, and comprehensive remote management.

[![Go Version](https://img.shields.io/badge/Go-1.22.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

## ✨ Features

### 🔧 Repository Operations
- **Initialize** - Create new Git repositories with advanced options
- **Clone** - Clone repositories with depth control, branch selection, and submodule support
- **Repository Info** - Get detailed repository information and validation

### 🌿 Branch Management
- **Create/Delete** - Full branch lifecycle management
- **Switch/Checkout** - Navigate between branches efficiently
- **List** - View all branches with filtering and sorting options
- **Tracking** - Set up and manage upstream tracking relationships
- **Comparison** - Analyze ahead/behind status and commit differences

### 💾 Staging & Commits
- **Smart Staging** - Add files with pattern matching and interactive selection
- **Flexible Commits** - Support for standard commits, amendments, and sign-offs
- **History Management** - Comprehensive commit history with advanced filtering
- **Reset Operations** - Soft, mixed, and hard reset capabilities
- **Cherry-picking & Reverting** - Advanced commit manipulation

### 🔄 Remote Operations
- **Remote Management** - Add, remove, rename, and configure remotes
- **Push/Pull** - Advanced push and pull operations with force options
- **Fetch** - Granular fetching with pruning and tag management
- **Tracking Branches** - Comprehensive remote branch tracking

### 📦 Stash Operations
- **Save/Apply** - Flexible stashing with message support and path selection
- **List/Show** - Browse and inspect stash entries
- **Branch Creation** - Create branches from stash entries
- **Advanced Options** - Include untracked files, keep index, and selective stashing

### 🏷️ Tag Management
- **Lightweight Tags** - Simple reference tags
- **Annotated Tags** - Rich tags with messages and signing
- **Tag Operations** - Push, fetch, and delete tags
- **Verification** - GPG signature verification for signed tags

### 🔄 Advanced Workflows
- **Interactive Rebase** - Sophisticated commit rewriting and squashing
- **Merge Strategies** - Multiple merge strategies and conflict resolution
- **Conflict Resolution** - Tools for handling merge and rebase conflicts
- **Workflow Automation** - Automated common Git workflows

## 📁 Project Structure

```
githubber/
├── cmd/                    # Application entry points
│   └── main.go
├── pkg/                    # Public packages
│   ├── git/               # Core Git functionality
│   │   ├── client.go      # Main Git client
│   │   ├── repository.go  # Repository operations
│   │   ├── branches.go    # Branch management
│   │   ├── commits.go     # Commit operations
│   │   ├── staging.go     # Staging operations
│   │   ├── remotes.go     # Remote management
│   │   ├── stash.go       # Stash operations
│   │   ├── tags.go        # Tag management
│   │   ├── rebase.go      # Rebase operations
│   │   └── executor.go    # Command execution
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── logger/            # Logging utilities
│   │   └── logger.go
│   ├── types/             # Type definitions
│   │   └── git.go
│   └── utils/             # Utility functions
│       └── parsers.go
├── internal/              # Private packages
│   ├── commands/          # CLI commands
│   │   ├── root.go
│   │   ├── init.go
│   │   ├── clone.go
│   │   └── status.go
│   ├── handlers/          # Command handlers
│   ├── interactive/       # Interactive UI components
│   └── middleware/        # Command middleware
├── tests/                 # Test files
│   ├── unit/             # Unit tests
│   ├── integration/      # Integration tests
│   └── fixtures/         # Test fixtures
├── configs/              # Configuration files
│   └── default.yaml
├── docs/                 # Documentation
├── scripts/              # Build and utility scripts
├── examples/             # Usage examples
├── Makefile             # Build automation
├── .golangci.yml        # Linter configuration
└── go.mod               # Go module definition
```

## 🚀 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ritankarsaha/githubber.git
cd githubber

# Build the binary
make build

# Install to system PATH
make install-system
```

### Using Go Install

```bash
go install github.com/ritankarsaha/githubber/cmd@latest
```

### Pre-built Binaries

Download the latest release from the [releases page](https://github.com/ritankarsaha/githubber/releases).

## 🎯 Usage

### Basic Commands

```bash
# Initialize a new repository
githubber init [directory]

# Clone a repository
githubber clone <repository-url> [directory]

# Check repository status
githubber status

# View detailed help
githubber --help
```

### Advanced Examples

```bash
# Clone with specific options
githubber clone --depth 10 --branch main https://github.com/user/repo.git

# Interactive status with porcelain output
githubber status --short

# Initialize bare repository
githubber init --bare /path/to/repo.git
```

## ⚙️ Configuration

GitHubber uses a configuration file located at `~/.githubber/config.yaml`. You can customize various settings:

```yaml
app:
  name: "GitHubber"
  version: "2.0.0"
  debug: false

git:
  default_remote: "origin"
  default_branch: "main"
  auto_push: false
  sign_commits: false
  pretty_format: "%h %s (%an, %ar)"

github:
  token: ""  # GitHub API token for enhanced features
  api_base_url: "https://api.github.com"

ui:
  theme: "default"
  show_icons: true
  show_spinner: true
  confirm_actions: true
  page_size: 10

logging:
  level: "info"
  format: "text"
  output_file: ""  # Leave empty for stdout
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Run all checks
make check
```

## 🔧 Development

### Prerequisites

- Go 1.22.4 or later
- Git 2.20 or later
- Make (for build automation)

### Setup Development Environment

```bash
# Clone and setup
git clone https://github.com/ritankarsaha/githubber.git
cd githubber

# Install development dependencies
make dev-setup

# Run in development mode with auto-reload
make watch
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create release packages
make release
```

## 📊 Comprehensive Git Operations Coverage

### Basic Operations ✅
- Repository initialization and cloning
- File staging and committing
- Branch creation and switching
- Basic push/pull operations

### Intermediate Operations ✅
- Remote management
- Stash operations
- Tag creation and management
- Merge operations
- Conflict resolution

### Advanced Operations ✅
- Interactive rebasing
- Commit squashing and rewriting
- Cherry-picking and reverting
- Advanced merge strategies
- GPG signing and verification
- Submodule management
- Worktree operations
- Bisect operations
- Reflog management
- Garbage collection
- Repository maintenance

### GitHub Integration 🚧
- Pull request management
- Issue tracking
- GitHub Actions integration
- Release management

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Git community for the amazing version control system
- Go community for the excellent tooling and libraries
- Contributors and testers who help make this project better

## 🐛 Bug Reports & Feature Requests

Please use the [GitHub Issues](https://github.com/ritankarsaha/githubber/issues) page to report bugs or request features.

## 📈 Roadmap

- [ ] GitHub CLI integration
- [ ] GitLab support
- [ ] Bitbucket integration
- [ ] Advanced conflict resolution UI
- [ ] Plugin system
- [ ] Web interface
- [ ] Team collaboration features
- [ ] Advanced analytics and reporting

---

**GitHubber** - Making Git operations intuitive and powerful! 🚀