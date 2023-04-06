package commit

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

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

type Message struct {
	Type    string
	Scope   string
	Summary string
	Body    string
	Footer  string
}

func Type() string {
	ctype, err := p.Ask("Select the type of change:").AdvancedChoose(
		config.CommitTypes(),
		choose.WithTheme(choose.ThemeArrow),
		choose.WithHelp(true),
	)
	CheckErr(err)
	return ctype
}

func Scope() string {
	if config.CommitDisableScope() {
		return ""
	}

	scope, err := p.Ask("Input the scope of change: (skip if empty)").Input(
		"",
		input.WithHelp(true),
	)
	CheckErr(err)

	scope = strings.Trim(scope, " ")
	return scope
}

func Summary() string {
	summary, err := p.Ask("Input the summary of change:").Input(
		"",
		input.WithHelp(true),
	)
	CheckErr(err)
	return summary
}

func Body() string {
	if config.CommitDisableBody() {
		return ""
	}

	body, err := p.Ask("Input the message body of change: (skip if empty)").Write(
		"",
		write.WithHelp(true),
	)
	CheckErr(err)

	body = strings.Trim(body, " \n")
	return body
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
	if config.CommitDisableFooter() {
		return ""
	}

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
	// Footer
	footer := Issues()
	if footer != "" {
		footer = "Closes " + footer
	}

	msg := Message{
		Type:    Type(),
		Scope:   Scope(),
		Summary: Summary(),
		Body:    Body(),
		Footer:  footer,
	}

	tmpl, err := template.New("message").
		Parse(config.CommitMessageTemplate())
	if err != nil {
		return err
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, msg)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "commit", "-m", buf.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
