package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ritankarsaha/githubber/pkg/git"
)

func TestRepository(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Init", func(t *testing.T) {
		executor.SetResponse("git init", "Initialized empty Git repository")

		err := client.Repository.Init("", false)
		assert.NoError(t, err)

		commands := executor.GetExecutedCommands()
		assert.Contains(t, commands, "git init")
	})

	t.Run("IsGitRepository", func(t *testing.T) {
		executor.SetResponse("git rev-parse --git-dir", ".git")

		isRepo := client.Repository.IsGitRepository()
		assert.True(t, isRepo)
	})

	t.Run("GetCurrentBranch", func(t *testing.T) {
		executor.SetResponse("git rev-parse --abbrev-ref HEAD", "main")

		branch, err := client.Repository.GetCurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "main", branch)
	})
}

func TestBranchManager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Create", func(t *testing.T) {
		executor.SetResponse("git checkout -b feature-branch", "Switched to a new branch 'feature-branch'")

		err := client.Branches.Create("feature-branch", "")
		assert.NoError(t, err)
	})

	t.Run("Switch", func(t *testing.T) {
		executor.SetResponse("git checkout main", "Switched to branch 'main'")

		err := client.Branches.Switch("main")
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		executor.SetResponse("git branch -d feature-branch", "Deleted branch feature-branch")

		err := client.Branches.Delete("feature-branch", false)
		assert.NoError(t, err)
	})
}

func TestCommitManager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Commit", func(t *testing.T) {
		executor.SetResponse("git commit -m Test commit", "[main abc1234] Test commit")

		options := git.CommitOptions{}
		err := client.Commits.Commit("Test commit", options)
		assert.NoError(t, err)
	})

	t.Run("AmendLastCommit", func(t *testing.T) {
		executor.SetResponse("git commit --amend --no-edit", "[main abc1234] Test commit")

		err := client.Commits.AmendLastCommit("", true)
		assert.NoError(t, err)
	})
}

func TestStager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("AddFiles", func(t *testing.T) {
		executor.SetResponse("git add file1.go file2.go", "")

		err := client.Staging.AddFiles([]string{"file1.go", "file2.go"})
		assert.NoError(t, err)
	})

	t.Run("AddAll", func(t *testing.T) {
		executor.SetResponse("git add .", "")

		err := client.Staging.AddAll()
		assert.NoError(t, err)
	})

	t.Run("IsClean", func(t *testing.T) {
		executor.SetResponse("git status --porcelain", "")

		clean, err := client.Staging.IsClean()
		assert.NoError(t, err)
		assert.True(t, clean)
	})
}

func TestRemoteManager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Add", func(t *testing.T) {
		executor.SetResponse("git remote add origin https://github.com/user/repo.git", "")

		err := client.Remotes.Add("origin", "https://github.com/user/repo.git")
		assert.NoError(t, err)
	})

	t.Run("Remove", func(t *testing.T) {
		executor.SetResponse("git remote remove origin", "")

		err := client.Remotes.Remove("origin")
		assert.NoError(t, err)
	})

	t.Run("GetURL", func(t *testing.T) {
		executor.SetResponse("git remote get-url origin", "https://github.com/user/repo.git")

		url, err := client.Remotes.GetURL("origin", false)
		assert.NoError(t, err)
		assert.Equal(t, "https://github.com/user/repo.git", url)
	})
}

func TestStashManager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Save", func(t *testing.T) {
		executor.SetResponse("git stash push -m Test stash", "Saved working directory and index state")

		options := git.StashSaveOptions{}
		err := client.Stash.Save("Test stash", options)
		assert.NoError(t, err)
	})

	t.Run("Pop", func(t *testing.T) {
		executor.SetResponse("git stash pop", "Applied stash@{0}")

		err := client.Stash.Pop("")
		assert.NoError(t, err)
	})

	t.Run("List", func(t *testing.T) {
		executor.SetResponse("git stash list", "stash@{0}: WIP on main: abc1234 Test commit")

		stashes, err := client.Stash.List()
		assert.NoError(t, err)
		assert.Len(t, stashes, 1)
		assert.Equal(t, 0, stashes[0].Index)
	})
}

func TestTagManager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("CreateLightweight", func(t *testing.T) {
		executor.SetResponse("git tag v1.0.0", "")

		err := client.Tags.CreateLightweight("v1.0.0", "", false)
		assert.NoError(t, err)
	})

	t.Run("CreateAnnotated", func(t *testing.T) {
		executor.SetResponse("git tag -a -m Version 1.0.0 v1.0.0", "")

		err := client.Tags.CreateAnnotated("v1.0.0", "Version 1.0.0", "", false)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		executor.SetResponse("git tag -d v1.0.0", "Deleted tag 'v1.0.0'")

		err := client.Tags.Delete([]string{"v1.0.0"})
		assert.NoError(t, err)
	})
}

func TestRebaseManager(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Rebase", func(t *testing.T) {
		executor.SetResponse("git rebase main", "Successfully rebased and updated refs/heads/feature")

		options := git.RebaseOptions{}
		err := client.Rebase.Rebase("main", options)
		assert.NoError(t, err)
	})

	t.Run("InteractiveRebase", func(t *testing.T) {
		executor.SetResponse("git rebase -i main", "Successfully rebased and updated refs/heads/feature")

		options := git.InteractiveOptions{}
		err := client.Rebase.InteractiveRebase("main", options)
		assert.NoError(t, err)
	})

	t.Run("Continue", func(t *testing.T) {
		executor.SetResponse("git rebase --continue", "")

		err := client.Rebase.Continue()
		assert.NoError(t, err)
	})

	t.Run("Abort", func(t *testing.T) {
		executor.SetResponse("git rebase --abort", "")

		err := client.Rebase.Abort()
		assert.NoError(t, err)
	})
}

func TestGitClient(t *testing.T) {
	client, executor := git.NewTestClient()

	t.Run("Version", func(t *testing.T) {
		executor.SetResponse("git --version", "git version 2.34.1")

		version, err := client.Version()
		assert.NoError(t, err)
		assert.Contains(t, version, "git version")
	})

	t.Run("IsGitRepository", func(t *testing.T) {
		executor.SetResponse("git rev-parse --git-dir", ".git")

		isRepo := client.IsGitRepository()
		assert.True(t, isRepo)
	})

	t.Run("ValidateGitInstallation", func(t *testing.T) {
		executor.SetResponse("git --version", "git version 2.34.1")

		err := client.ValidateGitInstallation()
		assert.NoError(t, err)
	})
}