# 🤝 Contributing to GitHubber

Thank you for your interest in contributing to GitHubber! This document provides comprehensive guidelines for contributors to help you get started and ensure a smooth contribution process.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Development Workflow](#development-workflow)
- [Release Process](#release-process)

## 📜 Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please treat all contributors with respect and create a welcoming environment for everyone.

## 🚀 Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.23 or higher**: [Download Go](https://golang.org/dl/)
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Make**: Usually pre-installed on Unix systems
- **GitHub CLI** (optional but recommended): [Install gh](https://cli.github.com/)

### Verify Installation

```bash
# Check Go version
go version

# Check Git version
git --version

# Check Make
make --version
```

## 🛠 Development Setup

### 1. Fork and Clone the Repository

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/GitHubber.git
cd GitHubber

# Add upstream remote
git remote add upstream https://github.com/ritankarsaha/GitHubber.git

# Verify remotes
git remote -v
```

### 2. Set Up Development Environment

```bash
# Download dependencies
make deps

# Build the application
make build

# Run tests to ensure everything works
make test

# Install development tools
make dev-setup  # This will install golangci-lint, air, etc.
```

### 3. Verify Setup

```bash
# Run the application
make run

# Run in development mode (with hot reload)
make dev

# Run linting
make lint

# Run tests with coverage
make test-coverage
```

## 🔄 Contributing Process

### 1. Choose What to Work On

- **Issues**: Look for issues labeled `good first issue`, `help wanted`, or `bug`
- **Features**: Check the roadmap in issues or propose new features
- **Documentation**: Improve existing docs or add missing documentation

### 2. Create a New Branch

```bash
# Update your local main branch
git checkout main
git pull upstream main

# Create a new feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/bug-description
```

### 3. Make Your Changes

- Write clean, maintainable code
- Follow the existing code style
- Add tests for new functionality
- Update documentation as needed

### 4. Test Your Changes

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linting
make lint

# Format your code
make fmt
```

### 5. Commit Your Changes

```bash
# Stage your changes
git add .

# Commit with a descriptive message
git commit -m "feat: add user authentication system"

# Or for bug fixes
git commit -m "fix: resolve nil pointer in git status command"
```

## 📝 Coding Standards

### Go Style Guide

We follow the standard Go conventions and best practices:

#### Code Formatting
- Use `gofmt` and `goimports` for formatting (available via `make fmt`)
- Use meaningful variable and function names
- Keep functions small and focused
- Add comments for exported functions and complex logic

#### Naming Conventions
```go
// ✅ Good
func GetRepositoryInfo() (*RepositoryInfo, error) {}
var githubToken string
const DefaultTimeout = 30 * time.Second

// ❌ Bad  
func get_repo_info() (*repoInfo, error) {}
var gt string
const timeout = 30
```

#### Error Handling
```go
// ✅ Good
if err != nil {
    return fmt.Errorf("failed to get repository info: %w", err)
}

// ❌ Bad
if err != nil {
    panic(err)
}
```

#### Package Structure
- Keep packages focused and cohesive
- Use internal packages for implementation details
- Export only what's necessary for the public API

### Directory Structure Guidelines

```
GitHubber/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── cli/              # Command-line interface
│   ├── git/              # Git operations
│   ├── github/           # GitHub API integration
│   ├── config/           # Configuration management
│   ├── ui/               # Terminal UI components
│   ├── logging/          # Logging infrastructure
│   ├── plugins/          # Plugin system
│   └── providers/        # External service providers
├── docs/                 # Documentation
├── scripts/              # Utility scripts
├── tests/               # Test files and fixtures
└── examples/            # Usage examples
```

## 🧪 Testing Guidelines

### Writing Tests

- Write tests for all new functionality
- Use table-driven tests for multiple test cases
- Test both success and error cases
- Use meaningful test names

```go
func TestGetRepositoryInfo(t *testing.T) {
    tests := []struct {
        name    string
        setup   func()
        want    *RepositoryInfo
        wantErr bool
    }{
        {
            name: "valid git repository",
            setup: func() {
                // Setup test repository
            },
            want:    &RepositoryInfo{URL: "...", Branch: "main"},
            wantErr: false,
        },
        {
            name:    "not a git repository",
            setup:   func() {},
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            got, err := GetRepositoryInfo()
            if (err != nil) != tt.wantErr {
                t.Errorf("GetRepositoryInfo() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetRepositoryInfo() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test ./internal/git/

# Run specific test
go test -run TestGetRepositoryInfo ./internal/git/
```

## 📥 Pull Request Process

### Before Submitting

1. **Rebase your branch** on the latest main:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**:
   ```bash
   make test
   make lint
   make fmt
   ```

3. **Update documentation** if needed

### PR Description Template

Use this template for your Pull Request description:

```markdown
## Description
Brief description of changes made

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Tests pass locally with `make test`
- [ ] Linting passes with `make lint`
- [ ] Added tests for new functionality
- [ ] Updated documentation

## Screenshots (if applicable)
Add screenshots here

## Additional Notes
Any additional information or context
```

### Review Process

1. **Automated Checks**: CI will run tests and linting
2. **Code Review**: Maintainers will review your code
3. **Address Feedback**: Make requested changes
4. **Merge**: Once approved, your PR will be merged

## 🐛 Issue Guidelines

### Reporting Bugs

When reporting bugs, please include:

1. **Clear title** describing the issue
2. **Steps to reproduce** the bug
3. **Expected behavior**
4. **Actual behavior**
5. **Environment information**:
   - Go version
   - OS and version
   - GitHubber version

### Bug Report Template

```markdown
**Bug Description**
A clear description of the bug

**Steps to Reproduce**
1. Run command `githubber ...`
2. Select option '...'
3. See error

**Expected Behavior**
What should have happened

**Actual Behavior**
What actually happened

**Environment**
- OS: [e.g., macOS 12.0]
- Go Version: [e.g., 1.23.0]
- GitHubber Version: [e.g., v2.0.0]

**Additional Context**
Any other relevant information
```

### Feature Requests

For feature requests, please include:

1. **Clear description** of the feature
2. **Use case** explaining why it's needed
3. **Proposed solution** (if you have ideas)
4. **Alternatives considered**

## 🔧 Development Workflow

### Daily Development

```bash
# Start development mode (with hot reload)
make dev

# In another terminal, run tests continuously
make watch-test  # If available

# Format code before committing
make fmt

# Check your code
make lint
```

### Common Commands

```bash
# Development
make dev           # Development mode with hot reload
make run           # Run the application once
make build         # Build binary
make clean         # Clean build artifacts

# Testing
make test          # Run all tests
make test-coverage # Run tests with coverage
make test-watch    # Watch for changes and re-run tests

# Code Quality
make lint          # Run linter
make fmt           # Format code
make vet           # Run go vet

# Dependencies
make deps          # Download dependencies
make deps-update   # Update dependencies

# Release
make build-all     # Cross-compile for all platforms
make release       # Create release archives
```

## 🚀 Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Workflow

1. **Update Version**: Update version in relevant files
2. **Update Changelog**: Add new features and fixes
3. **Create Tag**: `git tag -a v2.1.0 -m "Release v2.1.0"`
4. **Push Tag**: `git push upstream v2.1.0`
5. **GitHub Release**: Create release on GitHub
6. **Build Artifacts**: CI will build and attach binaries

## 💡 Tips for Contributors

### First-Time Contributors

- Start with issues labeled `good first issue`
- Read the existing code to understand patterns
- Ask questions in issues or discussions
- Don't hesitate to ask for help

### Code Review Tips

- Be respectful and constructive
- Explain the "why" behind your feedback
- Suggest improvements rather than just pointing out problems
- Test the changes locally when possible

### Git Best Practices

```bash
# Use descriptive commit messages
git commit -m "feat(cli): add interactive branch selection menu"
git commit -m "fix(git): handle empty repository case in status command"
git commit -m "docs: update installation instructions for macOS"

# Squash related commits before submitting PR
git rebase -i HEAD~3

# Keep your branch up to date
git fetch upstream
git rebase upstream/main
```

## 📞 Getting Help

If you need help or have questions:

1. **Check the documentation** in the `docs/` folder
2. **Search existing issues** for similar problems
3. **Create a new issue** with the `question` label
4. **Join discussions** on GitHub Discussions
5. **Contact the maintainer**: ritankar.saha786@gmail.com

## 🙏 Acknowledgments

Thank you for contributing to GitHubber! Your contributions help make this tool better for everyone in the Git and GitHub community.

---

**Happy Contributing! 🚀**