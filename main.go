package main

import (
	"github.com/cqroot/gear/cmd"
	"github.com/cqroot/gear/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	err := config.ReadConfig("./.gear.yml")
	cobra.CheckErr(err)

	cmd.Execute()
}
