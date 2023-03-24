package cmd

import (
	"github.com/cqroot/git-commit-helper/internal/commit"
	"github.com/spf13/cobra"
)

func newCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(commit.Run())
		},
	}
	return cmd
}
