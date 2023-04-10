package committer

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

type Committer struct {
	p   prompt.Prompt
	typ string
}

func checkErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}
}

func New() *Committer {
	c := Committer{
		p: *prompt.New(),
	}

	typ, err := c.p.Ask("Select the type of change:").AdvancedChoose(
		config.CommitTypes(),
		choose.WithTheme(choose.ThemeArrow),
	)
	checkErr(err)

	c.typ = typ
	return &c
}

func (c Committer) scope() string {
	if config.CommitDisableScope() {
		return ""
	}

	scope, err := c.p.Ask("Input the scope of change: (skip if empty)").Input(
		"",
		input.WithHelp(true),
	)
	checkErr(err)

	scope = strings.Trim(scope, " ")
	return scope
}

func (c Committer) summary() string {
	summary, err := c.p.Ask("Input the summary of change:").Input(
		"",
		input.WithHelp(true),
	)
	checkErr(err)
	return summary
}

func (c Committer) body() string {
	if config.CommitDisableBody() {
		return ""
	}

	body, err := c.p.Ask("Input the message body of change: (skip if empty)").Write(
		"",
		write.WithHelp(true),
	)
	checkErr(err)

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

func (c Committer) issues() string {
	if config.CommitDisableFooter() {
		return ""
	}

	issues, err := c.p.Ask("Input the issues you want to close: (Such as \"#1, #2\". skip if empty)").Input(
		"", input.WithHelp(true),
		input.WithValidateFunc(validateIssues),
	)
	if err != nil {
		return ""
	}
	return issues
}

func (c Committer) footer() string {
	footer := c.issues()
	if footer != "" {
		footer = "Closes " + footer
	}
	return footer
}

type Message struct {
	Type    string
	Scope   string
	Summary string
	Body    string
	Footer  string
}

func (c Committer) Run() error {
	msg := Message{
		Type:    c.typ,
		Scope:   c.scope(),
		Summary: c.summary(),
		Body:    c.body(),
		Footer:  c.footer(),
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
