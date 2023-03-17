package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cqroot/git-commit-helper/internal/commit"
)

func main() {
	ctype := commit.Type()
	scope := commit.Scope()
	if scope != "" {
		scope = "(" + scope + ")"
	}
	summary := commit.Summary()
	message := fmt.Sprintf("%s%s: %s", ctype, scope, summary)

	body := commit.Body()
	if body != "" {
		message = message + "\n\n" + body
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
