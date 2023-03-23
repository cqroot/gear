package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cqroot/git-commit-helper/internal/commit"
)

func main() {
	ctype := commit.Type()
	scope := strings.Trim(commit.Scope(), " ")
	if scope != "" {
		scope = "(" + scope + ")"
	}
	summary := commit.Summary()
	message := fmt.Sprintf("%s%s: %s", ctype, scope, summary)

	body := strings.Trim(commit.Body(), " \n")
	if body != "" {
		message = message + "\n\n" + body
	}

	issues := commit.Issues()
	if issues != "" {
		message = message + "\n\n" + "Closes " + issues
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
