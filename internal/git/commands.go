/*
 * GitHubber - Git Commands Module
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Core Git command implementations and wrappers
 */

package git

import (
	"fmt"
	"strings"
)

// Repository Operations
func Init() error {
	_, err := RunCommand("git init")
	return err
}

func Clone(url string) error {
	_, err := RunCommand(fmt.Sprintf("git clone %s", url))
	return err
}

// Branch Operations
func CreateBranch(name string) error {
	_, err := RunCommand(fmt.Sprintf("git checkout -b %s", name))
	return err
}

func DeleteBranch(name string) error {
	_, err := RunCommand(fmt.Sprintf("git branch -D %s", name))
	return err
}

func SwitchBranch(name string) error {
	_, err := RunCommand(fmt.Sprintf("git checkout %s", name))
	return err
}

func ListBranches() ([]string, error) {
	output, err := RunCommand("git branch")
	if err != nil {
		return nil, err
	}
	branches := strings.Split(output, "\n")
	return branches, nil
}

// Changes and Staging
func Status() (string, error) {
	return RunCommand("git status")
}

func AddFiles(files ...string) error {
	if len(files) == 0 {
		_, err := RunCommand("git add .")
		return err
	}
	_, err := RunCommand(fmt.Sprintf("git add %s", strings.Join(files, " ")))
	return err
}

func Commit(message string) error {
	_, err := RunCommand(fmt.Sprintf("git commit -m \"%s\"", message))
	return err
}

// Remote Operations
func Push(remote, branch string) error {
	_, err := RunCommand(fmt.Sprintf("git push %s %s", remote, branch))
	return err
}

func Pull(remote, branch string) error {
	_, err := RunCommand(fmt.Sprintf("git pull %s %s", remote, branch))
	return err
}

func Fetch(remote string) error {
	_, err := RunCommand(fmt.Sprintf("git fetch %s", remote))
	return err
}

// History and Diff
func Log(n int) (string, error) {
	return RunCommand(fmt.Sprintf("git log -%d --oneline", n))
}

func Diff(file string) (string, error) {
	return RunCommand(fmt.Sprintf("git diff %s", file))
}

// Stash Operations
func StashSave(message string) error {
	_, err := RunCommand(fmt.Sprintf("git stash push -m \"%s\"", message))
	return err
}

func StashPop() error {
	_, err := RunCommand("git stash pop")
	return err
}

func StashList() (string, error) {
	return RunCommand("git stash list")
}

// Tag Operations
func CreateTag(name, message string) error {
	_, err := RunCommand(fmt.Sprintf("git tag -a %s -m \"%s\"", name, message))
	return err
}

func DeleteTag(name string) error {
	_, err := RunCommand(fmt.Sprintf("git tag -d %s", name))
	return err
}

func ListTags() (string, error) {
	return RunCommand("git tag")
}
