package compose

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCommand_WithAlias(t *testing.T) {
	fixtureDir := setupFixture(t, "with-aliases")

	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	// Parse the alias "fresh: down -v && up --build -d"
	aliasDef := cfg.Aliases["fresh"]
	assert.Equal(t, "down -v && up --build -d", aliasDef)

	// Simulate what the alias parser would do
	commands := strings.Split(aliasDef, "&&")
	assert.Len(t, commands, 2)
}

func TestBuildCommand_AliasFresh(t *testing.T) {
	fixtureDir := setupFixture(t, "with-aliases")

	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "")

	// The "fresh" alias expands to "down -v && up --build -d"
	// First command: down -v
	downCmd, err := b.BuildDown([]string{"-v"})
	require.NoError(t, err)
	assert.True(t, sliceContains(downCmd, "down"))
	assert.True(t, sliceContains(downCmd, "-v"))

	// Second command: up --build -d
	upCmd, err := b.BuildUp([]string{"--build", "-d"})
	require.NoError(t, err)
	assert.True(t, sliceContains(upCmd, "up"))
	assert.True(t, sliceContains(upCmd, "--build"))
	assert.True(t, sliceContains(upCmd, "-d"))
}

func TestBuildCommand_AliasRestartAll(t *testing.T) {
	fixtureDir := setupFixture(t, "with-aliases")

	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "")

	// The "restart-all" alias expands to "down && up -d"
	// First command: down
	downCmd, err := b.BuildDown([]string{})
	require.NoError(t, err)
	assert.True(t, sliceContains(downCmd, "down"))

	// Second command: up -d
	upCmd, err := b.BuildUp([]string{"-d"})
	require.NoError(t, err)
	assert.True(t, sliceContains(upCmd, "up"))
	assert.True(t, sliceContains(upCmd, "-d"))
}

func TestParseAlias_SimpleCommand(t *testing.T) {
	alias := "up -d"
	parts := parseAlias(alias)
	assert.Equal(t, []string{"up", "-d"}, parts)
}

func TestParseAlias_ChainedCommands(t *testing.T) {
	alias := "down -v && up --build"
	parts := parseAlias(alias)
	// Note: parseAlias should preserve the && for later splitting
	assert.Contains(t, parts, "down")
	assert.Contains(t, parts, "-v")
	assert.Contains(t, parts, "&&")
	assert.Contains(t, parts, "up")
	assert.Contains(t, parts, "--build")
}

// Helper function to parse alias into command parts
func parseAlias(alias string) []string {
	return strings.Fields(alias)
}
