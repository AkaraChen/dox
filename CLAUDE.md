# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

You are a manager and an agent orchestrator. You must never perform implementations yourself; instead, delegate everything to sub-agents or task agents. Break down tasks into ultra-granular steps and establish a PDCA cycle.

## Common Development Commands

### Building
```bash
make build          # Build binary for current platform to build/dox
make build-all      # Build for darwin/amd64/arm64, linux/amd64, windows/amd64
make install        # Install to $GOPATH/bin
```

### Testing
```bash
make test           # Run all tests with coverage (includes race detection)
make test-short     # Run tests without race detection (faster)
make test-ci        # CI tests with 80% coverage threshold
```

### Code Quality
```bash
make lint           # Run golangci-lint
make fmt            # Format code with gofmt
make deps           # Download and tidy dependencies
make clean          # Remove build artifacts and coverage files
```

### Running Single Tests
```bash
go test -v -race ./internal/config           # Test specific package
go test -v -run TestDiscoverFiles ./...      # Run specific test
go test -bench=. ./...                       # Run benchmarks
```

## Architecture Overview

This is a Go CLI wrapper for Docker Compose V2 that simplifies multi-file compose setups.

### Layer Structure (Clean Architecture)

```
cmd/ (CLI Layer)
  ├── root.go          # Root command, global flags (--verbose, --dry-run, --profile)
  ├── compose.go       # Main command group (builder + executor)
  ├── up.go, down.go   # Individual compose commands
  └── convenience.go   # dup, nuke, fresh commands
       ↓
internal/ (Business Logic)
  ├── config/          # Configuration: profiles, aliases, hooks, file discovery
  ├── compose/         # Docker operations: builder (command construction), executor (execution)
  └── project/         # Multi-project: global config, aliases, history
       ↓
pkg/ (Reusable components)
```

### Key Data Flow

```
User Input (dox c up -d)
  ↓
cmd/root.go (global flags)
  ↓
cmd/compose.go (getComposeBuilder → getComposeExecutor)
  ↓
internal/compose/builder.go (BuildUp → buildBase → resolveFiles)
  ↓
internal/config/ (DiscoverFiles → ResolveProfile)
  ↓
internal/compose/executor.go (RunCommand → docker compose)
```

### Configuration Discovery Priority

1. Explicit `-f` flags (highest priority)
2. `dox.yaml` profile configuration
3. Auto-discovery: `compose.yaml` + `compose.*.yaml` (alphabetical)

### Core Components

**Config Package** (`internal/config/`):
- `Config` struct: profiles, aliases, hooks, defaults
- `DiscoverFiles()`: Auto-discovers compose files
- `ResolveProfile()`: Merges slices with inheritance
- Test coverage: 93.9%

**Compose Package** (`internal/compose/`):
- `Builder`: Constructs docker compose commands (`BuildUp()`, `BuildDown()`, `BuildPs()`, `BuildLogs()`)
- `Executor`: Executes commands with dry-run mode (`RunCommand()`, `RunCommands()`)
- Test coverage: 79.3%

**Project Package** (`internal/project/`):
- `GlobalConfig`: Project aliases and global aliases
- `History`: Command execution tracking
- `RemoteProject`: @project syntax resolution
- Test coverage: 93.9%

### Test Fixtures

Located in `test/fixtures/` with 16+ scenarios:
- `simple/` - Basic compose.yaml
- `multi-slice/` - Multiple compose files
- `with-profiles/` - Profile configurations
- `with-aliases/` - Custom aliases
- `with-hooks/` - Pre/post hooks
- `edge-cases/` - Error scenarios

### Project Conventions

- Code style: `gofmt`, lowercase package names, explicit error handling
- Architecture: Clean architecture with cmd → internal → pkg separation
- Testing: TDD with >80% coverage goal, race detection enabled
- Git: Conventional commits (feat:, fix:, refactor:)
- External dependency: Docker Compose V2 (must be installed separately)

### Important Files

- `openspec/project.md` - Tech stack, conventions, constraints
- `openspec/AGENTS.md` - Spec-driven development workflow
- `dox.yaml` (project-local) - Profiles, slices, aliases, hooks
- `~/.config/dox/config.yaml` - Global project aliases
- `~/.cache/dox/history.yaml` - Command execution history
