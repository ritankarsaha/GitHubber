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

// Advanced Git Operations

// Interactive Rebase
func InteractiveRebase(base string) error {
	_, err := RunCommand(fmt.Sprintf("git rebase -i %s", base))
	return err
}

func RebaseOnto(upstream, onto string) error {
	_, err := RunCommand(fmt.Sprintf("git rebase --onto %s %s", onto, upstream))
	return err
}

func RebaseContinue() error {
	_, err := RunCommand("git rebase --continue")
	return err
}

func RebaseAbort() error {
	_, err := RunCommand("git rebase --abort")
	return err
}

func RebaseSkip() error {
	_, err := RunCommand("git rebase --skip")
	return err
}

// Cherry Pick Operations
func CherryPick(commitHash string) error {
	_, err := RunCommand(fmt.Sprintf("git cherry-pick %s", commitHash))
	return err
}

func CherryPickRange(startCommit, endCommit string) error {
	_, err := RunCommand(fmt.Sprintf("git cherry-pick %s..%s", startCommit, endCommit))
	return err
}

func CherryPickContinue() error {
	_, err := RunCommand("git cherry-pick --continue")
	return err
}

func CherryPickAbort() error {
	_, err := RunCommand("git cherry-pick --abort")
	return err
}

// Reset Operations
func ResetSoft(commit string) error {
	_, err := RunCommand(fmt.Sprintf("git reset --soft %s", commit))
	return err
}

func ResetMixed(commit string) error {
	_, err := RunCommand(fmt.Sprintf("git reset --mixed %s", commit))
	return err
}

func ResetHard(commit string) error {
	_, err := RunCommand(fmt.Sprintf("git reset --hard %s", commit))
	return err
}

func ResetFile(file string) error {
	_, err := RunCommand(fmt.Sprintf("git reset HEAD %s", file))
	return err
}

// Revert Operations
func Revert(commitHash string) error {
	_, err := RunCommand(fmt.Sprintf("git revert %s", commitHash))
	return err
}

func RevertNoCommit(commitHash string) error {
	_, err := RunCommand(fmt.Sprintf("git revert --no-commit %s", commitHash))
	return err
}

// Merge Operations
func Merge(branch string) error {
	_, err := RunCommand(fmt.Sprintf("git merge %s", branch))
	return err
}

func MergeNoFF(branch string) error {
	_, err := RunCommand(fmt.Sprintf("git merge --no-ff %s", branch))
	return err
}

func MergeSquash(branch string) error {
	_, err := RunCommand(fmt.Sprintf("git merge --squash %s", branch))
	return err
}

func MergeAbort() error {
	_, err := RunCommand("git merge --abort")
	return err
}

func MergeContinue() error {
	_, err := RunCommand("git merge --continue")
	return err
}

// Bisect Operations
func BisectStart() error {
	_, err := RunCommand("git bisect start")
	return err
}

func BisectBad(commit string) error {
	if commit == "" {
		_, err := RunCommand("git bisect bad")
		return err
	}
	_, err := RunCommand(fmt.Sprintf("git bisect bad %s", commit))
	return err
}

func BisectGood(commit string) error {
	if commit == "" {
		_, err := RunCommand("git bisect good")
		return err
	}
	_, err := RunCommand(fmt.Sprintf("git bisect good %s", commit))
	return err
}

func BisectReset() error {
	_, err := RunCommand("git bisect reset")
	return err
}

func BisectSkip() error {
	_, err := RunCommand("git bisect skip")
	return err
}

// Remote Management
func AddRemote(name, url string) error {
	_, err := RunCommand(fmt.Sprintf("git remote add %s %s", name, url))
	return err
}

func RemoveRemote(name string) error {
	_, err := RunCommand(fmt.Sprintf("git remote remove %s", name))
	return err
}

func RenameRemote(oldName, newName string) error {
	_, err := RunCommand(fmt.Sprintf("git remote rename %s %s", oldName, newName))
	return err
}

func ListRemotes() (string, error) {
	return RunCommand("git remote -v")
}

func SetRemoteURL(name, url string) error {
	_, err := RunCommand(fmt.Sprintf("git remote set-url %s %s", name, url))
	return err
}

// Working Directory Operations
func CheckoutFile(file string) error {
	_, err := RunCommand(fmt.Sprintf("git checkout -- %s", file))
	return err
}

func CheckoutCommit(commit string) error {
	_, err := RunCommand(fmt.Sprintf("git checkout %s", commit))
	return err
}

func CheckoutNewBranch(branchName, startPoint string) error {
	if startPoint == "" {
		_, err := RunCommand(fmt.Sprintf("git checkout -b %s", branchName))
		return err
	}
	_, err := RunCommand(fmt.Sprintf("git checkout -b %s %s", branchName, startPoint))
	return err
}

// Clean Operations
func Clean() error {
	_, err := RunCommand("git clean -fd")
	return err
}

func CleanDryRun() (string, error) {
	return RunCommand("git clean -fd --dry-run")
}

// Configuration Operations
func SetConfig(key, value string, global bool) error {
	scope := "--local"
	if global {
		scope = "--global"
	}
	_, err := RunCommand(fmt.Sprintf("git config %s %s \"%s\"", scope, key, value))
	return err
}

func GetConfig(key string, global bool) (string, error) {
	scope := "--local"
	if global {
		scope = "--global"
	}
	return RunCommand(fmt.Sprintf("git config %s %s", scope, key))
}

// Submodule Operations
func AddSubmodule(url, path string) error {
	_, err := RunCommand(fmt.Sprintf("git submodule add %s %s", url, path))
	return err
}

func UpdateSubmodules() error {
	_, err := RunCommand("git submodule update --init --recursive")
	return err
}

func RemoveSubmodule(path string) error {
	// Remove submodule (requires multiple steps)
	commands := []string{
		fmt.Sprintf("git submodule deinit -f %s", path),
		fmt.Sprintf("rm -rf .git/modules/%s", path),
		fmt.Sprintf("git rm -f %s", path),
	}

	for _, cmd := range commands {
		if _, err := RunCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

// Archive Operations
func CreateArchive(format, output string) error {
	_, err := RunCommand(fmt.Sprintf("git archive --format=%s --output=%s HEAD", format, output))
	return err
}

// Show Operations
func ShowCommit(commit string) (string, error) {
	return RunCommand(fmt.Sprintf("git show %s", commit))
}

func ShowBranch() (string, error) {
	return RunCommand("git show-branch")
}

// Reflog Operations
func Reflog() (string, error) {
	return RunCommand("git reflog")
}

func ReflogExpire() error {
	_, err := RunCommand("git reflog expire --expire=now --all")
	return err
}

// Advanced History and Analysis
func InteractiveLog() (string, error) {
	return RunCommand("git log --oneline --graph --decorate --all")
}

func FileHistory(file string) (string, error) {
	return RunCommand(fmt.Sprintf("git log --follow --patch -- %s", file))
}

func BlameFile(file string) (string, error) {
	return RunCommand(fmt.Sprintf("git blame %s", file))
}

func ShowCommitDetails(commit string) (string, error) {
	return RunCommand(fmt.Sprintf("git show %s --stat", commit))
}

func CompareBranches(branch1, branch2 string) (string, error) {
	return RunCommand(fmt.Sprintf("git log --oneline %s..%s", branch1, branch2))
}

func FindCommitsByAuthor(author string) (string, error) {
	return RunCommand(fmt.Sprintf("git log --author=\"%s\" --oneline", author))
}

func FindCommitsByMessage(message string) (string, error) {
	return RunCommand(fmt.Sprintf("git log --grep=\"%s\" --oneline", message))
}

// Patch Operations
func CreatePatch(outputFile string) error {
	_, err := RunCommand(fmt.Sprintf("git diff > %s", outputFile))
	return err
}

func CreatePatchFromCommit(commit, outputFile string) error {
	_, err := RunCommand(fmt.Sprintf("git format-patch -1 %s --stdout > %s", commit, outputFile))
	return err
}

func ApplyPatch(patchFile string) error {
	_, err := RunCommand(fmt.Sprintf("git apply %s", patchFile))
	return err
}

func FormatPatchForEmail(since string) (string, error) {
	return RunCommand(fmt.Sprintf("git format-patch %s", since))
}

// Bundle Operations
func CreateBundle(bundleFile, refSpec string) error {
	_, err := RunCommand(fmt.Sprintf("git bundle create %s %s", bundleFile, refSpec))
	return err
}

func VerifyBundle(bundleFile string) (string, error) {
	return RunCommand(fmt.Sprintf("git bundle verify %s", bundleFile))
}

func ListBundleRefs(bundleFile string) (string, error) {
	return RunCommand(fmt.Sprintf("git bundle list-heads %s", bundleFile))
}

func CloneFromBundle(bundleFile, directory string) error {
	_, err := RunCommand(fmt.Sprintf("git clone %s %s", bundleFile, directory))
	return err
}

// Worktree Operations
func ListWorktrees() (string, error) {
	return RunCommand("git worktree list")
}

func AddWorktree(path, branch string) error {
	if branch == "" {
		_, err := RunCommand(fmt.Sprintf("git worktree add %s", path))
		return err
	}
	_, err := RunCommand(fmt.Sprintf("git worktree add %s %s", path, branch))
	return err
}

func RemoveWorktree(path string) error {
	_, err := RunCommand(fmt.Sprintf("git worktree remove %s", path))
	return err
}

func MoveWorktree(oldPath, newPath string) error {
	_, err := RunCommand(fmt.Sprintf("git worktree move %s %s", oldPath, newPath))
	return err
}

func PruneWorktrees() error {
	_, err := RunCommand("git worktree prune")
	return err
}

// Repository Maintenance
func GarbageCollect() error {
	_, err := RunCommand("git gc --aggressive")
	return err
}

func VerifyRepository() (string, error) {
	return RunCommand("git fsck --full")
}

func OptimizeRepository() error {
	commands := []string{
		"git gc --aggressive",
		"git repack -a -d --depth=250 --window=250",
		"git prune",
	}

	for _, cmd := range commands {
		if _, err := RunCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func RepositoryStatistics() (string, error) {
	stats, err := RunCommand("git count-objects -v")
	if err != nil {
		return "", err
	}

	size, _ := RunCommand("du -sh .git")
	branches, _ := RunCommand("git branch -a | wc -l")
	tags, _ := RunCommand("git tag | wc -l")
	commits, _ := RunCommand("git rev-list --count HEAD")

	result := fmt.Sprintf("Repository Statistics:\n%s\n\nRepository size: %s\nBranches: %s\nTags: %s\nTotal commits: %s",
		stats, strings.TrimSpace(size), strings.TrimSpace(branches), strings.TrimSpace(tags), strings.TrimSpace(commits))

	return result, nil
}

func PruneRemoteBranches(remote string) error {
	if remote == "" {
		remote = "origin"
	}
	_, err := RunCommand(fmt.Sprintf("git remote prune %s", remote))
	return err
}

func ReflogShow(ref string) (string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	return RunCommand(fmt.Sprintf("git reflog show %s", ref))
}

// Smart Git Operations
func InteractiveAdd() error {
	_, err := RunCommand("git add -p")
	return err
}

func PartialCommit(message string) error {
	// First do interactive add, then commit
	fmt.Println("Starting interactive add mode...")
	if err := InteractiveAdd(); err != nil {
		return err
	}
	return Commit(message)
}

func AmendLastCommit(message string) error {
	if message == "" {
		_, err := RunCommand("git commit --amend --no-edit")
		return err
	}
	_, err := RunCommand(fmt.Sprintf("git commit --amend -m \"%s\"", message))
	return err
}

func BranchComparison(branch1, branch2 string) (string, error) {
	ahead, err := RunCommand(fmt.Sprintf("git rev-list --count %s..%s", branch2, branch1))
	if err != nil {
		return "", err
	}

	behind, err := RunCommand(fmt.Sprintf("git rev-list --count %s..%s", branch1, branch2))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Branch Comparison:\n%s is %s commits ahead and %s commits behind %s",
		branch1, strings.TrimSpace(ahead), strings.TrimSpace(behind), branch2), nil
}

func ConflictPreventionCheck(branch string) (string, error) {
	// Check if merge would cause conflicts without actually merging
	output, err := RunCommand(fmt.Sprintf("git merge-tree $(git merge-base HEAD %s) HEAD %s", branch, branch))
	if err != nil {
		return "", err
	}

	if strings.Contains(output, "<<<<<<<") {
		return "⚠️  Merge conflicts detected! Files that would conflict:\n" + output, nil
	}

	return "✅ No merge conflicts detected. Safe to merge!", nil
}

// Sync Operations
func SyncWithAllRemotes() error {
	remotes, err := RunCommand("git remote")
	if err != nil {
		return err
	}

	for _, remote := range strings.Split(strings.TrimSpace(remotes), "\n") {
		if remote != "" {
			if err := Fetch(remote); err != nil {
				fmt.Printf("Warning: Failed to fetch from %s: %v\n", remote, err)
			}
		}
	}
	return nil
}

// Enhanced Diff Operations
func DiffCached() (string, error) {
	return RunCommand("git diff --cached")
}

func DiffBetweenCommits(commit1, commit2 string) (string, error) {
	return RunCommand(fmt.Sprintf("git diff %s..%s", commit1, commit2))
}

func DiffStats() (string, error) {
	return RunCommand("git diff --stat")
}

func DiffWordLevel(file string) (string, error) {
	if file == "" {
		return RunCommand("git diff --word-diff")
	}
	return RunCommand(fmt.Sprintf("git diff --word-diff %s", file))
}
