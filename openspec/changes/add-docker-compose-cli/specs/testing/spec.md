# Testing Specification

## ADDED Requirements

### Requirement: Test-Driven Development Approach
The project SHALL follow TDD methodology for all core functionality.

#### Scenario: Write test before implementation
- **WHEN** implementing a new feature or command
- **THEN** write failing tests first (Red phase)
- **AND** write minimum code to pass tests (Green phase)
- **AND** refactor while keeping tests green (Refactor phase)

#### Scenario: Test coverage requirement
- **WHEN** running `go test -cover ./...`
- **THEN** overall coverage must be at least 80%
- **AND** core packages (compose, config) must have at least 90% coverage

### Requirement: Test Fixtures
The system SHALL provide a comprehensive set of test fixtures for E2E testing.

#### Scenario: Simple project fixture
- **GIVEN** fixture directory `fixtures/simple/`
- **AND** contains only `compose.yaml` with basic web service
- **WHEN** tests run against this fixture
- **THEN** `do c up --dry-run` outputs `docker compose -f compose.yaml up`

#### Scenario: Multi-slice project fixture
- **GIVEN** fixture directory `fixtures/multi-slice/`
- **AND** contains `compose.yaml`, `compose.dev.yaml`, `compose.prod.yaml`, `compose.db.yaml`
- **WHEN** tests run `do c up --dry-run` with profile using dev and db
- **THEN** output includes files in correct order: `-f compose.yaml -f compose.db.yaml -f compose.dev.yaml`
- **AND** files are ordered alphabetically after base

#### Scenario: Profile-based fixture
- **GIVEN** fixture directory `fixtures/with-profiles/`
- **AND** contains `dox.yaml` with profiles: dev (slices: [dev]), prod (slices: [prod])
- **WHEN** test runs `do c up -p dev --dry-run`
- **THEN** output includes `compose.yaml` and `compose.dev.yaml`
- **AND** when running `do c up -p prod --dry-run`, includes `compose.prod.yaml`

#### Scenario: Env file fixture
- **GIVEN** fixture directory `fixtures/with-env/`
- **AND** contains `.env.dev`, `.env.prod`
- **AND** `dox.yaml` defines profiles with env_file mapping
- **WHEN** test runs `do c up -p dev --dry-run`
- **THEN** output includes `--env-file .env.dev`

#### Scenario: Aliases fixture
- **GIVEN** fixture directory `fixtures/with-aliases/`
- **AND** `dox.yaml` defines `aliases.fresh: "down -v && up --build -d"`
- **WHEN** test runs `do c fresh --dry-run`
- **THEN** output shows two commands: `docker compose ... down -v` and `docker compose ... up --build -d`

#### Scenario: Hooks fixture
- **GIVEN** fixture directory `fixtures/with-hooks/`
- **AND** `dox.yaml` defines pre_up and post_up hooks
- **WHEN** test runs `do c up --dry-run`
- **THEN** output includes hook commands before and after docker compose command

#### Scenario: Edge cases fixture
- **GIVEN** fixture directory `fixtures/edge-cases/`
- **AND** includes scenarios: no compose files, invalid yaml, missing slices
- **WHEN** tests run against each scenario
- **THEN** appropriate error messages are displayed
- **AND** exit codes are non-zero

### Requirement: Dry-Run Command Testing
The system SHALL verify all assembled commands via dry-run testing.

#### Scenario: Up command dry-run
- **GIVEN** fixture with compose.yaml and compose.dev.yaml
- **WHEN** running `do c up --dry-run`
- **THEN** output shows exact command: `docker compose -f compose.yaml -f compose.dev.yaml up`
- **AND** command is NOT executed

#### Scenario: Up with detached flag dry-run
- **GIVEN** any fixture with valid compose files
- **WHEN** running `do c up -d --dry-run`
- **THEN** output shows command with `-d` flag: `docker compose -f ... up -d`

#### Scenario: Down with flags dry-run
- **GIVEN** any fixture with valid compose files
- **WHEN** running `do c down -v --remove-orphans --dry-run`
- **THEN** output shows: `docker compose -f ... down -v --remove-orphans`

#### Scenario: Logs with service and flags dry-run
- **GIVEN** fixture with compose files
- **WHEN** running `do c logs api -f --tail 50 --dry-run`
- **THEN** output shows: `docker compose -f ... logs -f --tail 50 api`

#### Scenario: Service restart dry-run
- **GIVEN** fixture with compose files
- **WHEN** running `do c restart api --dry-run`
- **THEN** output shows: `docker compose -f ... restart api`

#### Scenario: Exec command dry-run
- **GIVEN** fixture with compose files
- **WHEN** running `do c exec api bash --dry-run`
- **THEN** output shows: `docker compose -f ... exec api bash`

#### Scenario: Convenience command dup dry-run
- **GIVEN** fixture with compose files
- **WHEN** running `do c dup --dry-run`
- **THEN** output shows two commands in sequence:
  - `docker compose -f ... down`
  - `docker compose -f ... up`

#### Scenario: Convenience command nuke dry-run
- **GIVEN** fixture with compose files
- **WHEN** running `do c nuke --dry-run`
- **THEN** output shows: `docker compose -f ... down -v --remove-orphans`

#### Scenario: Convenience command fresh dry-run
- **GIVEN** fixture with compose files
- **WHEN** running `do c fresh --dry-run`
- **THEN** output shows two commands:
  - `docker compose -f ... down -v`
  - `docker compose -f ... up --build`

#### Scenario: Profile with env file dry-run
- **GIVEN** fixture with dox.yaml profile defining env_file: .env.dev
- **WHEN** running `do c up -p dev --dry-run`
- **THEN** output includes: `docker compose -f ... -f compose.dev.yaml --env-file .env.dev up`

#### Scenario: Alias expansion dry-run
- **GIVEN** fixture with alias defined as "fresh: down -v && up --build"
- **WHEN** running `do c fresh --dry-run`
- **THEN** output shows expanded steps of the alias

### Requirement: Table-Driven Tests
The system SHALL use table-driven tests for command variants.

#### Scenario: Command variant test table
- **GIVEN** test file for command building
- **WHEN** defining test cases
- **THEN** use table-driven approach with test cases for:
  - No flags
  - With -d flag
  - With --build flag
  - With multiple flags
  - With service argument
  - With profile selection
- **AND** each test case verifies dry-run output matches expected command

#### Scenario: Error condition test table
- **GIVEN** test file for error handling
- **WHEN** defining error cases
- **THEN** use table-driven approach for:
  - No compose files found
  - Invalid dox.yaml syntax
  - Missing profile reference
  - Missing slice file
  - Docker compose not found
- **AND** each case verifies error message and exit code

### Requirement: Unit Tests for Config Parsing
The system SHALL have comprehensive unit tests for configuration parsing.

#### Scenario: Parse valid dox.yaml
- **GIVEN** a valid dox.yaml with profiles, aliases, hooks
- **WHEN** parsing the config
- **THEN** return struct with all fields populated
- **AND** no errors

#### Scenario: Parse invalid yaml
- **GIVEN** a dox.yaml with invalid YAML syntax
- **WHEN** parsing the config
- **THEN** return error with line number
- **AND** error message describes the issue

#### Scenario: Parse empty config
- **GIVEN** a dox.yaml with only `version: 1`
- **WHEN** parsing the config
- **THEN** return default config struct
- **AND** no errors

#### Scenario: Parse profile inheritance
- **GIVEN** a dox.yaml with profile extending another
- **WHEN** parsing and resolving profiles
- **THEN** merged profile includes base and extended slices
- **AND** no duplicate slices

#### Scenario: Parse with variable substitution
- **GIVEN** a dox.yaml with ${VAR} references
- **AND** environment variables set
- **WHEN** parsing the config
- **THEN** variables are substituted with actual values
- **AND** undefined variables return error

### Requirement: Unit Tests for File Discovery
The system SHALL have unit tests for compose file auto-discovery.

#### Scenario: Discover standard layout
- **GIVEN** directory with compose.yaml and slice files
- **WHEN** running discovery
- **THEN** return base file and all slice files
- **AND** slices are sorted alphabetically

#### Scenario: Discover with docker-compose.yaml
- **GIVEN** directory with docker-compose.yaml (legacy name)
- **WHEN** running discovery
- **THEN** treat docker-compose.yaml as base file

#### Scenario: Discover prefers .yaml over .yml
- **GIVEN** directory with both compose.yaml and compose.yml
- **WHEN** running discovery
- **THEN** use compose.yaml, ignore compose.yml

#### Scenario: Discover with no files
- **GIVEN** empty directory or directory without compose files
- **WHEN** running discovery
- **THEN** return empty result
- **AND** error indicating no files found

### Requirement: Integration Tests
The system SHALL have integration tests for full command workflows.

#### Scenario: Full up and down workflow
- **GIVEN** fixture with simple compose.yaml
- **WHEN** running `do c up --dry-run` then `do c down --dry-run`
- **THEN** both commands output correctly
- **AND** files are identical between commands

#### Scenario: Profile switch workflow
- **GIVEN** fixture with dev and prod profiles
- **WHEN** running `do c up -p dev --dry-run`
- **AND** then running `do c up -p prod --dry-run`
- **THEN** outputs use different slice files
- **AND** each profile uses correct compose files

#### Scenario: Alias execution workflow
- **GIVEN** fixture with defined aliases
- **WHEN** running alias command with --dry-run
- **THEN** all steps in alias are output
- **AND** steps use correct compose files

### Requirement: Benchmark Tests
The system SHALL include benchmark tests for performance-critical code.

#### Scenario: Benchmark command building
- **GIVEN** a function that builds docker compose commands
- **WHEN** running `go test -bench=.`
- **THEN** benchmark reports time per operation
- **AND** target: < 1ms per command build

#### Scenario: Benchmark config parsing
- **GIVEN** a function that parses dox.yaml
- **WHEN** running benchmarks
- **THEN** benchmark reports time for various config sizes
- **AND** target: < 10ms for typical config

### Requirement: Race Detection
The system SHALL pass all race detection tests.

#### Scenario: Run tests with race detector
- **WHEN** running `go test -race ./...`
- **THEN** no data races are reported
- **AND** all tests pass

### Requirement: Test Execution in CI
All tests MUST pass in CI before code is merged.

#### Scenario: CI test pipeline
- **GIVEN** a pull request with changes
- **WHEN** CI pipeline runs
- **THEN** execute `go test -v -race -cover ./...`
- **AND** fail pipeline if any test fails
- **AND** fail pipeline if coverage < 80%

#### Scenario: Pre-commit test hook
- **GIVEN** developer attempts to commit code
- **WHEN** pre-commit hook runs
- **THEN** run quick subset of tests
- **AND** block commit if tests fail
