# Project Context

## Purpose

A CLI wrapper tool for Docker Compose that simplifies multi-file compose setups. The `dox` command eliminates the need for repetitive `-f` flags when working with compose file slices for different environments (dev, prod, db, etc.).

## Tech Stack

- **Language**: Go 1.21+
- **CLI Framework**: spf13/cobra
- **Configuration**: spf13/viper, gopkg.in/yaml.v3
- **Target**: Single binary distribution (Linux, macOS, Windows)
- **Dependency**: Docker Compose V2 (external, not bundled)

## Project Conventions

### Code Style

- Go standard formatting (`gofmt`)
- Package names: lowercase, single word when possible
- Exported functions: Go naming conventions (PascalCase)
- Error handling: explicit, wrap with context using `fmt.Errorf`
- File naming: `snake_case.go` for non-exported, `PascalCase.go` may contain exported types

### Architecture Patterns

- Clean architecture: cmd (CLI layer), internal (business logic), pkg (reusable)
- Dependency injection via interfaces where beneficial
- Configuration-driven behavior (dox.yaml)
- Explicit is better than implicit (fallbacks with clear messages)

### Testing Strategy

- Unit tests for config parsing and discovery logic
- Integration tests for command building (with docker compose mock)
- Manual testing for actual docker compose interaction

### Git Workflow

- Feature branches: `feature/feature-name`
- Commit format: Conventional Commits (feat:, fix:, refactor:)
- PR description required for changes
- Auto-update openspec on merge (if applicable)

## Domain Context

Docker Compose uses multiple files for environment-specific configuration:
- Base file: `compose.yaml` or `docker-compose.yaml` (common services)
- Slice files: `compose.dev.yaml`, `compose.prod.yaml` (environment overrides)
- Override file: `compose.override.yaml` (local development)

The `do` tool auto-discovers these files and builds the appropriate docker compose command.

## Important Constraints

- **Must not replace Docker Compose**: wrapper only, calls external binary
- **Cross-platform**: Support Linux, macOS, Windows (path handling differences)
- **Backward compatible**: work with Docker Compose V2+ only
- **No state**: don't persist running state (except optional profile cache)
- **Exit codes**: propagate docker compose exit codes

## External Dependencies

- **Docker Compose V2**: must be installed and in PATH
- **Git**: optional, for auto-detect feature
- **YAML files**: project-local dox.yaml, compose.*.yaml files
