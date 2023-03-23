package commit

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/cqroot/prompt/write"
)

var (
	p           = prompt.New()
	CommitTypes = []choose.Choice{
		{Text: "feat", Note: "A new feature"},
		{Text: "fix", Note: "A bug fix"},
		{Text: "docs", Note: "Documentation only changes"},
		{Text: "refactor", Note: "A code change that neither fixes a bug nor adds a feature"},
		{Text: "test", Note: "Adding missing tests or correcting existing tests"},
		{Text: "build", Note: "Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)"},
		{Text: "ci", Note: "Changes to our CI configuration files and scripts (examples: CircleCi, SauceLabs)"},
		{Text: "perf", Note: "A code change that improves performance"},
	}
)

func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

func Type() string {
	ctype, err := p.Ask("Select the type of change:").AdvancedChoose(
		CommitTypes,
		choose.WithHelp(true),
	)
	CheckErr(err)
	return ctype
}

func Scope() string {
	scope, err := p.Ask("Input the scope of change: (skip if empty)").Input(
		"",
		input.WithHelp(true),
	)
	CheckErr(err)
	return scope
}

func Body() string {
	scope, err := p.Ask("Input the summary of change:").Write(
		"",
		write.WithHelp(true),
	)
	CheckErr(err)
	return scope
}

func Summary() string {
	scope, err := p.Ask("Input the message body of change: (skip if empty)").Input(
		"",
		input.WithHelp(true),
	)
	CheckErr(err)
	return scope
}

func validateIssues(text string) error {
	if len(text) == 0 {
		return nil
	}

	issues := strings.Split(text, ", ")
	for _, issue := range issues {
		if len(issue) == 0 {
			return errors.New("Empty issues are not allowed")
		}
		if issue[0] != '#' {
			return errors.New("Issue must start with #")
		}
		if _, err := strconv.Atoi(issue[1:]); err != nil {
			return errors.New("Issue must be like \"#number\"")
		}
	}
	return nil
}

func Issues() string {
	issues, err := p.Ask("Input the issues you want to close: (Such as \"#1, #2\". skip if empty)").Input(
		"", input.WithHelp(true),
		input.WithValidateFunc(validateIssues),
	)
	if err != nil {
		return ""
	}
	return issues
}
