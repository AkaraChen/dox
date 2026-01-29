# Compose Configuration Specification

## ADDED Requirements

### Requirement: Auto-Discovery of Compose Files
The system SHALL automatically discover Docker Compose files in the current directory.

#### Scenario: Discover base and slice files
- **WHEN** directory contains `compose.yaml` and `compose.dev.yaml`
- **THEN** auto-discovery finds both files
- **AND** base file is `compose.yaml`
- **AND** slice file is `compose.dev.yaml`

#### Scenario: Discover with .yml extension
- **WHEN** directory contains `compose.yml` and `compose.dev.yml`
- **THEN** auto-discovery finds both files
- **AND** treats `.yml` same as `.yaml`

#### Scenario: Multiple slice files
- **WHEN** directory contains `compose.yaml`, `compose.dev.yaml`, `compose.db.yaml`, `compose.redis.yaml`
- **THEN** auto-discovery finds all files
- **AND** identifies slices: dev, db, redis

#### Scenario: No compose files found
- **WHEN** directory contains no compose.* files
- **THEN** return empty discovery result
- **AND** suggest user create compose.yaml or specify -f flag

#### Scenario: Both .yaml and .yml present
- **WHEN** directory contains both `compose.yaml` and `compose.yml`
- **THEN** prefer `.yaml` extension
- **AND** ignore `.yml` variant

### Requirement: Compose File Ordering
The system SHALL order discovered compose files correctly for Docker Compose.

#### Scenario: Base file first
- **WHEN** building command with auto-discovered files
- **THEN** base file (`compose.yaml` or `docker-compose.yaml`) is always first
- **AND** followed by slice files in alphabetical order

#### Scenario: Slice file ordering
- **WHEN** multiple slice files exist (dev, db, redis)
- **THEN** order slices alphabetically: db, dev, redis
- **AND** build command: `-f compose.yaml -f compose.db.yaml -f compose.dev.yaml -f compose.redis.yaml`

### Requirement: dox.yaml Configuration File
The system SHALL support a project-local `dox.yaml` configuration file.

#### Scenario: Parse minimal config
- **WHEN** `dox.yaml` contains only `version: 1`
- **THEN** parse successfully
- **AND** use default behavior for all settings

#### Scenario: Parse profiles config
- **WHEN** `dox.yaml` contains profiles section with dev and prod
- **THEN** parse profiles into memory
- **AND** make available for profile resolution

#### Scenario: Parse aliases config
- **WHEN** `dox.yaml` contains aliases section
- **THEN** parse aliases into memory
- **AND** map alias names to command templates

#### Scenario: Parse hooks config
- **WHEN** `dox.yaml` contains hooks section with pre_up and post_up
- **THEN** parse hooks into memory
- **AND** associate hooks with appropriate lifecycle events

#### Scenario: Invalid YAML syntax
- **WHEN** `dox.yaml` contains invalid YAML syntax
- **THEN** display helpful error message with line number
- **AND** fall back to default behavior
- **AND** exit with non-zero status if critical config is invalid

### Requirement: Discovery Configuration
The system SHALL allow configuration of auto-discovery behavior via dox.yaml.

#### Scenario: Custom pattern
- **WHEN** dox.yaml specifies `discovery.pattern: "docker-compose.*.yml"`
- **THEN** use custom pattern for file discovery
- **AND** match files against custom pattern

#### Scenario: Custom base file
- **WHEN** dox.yaml specifies `discovery.base: "docker-compose.yml"`
- **THEN** use custom base file
- **AND** treat specified file as base (not as slice)

#### Scenario: Disable auto-discovery
- **WHEN** dox.yaml specifies `discovery.enabled: false`
- **THEN** skip auto-discovery
- **AND** require explicit file specification or profile

### Requirement: Default Behavior Configuration
The system SHALL support default behavior settings in dox.yaml.

#### Scenario: Default profile
- **WHEN** dox.yaml specifies `defaults.profile: dev`
- **AND** user runs `do c up` without `-p` flag
- **THEN** use dev profile automatically

#### Scenario: Default slice
- **WHEN** dox.yaml specifies `defaults.slice: dev`
- **AND** no profile is specified
- **THEN** use base file + dev slice

#### Scenario: No default configured
- **WHEN** dox.yaml has no defaults section
- **AND** user runs `do c up` without specifying profile/slice
- **THEN** use only base file (compose.yaml)
- **AND** warn user that no profile is active

### Requirement: Environment Variable File Mapping
The system SHALL support environment file configuration per profile.

#### Scenario: Env file in profile
- **WHEN** profile 'dev' specifies `env_file: .env.dev`
- **AND** user runs `do c up -p dev`
- **THEN** pass `--env-file .env.dev` to docker compose
- **AND** load environment variables from specified file

#### Scenario: Global env_files mapping
- **WHEN** dox.yaml contains `env_files.dev: .env.dev` and `env_files.prod: .env.prod`
- **AND** profile references env: dev
- **THEN** resolve env file from global mapping
- **AND** pass to docker compose command

#### Scenario: Missing env file
- **WHEN** specified env file does not exist
- **THEN** display warning message
- **AND** continue execution (docker compose will handle missing file)

### Requirement: Aliases Configuration
The system SHALL support custom command aliases in dox.yaml.

#### Scenario: Simple alias
- **WHEN** dox.yaml defines `aliases.fresh: "down -v && up --build -d"`
- **AND** user runs `do c fresh`
- **THEN** execute down with volumes
- **AND** then execute up with build and detach flags

#### Scenario: Alias with profile
- **WHEN** alias contains profile reference `-p prod`
- **AND** user runs the alias
- **THEN** execute command with prod profile active

#### Scenario: Alias expansion in dry-run
- **WHEN** user runs `do c fresh --dry-run`
- **THEN** display expanded command(s) that would execute
- **AND** show each step of the alias

### Requirement: Hooks Configuration
The system SHALL support lifecycle hooks for command execution.

#### Scenario: Pre-up hook
- **WHEN** dox.yaml defines `hooks.pre_up: ["echo 'Starting...'"]`
- **AND** user runs `do c up`
- **THEN** execute pre_up hook before docker compose command
- **AND** display hook output
- **AND** proceed to up command if hook succeeds
- **AND** abort if hook fails (non-zero exit)

#### Scenario: Post-up hook
- **WHEN** dox.yaml defines `hooks.post_up: ["echo 'Ready!'"]`
- **AND** user runs `do c up`
- **THEN** execute docker compose up command
- **AND** after completion, execute post_up hook
- **AND** display hook output

#### Scenario: Multiple hooks
- **WHEN** dox.yaml defines multiple pre_up hooks
- **THEN** execute hooks in order defined
- **AND** stop on first failure
- **AND** report which hook failed

#### Scenario: Hook with service reference
- **WHEN** hook contains command using docker-compose service
- **THEN** execute hook with discovered compose files
- **AND** apply same profile/slice context

### Requirement: User Config Directory
The system SHALL support global configuration in `~/.config/dox/config.yaml`.

#### Scenario: User config for project aliases
- **WHEN** `~/.config/dox/config.yaml` contains projects section
- **THEN** parse global config on startup
- **AND** merge with project-local dox.yaml (project takes precedence)

#### Scenario: Missing user config directory
- **WHEN** `~/.config/dox/` does not exist
- **THEN** create directory structure
- **AND** use default configuration

#### Scenario: User config for defaults
- **WHEN** user config specifies global defaults
- **AND** project config doesn't override
- **THEN** apply user config defaults

### Requirement: Config Validation
The system SHALL validate configuration files and report errors clearly.

#### Scenario: Invalid profile reference
- **WHEN** dox.yaml references non-existent profile in defaults
- **THEN** display error: "Profile 'xyz' not found in configuration"
- **AND** list available profiles
- **AND** exit with non-zero status

#### Scenario: Invalid slice reference
- **WHEN** profile references slice that doesn't exist in discovery
- **THEN** display error: "Slice file for 'xyz' not found"
- **AND** list discovered slices
- **AND** exit with non-zero status

#### Scenario: Unknown config keys
- **WHEN** dox.yaml contains keys not recognized by schema
- **THEN** display warning: "Unknown configuration key: 'xyz'"
- **AND** continue with valid configuration
