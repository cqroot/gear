package config

import (
	"github.com/cqroot/prompt/choose"
)

type CommitType struct {
	Text            string `yaml:"text"`
	Note            string `yaml:"note"`
	EnableScope     string `yaml:"enable-scope"`
	EnableBody      string `yaml:"enable-body"`
	EnableFooter    string `yaml:"enable-footer"`
	MessageTemplate string `yaml:"message-template"`
}

type CommitConfig struct {
	Types           []CommitType `yaml:"types"`
	EnableScope     string       `yaml:"enable-scope"`
	EnableBody      string       `yaml:"enable-body"`
	EnableFooter    string       `yaml:"enable-footer"`
	MessageTemplate string       `yaml:"message-template"`
}

var (
	commitEnableScope     = "false"
	commitEnableBody      = "false"
	commitEnableFooter    = "false"
	commitMessageTemplate = `{{.Type}}{{if .Scope}}({{.Scope}}){{end}}: {{.Summary}}{{if .Body}}

{{.Body}}{{end}}{{if .Footer}}

{{.Footer}}{{end}}`
)

func CommitTypes() []choose.Choice {
	var choices []choose.Choice

	for _, typ := range conf.Commit.Types {
		choices = append(choices, choose.Choice{
			Text: typ.Text,
			Note: typ.Note,
		})
	}

	return choices
}

func InitCommitConfig(typ string) {
	set := func(val string, def string) string {
		if val == "" {
			return def
		}
		return val
	}

	for _, t := range conf.Commit.Types {
		if t.Text != typ {
			continue
		}

		commitEnableScope = set(t.EnableScope, conf.Commit.EnableScope)
		commitEnableBody = set(t.EnableBody, conf.Commit.EnableScope)
		commitEnableFooter = set(t.EnableFooter, conf.Commit.EnableScope)
		commitMessageTemplate = set(t.MessageTemplate, conf.Commit.MessageTemplate)
	}
}

func CommitEnableScope() bool {
	return toBool(commitEnableScope)
}

func CommitEnableBody() bool {
	return toBool(commitEnableBody)
}

func CommitEnableFooter() bool {
	return toBool(commitEnableFooter)
}

func CommitMessageTemplate() string {
	return commitMessageTemplate
}
