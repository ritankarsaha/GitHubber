/*
 * GitHubber - Conflict Resolution Interface
 * Author: Ritankar Saha <ritankar.saha786@gmail.com>
 * Description: Interactive conflict resolution tools
 */

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ritankarsaha/git-tool/internal/git"
	"github.com/ritankarsaha/git-tool/internal/ui"
)

// ConflictFile represents a file with merge conflicts
type ConflictFile struct {
	Path     string
	Content  string
	Resolved bool
}

// ConflictResolver provides interactive conflict resolution
type ConflictResolver struct {
	Files []ConflictFile
}

// NewConflictResolver creates a new conflict resolver
func NewConflictResolver() (*ConflictResolver, error) {
	files, err := getConflictedFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get conflicted files: %w", err)
	}

	var conflictFiles []ConflictFile
	for _, file := range files {
		content, err := readFileContent(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", file, err)
		}

		conflictFiles = append(conflictFiles, ConflictFile{
			Path:     file,
			Content:  content,
			Resolved: false,
		})
	}

	return &ConflictResolver{
		Files: conflictFiles,
	}, nil
}

// StartResolution starts the interactive conflict resolution process
func (cr *ConflictResolver) StartResolution() error {
	if len(cr.Files) == 0 {
		fmt.Println(ui.FormatSuccess("No conflicts to resolve!"))
		return nil
	}

	fmt.Println(ui.FormatTitle("🔧 Conflict Resolution Tool"))
	fmt.Printf("Found %d files with conflicts\n\n", len(cr.Files))

	for i := range cr.Files {
		if err := cr.resolveFile(&cr.Files[i]); err != nil {
			return fmt.Errorf("failed to resolve file %s: %w", cr.Files[i].Path, err)
		}
	}

	// Check if all conflicts are resolved
	allResolved := true
	for _, file := range cr.Files {
		if !file.Resolved {
			allResolved = false
			break
		}
	}

	if allResolved {
		fmt.Println(ui.FormatSuccess("All conflicts resolved successfully!"))

		// Ask if user wants to continue with merge/rebase
		choice := GetInput(ui.FormatPrompt("Continue with merge/rebase? (y/n): "))
		if strings.ToLower(choice) == "y" || strings.ToLower(choice) == "yes" {
			return cr.continueMergeRebase()
		}
	} else {
		fmt.Println(ui.FormatWarning("Some conflicts remain unresolved"))
		cr.showUnresolvedFiles()
	}

	return nil
}

func (cr *ConflictResolver) resolveFile(file *ConflictFile) error {
	fmt.Printf("\n%s Resolving conflicts in: %s\n", ui.IconInfo, file.Path)

	// Show conflict markers and content
	conflicts := extractConflicts(file.Content)
	if len(conflicts) == 0 {
		file.Resolved = true
		return nil
	}

	fmt.Printf("Found %d conflict(s) in this file\n", len(conflicts))

	for {
		fmt.Printf("\nOptions for %s:\n", file.Path)
		fmt.Println("1. Show conflicts")
		fmt.Println("2. Accept current changes (HEAD)")
		fmt.Println("3. Accept incoming changes")
		fmt.Println("4. Edit manually")
		fmt.Println("5. Open in editor")
		fmt.Println("6. Mark as resolved")
		fmt.Println("7. Skip this file")

		choice := GetInput(ui.FormatPrompt("Choose an option (1-7): "))

		switch choice {
		case "1":
			cr.showConflicts(file.Path, conflicts)
		case "2":
			if err := cr.acceptCurrent(file); err != nil {
				return err
			}
		case "3":
			if err := cr.acceptIncoming(file); err != nil {
				return err
			}
		case "4":
			if err := cr.editManually(file); err != nil {
				return err
			}
		case "5":
			if err := cr.openInEditor(file.Path); err != nil {
				return err
			}
		case "6":
			file.Resolved = true
			fmt.Println(ui.FormatSuccess("File marked as resolved"))
			return nil
		case "7":
			fmt.Println(ui.FormatWarning("Skipping file - conflicts remain unresolved"))
			return nil
		default:
			fmt.Println(ui.FormatError("Invalid choice. Please try again."))
		}
	}
}

func (cr *ConflictResolver) showConflicts(filePath string, conflicts []ConflictSection) {
	fmt.Printf("\n%s Conflicts in %s:\n", ui.IconCommit, filePath)

	for i, conflict := range conflicts {
		fmt.Printf("\n--- Conflict %d ---\n", i+1)
		fmt.Printf("📥 Current (HEAD):\n%s\n", ui.FormatCode(conflict.Current))
		fmt.Printf("📤 Incoming:\n%s\n", ui.FormatCode(conflict.Incoming))
		if conflict.Base != "" {
			fmt.Printf("🔀 Base:\n%s\n", ui.FormatCode(conflict.Base))
		}
	}
}

func (cr *ConflictResolver) acceptCurrent(file *ConflictFile) error {
	resolved := resolveConflicts(file.Content, "current")
	if err := writeFileContent(file.Path, resolved); err != nil {
		return err
	}

	file.Content = resolved
	file.Resolved = true
	fmt.Println(ui.FormatSuccess("Accepted current changes"))
	return nil
}

func (cr *ConflictResolver) acceptIncoming(file *ConflictFile) error {
	resolved := resolveConflicts(file.Content, "incoming")
	if err := writeFileContent(file.Path, resolved); err != nil {
		return err
	}

	file.Content = resolved
	file.Resolved = true
	fmt.Println(ui.FormatSuccess("Accepted incoming changes"))
	return nil
}

func (cr *ConflictResolver) editManually(file *ConflictFile) error {
	fmt.Println(ui.FormatInfo("Manual editing mode - resolve conflicts and save the file"))
	fmt.Printf("File: %s\n", file.Path)
	fmt.Println("Conflict markers:")
	fmt.Println("  <<<<<<< HEAD (your changes)")
	fmt.Println("  ======= (separator)")
	fmt.Println("  >>>>>>> branch (incoming changes)")

	GetInput("Press Enter when you've finished editing the file manually...")

	// Re-read the file to check if conflicts are resolved
	content, err := readFileContent(file.Path)
	if err != nil {
		return err
	}

	file.Content = content
	conflicts := extractConflicts(content)

	if len(conflicts) == 0 {
		file.Resolved = true
		fmt.Println(ui.FormatSuccess("All conflicts resolved in this file!"))
	} else {
		fmt.Printf(ui.FormatWarning("File still has %d unresolved conflicts"), len(conflicts))
	}

	return nil
}

func (cr *ConflictResolver) openInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // Default fallback
	}

	fmt.Printf("Opening %s in %s...\n", filePath, editor)

	_, err := git.RunCommand(fmt.Sprintf("%s %s", editor, filePath))
	if err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Re-check for conflicts after editing
	content, err := readFileContent(filePath)
	if err != nil {
		return err
	}

	conflicts := extractConflicts(content)
	if len(conflicts) == 0 {
		fmt.Println(ui.FormatSuccess("All conflicts resolved!"))
	} else {
		fmt.Printf(ui.FormatWarning("File still has %d unresolved conflicts"), len(conflicts))
	}

	return nil
}

func (cr *ConflictResolver) continueMergeRebase() error {
	// Check if we're in a merge or rebase state
	if isInMergeState() {
		fmt.Println(ui.FormatInfo("Continuing merge..."))
		_, err := git.RunCommand("git commit")
		return err
	}

	if isInRebaseState() {
		fmt.Println(ui.FormatInfo("Continuing rebase..."))
		return git.RebaseContinue()
	}

	if isInCherryPickState() {
		fmt.Println(ui.FormatInfo("Continuing cherry-pick..."))
		return git.CherryPickContinue()
	}

	return fmt.Errorf("not in a merge/rebase/cherry-pick state")
}

func (cr *ConflictResolver) showUnresolvedFiles() {
	fmt.Println(ui.FormatWarning("Unresolved files:"))
	for _, file := range cr.Files {
		if !file.Resolved {
			fmt.Printf("  - %s\n", file.Path)
		}
	}
}

// ConflictSection represents a single conflict section
type ConflictSection struct {
	Current  string
	Incoming string
	Base     string
}

// Helper functions

func getConflictedFiles() ([]string, error) {
	output, err := git.RunCommand("git diff --name-only --diff-filter=U")
	if err != nil {
		return nil, err
	}

	if output == "" {
		return []string{}, nil
	}

	return strings.Split(output, "\n"), nil
}

func readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeFileContent(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func extractConflicts(content string) []ConflictSection {
	var conflicts []ConflictSection
	lines := strings.Split(content, "\n")

	i := 0
	for i < len(lines) {
		if strings.HasPrefix(lines[i], "<<<<<<<") {
			conflict := ConflictSection{}
			i++ // Skip conflict marker

			// Read current (HEAD) section
			var currentLines []string
			for i < len(lines) && !strings.HasPrefix(lines[i], "=======") {
				currentLines = append(currentLines, lines[i])
				i++
			}
			conflict.Current = strings.Join(currentLines, "\n")

			if i < len(lines) {
				i++ // Skip separator
			}

			// Read incoming section
			var incomingLines []string
			for i < len(lines) && !strings.HasPrefix(lines[i], ">>>>>>>") {
				incomingLines = append(incomingLines, lines[i])
				i++
			}
			conflict.Incoming = strings.Join(incomingLines, "\n")

			conflicts = append(conflicts, conflict)
		}
		i++
	}

	return conflicts
}

func resolveConflicts(content, resolution string) string {
	lines := strings.Split(content, "\n")
	var result []string

	i := 0
	for i < len(lines) {
		if strings.HasPrefix(lines[i], "<<<<<<<") {
			i++ // Skip conflict marker

			if resolution == "current" {
				// Keep current (HEAD) section
				for i < len(lines) && !strings.HasPrefix(lines[i], "=======") {
					result = append(result, lines[i])
					i++
				}
				// Skip to end of conflict
				for i < len(lines) && !strings.HasPrefix(lines[i], ">>>>>>>") {
					i++
				}
			} else if resolution == "incoming" {
				// Skip current (HEAD) section
				for i < len(lines) && !strings.HasPrefix(lines[i], "=======") {
					i++
				}
				if i < len(lines) {
					i++ // Skip separator
				}
				// Keep incoming section
				for i < len(lines) && !strings.HasPrefix(lines[i], ">>>>>>>") {
					result = append(result, lines[i])
					i++
				}
			}
		} else {
			result = append(result, lines[i])
		}
		i++
	}

	return strings.Join(result, "\n")
}

func isInMergeState() bool {
	_, err := os.Stat(filepath.Join(".git", "MERGE_HEAD"))
	return err == nil
}

func isInRebaseState() bool {
	_, err := os.Stat(filepath.Join(".git", "rebase-apply"))
	if err == nil {
		return true
	}
	_, err = os.Stat(filepath.Join(".git", "rebase-merge"))
	return err == nil
}

func isInCherryPickState() bool {
	_, err := os.Stat(filepath.Join(".git", "CHERRY_PICK_HEAD"))
	return err == nil
}

// StartConflictResolution starts the conflict resolution process
func StartConflictResolution() error {
	resolver, err := NewConflictResolver()
	if err != nil {
		return err
	}

	return resolver.StartResolution()
}
