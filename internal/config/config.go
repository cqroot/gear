package config

import (
	"os"
	"strings"

	"github.com/cqroot/prompt/choose"
	"gopkg.in/yaml.v3"
)

type Config struct {
	content struct {
		Commit CommitContent `yaml:"commit"`
	}
	commitType   string
	commitTypes  []choose.Choice
	commitConfig map[string]CommitConfig
}

func New(name string) (*Config, error) {
	c := Config{
		commitConfig: make(map[string]CommitConfig),
	}

	c.content.Commit = CommitContent{
		Types: []CommitType{
			{Text: "feat", Note: "A new feature"},
			{Text: "fix", Note: "A bug fix"},
			{Text: "docs", Note: "Documentation only changes"},
			{Text: "refactor", Note: "A code change that neither fixes a bug nor adds a feature"},
			{Text: "test", Note: "Adding missing tests or correcting existing tests"},
			{Text: "build", Note: "Changes that affect the build system or external dependencies"},
			{Text: "ci", Note: "Changes to our CI configuration files and scripts"},
			{Text: "perf", Note: "A code change that improves performance"},
		},
		EnableScope:  "false",
		EnableBody:   "false",
		EnableFooter: "false",
		MessageTemplate: `{{.Type}}{{if .Scope}}({{.Scope}}){{end}}: {{.Summary}}{{if .Body}}

{{.Body}}{{end}}{{if .Footer}}

{{.Footer}}{{end}}`,
	}

	var err error
	if name != "" {
		err = c.read(name)
	}

	c.commitConfig[""] = CommitConfig{
		Scopes:          c.content.Commit.Scopes,
		EnableScope:     toBool(c.content.Commit.EnableScope),
		EnableBody:      toBool(c.content.Commit.EnableBody),
		EnableFooter:    toBool(c.content.Commit.EnableFooter),
		MessageTemplate: c.content.Commit.MessageTemplate,
	}

	for _, typ := range c.content.Commit.Types {
		c.commitTypes = append(c.commitTypes, choose.Choice{
			Text: typ.Text,
			Note: typ.Note,
		})

		c.commitConfig[typ.Text] = CommitConfig{
			Scopes: func() []string {
				if len(typ.Scopes) == 0 {
					return c.commitConfig[""].Scopes
				}
				return typ.Scopes
			}(),
			EnableScope: func() bool {
				if typ.EnableScope == "" {
					return c.commitConfig[""].EnableScope
				}
				return toBool(typ.EnableScope)
			}(),
			EnableBody: func() bool {
				if typ.EnableBody == "" {
					return c.commitConfig[""].EnableBody
				}
				return toBool(typ.EnableBody)
			}(),
			EnableFooter: func() bool {
				if typ.EnableFooter == "" {
					return c.commitConfig[""].EnableFooter
				}
				return toBool(typ.EnableFooter)
			}(),
			MessageTemplate: func() string {
				if typ.MessageTemplate == "" {
					return c.commitConfig[""].MessageTemplate
				}
				return typ.MessageTemplate
			}(),
		}
	}

	return &c, err
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func toBool(val string) bool {
	return strings.ToLower(val) == "true"
}

func (c *Config) read(name string) error {
	if !fileExists(name) {
		return nil
	}

	content, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &c.content)
	if err != nil {
		return err
	}

	return nil
}
