// Copyright (C) 2026 Keith Chu <cqroot@outlook.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
