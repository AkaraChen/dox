# Design: Docker Compose CLI Tool

## Context

Docker Compose projects often span multiple files for different environments (dev, prod, db, etc.). The native `docker compose -f file1.yaml -f file2.yaml ...` syntax becomes unwieldy with 3+ files. This tool aims to simplify the developer experience while maintaining full compatibility with Docker Compose.

## Goals / Non-Goals

**Goals:**
- Reduce typing for common docker compose operations
- Provide project-local configuration (do.yaml)
- Support profile-based compose file combinations
- Enable dry-run preview of commands
- Work with existing Docker Compose installations (no replacement)

**Non-Goals:**
- Replacing Docker Compose (wrapper only, not a reimplementation)
- Managing Docker daemon or containers directly
- Cross-machine orchestration
- GUI or web interface

## Decisions

### Language: Go with Cobra
- **What**: Go implementation using spf13/cobra CLI framework
- **Why**: Single binary distribution, fast execution, excellent CLI library ecosystem
- **Alternatives considered**:
  - Rust (also good, but Cobra is more mature for CLI apps)
  - Python/bash script (simpler but requires runtime, harder to distribute)
  - Docker Compose plugin (limited by compose's extension API)

### Configuration Format: YAML
- **What**: Project-local `do.yaml` for profiles, slices, aliases
- **Why**: Familiar to Docker users, easy to parse, human-editable
- **Alternatives considered**: TOML (less common in Docker ecosystem), JSON (not human-friendly)

### Auto-Discovery Strategy
- **What**: Scan current directory for `compose.yaml` + `compose.*.yaml` patterns
- **Why**: Convention over configuration - zero config for simple cases
- **Fallback**: Explicit file list in do.yaml when auto-detection insufficient

### Command Structure: `do c <subcommand>`
- **What**: `c` as shorthand for compose, followed by familiar docker compose verbs
- **Why**: Keeps namespace open for future tools, mirrors docker compose mental model
- **Alternatives considered**: `do compose` (verbose), `do up` (no namespace for other tools)

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  Root (do) → Compose (c) → Commands (up, down, ps, logs...)  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Config Layer                            │
│  do.yaml parser + Auto-discovery + Profile resolver         │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Command Builder                           │
│  Assembles: docker compose -f <base> -f <slice1> ... <cmd>  │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Executor Layer                            │
│  os/exec for docker compose, output formatting, dry-run     │
└─────────────────────────────────────────────────────────────┘
```

### Package Structure

```
do/
├── cmd/
│   └── root.go          # Cobra root command
│   └── compose.go       # Compose command group
├── internal/
│   ├── config/
│   │   ├── parser.go    # do.yaml parsing
│   │   ├── discovery.go # Auto-discover compose files
│   │   └── profile.go   # Profile resolution
│   ├── compose/
│   │   ├── builder.go   # Build docker compose commands
│   │   └── executor.go  # Execute commands
│   └── project/
│       └── alias.go     # Project alias system
├── pkg/
│   └── doctl/           # Public API (optional)
└── go.mod
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Docker Compose version compatibility | Validate compose binary, feature detection |
| Shell injection from user input | Use safe exec API, validate all inputs |
| Config file conflicts | Clear precedence: CLI flags > do.yaml > auto-discovery |
| Breaking changes in Docker Compose | Pin minimum version, update on major releases |
| Complex config scenarios | Always allow explicit `-f` override |

## Migration Plan

### Phase 1: Core Wrapper
- Basic up/down/ps/logs with auto-discovery
- No config file required
- Single profile (default)

### Phase 2: Profiles & Config
- do.yaml support
- Profile system
- Environment files

### Phase 3: Advanced Features
- Project aliases
- Hooks
- History tracking

### Rollback
- Tool is optional wrapper; users can always fall back to direct `docker compose` commands
- No state modification except docker compose operations (which are idempotent/reversible)

## Open Questions

1. Should we support Docker Compose V1 (deprecated Python version)?
   - **Decision**: No, V2 only (Go version shipped with Docker Desktop)

2. How to handle Windows paths (backslashes vs forward slashes)?
   - **Decision**: Use filepath.Join() for cross-platform compatibility

3. Should we support `.yml` extension or only `.yaml`?
   - **Decision**: Support both, prefer `.yaml` (Docker convention)

4. How to handle project config vs user config?
   - **Decision**: Project-local `do.yaml` takes precedence, `~/.config/do/config.yaml` for global aliases
