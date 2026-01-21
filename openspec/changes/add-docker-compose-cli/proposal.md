# Change: Add Docker Compose CLI Tool

## Why

Managing multi-file Docker Compose setups with multiple `-f` flags is tedious and error-prone. Users with complex compose stacks (base + dev/prod slices + overrides) need a streamlined way to compose, launch, and manage their containers without repetitive command-line arguments.

## What Changes

- Add a Go-based CLI tool `do` using Cobra framework
- Auto-discover compose files (compose.yaml, compose.*.yaml) in current directory
- Support slice profiles for pre-defined compose file combinations
- Add shorthand commands for common operations (up, down, ps, logs, exec, restart)
- Support environment variable file management per profile
- Add dry-run mode to preview generated commands
- Add convenience commands (dup, nuke, fresh)
- Support project aliases for working with multiple compose projects
- **TDD approach with comprehensive test fixtures and dry-run verification**

## Impact

- Affected specs: New capabilities (core-cli, compose-config, profiles, testing)
- Affected code: New project - initial Go module structure
- Dependencies: Cobra (CLI framework), viper (config), yaml handling, testify (testing)
- Test fixtures: 6+ fixture directories covering simple, multi-slice, profiles, env files, aliases, hooks, edge cases
- Coverage requirement: >80% overall, >90% for core packages
