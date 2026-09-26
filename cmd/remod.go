package cmd

import (
	"github.com/cqroot/gear/pkg/remod"
	"github.com/spf13/cobra"
)

func newRemodCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "remod",
		Short: "Rebuild go.mod and go.sum in place",
		Long: `remod reconstructs the Go module files for the current directory.

It will:
  1. Detect the project root (the working directory).
  2. Read the existing go.mod to capture the declared module path.
  3. Delete the current go.mod and go.sum files.
  4. Re-run 'go mod init <module>' followed by 'go mod tidy'.

The module path is preserved so existing import paths keep working.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return remod.Run(
				"", // empty => use os.Getwd() inside pkg/remod
				cmd.OutOrStdout(),
				cmd.ErrOrStderr(),
			)
		},
		SilenceUsage: true,
	}
	return c
}
