/*
 * GitHubber - Shell Completion Generator
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Generate shell completion scripts for bash, zsh, and fish
 */

package cli

import (
	"fmt"
	"strings"
)

// GenerateCompletion generates shell completion scripts
func GenerateCompletion(shell string) (string, error) {
	switch strings.ToLower(shell) {
	case "bash":
		return generateBashCompletion(), nil
	case "zsh":
		return generateZshCompletion(), nil
	case "fish":
		return generateFishCompletion(), nil
	default:
		return "", fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish)", shell)
	}
}

// generateBashCompletion generates bash completion script
func generateBashCompletion() string {
	commands := GetCommands()
	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}

	return fmt.Sprintf(`#!/bin/bash
# GitHubber bash completion script

_githubber_complete() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # Main commands
    local main_commands="%s"
    
    # Branch operations
    local branch_actions="create delete list switch"
    
    # Stash operations
    local stash_actions="push pop list show drop"
    
    # Tag operations  
    local tag_actions="create delete list"
    
    # GitHub operations
    local github_actions="repo pr issue"
    
    # PR operations
    local pr_actions="create list view close merge"
    
    # Issue operations
    local issue_actions="create list view close"
    
    # Reset options
    local reset_options="--soft --mixed --hard"
    
    # Merge options
    local merge_options="--no-ff --squash"

    case "${prev}" in
        githubber)
            COMPREPLY=( $(compgen -W "${main_commands}" -- ${cur}) )
            return 0
            ;;
        branch)
            COMPREPLY=( $(compgen -W "${branch_actions}" -- ${cur}) )
            return 0
            ;;
        stash)
            COMPREPLY=( $(compgen -W "${stash_actions}" -- ${cur}) )
            return 0
            ;;
        tag)
            COMPREPLY=( $(compgen -W "${tag_actions}" -- ${cur}) )
            return 0
            ;;
        github)
            COMPREPLY=( $(compgen -W "${github_actions}" -- ${cur}) )
            return 0
            ;;
        pr)
            COMPREPLY=( $(compgen -W "${pr_actions}" -- ${cur}) )
            return 0
            ;;
        issue)
            COMPREPLY=( $(compgen -W "${issue_actions}" -- ${cur}) )
            return 0
            ;;
        reset)
            COMPREPLY=( $(compgen -W "${reset_options}" -- ${cur}) )
            return 0
            ;;
        merge)
            COMPREPLY=( $(compgen -W "${merge_options}" -- ${cur}) )
            return 0
            ;;
        checkout|switch)
            # Complete with branch names
            local branches=$(git branch 2>/dev/null | sed 's/^..//; s/ *$//')
            COMPREPLY=( $(compgen -W "${branches}" -- ${cur}) )
            return 0
            ;;
        add)
            # Complete with modified files
            local files=$(git status --porcelain 2>/dev/null | awk '{print $2}')
            COMPREPLY=( $(compgen -W "${files}" -- ${cur}) )
            return 0
            ;;
        diff)
            # Complete with modified files  
            local files=$(git status --porcelain 2>/dev/null | awk '{print $2}')
            COMPREPLY=( $(compgen -f -W "${files}" -- ${cur}) )
            return 0
            ;;
        *)
            COMPREPLY=( $(compgen -W "${main_commands}" -- ${cur}) )
            return 0
            ;;
    esac
}

complete -F _githubber_complete githubber

# Installation instructions:
# 1. Save this script to a file (e.g., githubber-completion.bash)
# 2. Source it in your ~/.bashrc: source /path/to/githubber-completion.bash
# 3. Or copy it to /etc/bash_completion.d/ (requires sudo)
`, strings.Join(commandNames, " "))
}

// generateZshCompletion generates zsh completion script
func generateZshCompletion() string {
	commands := GetCommands()
	var commandDescriptions []string
	for name, cmd := range commands {
		commandDescriptions = append(commandDescriptions, fmt.Sprintf("    '%s:%s'", name, cmd.Description))
	}

	return fmt.Sprintf(`#compdef githubber
# GitHubber zsh completion script

_githubber() {
    local context state line
    typeset -A opt_args

    _arguments -C \
        '1: :_githubber_commands' \
        '*::arg:->args'

    case $line[1] in
        branch)
            _githubber_branch
            ;;
        stash)
            _githubber_stash
            ;;
        tag)
            _githubber_tag
            ;;
        github)
            _githubber_github
            ;;
        pr)
            _githubber_pr
            ;;
        issue)
            _githubber_issue
            ;;
        checkout|switch)
            _githubber_branches
            ;;
        add)
            _githubber_modified_files
            ;;
        diff)
            _githubber_files
            ;;
        reset)
            _githubber_reset
            ;;
        merge)
            _githubber_merge
            ;;
    esac
}

_githubber_commands() {
    local commands; commands=(
%s
    )
    _describe 'commands' commands
}

_githubber_branch() {
    local actions; actions=(
        'create:Create a new branch'
        'delete:Delete a branch'
        'list:List all branches'
        'switch:Switch to a branch'
    )
    _describe 'branch actions' actions
}

_githubber_stash() {
    local actions; actions=(
        'push:Stash current changes'
        'pop:Apply and remove latest stash'
        'list:List all stashes'
        'show:Show stash content'
        'drop:Delete a stash'
    )
    _describe 'stash actions' actions
}

_githubber_tag() {
    local actions; actions=(
        'create:Create a new tag'
        'delete:Delete a tag'
        'list:List all tags'
    )
    _describe 'tag actions' actions
}

_githubber_github() {
    local actions; actions=(
        'repo:Repository operations'
        'pr:Pull request operations'
        'issue:Issue operations'
    )
    _describe 'github actions' actions
}

_githubber_pr() {
    local actions; actions=(
        'create:Create a pull request'
        'list:List pull requests'
        'view:View pull request details'
        'close:Close a pull request'
        'merge:Merge a pull request'
    )
    _describe 'pr actions' actions
}

_githubber_issue() {
    local actions; actions=(
        'create:Create an issue'
        'list:List issues'
        'view:View issue details'
        'close:Close an issue'
    )
    _describe 'issue actions' actions
}

_githubber_branches() {
    local branches
    branches=(${(f)"$(git branch 2>/dev/null | sed 's/^..//')"})
    _describe 'branches' branches
}

_githubber_modified_files() {
    local files
    files=(${(f)"$(git status --porcelain 2>/dev/null | awk '{print $2}')"})
    _describe 'modified files' files
}

_githubber_files() {
    _files
}

_githubber_reset() {
    local options; options=(
        '--soft:Reset HEAD only'
        '--mixed:Reset HEAD and index (default)'
        '--hard:Reset HEAD, index and working tree'
    )
    _describe 'reset options' options
}

_githubber_merge() {
    local options; options=(
        '--no-ff:Create merge commit even for fast-forward'
        '--squash:Squash commits into single commit'
    )
    _describe 'merge options' options
}

_githubber "$@"

# Installation instructions:
# 1. Save this script to a file in your fpath (e.g., _githubber)
# 2. Make sure the directory is in your fpath
# 3. Run: compinit
`, strings.Join(commandDescriptions, "\n"))
}

// generateFishCompletion generates fish completion script
func generateFishCompletion() string {
	commands := GetCommands()
	var completions []string

	for name, cmd := range commands {
		completions = append(completions,
			fmt.Sprintf("complete -c githubber -f -n '__fish_use_subcommand' -a '%s' -d '%s'",
				name, strings.ReplaceAll(cmd.Description, "'", "\\'")))
	}

	return fmt.Sprintf(`# GitHubber fish completion script

# Main commands
%s

# Branch operations
complete -c githubber -f -n '__fish_seen_subcommand_from branch' -a 'create delete list switch' -d 'Branch operations'

# Stash operations  
complete -c githubber -f -n '__fish_seen_subcommand_from stash' -a 'push pop list show drop' -d 'Stash operations'

# Tag operations
complete -c githubber -f -n '__fish_seen_subcommand_from tag' -a 'create delete list' -d 'Tag operations'

# GitHub operations
complete -c githubber -f -n '__fish_seen_subcommand_from github' -a 'repo pr issue' -d 'GitHub operations'

# PR operations
complete -c githubber -f -n '__fish_seen_subcommand_from pr' -a 'create list view close merge' -d 'Pull request operations'

# Issue operations
complete -c githubber -f -n '__fish_seen_subcommand_from issue' -a 'create list view close' -d 'Issue operations'

# Reset options
complete -c githubber -f -n '__fish_seen_subcommand_from reset' -l soft -d 'Reset HEAD only'
complete -c githubber -f -n '__fish_seen_subcommand_from reset' -l mixed -d 'Reset HEAD and index'
complete -c githubber -f -n '__fish_seen_subcommand_from reset' -l hard -d 'Reset HEAD, index and working tree'

# Merge options
complete -c githubber -f -n '__fish_seen_subcommand_from merge' -l no-ff -d 'Create merge commit even for fast-forward'
complete -c githubber -f -n '__fish_seen_subcommand_from merge' -l squash -d 'Squash commits into single commit'

# Branch completion for checkout/switch
complete -c githubber -f -n '__fish_seen_subcommand_from checkout switch' -a '(__fish_git_branches)'

# File completion for add/diff
complete -c githubber -f -n '__fish_seen_subcommand_from add diff' -a '(__fish_git_modified_files)'

# Functions for git-aware completions
function __fish_git_branches
    git branch 2>/dev/null | string replace -r '^..' '' | string trim
end

function __fish_git_modified_files
    git status --porcelain 2>/dev/null | awk '{print $2}'
end

# Installation instructions:
# 1. Save this script to ~/.config/fish/completions/githubber.fish
# 2. Fish will automatically load it when you start a new shell session
`, strings.Join(completions, "\n"))
}

// ShowCompletionInstructions shows installation instructions for shell completion
func ShowCompletionInstructions(shell string) string {
	switch strings.ToLower(shell) {
	case "bash":
		return `
To install bash completion:

1. Generate and save the completion script:
   githubber completion bash > githubber-completion.bash

2. Source it in your ~/.bashrc:
   echo "source /path/to/githubber-completion.bash" >> ~/.bashrc

3. Or install system-wide (requires sudo):
   sudo cp githubber-completion.bash /etc/bash_completion.d/

4. Restart your shell or run: source ~/.bashrc
`
	case "zsh":
		return `
To install zsh completion:

1. Generate and save the completion script:
   githubber completion zsh > _githubber

2. Move it to a directory in your fpath:
   mkdir -p ~/.local/share/zsh/site-functions
   mv _githubber ~/.local/share/zsh/site-functions/

3. Add the directory to your fpath in ~/.zshrc:
   fpath=(~/.local/share/zsh/site-functions $fpath)

4. Regenerate completions: compinit

5. Restart your shell
`
	case "fish":
		return `
To install fish completion:

1. Generate and save the completion script:
   githubber completion fish > ~/.config/fish/completions/githubber.fish

2. Fish will automatically load it in new shell sessions

3. Or load it immediately: source ~/.config/fish/completions/githubber.fish
`
	default:
		return fmt.Sprintf("Unsupported shell: %s", shell)
	}
}

// GetAvailableShells returns list of supported shells for completion
func GetAvailableShells() []string {
	return []string{"bash", "zsh", "fish"}
}
