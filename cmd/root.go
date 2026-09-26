package cmd

import (
	"github.com/cqroot/gear/internal/version"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "gear",
		Short: "Lightweight utilities around the Go toolchain.",
		Long:  "gear is a thin CLI wrapping common Go development tasks.",
	}
	c.AddCommand(newRemodCmd())
	c.Version = version.Get().String()
	return &c
}

func Execute() {
	c := newRootCmd()
	err := c.Execute()
	cobra.CheckErr(err)
}
