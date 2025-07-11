package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ritankarsaha/githubber/pkg/config"
	"github.com/ritankarsaha/githubber/pkg/git"
)

func NewCloneCommand(cfg *config.Config, gitClient *git.Client) *cobra.Command {
	var (
		branch       string
		depth        int
		singleBranch bool
		recursive    bool
		mirror       bool
		bare         bool
	)

	cmd := &cobra.Command{
		Use:   "clone <repository> [directory]",
		Short: "Clone a repository into a new directory",
		Long:  "Clone a repository into a newly created directory.",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			var destination string
			if len(args) > 1 {
				destination = args[1]
			}

			options := git.CloneOptions{
				Branch:       branch,
				Depth:        depth,
				SingleBranch: singleBranch,
				Recursive:    recursive,
				Mirror:       mirror,
				Bare:         bare,
			}

			err := gitClient.Repository.Clone(url, destination, options)
			if err != nil {
				return fmt.Errorf("failed to clone repository: %w", err)
			}

			if destination != "" {
				fmt.Printf("✅ Repository cloned into %s\n", destination)
			} else {
				fmt.Println("✅ Repository cloned successfully")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&branch, "branch", "b", "", "Clone a specific branch instead of HEAD")
	cmd.Flags().IntVar(&depth, "depth", 0, "Create a shallow clone with history truncated to the specified number of commits")
	cmd.Flags().BoolVar(&singleBranch, "single-branch", false, "Clone only one branch")
	cmd.Flags().BoolVar(&recursive, "recursive", false, "Initialize submodules in the clone")
	cmd.Flags().BoolVar(&mirror, "mirror", false, "Set up a mirror of the source repository")
	cmd.Flags().BoolVar(&bare, "bare", false, "Make a bare Git repository")

	return cmd
}