package compose

import (
	"path/filepath"
	"testing"

	configpkg "github.com/AkaraChen/dox/internal/config"
)

func BenchmarkBuildCommand_SimpleUp(b *testing.B) {
	fixtureDir := setupTestFixture(b, "simple")
	builder := NewBuilder(fixtureDir, nil, "")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = builder.BuildUp([]string{})
	}
}

func BenchmarkBuildCommand_MultiSliceUp(b *testing.B) {
	fixtureDir := setupTestFixture(b, "multi-slice")
	builder := NewBuilder(fixtureDir, nil, "")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = builder.BuildUp([]string{})
	}
}

func BenchmarkBuildCommand_WithProfile(b *testing.B) {
	fixtureDir := setupTestFixture(b, "with-profiles")
	cfg, _, _ := configpkg.LoadConfigFromDirectory(fixtureDir)
	builder := NewBuilder(fixtureDir, cfg, "dev")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = builder.BuildUp([]string{})
	}
}

func BenchmarkBuildCommand_Down(b *testing.B) {
	fixtureDir := setupTestFixture(b, "simple")
	builder := NewBuilder(fixtureDir, nil, "")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = builder.BuildDown([]string{"-v"})
	}
}

func BenchmarkBuildCommand_Logs(b *testing.B) {
	fixtureDir := setupTestFixture(b, "simple")
	builder := NewBuilder(fixtureDir, nil, "")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = builder.BuildLogs([]string{"-f", "--tail", "100"})
	}
}

func BenchmarkDiscoverFiles_Simple(b *testing.B) {
	fixtureDir := setupTestFixture(b, "simple")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = configpkg.DiscoverFiles(fixtureDir)
	}
}

func BenchmarkDiscoverFiles_MultiSlice(b *testing.B) {
	fixtureDir := setupTestFixture(b, "multi-slice")

	b.ResetTimer()
	for b.Loop() {
	 _, _ = configpkg.DiscoverFiles(fixtureDir)
	}
}

func BenchmarkResolveProfile_SingleSlice(b *testing.B) {
	fixtureDir := setupTestFixture(b, "with-profiles")
	cfg, _, _ := configpkg.LoadConfigFromDirectory(fixtureDir)
	discovery, _ := configpkg.DiscoverFiles(fixtureDir)

	b.ResetTimer()
	for b.Loop() {
	 _, _, _ = cfg.ResolveProfile("dev", discovery)
	}
}

func BenchmarkResolveProfile_MultiSlice(b *testing.B) {
	fixtureDir := setupTestFixture(b, "with-profiles")
	cfg, _, _ := configpkg.LoadConfigFromDirectory(fixtureDir)
	discovery, _ := configpkg.DiscoverFiles(fixtureDir)

	b.ResetTimer()
	for b.Loop() {
	 _, _, _ = cfg.ResolveProfile("full", discovery)
	}
}

func BenchmarkFormatCommand(b *testing.B) {
	cmd := []string{"docker", "compose", "-f", "compose.yaml", "-f", "compose.dev.yaml", "up", "-d"}

	b.ResetTimer()
	for b.Loop() {
	 _ = FormatCommand(cmd)
	}
}

func BenchmarkStringCommand(b *testing.B) {
	builder := NewBuilder("", nil, "")
	cmd := []string{"docker", "compose", "up"}

	b.ResetTimer()
	for b.Loop() {
	 _ = builder.String(cmd)
	}
}

// Helper setup function for benchmarks
func setupTestFixture(b *testing.B, name string) string {
	b.Helper()
	return filepath.Join("..", "..", "test", "fixtures", name)
}
