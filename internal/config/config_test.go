package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverFiles_SimpleProject(t *testing.T) {
	// Setup: Use simple fixture
	fixtureDir := filepath.Join("..", "..", "test", "fixtures", "simple")

	d, err := DiscoverFiles(fixtureDir)
	require.NoError(t, err)
	assert.NotNil(t, d)

	// Should find base file only
	assert.Contains(t, d.BaseFile, "fixtures"+string(filepath.Separator)+"simple")
	assert.Contains(t, d.BaseFile, "compose.yaml")
	assert.Empty(t, d.Slices)
	assert.Len(t, d.Files, 1)
	assert.Equal(t, d.BaseFile, d.Files[0])
}

func TestDiscoverFiles_MultiSliceProject(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "test", "fixtures", "multi-slice")

	d, err := DiscoverFiles(fixtureDir)
	require.NoError(t, err)

	// Should find base file
	assert.Contains(t, d.BaseFile, "compose.yaml")

	// Should find all slices
	assert.Len(t, d.Slices, 4)
	assert.Contains(t, d.Slices, "dev")
	assert.Contains(t, d.Slices, "prod")
	assert.Contains(t, d.Slices, "db")
	assert.Contains(t, d.Slices, "redis")

	// Files should be ordered: base, then slices alphabetically
	assert.Len(t, d.Files, 5)
	assert.Equal(t, d.BaseFile, d.Files[0])

	// Check alphabetical order of slices (db, dev, prod, redis)
	assert.Contains(t, d.Files[1], "compose.db.yaml")
	assert.Contains(t, d.Files[2], "compose.dev.yaml")
	assert.Contains(t, d.Files[3], "compose.prod.yaml")
	assert.Contains(t, d.Files[4], "compose.redis.yaml")
}

func TestDiscoverFiles_NoComposeFiles(t *testing.T) {
	// Create temp dir with no compose files
	tmpDir := t.TempDir()

	d, err := DiscoverFiles(tmpDir)
	require.NoError(t, err)

	assert.Empty(t, d.BaseFile)
	assert.Empty(t, d.Slices)
	assert.Empty(t, d.Files)
}

func TestResolveProfile_SingleSlice(t *testing.T) {
	// Setup
	c := &Config{
	 Profiles: map[string]Profile{
   "dev": {
    Slices: []string{"dev"},
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "dev": "compose.dev.yaml",
	 },
	}

	files, envFile, err := c.ResolveProfile("dev", d)
	require.NoError(t, err)
	assert.Empty(t, envFile)

	// Should have base + dev slice
	assert.Len(t, files, 2)
	assert.Equal(t, "compose.yaml", files[0])
	assert.Equal(t, "compose.dev.yaml", files[1])
}

func TestResolveProfile_MultiSlice(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{
   "full": {
    Slices: []string{"dev", "db", "redis"},
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "dev":   "compose.dev.yaml",
   "db":    "compose.db.yaml",
   "redis": "compose.redis.yaml",
	 },
	}

	files, _, err := c.ResolveProfile("full", d)
	require.NoError(t, err)

	// Should have base + all slices in order
	assert.Len(t, files, 4)
	assert.Equal(t, "compose.yaml", files[0])

	// Slices should be in the order specified in profile
	assert.Equal(t, "compose.dev.yaml", files[1])
	assert.Equal(t, "compose.db.yaml", files[2])
	assert.Equal(t, "compose.redis.yaml", files[3])
}

func TestResolveProfile_ProfileNotFound(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{},
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	}

	_, _, err := c.ResolveProfile("missing", d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "profile 'missing' not found")
}

func TestResolveProfile_SliceNotFound(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{
   "dev": {
    Slices: []string{"nonexistent"},
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "dev": "compose.dev.yaml",
	 },
	}

	_, _, err := c.ResolveProfile("dev", d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "slice file 'compose.nonexistent.yaml' not found")
}

func TestResolveProfile_WithEnvFile(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{
   "dev": {
    Slices:  []string{"dev"},
    EnvFile: ".env.dev",
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "dev": "compose.dev.yaml",
	 },
	}

	files, envFile, err := c.ResolveProfile("dev", d)
	require.NoError(t, err)

	assert.Len(t, files, 2)
	assert.Equal(t, ".env.dev", envFile)
}

func TestResolveProfile_WithEnvReference(t *testing.T) {
	c := &Config{
	 EnvFiles: map[string]string{
	 "dev": ".env.dev",
	 },
	 Profiles: map[string]Profile{
   "dev": {
    Slices: []string{"dev"},
    Env:    "dev",
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "dev": "compose.dev.yaml",
	 },
	}

	files, envFile, err := c.ResolveProfile("dev", d)
	require.NoError(t, err)

	assert.Len(t, files, 2)
	assert.Equal(t, ".env.dev", envFile)
}

func TestResolveProfile_Inheritance(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{
   "base": {
    Slices: []string{"db"},
   },
   "full": {
    Slices:  []string{"dev"},
    Extends: "base",
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "db":  "compose.db.yaml",
   "dev": "compose.dev.yaml",
	 },
	}

	files, _, err := c.ResolveProfile("full", d)
	require.NoError(t, err)

	// Should have base + inherited db + dev
	assert.Len(t, files, 3)
	assert.Equal(t, "compose.yaml", files[0])
	assert.Equal(t, "compose.db.yaml", files[1])
	assert.Equal(t, "compose.dev.yaml", files[2])
}

func TestResolveProfile_CircularInheritance(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{
   "a": {Slices: []string{}, Extends: "b"},
   "b": {Slices: []string{}, Extends: "a"},
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	}

	_, _, err := c.ResolveProfile("a", d)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circular")
}

func TestResolveProfile_DeduplicateSlices(t *testing.T) {
	c := &Config{
	 Profiles: map[string]Profile{
   "base": {Slices: []string{"db"}},
   "full": {
    Slices:  []string{"db", "dev"},
    Extends: "base",
   },
	 },
	}

	d := &Discovery{
	 BaseFile: "compose.yaml",
	 Slices: map[string]string{
   "db":  "compose.db.yaml",
   "dev": "compose.dev.yaml",
	 },
	}

	files, _, err := c.ResolveProfile("full", d)
	require.NoError(t, err)

	// Should deduplicate db
	assert.Len(t, files, 3)
}

func TestDiscoverFiles_PrefersYamlOverYml(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both .yaml and .yml versions
	composeYaml := filepath.Join(tmpDir, "compose.yaml")
	composeYml := filepath.Join(tmpDir, "compose.yml")
	devYaml := filepath.Join(tmpDir, "compose.dev.yaml")
	devYml := filepath.Join(tmpDir, "compose.dev.yml")

	// Create .yml first
	err := os.WriteFile(composeYml, []byte("services: {}"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(devYml, []byte("services: {}"), 0644)
	require.NoError(t, err)

	// Then create .yaml
	err = os.WriteFile(composeYaml, []byte("services: {}"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(devYaml, []byte("services: {}"), 0644)
	require.NoError(t, err)

	d, err := DiscoverFiles(tmpDir)
	require.NoError(t, err)

	// Should prefer .yaml
	assert.Equal(t, composeYaml, d.BaseFile)
	assert.Equal(t, devYaml, d.Slices["dev"])
}

func TestGetDefaultProfile(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
	 {
		name: "with default profile",
		config: &Config{
   Defaults: Defaults{Profile: "dev"},
	 },
	 expected: "dev",
	 },
	 {
		name:     "no default profile",
	 config:   &Config{},
	 expected: "",
	 },
	}

	for _, tt := range tests {
	 t.Run(tt.name, func(t *testing.T) {
   result := tt.config.GetDefaultProfile()
   assert.Equal(t, tt.expected, result)
	 })
	}
}
