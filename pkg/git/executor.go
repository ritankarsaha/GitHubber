package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ritankarsaha/githubber/pkg/logger"
)

type CommandExecutor interface {
	Execute(command string, args ...string) (string, error)
	ExecuteWithContext(ctx context.Context, command string, args ...string) (string, error)
	ExecuteInDir(dir string, command string, args ...string) (string, error)
	SetWorkingDirectory(dir string)
	GetWorkingDirectory() string
}

type DefaultCommandExecutor struct {
	workingDir string
	timeout    time.Duration
}

func NewCommandExecutor() CommandExecutor {
	return &DefaultCommandExecutor{
		timeout: 30 * time.Second,
	}
}

func (e *DefaultCommandExecutor) Execute(command string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	
	return e.ExecuteWithContext(ctx, command, args...)
}

func (e *DefaultCommandExecutor) ExecuteWithContext(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	
	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}
	
	// Log the command being executed
	logger.Debugf("Executing command: %s %s", command, strings.Join(args, " "))
	if cmd.Dir != "" {
		logger.Debugf("Working directory: %s", cmd.Dir)
	}
	
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))
	
	if err != nil {
		logger.Errorf("Command failed: %s %s", command, strings.Join(args, " "))
		logger.Errorf("Error: %v", err)
		logger.Errorf("Output: %s", outputStr)
		return outputStr, fmt.Errorf("command failed: %w, output: %s", err, outputStr)
	}
	
	logger.Debugf("Command output: %s", outputStr)
	return outputStr, nil
}

func (e *DefaultCommandExecutor) ExecuteInDir(dir string, command string, args ...string) (string, error) {
	originalDir := e.workingDir
	e.workingDir = dir
	defer func() {
		e.workingDir = originalDir
	}()
	
	return e.Execute(command, args...)
}

func (e *DefaultCommandExecutor) SetWorkingDirectory(dir string) {
	e.workingDir = dir
}

func (e *DefaultCommandExecutor) GetWorkingDirectory() string {
	if e.workingDir != "" {
		return e.workingDir
	}
	
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

// TestCommandExecutor for unit testing
type TestCommandExecutor struct {
	responses map[string]string
	errors    map[string]error
	executed  []string
}

func NewTestCommandExecutor() *TestCommandExecutor {
	return &TestCommandExecutor{
		responses: make(map[string]string),
		errors:    make(map[string]error),
		executed:  []string{},
	}
}

func (e *TestCommandExecutor) SetResponse(command string, response string) {
	e.responses[command] = response
}

func (e *TestCommandExecutor) SetError(command string, err error) {
	e.errors[command] = err
}

func (e *TestCommandExecutor) GetExecutedCommands() []string {
	return e.executed
}

func (e *TestCommandExecutor) Execute(command string, args ...string) (string, error) {
	return e.ExecuteWithContext(context.Background(), command, args...)
}

func (e *TestCommandExecutor) ExecuteWithContext(ctx context.Context, command string, args ...string) (string, error) {
	fullCommand := command + " " + strings.Join(args, " ")
	e.executed = append(e.executed, fullCommand)
	
	if err, exists := e.errors[fullCommand]; exists {
		return "", err
	}
	
	if response, exists := e.responses[fullCommand]; exists {
		return response, nil
	}
	
	return "", fmt.Errorf("no response configured for command: %s", fullCommand)
}

func (e *TestCommandExecutor) ExecuteInDir(dir string, command string, args ...string) (string, error) {
	return e.Execute(command, args...)
}

func (e *TestCommandExecutor) SetWorkingDirectory(dir string) {
	// No-op for test executor
}

func (e *TestCommandExecutor) GetWorkingDirectory() string {
	return "/test/dir"
}