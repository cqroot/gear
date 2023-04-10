package cmd

import (
	"github.com/cqroot/gear/internal/committer"
	"github.com/spf13/cobra"
)

func newCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Write conventional git commit messages.",
		Long:  "Write conventional git commit messages.",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(committer.New().Run())
		},
	}
	return cmd
}
