package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Commit CommitConfig `yaml:"commit"`
}

var conf = Config{
	Commit: CommitConfig{
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
	},
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func toBool(val string) bool {
	return strings.ToLower(val) == "true"
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
