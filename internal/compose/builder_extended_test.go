package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildDown_AllFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	tests := []struct {
		name string
		args []string
	}{
		{"with volumes", []string{"-v"}},
		{"with remove orphans", []string{"--remove-orphans"}},
		{"with both", []string{"-v", "--remove-orphans"}},
		{"with timeout", []string{"-t", "10"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := b.BuildDown(tt.args)
			require.NoError(t, err)
			assert.Contains(t, cmd, "down")
		})
	}
}

func TestBuildPs_WithServiceArgs(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	cmd, err := b.BuildPs([]string{"web", "db"})
	require.NoError(t, err)
	assert.True(t, sliceContains(cmd, "ps"))
	assert.True(t, sliceContains(cmd, "web"))
	assert.True(t, sliceContains(cmd, "db"))
}

func TestBuildPs_AllFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	tests := []struct {
		name string
		args []string
	}{
		{"quiet", []string{"-q"}},
		{"services", []string{"--services"}},
		{"filter status", []string{"--filter", "status=running"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := b.BuildPs(tt.args)
			require.NoError(t, err)
			assert.Contains(t, cmd, "ps")
		})
	}
}

func TestBuildLogs_AllFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	tests := []struct {
		name string
		args []string
	}{
		{"follow", []string{"-f"}},
		{"tail", []string{"--tail", "100"}},
		{"since", []string{"--since", "1h"}},
		{"until", []string{"--until", "2024-01-01"}},
		{"with service", []string{"web"}},
		{"combined", []string{"-f", "--tail", "50", "web"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := b.BuildLogs(tt.args)
			require.NoError(t, err)
			assert.True(t, sliceContains(cmd, "logs"))
		})
	}
}

func TestBuildRestart_WithServices(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	cmd, err := b.BuildRestart([]string{"web", "db"})
	require.NoError(t, err)
	assert.True(t, sliceContains(cmd, "restart"))
	assert.True(t, sliceContains(cmd, "web"))
	assert.True(t, sliceContains(cmd, "db"))
}

func TestBuildRestart_WithFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	cmd, err := b.BuildRestart([]string{"-t", "10", "web"})
	require.NoError(t, err)
	assert.True(t, sliceContains(cmd, "restart"))
	assert.True(t, sliceContains(cmd, "-t"))
	assert.True(t, sliceContains(cmd, "10"))
	assert.True(t, sliceContains(cmd, "web"))
}

func TestBuildExec_Variations(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{"simple", []string{"web", "ls"}, []string{"exec", "web", "ls"}},
		{"interactive", []string{"-it", "web", "bash"}, []string{"exec", "-it", "web", "bash"}},
		{"with env", []string{"-e", "DEBUG=1", "web", "sh"}, []string{"exec", "-e", "DEBUG=1", "web", "sh"}},
		{"privileged", []string{"--privileged", "web", "cmd"}, []string{"exec", "--privileged", "web", "cmd"}},
		{"user", []string{"--user", "root", "web", "id"}, []string{"exec", "--user", "root", "web", "id"}},
		{"work dir", []string{"-w", "/tmp", "web", "pwd"}, []string{"exec", "-w", "/tmp", "web", "pwd"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := b.BuildExec(tt.args)
			require.NoError(t, err)
			for _, c := range tt.contains {
				assert.True(t, sliceContains(cmd, c))
			}
		})
	}
}

func TestBuildBuild_Variations(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	tests := []struct {
		name string
		args []string
	}{
		{"simple service", []string{"web"}},
		{"multiple services", []string{"web", "api"}},
		{"with quiet", []string{"-q", "web"}},
		{"with no cache", []string{"--no-cache", "web"}},
		{"with build args", []string{"--build-arg", "VERSION=1.0", "web"}},
		{"with parallel", []string{"--parallel", "web"}},
		{"with progress", []string{"--progress", "plain", "web"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := b.BuildBuild(tt.args)
			require.NoError(t, err)
			assert.True(t, sliceContains(cmd, "build"))
		})
	}
}

func TestBuildNuke_ReturnsCommands(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	commands, err := b.BuildNuke()
	require.NoError(t, err)
	assert.Len(t, commands, 1)

	// Should contain down with -v and --remove-orphans
	assert.True(t, sliceContains(commands[0], "down"))
	assert.True(t, sliceContains(commands[0], "-v"))
	assert.True(t, sliceContains(commands[0], "--remove-orphans"))
}

func TestBuildFresh_ReturnsCommands(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	commands, err := b.BuildFresh()
	require.NoError(t, err)
	assert.Len(t, commands, 2)

	// First command should be down -v
	assert.True(t, sliceContains(commands[0], "down"))
	assert.True(t, sliceContains(commands[0], "-v"))

	// Second command should be up --build
	assert.True(t, sliceContains(commands[1], "up"))
	assert.True(t, sliceContains(commands[1], "--build"))
}

func TestBuildDup_ReturnsCommands(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	commands, err := b.BuildDup()
	require.NoError(t, err)
	assert.Len(t, commands, 2)

	// First should be down
	assert.True(t, sliceContains(commands[0], "down"))

	// Second should be up
	assert.True(t, sliceContains(commands[1], "up"))
}

func TestBuildStatus_AllFlags(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	b := NewBuilder(fixtureDir, nil, "")

	tests := []struct {
		name string
		args []string
	}{
		{"watch", []string{"--watch"}},
		{"format json", []string{"--format", "json"}},
		{"all", []string{"--all"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := b.BuildStatus(tt.args)
			require.NoError(t, err)
			assert.True(t, sliceContains(cmd, "ps"))
		})
	}
}

func TestBuilder_String(t *testing.T) {
	cmd := []string{"docker", "compose", "-f", "test.yaml", "up"}
	b := &Builder{}

	result := b.String(cmd)
	assert.Equal(t, "docker compose -f test.yaml up", result)
}

func TestFormatCommand_Escaping(t *testing.T) {
	cmd := []string{"docker", "compose", "exec", "web", "echo", "hello world"}
	result := FormatCommand(cmd)
	assert.Contains(t, result, "docker compose exec web")
}

func TestNewBuilder_WithConfig(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")
	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "dev")
	assert.NotNil(t, b)
}

func TestNewBuilder_WithEnvFile(t *testing.T) {
	fixtureDir := setupFixture(t, "with-env")
	cfg, _, err := loadConfigFromDir(fixtureDir)
	require.NoError(t, err)

	b := NewBuilder(fixtureDir, cfg, "dev")
	assert.NotNil(t, b)
}
