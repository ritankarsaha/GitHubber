package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ritankarsaha/githubber/pkg/config"
	"github.com/ritankarsaha/githubber/pkg/git"
)

func NewStatusCommand(cfg *config.Config, gitClient *git.Client) *cobra.Command {
	var (
		short   bool
		porcelain bool
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the working tree status",
		Long:  "Display paths that have differences between the index file and the current HEAD commit.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !gitClient.IsGitRepository() {
				return fmt.Errorf("not in a git repository")
			}

			if short || porcelain {
				status, err := gitClient.Staging.GetStatus()
				if err != nil {
					return fmt.Errorf("failed to get status: %w", err)
				}
				
				fmt.Printf("## %s\n", status.Branch)
				if status.Ahead > 0 || status.Behind > 0 {
					fmt.Printf("## %s [ahead %d, behind %d]\n", status.Branch, status.Ahead, status.Behind)
				}
				
				for _, file := range status.StagedFiles {
					fmt.Printf("%s %s\n", file.Status, file.Path)
				}
				for _, file := range status.ModifiedFiles {
					fmt.Printf(" M %s\n", file.Path)
				}
				for _, file := range status.UntrackedFiles {
					fmt.Printf("?? %s\n", file)
				}
				
				return nil
			}

			output, err := gitClient.Staging.GetDetailedStatus()
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}

			fmt.Print(output)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&short, "short", "s", false, "Give the output in the short-format")
	cmd.Flags().BoolVar(&porcelain, "porcelain", false, "Give the output in an easy-to-parse format")

	return cmd
}