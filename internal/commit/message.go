package commit

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cqroot/gear/internal/config"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/cqroot/prompt/write"
)

var p = prompt.New()

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
		config.CommitTypes(),
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

func Run() error {
	// Header
	ctype := Type()

	scope := ""
	if !config.CommitDisableScope() {
		scope = strings.Trim(Scope(), " ")

		if scope != "" {
			scope = "(" + scope + ")"
		}
	}

	summary := Summary()
	message := fmt.Sprintf("%s%s: %s", ctype, scope, summary)

	// Body
	if !config.CommitDisableBody() {
		body := strings.Trim(Body(), " \n")
		if body != "" {
			message = message + "\n\n" + body
		}
	}

	// Footer
	if !config.CommitDisableFooter() {
		issues := Issues()
		if issues != "" {
			message = message + "\n\n" + "Closes " + issues
		}
	}

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
