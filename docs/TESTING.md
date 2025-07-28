# 🧪 Testing Guide for GitHubber

This document provides comprehensive information about testing in the GitHubber project, including test coverage, running tests, and writing new tests.

## 📋 Table of Contents

- [Test Coverage Overview](#test-coverage-overview)
- [Running Tests](#running-tests)
- [Test Structure](#test-structure)
- [Writing Tests](#writing-tests)
- [Coverage Reports](#coverage-reports)
- [Continuous Integration](#continuous-integration)

## 📊 Test Coverage Overview

### Current Test Coverage by Package

| Package | Coverage | Status | Test Files |
|---------|----------|--------|------------|
| `internal/ui` | 92.3% | ✅ Excellent | `styles_test.go` |
| `internal/config` | 85%+ | ✅ Good | `config_test.go` |
| `internal/git` | 30.3% | ⚠️ Needs Work | `commands_test.go`, `squash_test.go`, `utils_test.go` |
| `internal/github` | 25%+ | ⚠️ Basic | `client_test.go` |
| `internal/cli` | 15%+ | ⚠️ Basic | `input_test.go`, `menu_test.go` |
| `cmd/` | 10%+ | ⚠️ Basic | `main_test.go` |
| `internal/plugins` | 5%+ | 🔄 Stub Tests | `types_test.go` |
| `internal/providers` | 5%+ | 🔄 Stub Tests | `registry_test.go` |
| `internal/logging` | 5%+ | 🔄 Stub Tests | `types_test.go` |

### Test Categories

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test interactions between components
- **Structure Tests**: Test data structures and types
- **Mock Tests**: Test with mocked dependencies

## 🚀 Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run detailed coverage analysis
make test-coverage-detailed

# Run tests for specific package
make test-package PKG=internal/ui

# Run benchmark tests
make test-bench

# Generate coverage badge
make test-badge
```

### Manual Test Commands

```bash
# Run tests for specific package
go test -v ./internal/ui/

# Run tests with coverage
go test -v -cover ./internal/config/

# Run tests with coverage profile
go test -v -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage by function
go tool cover -func=coverage.out
```

### Test Coverage Script

The project includes a custom test coverage script (`test_coverage.sh`) that:

1. Runs tests for all working packages
2. Generates individual coverage reports
3. Merges coverage data
4. Creates HTML coverage reports
5. Provides coverage summaries

```bash
./test_coverage.sh
```

## 🏗 Test Structure

### Directory Organization

```
GitHubber/
├── cmd/
│   └── main_test.go          # CLI application tests
├── internal/
│   ├── cli/
│   │   ├── input_test.go     # Input handling tests
│   │   └── menu_test.go      # Menu system tests
│   ├── config/
│   │   └── config_test.go    # Configuration tests
│   ├── git/
│   │   ├── commands_test.go  # Git command tests
│   │   ├── squash_test.go    # Commit squashing tests
│   │   ├── utils_test.go     # Git utility tests
│   │   └── test_helpers.go   # Test helper functions
│   ├── github/
│   │   └── client_test.go    # GitHub API tests
│   ├── ui/
│   │   └── styles_test.go    # UI styling tests
│   └── [other packages]/
│       └── *_test.go         # Package-specific tests
├── coverage/                 # Coverage reports
├── test_coverage.sh         # Test coverage script
└── docs/
    └── TESTING.md           # This document
```

### Test Naming Conventions

- Test files: `*_test.go`
- Test functions: `TestFunctionName(t *testing.T)`
- Benchmark functions: `BenchmarkFunctionName(b *testing.B)`
- Example functions: `ExampleFunctionName()`

## ✍️ Writing Tests

### Test Structure Template

```go
package packagename

import (
    "testing"
    // other imports
)

func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid case",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "error case",
            input:    invalidInput,
            expected: zeroValue,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionToTest(tt.input)
            
            if tt.wantErr && err == nil {
                t.Errorf("Expected error but got none")
            }
            
            if !tt.wantErr && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### Testing Patterns

#### 1. Table-Driven Tests
```go
func TestParseURL(t *testing.T) {
    tests := []struct {
        name      string
        url       string
        wantOwner string
        wantRepo  string
        wantErr   bool
    }{
        {"github https", "https://github.com/user/repo", "user", "repo", false},
        {"github ssh", "git@github.com:user/repo.git", "user", "repo", false},
        {"invalid url", "not-a-url", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            owner, repo, err := ParseURL(tt.url)
            // assertions...
        })
    }
}
```

#### 2. Mock Testing
```go
func TestWithMock(t *testing.T) {
    // Create mock
    mockClient := &MockClient{
        response: expectedResponse,
    }
    
    // Test with mock
    result := ServiceUnderTest(mockClient)
    
    // Verify expectations
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

#### 3. Temporary Directory Testing
```go
func TestFileOperations(t *testing.T) {
    tmpDir, err := os.MkdirTemp("", "test-*")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tmpDir)
    
    // Test file operations in tmpDir
}
```

#### 4. Git Repository Testing
```go
func TestGitOperations(t *testing.T) {
    // Use test_helpers.go functions
    tmpDir := setupTestRepo(t)
    defer cleanupTestRepo(tmpDir)
    
    // Test git operations
}
```

### Test Helpers

The project includes test helper functions in `internal/git/test_helpers.go`:

```go
// Setup test Git repository
func setupTestRepo(t *testing.T) string

// Cleanup test repository  
func cleanupTestRepo(dir string)

// Create test commits
func createTestCommit(t *testing.T, message string)

// Assert directory exists
func assertDirExists(t *testing.T, path string)

// Assert file exists
func assertFileExists(t *testing.T, path string)
```

## 📈 Coverage Reports

### Viewing Coverage

1. **HTML Report**: Open `coverage/coverage.html` in a browser
2. **Terminal Summary**: Use `go tool cover -func=coverage.out`
3. **VS Code Extension**: Use Go extension's coverage features

### Coverage Targets

| Component | Target | Current | Priority |
|-----------|--------|---------|----------|
| Core Git Operations | 80%+ | 30% | High |
| Configuration System | 90%+ | 85% | Medium |
| UI Components | 85%+ | 92% | ✅ Done |
| GitHub Integration | 70%+ | 25% | High |
| CLI Interface | 60%+ | 15% | Medium |

### Improving Coverage

1. **Identify Gaps**: Use coverage reports to find untested code
2. **Add Tests**: Write tests for uncovered functions
3. **Edge Cases**: Test error conditions and edge cases
4. **Integration**: Add integration tests for component interactions

## 🔄 Continuous Integration

### GitHub Actions

The project can be configured with GitHub Actions for automated testing:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.23
    - run: make test-coverage
    - uses: actions/upload-artifact@v2
      with:
        name: coverage-report
        path: coverage/
```

### Pre-commit Hooks

Set up pre-commit hooks to run tests before commits:

```bash
#!/bin/bash
# .git/hooks/pre-commit
make test-coverage
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

## 🔧 Testing Best Practices

### 1. Test Organization
- Group related tests in the same file
- Use descriptive test names
- Keep tests focused and independent

### 2. Test Data
- Use table-driven tests for multiple scenarios  
- Create test fixtures for complex data
- Use temporary directories for file operations

### 3. Error Testing
- Test both success and error paths
- Verify error messages are helpful
- Test edge cases and boundary conditions

### 4. Performance Testing
- Use benchmark tests for critical paths
- Monitor test execution time
- Profile memory usage for large operations

### 5. Maintainability
- Keep tests simple and readable
- Avoid testing implementation details
- Update tests when changing functionality

## 🐛 Debugging Tests

### Common Issues

1. **Failing Git Tests**: Ensure Git is configured properly
2. **File Path Issues**: Use absolute paths in tests
3. **Network Tests**: Mock external API calls
4. **Race Conditions**: Use proper synchronization

### Debugging Commands

```bash
# Run specific test with verbose output
go test -v -run TestSpecificFunction ./internal/package/

# Run tests with race detection
go test -race ./...

# Run tests with memory profiling
go test -memprofile=mem.prof ./...

# Run tests with CPU profiling
go test -cpuprofile=cpu.prof ./...
```

## 📚 Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TestingTesting)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Advanced Go Testing](https://segment.com/blog/5-advanced-testing-techniques-in-go/)

## 🎯 Next Steps

1. **Increase Git Package Coverage**: Add more comprehensive Git operation tests
2. **GitHub API Testing**: Implement proper mocking for API tests  
3. **Integration Tests**: Add end-to-end testing scenarios
4. **Performance Testing**: Add benchmark tests for critical operations
5. **CI/CD Integration**: Set up automated testing pipeline

---

**Happy Testing! 🚀**

*This testing guide is maintained alongside the codebase. Please update it when adding new test patterns or changing test structure.*