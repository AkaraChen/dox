package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Discovery holds auto-discovered compose files
type Discovery struct {
	BaseFile string
	Slices   map[string]string // slice name -> file path
	Files    []string          // ordered list of all files
}

// Config is the main dox.yaml configuration
type Config struct {
	Version    int                    `yaml:"version"`
	Discovery  DiscoveryConfig        `yaml:"discovery"`
	Profiles   map[string]Profile     `yaml:"profiles"`
	EnvFiles   map[string]string      `yaml:"env_files"`
	Defaults   Defaults               `yaml:"defaults"`
	Aliases    map[string]string      `yaml:"aliases"`
	Hooks      map[string][]string    `yaml:"hooks"`
}

// DiscoveryConfig configures auto-discovery behavior
type DiscoveryConfig struct {
	Enabled bool   `yaml:"enabled"`
	Pattern string `yaml:"pattern"`
	Base    string `yaml:"base"`
}

// Profile defines a set of compose slices
type Profile struct {
	Slices   []string `yaml:"slices"`
	EnvFile  string   `yaml:"env_file"`
	Env      string   `yaml:"env"`
	Extends  string   `yaml:"extends"`
}

// Defaults defines default behavior
type Defaults struct {
	Profile     string `yaml:"profile"`
	Slice       string `yaml:"slice"`
	AutoDetect  bool   `yaml:"auto_detect"`
}

// DiscoverFiles finds compose files in the current directory
func DiscoverFiles(dir string) (*Discovery, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	d := &Discovery{
		Slices: make(map[string]string),
	}

	// Base files to look for (in order of preference)
	baseFiles := []string{"compose.yaml", "docker-compose.yaml", "compose.yml", "docker-compose.yml"}

	for _, base := range baseFiles {
	 fullPath := filepath.Join(dir, base)
	 if _, err := os.Stat(fullPath); err == nil {
   d.BaseFile = fullPath
   break
	 }
	}

	// Find slice files (compose.*.yaml or compose.*.yml)
	for _, entry := range entries {
	 if entry.IsDir() {
   continue
	 }

	 name := entry.Name()
	 var sliceName string

	 // Prefer .yaml over .yml
	 if strings.HasPrefix(name, "compose.") && strings.HasSuffix(name, ".yaml") {
   sliceName = strings.TrimPrefix(name, "compose.")
   sliceName = strings.TrimSuffix(sliceName, ".yaml")
   // Skip if it's the base file
   if sliceName == "" || name == "compose.yaml" {
    continue
   }
   d.Slices[sliceName] = filepath.Join(dir, name)
	 } else if strings.HasPrefix(name, "compose.") && strings.HasSuffix(name, ".yml") {
   // Only use .yml if .yaml doesn't exist
   sliceName = strings.TrimPrefix(name, "compose.")
   sliceName = strings.TrimSuffix(sliceName, ".yml")
   if sliceName == "" || name == "compose.yml" {
    continue
   }
   // Check if .yaml version exists
   yamlExists := false
   for _, existingName := range d.Slices {
    if strings.HasSuffix(existingName, ".yaml") {
     yamlExists = true
     break
    }
   }
   if !yamlExists {
    d.Slices[sliceName] = filepath.Join(dir, name)
   }
	 }
	}

	// Build ordered file list
	if d.BaseFile != "" {
	 d.Files = append(d.Files, d.BaseFile)
	}

	// Sort slice names for deterministic ordering
	sliceNames := make([]string, 0, len(d.Slices))
	for name := range d.Slices {
	 sliceNames = append(sliceNames, name)
	}
	sort.Strings(sliceNames)

	for _, name := range sliceNames {
	 d.Files = append(d.Files, d.Slices[name])
	}

	return d, nil
}

// ResolveProfile resolves a profile to a list of compose files
func (c *Config) ResolveProfile(profileName string, discovery *Discovery) ([]string, string, error) {
	profile, exists := c.Profiles[profileName]
	if !exists {
	 return nil, "", fmt.Errorf("profile '%s' not found", profileName)
	}

	// Handle inheritance
	slices := append([]string{}, profile.Slices...)
	currentExtends := profile.Extends
	visited := map[string]bool{profileName: true}

	for currentExtends != "" {
	 if visited[currentExtends] {
   return nil, "", fmt.Errorf("circular profile inheritance detected")
	 }
	 visited[currentExtends] = true

	 parentProfile, exists := c.Profiles[currentExtends]
	 if !exists {
   return nil, "", fmt.Errorf("profile '%s' extends non-existent profile '%s'", profileName, currentExtends)
	 }
	 // Prepend parent slices
	 slices = append(parentProfile.Slices, slices...)
	 currentExtends = parentProfile.Extends
	}

	files := []string{}
	if discovery.BaseFile != "" {
	 files = append(files, discovery.BaseFile)
	}

	// Resolve slices to files, de-duplicating
	seen := map[string]bool{}
	for _, sliceName := range slices {
	 if seen[sliceName] {
   continue
	 }
	 seen[sliceName] = true

	 sliceFile, exists := discovery.Slices[sliceName]
	 if !exists {
   return nil, "", fmt.Errorf("slice file 'compose.%s.yaml' not found for profile '%s'", sliceName, profileName)
	 }
	 files = append(files, sliceFile)
	}

	// Resolve env file
	envFile := ""
	if profile.EnvFile != "" {
	 envFile = profile.EnvFile
	} else if profile.Env != "" {
	 if env, ok := c.EnvFiles[profile.Env]; ok {
   envFile = env
	 }
	}

	return files, envFile, nil
}

// GetDefaultProfile returns the default profile name
func (c *Config) GetDefaultProfile() string {
	if c.Defaults.Profile != "" {
	 return c.Defaults.Profile
	}
	return ""
}
