package config_test

import (
	"testing"

	"github.com/cqroot/gear/internal/config"
	"github.com/cqroot/prompt/choose"
	"github.com/stretchr/testify/require"
)

func TestCommitTypes(t *testing.T) {
	defConf, err := config.New("")
	require.Nil(t, err)

	require.Equal(t, []choose.Choice{
		{Text: "feat", Note: "A new feature"},
		{Text: "fix", Note: "A bug fix"},
		{Text: "docs", Note: "Documentation only changes"},
		{Text: "refactor", Note: "A code change that neither fixes a bug nor adds a feature"},
		{Text: "test", Note: "Adding missing tests or correcting existing tests"},
		{Text: "build", Note: "Changes that affect the build system or external dependencies"},
		{Text: "ci", Note: "Changes to our CI configuration files and scripts"},
		{Text: "perf", Note: "A code change that improves performance"},
	}, defConf.CommitTypes())

	conf, err := config.New("./testdata/gear_1.yml")
	require.Nil(t, err)

	require.Equal(t, []choose.Choice{
		{Text: "✨", Note: "feat: A new feature"},
		{Text: "🐛", Note: "fix: A bug fix"},
		{Text: "🔧", Note: "build: Changes that affect the build system or external dependencies"},
		{Text: "📝", Note: "docs: Documentation only changes"},
		{Text: "🎨", Note: "refactor: A code change that neither fixes a bug nor adds a feature"},
		{Text: "🧪", Note: "test: Adding missing tests or correcting existing tests"},
		{Text: "👷", Note: "ci: Changes to our CI configuration files and scripts"},
		{Text: "⚡️", Note: "perf: A code change that improves performance"},
	}, conf.CommitTypes())
}

func TestScopes(t *testing.T) {
	defConf, err := config.New("")
	require.Nil(t, err)
	require.Equal(t, []string(nil), defConf.CommitScopes())

	conf1, err := config.New("./testdata/gear_1.yml")
	require.Nil(t, err)
	require.Equal(t, []string{
		"something", "others",
	}, conf1.CommitScopes())

	conf2, err := config.New("./testdata/gear_2.yml")
	require.Nil(t, err)
	conf2.SetCommitType("✨")
	require.Equal(t, []string{
		"something", "others",
	}, conf2.CommitScopes())
}

func TestCommitEnableScope(t *testing.T) {
	defConf, err := config.New("")
	require.Nil(t, err)
	require.Equal(t, false, defConf.CommitEnableScope())

	conf1, err := config.New("./testdata/gear_1.yml")
	require.Nil(t, err)
	require.Equal(t, true, conf1.CommitEnableScope())

	conf2, err := config.New("./testdata/gear_2.yml")
	require.Nil(t, err)
	conf2.SetCommitType("✨")
	require.Equal(t, true, conf2.CommitEnableScope())
}

func TestCommitEnableBody(t *testing.T) {
	defConf, err := config.New("")
	require.Nil(t, err)
	require.Equal(t, false, defConf.CommitEnableBody())

	conf1, err := config.New("./testdata/gear_1.yml")
	require.Nil(t, err)
	require.Equal(t, true, conf1.CommitEnableBody())

	conf2, err := config.New("./testdata/gear_2.yml")
	require.Nil(t, err)
	conf2.SetCommitType("✨")
	require.Equal(t, true, conf2.CommitEnableBody())
}

func TestCommitEnableFooter(t *testing.T) {
	defConf, err := config.New("")
	require.Nil(t, err)
	require.Equal(t, false, defConf.CommitEnableFooter())

	conf1, err := config.New("./testdata/gear_1.yml")
	require.Nil(t, err)
	require.Equal(t, true, conf1.CommitEnableFooter())

	conf2, err := config.New("./testdata/gear_2.yml")
	require.Nil(t, err)
	conf2.SetCommitType("✨")
	require.Equal(t, true, conf2.CommitEnableFooter())
}

func TestCommitMessageTemplate(t *testing.T) {
	defConf, err := config.New("")
	require.Nil(t, err)
	require.Equal(t, `{{.Type}}{{if .Scope}}({{.Scope}}){{end}}: {{.Summary}}{{if .Body}}

{{.Body}}{{end}}{{if .Footer}}

{{.Footer}}{{end}}`, defConf.CommitMessageTemplate())

	conf1, err := config.New("./testdata/gear_1.yml")
	require.Nil(t, err)
	require.Equal(t, "{{.Type}} {{if .Scope}}({{.Scope}}): {{end}}{{.Summary}}", conf1.CommitMessageTemplate())

	conf2, err := config.New("./testdata/gear_2.yml")
	require.Nil(t, err)
	conf2.SetCommitType("✨")
	require.Equal(t, "{{.Type}} {{if .Scope}}({{.Scope}}): {{end}}{{.Summary}}", conf2.CommitMessageTemplate())
}
