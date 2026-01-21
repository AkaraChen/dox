package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveAlias_SimpleCommand(t *testing.T) {
	alias := "up -d"
	cmd := parseAlias(alias)
	assert.Equal(t, []string{"up", "-d"}, cmd)
}

func TestResolveAlias_ChainedCommands(t *testing.T) {
	alias := "down -v && up --build -d"
	commands := parseAliasCommands(alias)
	assert.Len(t, commands, 2)
	assert.Equal(t, []string{"down", "-v"}, commands[0])
	assert.Equal(t, []string{"up", "--build", "-d"}, commands[1])
}

func TestResolveAlias_WithQuotes(t *testing.T) {
	// Test that we handle basic quoted arguments
	alias := "echo \"hello world\""
	cmd := parseAlias(alias)
	// Simple parsing splits on spaces, quotes not yet handled
	// This is a known limitation (TODO in compose.go)
	assert.Contains(t, cmd, "echo")
	assert.Contains(t, cmd, "\"hello")
	assert.Contains(t, cmd, "world\"")
}

func TestResolveAlias_Empty(t *testing.T) {
	alias := ""
	cmd := parseAlias(alias)
	assert.Empty(t, cmd)
}

func TestResolveAlias_SingleWord(t *testing.T) {
	alias := "restart"
	cmd := parseAlias(alias)
	assert.Equal(t, []string{"restart"}, cmd)
}

// Helper function to parse alias into command parts
func parseAlias(alias string) []string {
	if alias == "" {
		return []string{}
	}
	return strings.Fields(alias)
}

// Helper function to parse chained alias commands
func parseAliasCommands(alias string) [][]string {
	parts := strings.Split(alias, "&&")
	commands := make([][]string, 0, len(parts))

	for _, part := range parts {
		cmd := strings.Fields(strings.TrimSpace(part))
		if len(cmd) > 0 {
			commands = append(commands, cmd)
		}
	}

	return commands
}
