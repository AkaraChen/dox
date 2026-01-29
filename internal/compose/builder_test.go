package compose

import (
	"path/filepath"
	"testing"

	"github.com/akrc/dox/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCommand_SimpleUp(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildUp([]string{})
	require.NoError(t, err)

	// Command structure: docker compose -f file1.yaml -f file2.yaml up
	assert.Equal(t, "docker", cmd[0])
	assert.Equal(t, "compose", cmd[1])
	// Find -f indices
	fIndex := -1
	for i, c := range cmd {
	 if c == "-f" {
   fIndex = i
   break
	 }
	}
	assert.Greater(t, fIndex, -1, "Should have -f flag")
	assert.Contains(t, cmd[fIndex+1], "compose.yaml")

	// Last argument should be "up"
	assert.Equal(t, "up", cmd[len(cmd)-1])
}

func TestBuildCommand_UpWithDetached(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildUp([]string{"-d"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "-d")
	assert.Contains(t, cmd, "up")
}

func TestBuildCommand_UpWithFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildUp([]string{"-d", "--build", "--force-recreate"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "-d")
	assert.Contains(t, cmd, "--build")
	assert.Contains(t, cmd, "--force-recreate")
	assert.Contains(t, cmd, "up")
}

func TestBuildCommand_UpWithProfile(t *testing.T) {
	fixtureDir := setupFixture(t, "with-profiles")

	// Need to load config for profile
	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "dev")
	cmd, err := b.BuildUp([]string{})
	require.NoError(t, err)

	// Should include compose.yaml and compose.dev.yaml
	fCount := countFlag(cmd, "-f")
	assert.Equal(t, 2, fCount)

	// Verify both files are present (using substring since paths are absolute)
	cmdStr := b.String(cmd)
	assert.Contains(t, cmdStr, "compose.yaml")
	assert.Contains(t, cmdStr, "compose.dev.yaml")
	assert.Equal(t, "up", cmd[len(cmd)-1])
}

func TestBuildCommand_UpWithProfileAndEnv(t *testing.T) {
	fixtureDir := setupFixture(t, "with-env")

	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "dev")
	cmd, err := b.BuildUp([]string{})
	require.NoError(t, err)

	// Should have --env-file flag
	assert.True(t, sliceContains(cmd, "--env-file"))
	assert.True(t, sliceContains(cmd, ".env.dev"))
}

func TestBuildCommand_Down(t *testing.T) {
	fixtureDir := setupFixture(t, "multi-slice")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildDown([]string{})
	require.NoError(t, err)

	assert.Contains(t, cmd, "docker")
	assert.Contains(t, cmd, "compose")
	assert.Contains(t, cmd, "down")
}

func TestBuildCommand_DownWithFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "multi-slice")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildDown([]string{"-v", "--remove-orphans"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "-v")
	assert.Contains(t, cmd, "--remove-orphans")
	assert.Contains(t, cmd, "down")
}

func TestBuildCommand_Ps(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildPs([]string{})
	require.NoError(t, err)

	assert.Contains(t, cmd, "ps")
}

func TestBuildCommand_Logs(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildLogs([]string{})
	require.NoError(t, err)

	assert.Contains(t, cmd, "logs")
}

func TestBuildCommand_LogsWithFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildLogs([]string{"-f", "--tail", "50"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "logs")
	assert.Contains(t, cmd, "-f")
	assert.Contains(t, cmd, "--tail")
	assert.Contains(t, cmd, "50")
}

func TestBuildCommand_LogsWithService(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildLogs([]string{"api"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "logs")
	assert.Contains(t, cmd, "api")
}

func TestBuildCommand_Restart(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildRestart([]string{"api"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "restart")
	assert.Contains(t, cmd, "api")
}

func TestBuildCommand_Exec(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildExec([]string{"api", "bash"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "exec")
	assert.Contains(t, cmd, "api")
	assert.Contains(t, cmd, "bash")
}

func TestBuildCommand_Build(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildBuild([]string{"api"})
	require.NoError(t, err)

	assert.Contains(t, cmd, "build")
	assert.Contains(t, cmd, "api")
}

func TestBuildCommand_String(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildUp([]string{"-d"})
	require.NoError(t, err)

	cmdStr := b.String(cmd)
	assert.Contains(t, cmdStr, "docker compose")
	assert.Contains(t, cmdStr, "-f")
	assert.Contains(t, cmdStr, "compose.yaml")
	assert.Contains(t, cmdStr, "up")
	assert.Contains(t, cmdStr, "-d")
}

func TestBuildCommand_Nuke(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	commands, err := b.BuildNuke()
	require.NoError(t, err)

	assert.Len(t, commands, 1)
	cmd := commands[0]
	assert.True(t, sliceContains(cmd, "down"))
	assert.True(t, sliceContains(cmd, "-v"))
	assert.True(t, sliceContains(cmd, "--remove-orphans"))
}

func TestBuildCommand_Fresh(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	commands, err := b.BuildFresh()
	require.NoError(t, err)

	assert.Len(t, commands, 2)
	// First: down -v
	assert.True(t, sliceContains(commands[0], "down"))
	assert.True(t, sliceContains(commands[0], "-v"))
	// Second: up --build
	assert.True(t, sliceContains(commands[1], "up"))
	assert.True(t, sliceContains(commands[1], "--build"))
}

func TestBuildCommand_Dup(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	commands, err := b.BuildDup()
	require.NoError(t, err)

	assert.Len(t, commands, 2)
	assert.True(t, sliceContains(commands[0], "down"))
	assert.True(t, sliceContains(commands[1], "up"))
}

func TestBuildCommand_FileOrder(t *testing.T) {
	fixtureDir := setupFixture(t, "multi-slice")

	// Without profile, should use only base
	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildUp([]string{})
	require.NoError(t, err)

	// Multi-slice auto-discovers all files
	fCount := countFlag(cmd, "-f")
	// Should have base + discovered slices
	assert.Greater(t, fCount, 0)
}

func TestBuildCommand_WithProfile_MultiSlice(t *testing.T) {
	fixtureDir := setupFixture(t, "with-profiles")

	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "full")
	cmd, err := b.BuildUp([]string{})
	require.NoError(t, err)

	// Should have base file + dev + prod
	fCount := countFlag(cmd, "-f")
	assert.Equal(t, 3, fCount)
}

func TestBuildCommand_NoComposeFiles(t *testing.T) {
	tmpDir := t.TempDir()

	b := NewBuilder(tmpDir, nil, "")
	_, err := b.BuildUp([]string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no compose files")
}

// Helper functions

func setupFixture(t *testing.T, name string) string {
	return filepath.Join("..", "..", "test", "fixtures", name)
}

func countFlag(slice []string, flag string) int {
	count := 0
	for i, s := range slice {
	 if s == flag && i < len(slice)-1 {
   count++
	 }
	}
	return count
}

func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
	 if s == item {
   return true
	 }
	}
	return false
}

func loadConfigFromDir(dir string) (*config.Config, string, error) {
	// Import config package locally to avoid circular import
	// This is a test helper so we'll inline the logic
	return config.LoadConfigFromDirectory(dir)
}
