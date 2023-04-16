package cmd

import (
	"fmt"
	"os"

	"github.com/cqroot/gear/internal/committer"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gear",
		Short: "Write conventional git commit messages.",
		Long:  "Write conventional git commit messages.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := committer.New()
			cobra.CheckErr(err)

			cobra.CheckErr(c.Run())
		},
	}

	return rootCmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
