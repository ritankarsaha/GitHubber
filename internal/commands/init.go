package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ritankarsaha/githubber/pkg/config"
	"github.com/ritankarsaha/githubber/pkg/git"
)

func NewInitCommand(cfg *config.Config, gitClient *git.Client) *cobra.Command {
	var (
		bare bool
		templateDir string
	)

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new Git repository",
		Long:  "Create an empty Git repository or reinitialize an existing one.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string
			if len(args) > 0 {
				path = args[0]
			}

			if path != "" {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return fmt.Errorf("failed to get absolute path: %w", err)
				}
				path = absPath
			}

			err := gitClient.Repository.Init(path, bare)
			if err != nil {
				return fmt.Errorf("failed to initialize repository: %w", err)
			}

			if path != "" {
				fmt.Printf("✅ Initialized empty Git repository in %s\n", path)
			} else {
				fmt.Println("✅ Initialized empty Git repository in current directory")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&bare, "bare", false, "Create a bare repository")
	cmd.Flags().StringVar(&templateDir, "template", "", "Directory from which templates will be used")

	return cmd
}