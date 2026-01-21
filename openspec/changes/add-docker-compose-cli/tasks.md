# Implementation Tasks

## 1. Project Setup & Test Infrastructure (TDD First)
- [x] 1.1 Initialize Go module with go.mod
- [x] 1.2 Set up test directory structure (internal/*/testdata, fixtures)
- [x] 1.3 Add test dependencies (testify/assert, testfixtures)
- [x] 1.4 Create basic project structure (cmd, internal, pkg)
- [x] 1.5 Set up CI test configuration (Makefile with test target)

## 2. Test Fixtures Creation (Before Implementation)
- [x] 2.1 Create fixture: simple project (compose.yaml only)
- [x] 2.2 Create fixture: multi-slice project (base + dev + prod + db)
- [x] 2.3 Create fixture: project with do.yaml profiles
- [x] 2.4 Create fixture: project with aliases and hooks
- [x] 2.5 Create fixture: project with env files per profile
- [x] 2.6 Create fixture: edge cases (no compose, invalid yaml, circular inheritance)

## 3. E2E Tests: Command Building (Dry-Run Focus)
- [x] 3.1 Write test: simple up command outputs correct docker compose
- [x] 3.2 Write test: multi-slice up orders files correctly
- [x] 3.3 Write test: profile selection includes correct slices
- [x] 3.4 Write test: env-file flag added when profile specifies
- [x] 3.5 Write test: down -v --remove-orphans outputs correctly
- [x] 3.6 Write test: logs with -f --tail flags passthrough
- [x] 3.7 Write test: service-specific commands (restart api, exec api bash)
- [x] 3.8 Write test: convenience commands (dup, nuke, fresh)
- [ ] 3.9 Write test: aliases expand to correct commands
- [ ] 3.10 Write test: hooks execute in correct order

## 4. Core CLI Commands (TDD-Driven)
- [x] 4.1 Implement `compose` command group (`do c`)
- [x] 4.2 Implement `up` command with auto-discovery (tests first)
- [x] 4.3 Implement `down` command (tests first)
- [x] 4.4 Implement `ps` command (tests first)
- [x] 4.5 Implement `logs` command with follow flag (tests first)
- [x] 4.6 Implement shorthand service commands (restart, exec, build) (tests first)

## 5. Configuration System (TDD-Driven)
- [x] 5.1 Implement do.yaml config parser (tests first)
- [x] 5.2 Add compose file auto-discovery logic (tests first)
- [x] 5.3 Add profile resolution (merge base + slices) (tests first)
- [x] 5.4 Add environment file support per profile (tests first)
- [ ] 5.5 Add project alias support (tests first)

## 6. Docker Command Builder (TDD-Driven)
- [x] 6.1 Build docker compose command with -f flags (tests first)
- [x] 6.2 Add dry-run mode flag (tests first)
- [x] 6.3 Implement command execution (tests first)
- [x] 6.4 Add output formatting (tests first)

## 7. Convenience Commands (TDD-Driven)
- [x] 7.1 Implement `dup` (down + up) (tests first)
- [x] 7.2 Implement `nuke` (down -v --remove-orphans) (tests first)
- [x] 7.3 Implement `fresh` (clean rebuild) (tests first)
- [ ] 7.4 Implement `status` command (enhanced ps) (tests first)

## 8. Project Features (TDD-Driven)
- [ ] 8.1 Add project alias system (~/.config/do/config.yaml) (tests first)
- [ ] 8.2 Implement `@project` syntax for remote project operations (tests first)
- [ ] 8.3 Add command history tracking (tests first)
- [ ] 8.4 Add hooks (pre_up, post_up, etc.) (tests first)

## 9. Test Coverage & Quality
- [x] 9.1 Achieve >80% code coverage (go test -cover) - config: 93.9%, compose: 73.8%
- [x] 9.2 Add table-driven tests for command variants
- [ ] 9.3 Add benchmarks for command building performance
- [ ] 9.4 Add race detection tests (go test -race)
- [x] 9.5 Verify all fixture scenarios pass dry-run tests

## 10. Docs & Release
- [ ] 10.1 Write README with usage examples
- [ ] 10.2 Add man page generation
- [x] 10.3 Add build script for multiple platforms (Makefile)
- [ ] 10.4 Set up GitHub Actions for CI/CD with test gating
- [ ] 10.5 Create Homebrew formula
