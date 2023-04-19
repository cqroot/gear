package config

import (
	"github.com/cqroot/prompt/choose"
)

type CommitType struct {
	Text            string   `yaml:"text"`
	Note            string   `yaml:"note"`
	Scopes          []string `yaml:"scopes"`
	EnableScope     string   `yaml:"enable-scope"`
	EnableBody      string   `yaml:"enable-body"`
	EnableFooter    string   `yaml:"enable-footer"`
	MessageTemplate string   `yaml:"message-template"`
}

type CommitContent struct {
	Types           []CommitType `yaml:"types"`
	Scopes          []string     `yaml:"scopes"`
	EnableScope     string       `yaml:"enable-scope"`
	EnableBody      string       `yaml:"enable-body"`
	EnableFooter    string       `yaml:"enable-footer"`
	MessageTemplate string       `yaml:"message-template"`
}

type CommitConfig struct {
	Scopes          []string
	EnableScope     bool
	EnableBody      bool
	EnableFooter    bool
	MessageTemplate string
}

func (c Config) CommitTypes() []choose.Choice {
	return c.commitTypes
}

func (c *Config) SetCommitType(typ string) {
	c.commitType = typ
}

func (c Config) CommitScopes() []string {
	return c.commitConfig[c.commitType].Scopes
}

func (c Config) CommitEnableScope() bool {
	return c.commitConfig[c.commitType].EnableScope
}

func (c Config) CommitEnableBody() bool {
	return c.commitConfig[c.commitType].EnableBody
}

func (c Config) CommitEnableFooter() bool {
	return c.commitConfig[c.commitType].EnableFooter
}

func (c Config) CommitMessageTemplate() string {
	return c.commitConfig[c.commitType].MessageTemplate
}
