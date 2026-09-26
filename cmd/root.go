package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	c := cobra.Command{
		Use:   "gear",
		Short: "Lightweight utilities around the Go toolchain.",
		Long:  "gear is a thin CLI wrapping common Go development tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello world.")
		},
	}
	c.AddCommand(newRemodCmd())
	return &c
}

func Execute() {
	c := newRootCmd()
	err := c.Execute()
	cobra.CheckErr(err)
}
