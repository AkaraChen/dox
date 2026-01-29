package compose

import (
	"fmt"
	"strings"

	"github.com/akrc/dox/internal/config"
)

// Builder builds docker compose commands
type Builder struct {
	dir        string
	config     *config.Config
	profile    string
	discovery  *config.Discovery
}

// NewBuilder creates a new command builder
func NewBuilder(dir string, cfg *config.Config, profile string) *Builder {
	b := &Builder{
	 dir:     dir,
	 config:  cfg,
	 profile: profile,
	}

	// Auto-discover compose files
	if discovery, err := config.DiscoverFiles(dir); err == nil {
	 b.discovery = discovery
	}

	return b
}

// resolveFiles resolves the compose files to use based on profile
func (b *Builder) resolveFiles() ([]string, error) {
	// If profile specified, use it
	if b.profile != "" && b.config != nil {
	 files, _, err := b.config.ResolveProfile(b.profile, b.discovery)
	 if err != nil {
   return nil, err
	 }
	 return files, nil
	}

	// Use default profile if configured
	if b.config != nil {
	 defaultProfile := b.config.GetDefaultProfile()
	 if defaultProfile != "" {
   files, _, err := b.config.ResolveProfile(defaultProfile, b.discovery)
   if err != nil {
    return nil, err
   }
   return files, nil
	 }
	}

	// Use auto-discovered files
	if b.discovery != nil && len(b.discovery.Files) > 0 {
	 return b.discovery.Files, nil
	}

	return nil, fmt.Errorf("no compose files found in %s", b.dir)
}

// resolveEnvFile returns the env file for the current profile
func (b *Builder) resolveEnvFile() string {
	if b.profile != "" && b.config != nil {
	 _, envFile, _ := b.config.ResolveProfile(b.profile, b.discovery)
	 return envFile
	}

	if b.config != nil {
	 defaultProfile := b.config.GetDefaultProfile()
	 if defaultProfile != "" {
   _, envFile, _ := b.config.ResolveProfile(defaultProfile, b.discovery)
   return envFile
	 }
	}

	return ""
}

// buildBase builds the base command with -f flags
func (b *Builder) buildBase() ([]string, error) {
	files, err := b.resolveFiles()
	if err != nil {
	 return nil, err
	}

	cmd := []string{"docker", "compose"}
	for _, file := range files {
	 cmd = append(cmd, "-f", file)
	}

	// Add env file if specified
	envFile := b.resolveEnvFile()
	if envFile != "" {
	 cmd = append(cmd, "--env-file", envFile)
	}

	return cmd, nil
}

// BuildUp builds the docker compose up command
func (b *Builder) BuildUp(args []string) ([]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "up")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildDown builds the docker compose down command
func (b *Builder) BuildDown(args []string) ([]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "down")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildPs builds the docker compose ps command
func (b *Builder) BuildPs(args []string) ([]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "ps")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildLogs builds the docker compose logs command
func (b *Builder) BuildLogs(args []string) ([]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "logs")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildRestart builds the docker compose restart command
func (b *Builder) BuildRestart(args []string) ([]string, error) {
	if len(args) == 0 {
	 return nil, fmt.Errorf("restart requires at least one service name")
	}

	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "restart")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildExec builds the docker compose exec command
func (b *Builder) BuildExec(args []string) ([]string, error) {
	if len(args) == 0 {
	 return nil, fmt.Errorf("exec requires at least a service name")
	}

	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "exec")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildBuild builds the docker compose build command
func (b *Builder) BuildBuild(args []string) ([]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "build")
	cmd = append(cmd, args...)
	return cmd, nil
}

// BuildNuke builds the nuke command (down -v --remove-orphans)
func (b *Builder) BuildNuke() ([][]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "down", "-v", "--remove-orphans")
	return [][]string{cmd}, nil
}

// BuildFresh builds the fresh command (down -v && up --build)
func (b *Builder) BuildFresh() ([][]string, error) {
	downCmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}
	downCmd = append(downCmd, "down", "-v")

	upCmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}
	upCmd = append(upCmd, "up", "--build")

	return [][]string{downCmd, upCmd}, nil
}

// BuildDup builds the dup command (down && up)
func (b *Builder) BuildDup() ([][]string, error) {
	downCmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}
	downCmd = append(downCmd, "down")

	upCmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}
	upCmd = append(upCmd, "up")

	return [][]string{downCmd, upCmd}, nil
}

// BuildStatus builds the status command (enhanced ps)
func (b *Builder) BuildStatus(args []string) ([]string, error) {
	cmd, err := b.buildBase()
	if err != nil {
	 return nil, err
	}

	cmd = append(cmd, "ps")
	cmd = append(cmd, args...)
	return cmd, nil
}

// String converts a command slice to a string
func (b *Builder) String(cmd []string) string {
	return strings.Join(cmd, " ")
}

// FormatCommand formats a command for display
func FormatCommand(cmd []string) string {
	return strings.Join(cmd, " ")
}
