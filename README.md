# do - Docker Compose CLI Wrapper

A simplified CLI wrapper for Docker Compose that eliminates the verbosity of managing multi-file compose stacks.

## Why do?

Docker Compose projects often span multiple files (base, dev, prod, db, etc.). The native `docker compose -f file1.yaml -f file2.yaml -f file3.yaml ...` syntax becomes unwieldy with 3+ files.

**do** simplifies this with:
- Auto-discovery of compose files
- Profile-based configuration
- Dry-run mode for command preview
- Convenient shortcuts for common workflows

## Installation

```bash
# Build from source
go install github.com/akrc/do@latest

# Or build locally
git clone https://github.com/akrc/do.git
cd do
make build
```

## Quick Start

```bash
# Start services (auto-discovers compose.yaml and compose.*.yaml)
do c up

# Start in detached mode
do c up -d

# View service status
do c ps

# View logs
do c logs -f

# Stop services
do c down
```

## Configuration

### Auto-Discovery (No Config Required)

do automatically discovers compose files in your directory:
- `compose.yaml` or `docker-compose.yaml`
- `compose.*.yaml` files (e.g., `compose.dev.yaml`, `compose.prod.yaml`)

Files are loaded in alphabetical order.

### Using do.yaml

Create a `do.yaml` in your project directory for advanced configuration:

```yaml
# Define profiles for different environments
default_profile: dev

profiles:
  dev:
    slices: [base, dev]
    env_file: .env.dev

  prod:
    slices: [base, prod]
    env_file: .env.prod

  full:
    slices: [base, dev, db, monitoring]
    env_file: .env.full

# Define custom file combinations
slices:
  base: ["compose.yaml"]
  dev: ["compose.dev.yaml"]
  prod: ["compose.prod.yaml"]
  db: ["compose.db.yaml"]
  monitoring: ["compose.monitoring.yaml"]

# Define command aliases
aliases:
  fresh: "down -v && up --build -d"
  restart-all: "restart && logs -f"
  rebuild: "down -v && up --build -d"

# Define hooks to run before/after commands
hooks:
  pre_up:
    - "echo 'Starting services...'"
    - "docker network prune -f"
  post_up:
    - "echo 'Services are ready!'"
  pre_down:
    - "echo 'Stopping services...'"
```

## Commands

### Core Compose Commands

```bash
# Start services
do c up
do c up -d                    # detached mode
do c up --build              # rebuild images

# Stop services
do c down
do c down -v                 # remove volumes
do c down --remove-orphans   # remove orphaned containers

# View status
do c ps
do c status                  # enhanced status view
do s                         # shorthand for status

# View logs
do c logs
do c logs -f                 # follow logs
do c logs --tail 100         # show last 100 lines
do c logs api                # logs for specific service

# Service management
do c restart api             # restart specific service
do c exec api bash           # execute command in container
do c build api               # rebuild specific service
```

### Convenience Commands

```bash
# dup: down then up
do c dup

# nuke: complete cleanup (down -v --remove-orphans)
do c nuke

# fresh: clean rebuild (down -v && up --build -d)
do c fresh
```

### Aliases

```bash
# List all aliases
do c alias

# Run an alias
do c alias fresh
```

### Global Flags

```bash
# Dry-run: preview commands without executing
do --dry-run c up

# Verbose: show debug information
do --verbose c up
do -v c up
```

## Profile Management

Switch between different compose configurations:

```bash
# Use a specific profile
do --profile prod c up

# Override default profile from do.yaml
do --profile full c up
```

## Project Aliases

Define project shortcuts in `~/.config/do/config.yaml`:

```yaml
projects:
  webapp:
    path: /home/user/projects/webapp
    description: "Web application"

  api:
    path: /home/user/projects/api
    description: "API service"

  microservices:
    path: /home/user/projects/microservices
    description: "Microservices architecture"

# Global aliases (available in all projects)
aliases:
  refresh: "down && up --build -d"
  clean: "down -v --remove-orphans"
```

Then use the `@project` syntax to run commands in any project:

```bash
# Run commands in a different project
do @webapp c up
do @api logs -f
do @microservices c status

# Works with any do command
do @webapp c nuke
do @api c fresh
```

## File Discovery

do discovers compose files using this precedence:

1. **Explicit `-f` flags** (highest priority)
2. **do.yaml profile configuration**
3. **Auto-discovery** in current directory

```bash
# Override auto-discovery with explicit files
do c up -f custom.yaml
do c up -f base.yaml -f override.yaml
```

## Environment Variables

Set environment files per profile:

```yaml
profiles:
  dev:
    slices: [base, dev]
    env_file: .env.dev
```

The env file will be passed to Docker Compose with `--env-file` flag.

## Hooks

Execute commands before or after Docker Compose operations:

```yaml
hooks:
  pre_up:
    - "echo 'Starting services...'"
    - "docker network prune -f"
  post_up:
    - "echo 'Services are ready!'"
  pre_down:
    - "echo 'Stopping services...'"
```

Hooks run in the order defined. If a hook fails, subsequent hooks and the main command are not executed.

## Examples

### Simple Project

```bash
# Directory with compose.yaml and compose.dev.yaml
myproject/
  compose.yaml
  compose.dev.yaml

# Just run
do c up
# Equivalent to: docker compose -f compose.yaml -f compose.dev.yaml up
```

### Complex Multi-Environment Setup

```bash
# do.yaml
profiles:
  dev:
    slices: [base, dev, db]
    env_file: .env.dev

  staging:
    slices: [base, staging, db]
    env_file: .env.staging

  prod:
    slices: [base, prod, db, monitoring]
    env_file: .env.prod

# Usage
do --profile dev c up
do --profile staging c up
do --profile prod c up
```

### Daily Workflow

```bash
# Morning: start everything fresh
do c fresh

# During development: rebuild and restart
do c dup

# Check status
do c ps

# View logs
do c logs -f

# End of day: clean shutdown
do c down
```

### Multi-Project Management

```bash
# ~/.config/do/config.yaml
projects:
  frontend:
    path: ~/projects/frontend
  backend:
    path: ~/projects/backend
  db:
    path: ~/projects/database

# Start all projects
do @frontend c up -d
do @backend c up -d
do @db c up -d

# Check all statuses
do @frontend c ps
do @backend c ps
do @db c ps
```

## File Locations

- **Project config**: `./do.yaml` (in your project directory)
- **Global config**: `~/.config/do/config.yaml`
- **Command history**: `~/.cache/do/history.yaml`

## Test Coverage

do follows TDD principles with comprehensive test coverage:

- Config: 93.9%
- Compose: 79.3%
- Overall: >80%

Run tests:
```bash
go test ./...
go test -race ./...    # race detection
go test -bench=. ./...  # benchmarks
```

## Requirements

- Go 1.21+
- Docker Compose V2

## License

MIT

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- New features include tests
- Code follows existing patterns
- Commit messages are clear
