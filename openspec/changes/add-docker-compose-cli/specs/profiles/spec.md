# Profiles Specification

## ADDED Requirements

### Requirement: Profile Definition
The system SHALL allow defining profiles in dox.yaml that group compose file slices.

#### Scenario: Define single-slice profile
- **WHEN** dox.yaml contains `profiles.dev.slices: [dev]`
- **THEN** profile 'dev' consists of base file + compose.dev.yaml

#### Scenario: Define multi-slice profile
- **WHEN** dox.yaml contains `profiles.fullstack.slices: [dev, db, redis]`
- **THEN** profile 'fullstack' includes base + all three slice files
- **AND** order slices: base, db, dev, redis (alphabetical after base)

#### Scenario: Define profile with env file
- **WHEN** dox.yaml contains `profiles.dev.env_file: .env.dev`
- **THEN** profile 'dev' loads environment from .env.dev
- **AND** pass `--env-file .env.dev` to docker compose

#### Scenario: Define profile with env reference
- **WHEN** dox.yaml contains `profiles.dev.env: dev`
- **AND** global `env_files.dev: .env.dev` is defined
- **THEN** resolve env file from global mapping
- **AND** use .env.dev for the profile

### Requirement: Profile Selection
The system SHALL allow selecting profiles via CLI flags.

#### Scenario: Select profile with -p flag
- **WHEN** user runs `do c up -p dev`
- **THEN** use profile 'dev' from dox.yaml
- **AND** build command with base + dev slice files

#### Scenario: Select profile with --profile flag
- **WHEN** user runs `do c up --profile prod`
- **THEN** use profile 'prod' from dox.yaml
- **AND** build command with base + prod slice files

#### Scenario: Profile not found
- **WHEN** user runs `do c up -p nonexistent`
- **AND** profile 'nonexistent' is not defined
- **THEN** display error: "Profile 'nonexistent' not found"
- **AND** list available profiles from dox.yaml
- **AND** exit with non-zero status

### Requirement: Profile Resolution
The system SHALL resolve profiles into compose file lists and env files.

#### Scenario: Resolve profile to compose files
- **WHEN** profile 'dev' defines slices: [dev, db]
- **AND** auto-discovery finds compose.yaml, compose.dev.yaml, compose.db.yaml
- **THEN** resolve to: compose.yaml, compose.db.yaml, compose.dev.yaml
- **AND** order with base first, then slices alphabetically

#### Scenario: Resolve profile with missing slice
- **WHEN** profile 'dev' defines slices: [dev, cache]
- **AND** auto-discovery finds compose.dev.yaml but not compose.cache.yaml
- **THEN** display error: "Slice file 'compose.cache.yaml' not found for profile 'dev'"
- **AND** exit with non-zero status

#### Scenario: Resolve profile with env file
- **WHEN** profile 'dev' defines env_file: .env.dev
- **THEN** resolved command includes `--env-file .env.dev`
- **AND** pass to docker compose command

### Requirement: Default Profile
The system SHALL support a default profile activated when no profile is specified.

#### Scenario: Use default profile
- **WHEN** dox.yaml defines `defaults.profile: dev`
- **AND** user runs `do c up` without -p flag
- **THEN** automatically use profile 'dev'
- **AND** display message: "Using default profile: dev"

#### Scenario: No default profile configured
- **WHEN** dox.yaml has no defaults.profile
- **AND** user runs `do c up` without -p flag
- **THEN** use base file only (compose.yaml)
- **AND** display info: "No profile specified, using base compose file only"

#### Scenario: Override default profile
- **WHEN** default profile is 'dev'
- **AND** user runs `do c up -p prod`
- **THEN** use specified profile 'prod' (ignore default)

### Requirement: Profile Auto-Detection
The system SHALL support automatic profile selection based on git branch.

#### Scenario: Auto-detect from git branch
- **WHEN** dox.yaml defines `defaults.auto_detect: true`
- **AND** current git branch is 'feature/xyz'
- **AND** profile 'feature' exists in dox.yaml
- **THEN** auto-select profile 'feature'
- **AND** display message: "Auto-detected profile 'feature' from branch 'feature/xyz'"

#### Scenario: No matching profile for branch
- **WHEN** auto_detect is enabled
- **AND** current branch is 'feature/xyz'
- **AND** no profile 'feature' exists
- **THEN** fall back to default profile or base file
- **AND** display info: "No profile matches branch 'feature/xyz', using default"

#### Scenario: Not a git repository
- **WHEN** auto_detect is enabled
- **AND** current directory is not a git repository
- **THEN** ignore auto_detect setting
- **AND** use default profile or base file

#### Scenario: Git not available
- **WHEN** auto_detect is enabled
- **AND** git binary is not installed
- **THEN** ignore auto_detect setting
- **AND** continue with default profile or base file

### Requirement: Profile Listing
The system SHALL allow listing available profiles.

#### Scenario: List profiles command
- **WHEN** user runs `do c profile list` or `do c profiles`
- **THEN** display all profiles defined in dox.yaml
- **AND** show slices included in each profile
- **AND** mark default profile with asterisk or (default)
- **AND** show env file if defined

#### Scenario: List profiles when none defined
- **WHEN** user runs `do c profiles`
- **AND** dox.yaml has no profiles section
- **THEN** display "No profiles defined"
- **AND** show available slices from auto-discovery

#### Scenario: Show current profile
- **WHEN** user runs `do c profile current`
- **THEN** display currently active profile (from default or -p flag)
- **OR** display "No active profile" if using base file only

### Requirement: Profile Inheritance
The system SHALL support profile inheritance for shared configuration.

#### Scenario: Profile extends another
- **WHEN** dox.yaml defines `profiles.api.extends: base`
- **AND** profile 'base' defines slices: [db]
- **THEN** profile 'api' includes slices from 'base' plus its own
- **AND** resolve to: base file + db slice + api-specific slices

#### Scenario: Multi-level inheritance
- **WHEN** profile A extends B, and B extends C
- **THEN** resolve slices: C's slices + B's slices + A's slices
- **AND** avoid duplicates (de-duplicate slice names)

#### Scenario: Circular inheritance
- **WHEN** profile A extends B, and B extends A
- **THEN** display error: "Circular profile inheritance detected"
- **AND** identify the circular path
- **AND** exit with non-zero status

### Requirement: Profile Variables
The system SHALL support variable substitution in profile definitions.

#### Scenario: Environment variable in slice name
- **WHEN** profile defines slices: [dev, ${EXTRA_SLICE}]
- **AND** environment variable EXTRA_SLICE=cache is set
- **THEN** resolve slices: [dev, cache]
- **AND** include compose.cache.yaml in command

#### Scenario: Undefined variable
- **WHEN** profile references ${UNDEFINED_VAR}
- **AND** variable is not set in environment
- **THEN** display error: "Environment variable 'UNDEFINED_VAR' not set"
- **AND** exit with non-zero status

#### Scenario: Default value for variable
- **WHEN** profile defines slices: [dev, ${EXTRA_SLICE:-cache}]
- **AND** EXTRA_SLICE is not set
- **THEN** use default value 'cache'
- **AND** resolve to: [dev, cache]

### Requirement: Profile Validation
The system SHALL validate profile configuration before execution.

#### Scenario: Validate on load
- **WHEN** dox.yaml is loaded
- **THEN** validate all profile definitions
- **AND** check slice references exist
- **AND** check env files exist (warning if missing)
- **AND** check inheritance is valid

#### Scenario: Invalid profile name
- **WHEN** profile name contains invalid characters (spaces, special chars)
- **THEN** display error: "Invalid profile name 'my profile'"
- **AND** explain naming rules (alphanumeric, hyphens, underscores)
- **AND** exit with non-zero status

### Requirement: Profile Per-Command
The system SHALL allow different profiles for different command types.

#### Scenario: Profile for up only
- **WHEN** user runs `do c up -p dev`
- **THEN** use dev profile for up command
- **AND** subsequent commands (logs, ps) auto-use same profile
- **AND** store profile selection in state file

#### Scenario: Override profile for command
- **WHEN** previous command used dev profile
- **AND** user runs `do c logs -p prod`
- **THEN** use prod profile for logs command
- **AND** update stored profile to prod

#### Scenario: Clear profile selection
- **WHEN** user runs `do c up --no-profile`
- **THEN** use base file only
- **AND** clear stored profile state
