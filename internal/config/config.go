package config

import (
	"os"

	"github.com/cqroot/prompt/choose"
	"gopkg.in/yaml.v3"
)

type CommitConfig struct {
	Types           []choose.Choice `yaml:"types"`
	DisableScope    bool            `yaml:"disable-scope"`
	DisableBody     bool            `yaml:"disable-body"`
	DisableFooter   bool            `yaml:"disable-footer"`
	MessageTemplate string          `yaml:"message-template"`
}

type Config struct {
	Commit CommitConfig `yaml:"commit"`
}

var conf = Config{
	Commit: CommitConfig{
		Types: []choose.Choice{
			{Text: "feat", Note: "A new feature"},
			{Text: "fix", Note: "A bug fix"},
			{Text: "docs", Note: "Documentation only changes"},
			{Text: "refactor", Note: "A code change that neither fixes a bug nor adds a feature"},
			{Text: "test", Note: "Adding missing tests or correcting existing tests"},
			{Text: "build", Note: "Changes that affect the build system or external dependencies"},
			{Text: "ci", Note: "Changes to our CI configuration files and scripts"},
			{Text: "perf", Note: "A code change that improves performance"},
		},
		MessageTemplate: `{{.Type}}{{if .Scope}}({{.Scope}}){{end}}: {{.Summary}}{{if .Body}}

{{.Body}}{{end}}{{if .Footer}}

{{.Footer}}{{end}}`,
	},
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func ReadConfig(name string) error {
	if !fileExists(name) {
		return nil
	}

	content, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &conf)
	if err != nil {
		return err
	}

	return nil
}

func CommitTypes() []choose.Choice {
	return conf.Commit.Types
}

func CommitDisableScope() bool {
	return conf.Commit.DisableScope
}

func CommitDisableBody() bool {
	return conf.Commit.DisableBody
}

func CommitDisableFooter() bool {
	return conf.Commit.DisableFooter
}

func CommitMessageTemplate() string {
	return conf.Commit.MessageTemplate
}
