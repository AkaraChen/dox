# dox - Docker Compose CLI Wrapper

A simplified CLI wrapper for Docker Compose that eliminates the verbosity of managing multi-file compose stacks.

## Why dox?

Docker Compose projects often span multiple files (base, dev, prod, db, etc.). The native `docker compose -f file1.yaml -f file2.yaml -f file3.yaml ...` syntax becomes unwieldy with 3+ files.

**dox** simplifies this with:
- Auto-discovery of compose files
- Profile-based configuration
- Dry-run mode for command preview
- Convenient shortcuts for common workflows

## Installation

```bash
# Build from source
go install github.com/akrc/dox@latest

# Or build locally
git clone https://github.com/akrc/dox.git
cd dox
make build
```

## Quick Start

```bash
# Start services (auto-discovers compose.yaml and compose.*.yaml)
dox c up

# Start in detached mode
dox c up -d

# View service status
dox c ps

# View logs
dox c logs -f

# Stop services
dox c down
```

## Configuration

### Auto-Discovery (No Config Required)

dox automatically discovers compose files in your directory:
- `compose.yaml` or `docker-compose.yaml`
- `compose.*.yaml` files (e.g., `compose.dev.yaml`, `compose.prod.yaml`)

Files are loaded in alphabetical order.

### Using dox.yaml

Create a `dox.yaml` in your project directory for advanced configuration:

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
dox c up
dox c up -d                    # detached mode
dox c up --build              # rebuild images

# Stop services
dox c down
dox c down -v                 # remove volumes
dox c down --remove-orphans   # remove orphaned containers

# View status
dox c ps
dox c status                  # enhanced status view
dox s                         # shorthand for status

# View logs
dox c logs
dox c logs -f                 # follow logs
dox c logs --tail 100         # show last 100 lines
dox c logs api                # logs for specific service

# Service management
dox c restart api             # restart specific service
dox c exec api bash           # execute command in container
dox c build api               # rebuild specific service
```

### Convenience Commands

```bash
# dup: down then up
dox c dup

# nuke: complete cleanup (down -v --remove-orphans)
dox c nuke

# fresh: clean rebuild (down -v && up --build -d)
dox c fresh
```

### Aliases

```bash
# List all aliases
dox c alias

# Run an alias
dox c alias fresh
```

### Global Flags

```bash
# Dry-run: preview commands without executing
dox --dry-run c up

# Verbose: show debug information
dox --verbose c up
dox -v c up
```

## Profile Management

Switch between different compose configurations:

```bash
# Use a specific profile
dox --profile prod c up

# Override default profile from dox.yaml
dox --profile full c up
```

## Project Aliases

Define project shortcuts in `~/.config/dox/config.yaml`:

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
dox @webapp c up
dox @api logs -f
dox @microservices c status

# Works with any dox command
dox @webapp c nuke
dox @api c fresh
```

## File Discovery

dox discovers compose files using this precedence:

1. **Explicit `-f` flags** (highest priority)
2. **dox.yaml profile configuration**
3. **Auto-discovery** in current directory

```bash
# Override auto-discovery with explicit files
dox c up -f custom.yaml
dox c up -f base.yaml -f override.yaml
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
dox c up
# Equivalent to: docker compose -f compose.yaml -f compose.dev.yaml up
```

### Complex Multi-Environment Setup

```bash
# dox.yaml
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
dox --profile dev c up
dox --profile staging c up
dox --profile prod c up
```

### Daily Workflow

```bash
# Morning: start everything fresh
dox c fresh

# During development: rebuild and restart
dox c dup

# Check status
dox c ps

# View logs
dox c logs -f

# End of day: clean shutdown
dox c down
```

### Multi-Project Management

```bash
# ~/.config/dox/config.yaml
projects:
  frontend:
    path: ~/projects/frontend
  backend:
    path: ~/projects/backend
  db:
    path: ~/projects/database

# Start all projects
dox @frontend c up -d
dox @backend c up -d
dox @db c up -d

# Check all statuses
dox @frontend c ps
dox @backend c ps
dox @db c ps
```

## File Locations

- **Project config**: `./dox.yaml` (in your project directory)
- **Global config**: `~/.config/dox/config.yaml`
- **Command history**: `~/.cache/dox/history.yaml`

## Test Coverage

dox follows TDD principles with comprehensive test coverage:

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
