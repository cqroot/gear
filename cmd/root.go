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
