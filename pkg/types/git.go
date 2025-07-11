package types

import "time"

type RepositoryInfo struct {
	URL           string `json:"url"`
	CurrentBranch string `json:"current_branch"`
	RemoteName    string `json:"remote_name"`
	IsGitRepo     bool   `json:"is_git_repo"`
	RootPath      string `json:"root_path"`
}

type BranchInfo struct {
	Name      string `json:"name"`
	IsCurrent bool   `json:"is_current"`
	IsRemote  bool   `json:"is_remote"`
	Upstream  string `json:"upstream,omitempty"`
}

type CommitInfo struct {
	Hash        string    `json:"hash"`
	ShortHash   string    `json:"short_hash"`
	Message     string    `json:"message"`
	Author      string    `json:"author"`
	AuthorEmail string    `json:"author_email"`
	Date        time.Time `json:"date"`
	Parents     []string  `json:"parents"`
}

type TagInfo struct {
	Name        string    `json:"name"`
	Message     string    `json:"message"`
	Hash        string    `json:"hash"`
	Date        time.Time `json:"date"`
	IsAnnotated bool      `json:"is_annotated"`
}

type StashInfo struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
	Branch  string `json:"branch"`
	Hash    string `json:"hash"`
}

type RemoteInfo struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	FetchURL string `json:"fetch_url"`
	PushURL  string `json:"push_url"`
}

type StatusInfo struct {
	Branch              string             `json:"branch"`
	Ahead               int                `json:"ahead"`
	Behind              int                `json:"behind"`
	StagedFiles         []FileStatusInfo   `json:"staged_files"`
	ModifiedFiles       []FileStatusInfo   `json:"modified_files"`
	UntrackedFiles      []string           `json:"untracked_files"`
	DeletedFiles        []FileStatusInfo   `json:"deleted_files"`
	RenamedFiles        []RenamedFileInfo  `json:"renamed_files"`
	ConflictedFiles     []FileStatusInfo   `json:"conflicted_files"`
	IsClean             bool               `json:"is_clean"`
}

type FileStatusInfo struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

type RenamedFileInfo struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	Status  string `json:"status"`
}

type DiffInfo struct {
	FilePath     string      `json:"file_path"`
	Additions    int         `json:"additions"`
	Deletions    int         `json:"deletions"`
	Changes      []LineChange `json:"changes"`
	IsBinary     bool        `json:"is_binary"`
}

type LineChange struct {
	LineNumber int    `json:"line_number"`
	Type       string `json:"type"` // "added", "deleted", "context"
	Content    string `json:"content"`
}

type PullRequestInfo struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BaseBranch  string `json:"base_branch"`
	HeadBranch  string `json:"head_branch"`
	State       string `json:"state"`
	URL         string `json:"url"`
}

type GitHubIssue struct {
	Number      int      `json:"number"`
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	State       string   `json:"state"`
	Labels      []string `json:"labels"`
	Assignees   []string `json:"assignees"`
	URL         string   `json:"url"`
}

type GitConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Core  struct {
		Editor   string `json:"editor"`
		Autocrlf string `json:"autocrlf"`
	} `json:"core"`
	Push struct {
		Default string `json:"default"`
	} `json:"push"`
}

type RebaseInfo struct {
	InProgress bool     `json:"in_progress"`
	Branch     string   `json:"branch"`
	Onto       string   `json:"onto"`
	Steps      []string `json:"steps"`
}

type MergeInfo struct {
	InProgress      bool     `json:"in_progress"`
	Branch          string   `json:"branch"`
	ConflictedFiles []string `json:"conflicted_files"`
}

type WorkflowInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
}