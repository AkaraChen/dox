# Core CLI Specification

## ADDED Requirements

### Requirement: Root Command Structure
The CLI SHALL provide a root command `do` using the Cobra framework.

#### Scenario: Help output displays available commands
- **WHEN** user runs `do --help` or `do -h`
- **THEN** display usage, available commands, and global flags
- **AND** show compose command group (`do c`)

#### Scenario: Version information
- **WHEN** user runs `do --version` or `do -v`
- **THEN** display version number and build information

### Requirement: Compose Command Group
The CLI SHALL provide a `c` subcommand as a namespace for Docker Compose operations.

#### Scenario: Compose help displays subcommands
- **WHEN** user runs `do c --help`
- **THEN** display available compose subcommands (up, down, ps, logs, etc.)

#### Scenario: Unknown command error
- **WHEN** user runs `do c invalid-command`
- **THEN** display helpful error message suggesting valid commands
- **AND** exit with non-zero status

### Requirement: Up Command
The CLI SHALL provide an `up` subcommand to start Docker Compose services.

#### Scenario: Basic up with auto-discovery
- **WHEN** user runs `do c up` in directory with `compose.yaml` and `compose.dev.yaml`
- **THEN** execute `docker compose -f compose.yaml -f compose.dev.yaml up`
- **AND** display output from docker compose

#### Scenario: Up with detached flag
- **WHEN** user runs `do c up -d` or `do c up --detach`
- **THEN** execute docker compose with `-d` flag
- **AND** run containers in detached mode

#### Scenario: Up with dry-run
- **WHEN** user runs `do c up --dry-run`
- **THEN** display the exact docker compose command that would be executed
- **AND** do NOT execute the command
- **AND** exit with status 0

#### Scenario: Up with explicit files
- **WHEN** user runs `do c up -f custom.yaml`
- **THEN** use only the specified file(s), ignoring auto-discovery
- **AND** execute `docker compose -f custom.yaml up`

### Requirement: Down Command
The CLI SHALL provide a `down` subcommand to stop and remove Docker Compose services.

#### Scenario: Basic down
- **WHEN** user runs `do c down`
- **THEN** execute `docker compose -f <discovered-files> down`
- **AND** stop and remove containers

#### Scenario: Down with volumes
- **WHEN** user runs `do c down -v`
- **THEN** execute docker compose down with `-v` flag
- **AND** remove named volumes declared in the compose file

#### Scenario: Down with remove-orphans
- **WHEN** user runs `do c down --remove-orphans`
- **THEN** execute docker compose down with `--remove-orphans` flag

### Requirement: PS Command
The CLI SHALL provide a `ps` subcommand to list running containers.

#### Scenario: List containers
- **WHEN** user runs `do c ps`
- **THEN** execute `docker compose -f <discovered-files> ps`
- **AND** display container status

### Requirement: Logs Command
The CLI SHALL provide a `logs` subcommand to view container logs.

#### Scenario: View all logs
- **WHEN** user runs `do c logs`
- **THEN** execute `docker compose -f <discovered-files> logs`
- **AND** display logs from all services

#### Scenario: View specific service logs
- **WHEN** user runs `do c logs api`
- **THEN** execute `docker compose -f <discovered-files> logs api`
- **AND** display logs only for the api service

#### Scenario: Follow logs
- **WHEN** user runs `do c logs -f` or `do c logs --follow`
- **THEN** execute docker compose logs with `-f` flag
- **AND** continuously stream new log output

#### Scenario: Tail logs
- **WHEN** user runs `do c logs --tail 50`
- **THEN** execute docker compose logs with `--tail 50` flag
- **AND** display last 50 lines per service

### Requirement: Service Shorthand Commands
The CLI SHALL provide shorthand commands for common service operations.

#### Scenario: Restart service
- **WHEN** user runs `do c restart api`
- **THEN** execute `docker compose -f <discovered-files> restart api`
- **AND** restart only the api service

#### Scenario: Execute command in container
- **WHEN** user runs `do c exec api bash`
- **THEN** execute `docker compose -f <discovered-files> exec api bash`
- **AND** open interactive shell in api container

#### Scenario: Build specific service
- **WHEN** user runs `do c build api`
- **THEN** execute `docker compose -f <discovered-files> build api`
- **AND** rebuild only the api service image

### Requirement: Convenience Commands
The CLI SHALL provide compound commands for common workflows.

#### Scenario: Dup (down then up)
- **WHEN** user runs `do c dup`
- **THEN** execute down command first
- **AND** after down completes, execute up command
- **AND** report failure if down fails (don't proceed to up)

#### Scenario: Nuke (complete cleanup)
- **WHEN** user runs `do c nuke`
- **THEN** execute `docker compose -f <discovered-files> down -v --remove-orphans`
- **AND** remove containers, volumes, and orphaned containers

#### Scenario: Fresh (clean rebuild)
- **WHEN** user runs `do c fresh`
- **THEN** execute down with volumes removed
- **AND** execute up with --build flag
- **AND** ensure fresh start with rebuilt images

### Requirement: Status Command
The CLI SHALL provide an enhanced `status` command (or `do s` shorthand) for container status.

#### Scenario: View status
- **WHEN** user runs `do s` or `do status`
- **THEN** display formatted container status with colors
- **AND** show service name, state, ports, and health

#### Scenario: Watch mode
- **WHEN** user runs `do s -w` or `do status --watch`
- **THEN** refresh status display every 2 seconds
- **AND** continue until user interrupts (Ctrl+C)

### Requirement: Global Flags
The CLI SHALL support global flags applicable to all commands.

#### Scenario: Dry-run globally
- **WHEN** user runs `do --dry-run c up`
- **THEN** display command without executing
- **AND** apply dry-run to any compose subcommand

#### Scenario: Verbose output
- **WHEN** user runs `do -v c up` or `do --verbose c up`
- **THEN** display debug information including discovered files
- **AND** show exact docker compose command being executed

### Requirement: Flag Passthrough
The CLI SHALL pass unrecognized flags through to Docker Compose.

#### Scenario: Passthrough of docker compose flags
- **WHEN** user runs `do c up --force-recreate --build`
- **THEN** include `--force-recreate` and `--build` in docker compose command
- **AND** execute with all passthrough flags

### Requirement: Error Handling
The CLI SHALL handle Docker Compose errors gracefully.

#### Scenario: Docker not found
- **WHEN** docker compose binary is not found
- **THEN** display helpful error message with installation instructions
- **AND** exit with status code indicating missing dependency

#### Scenario: No compose files found
- **WHEN** user runs `do c up` in directory without compose files
- **THEN** display error message explaining auto-discovery failure
- **AND** suggest creating compose.yaml or using -f flag

#### Scenario: Compose command fails
- **WHEN** docker compose command exits with error
- **THEN** propagate the exit code from docker compose
- **AND** display error output from docker compose

### Requirement: Test Coverage for CLI Commands
The system SHALL maintain comprehensive test coverage for all CLI commands using dry-run verification.

#### Scenario: Every command variant has dry-run test
- **GIVEN** any CLI command (up, down, ps, logs, restart, exec, build, dup, nuke, fresh)
- **WHEN** running with --dry-run flag
- **THEN** test verifies exact docker compose command output
- **AND** test covers all flag combinations

#### Scenario: Command output is deterministic
- **GIVEN** same fixture and command
- **WHEN** running dry-run multiple times
- **THEN** output is identical across runs
- **AND** file ordering is consistent

#### Scenario: Dry-run doesn't execute
- **GIVEN** any command with --dry-run
- **WHEN** executed
- **THEN** docker compose binary is NOT called
- **AND** exit status is 0
- **AND** only command string is printed

#### Scenario: Dry-run output format
- **GIVEN** any command with --dry-run
- **WHEN** executed
- **THEN** output format is: `docker compose -f <file1> -f <file2> ... <subcommand> <flags> <args>`
- **AND** files are in correct order
- **AND** flags appear before arguments
