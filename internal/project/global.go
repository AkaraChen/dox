package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// GlobalConfig represents the user's global dox configuration
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
	return filepath.Join(homeDir, ".config", "dox", "config.yaml")
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

// RemoteProject represents a resolved remote project reference
type RemoteProject struct {
	ProjectName     string
	ProjectPath     string
	RemainingCommand string
}

// ResolveRemoteProject parses and resolves an @project reference
// Returns nil if input is not an @project reference
// Returns error if project is not found in config
func (c *GlobalConfig) ResolveRemoteProject(input string) (*RemoteProject, error) {
	isRef, projectName, remainingCmd := ParseAtProjectReference(input)
	if !isRef {
		return nil, nil
	}

	projectPath, exists := c.ResolveProjectPath(projectName)
	if !exists {
		return nil, fmt.Errorf("project '%s' not found in global config", projectName)
	}

	return &RemoteProject{
		ProjectName:      projectName,
		ProjectPath:      projectPath,
		RemainingCommand: remainingCmd,
	}, nil
}

// ProjectInfo represents a project with its metadata
type ProjectInfo struct {
	Name        string
	Path        string
	Description string
}

// ListProjects returns all projects with their metadata
func (c *GlobalConfig) ListProjects() []ProjectInfo {
	projects := make([]ProjectInfo, 0, len(c.Projects))
	for name, entry := range c.Projects {
		projects = append(projects, ProjectInfo{
			Name:        name,
			Path:        entry.Path,
			Description: entry.Description,
		})
	}
	return projects
}
