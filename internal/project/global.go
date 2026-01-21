package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// GlobalConfig represents the user's global do configuration
type GlobalConfig struct {
	Projects map[string]ProjectEntry `yaml:"projects,omitempty"`
	Aliases  map[string]string       `yaml:"aliases,omitempty"`
}

// ProjectEntry represents a project alias in the global config
type ProjectEntry struct {
	Path        string `yaml:"path"`
	Description string `yaml:"description,omitempty"`
}

// GetGlobalConfigPath returns the default path for the global config file
func GetGlobalConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}
	return filepath.Join(homeDir, ".config", "do", "config.yaml")
}

// LoadGlobalConfig loads the global configuration from the specified path.
// Returns nil, nil if the file doesn't exist (not an error).
func LoadGlobalConfig(path string) (*GlobalConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read global config: %w", err)
	}

	var cfg GlobalConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse global config: %w", err)
	}

	// Initialize maps if nil
	if cfg.Projects == nil {
		cfg.Projects = make(map[string]ProjectEntry)
	}
	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
	}

	return &cfg, nil
}

// LoadGlobalConfigOrDefault loads the global config or returns an empty default if not found
func LoadGlobalConfigOrDefault(path string) (*GlobalConfig, error) {
	cfg, err := LoadGlobalConfig(path)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return &GlobalConfig{
			Projects: make(map[string]ProjectEntry),
			Aliases:  make(map[string]string),
		}, nil
	}
	return cfg, nil
}

// ResolveProjectPath looks up a project by name and returns its path
func (c *GlobalConfig) ResolveProjectPath(name string) (string, bool) {
	if proj, ok := c.Projects[name]; ok {
		return proj.Path, true
	}
	return "", false
}

// GetAlias looks up a global alias by name
func (c *GlobalConfig) GetAlias(name string) (string, bool) {
	alias, ok := c.Aliases[name]
	return alias, ok
}

// HasProject checks if a project exists in the config
func (c *GlobalConfig) HasProject(name string) bool {
	_, ok := c.Projects[name]
	return ok
}

// ProjectNames returns a list of all project names
func (c *GlobalConfig) ProjectNames() []string {
	names := make([]string, 0, len(c.Projects))
	for name := range c.Projects {
		names = append(names, name)
	}
	return names
}

// AliasNames returns a list of all alias names
func (c *GlobalConfig) AliasNames() []string {
	names := make([]string, 0, len(c.Aliases))
	for name := range c.Aliases {
		names = append(names, name)
	}
	return names
}

// Regular expression for matching @project references
// Matches: @projectname or @projectname command...
var atProjectRegex = regexp.MustCompile(`^@([a-zA-Z0-9_-]+)(?:\s+(.*))?$`)

// ParseAtProjectReference parses a string that starts with @project
// Returns (isAtRef, projectName, restOfCommand)
func ParseAtProjectReference(input string) (bool, string, string) {
	matches := atProjectRegex.FindStringSubmatch(input)
	if matches == nil {
		return false, "", ""
	}

	projectName := matches[1]
	rest := ""
	if len(matches) > 2 && matches[2] != "" {
		rest = strings.TrimSpace(matches[2])
	}

	return true, projectName, rest
}

// IsAtProjectReference checks if a string is an @project reference
func IsAtProjectReference(input string) bool {
	isRef, _, _ := ParseAtProjectReference(input)
	return isRef
}
