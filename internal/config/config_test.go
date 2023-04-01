package config_test

import (
	"testing"

	"github.com/cqroot/gear/internal/config"
	"github.com/cqroot/prompt/choose"
	"github.com/stretchr/testify/require"
)

func TestCommitTypes(t *testing.T) {
	require.Equal(t, []choose.Choice{
		{Text: "feat", Note: "A new feature"},
		{Text: "fix", Note: "A bug fix"},
		{Text: "docs", Note: "Documentation only changes"},
		{Text: "refactor", Note: "A code change that neither fixes a bug nor adds a feature"},
		{Text: "test", Note: "Adding missing tests or correcting existing tests"},
		{Text: "build", Note: "Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)"},
		{Text: "ci", Note: "Changes to our CI configuration files and scripts (examples: CircleCi, SauceLabs)"},
		{Text: "perf", Note: "A code change that improves performance"},
	}, config.CommitTypes())

	err := config.ReadConfig("./testdata/.gear.yml")
	require.Nil(t, err)

	require.Equal(t, []choose.Choice{
		{Text: "✨", Note: "feat: A new feature"},
		{Text: "🐛", Note: "fix: A bug fix"},
		{Text: "🔧", Note: "build: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)"},
		{Text: "📝", Note: "docs: Documentation only changes"},
		{Text: "🎨", Note: "refactor: A code change that neither fixes a bug nor adds a feature"},
		{Text: "🧪", Note: "test: Adding missing tests or correcting existing tests"},
		{Text: "👷", Note: "ci: Changes to our CI configuration files and scripts (examples: CircleCi, SauceLabs)"},
		{Text: "⚡️", Note: "perf: A code change that improves performance"},
	}, config.CommitTypes())
}
