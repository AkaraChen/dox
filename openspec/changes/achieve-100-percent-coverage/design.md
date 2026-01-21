# Design: Achieve 100% Test Coverage

## Coverage Analysis

### cmd Package (7.6% coverage)
**Uncovered functions:**
- `listAliases()` - alias listing
- `executeAlias()` - alias execution
- `getComposeBuilder()` - builder getter
- `getConfig()` - config getter
- `getComposeExecutor()` - executor getter
- `executeHooks()` - hooks execution
- `parseHookCommand()` - hook parsing
- `executeCommand()` - single command execution
- `executeCommands()` - multiple commands
- `printCommand()` - command printing
- `resolveFile()` - file path resolution
- `resolveAlias()` - alias resolution
- `isKnownCommand()` - command validation

**Strategy:**
- Create integration-style tests that exercise Cobra commands
- Use test doubles for executor/builder
- Test error paths and edge cases

### compose Package (79.3% coverage)
**Partially covered functions:**
- `resolveFiles()` - 86.7% (missing error paths)
- `BuildDown()` - 83.3% (missing flag combinations)
- `BuildPs()` - 83.3% (missing service args)
- `BuildLogs()` - 83.3% (missing flag variations)
- `BuildRestart()` - 75.0% (missing service args)
- `BuildExec()` - 75.0% (missing command variations)
- `BuildBuild()` - 83.3% (missing args)
- `BuildNuke()` - 80.0% (missing flags)
- `BuildFresh()` - 77.8% (missing variations)
- `BuildDup()` - 77.8% (missing variations)
- `BuildStatus()` - 83.3% (missing flags)
- `RunCommand()` - 94.1% (missing error paths)
- `RunCommandWithOutput()` - 61.5% (missing error paths)

**Uncovered:**
- `RunInteractive()` - 0% (exec with stdin)

**Strategy:**
- Add table-driven tests for all flag combinations
- Test error paths for executor
- Test RunInteractive with mock stdin

### config Package (93.9% coverage)
**Partially covered:**
- `DiscoverFiles()` - 92.9% (missing edge cases)
- `ResolveProfile()` - 97.1% (missing some paths)
- `LoadConfig()` - 88.9% (error paths)
- `LoadConfigFromDirectory()` - 85.7% (edge cases)

**Strategy:**
- Add tests for:
  - Non-existent directories
  - Permission errors
  - Empty slices
  - Invalid YAML structures

### project Package (75.5% coverage)
**Uncovered functions:**
- `GetAlias()` - 0%
- `HasProject()` - 0%
- `ProjectNames()` - 0%
- `AliasNames()` - 0%
- `IsAtProjectReference()` - 0%
- `GetGlobalConfigPath()` - 75% (error path)
- `LoadGlobalConfigOrDefault()` - 83.3%
- `GetHistoryPath()` - 75%
- `LoadHistory()` - 81.8%
- `Save()` - 66.7%
- `NewHistoryEntry()` - 0%

**Strategy:**
- Add straightforward unit tests for helper functions
- Test error paths for file operations

## New Test Fixtures Needed

1. **fixture: empty-profiles**
   - do.yaml with empty profiles map
   - Tests nil/empty map handling

2. **fixture: duplicate-services**
   - Same service in multiple files
   - Tests override behavior

3. **fixture: nested-includes**
   - Complex include chains
   - Tests circular reference detection

4. **fixture: large-project**
   - 10+ compose files
   - Tests performance and ordering

5. **fixture: special-characters**
   - Paths with spaces, unicode
   - Tests path handling

6. **fixture: mixed-extensions**
   - Both .yaml and .yml files
   - Tests extension preference

7. **fixture: complex-aliases**
   - Nested aliases, quoted commands
   - Tests alias parsing edge cases

8. **fixture: hook-failures**
   - Hooks that fail
   - Tests error propagation

## Test Organization

```
cmd/
  cmd_test.go           # Integration tests for commands
  alias_test.go         # Alias command tests
  compose_test.go       # Compose command tests

internal/compose/
  builder_test.go       # Existing, expand
  executor_test.go      # Existing, expand error paths
  integration_test.go   # New: end-to-end tests

internal/config/
  config_test.go        # Existing, expand edge cases
  parser_test.go        # Existing, expand

internal/project/
  global_test.go        # Existing, add missing tests
  history_test.go       # Existing, add missing tests
  remote_test.go        # Existing
```

## Testing Strategy

### 1. Unit Tests
- Test individual functions in isolation
- Use table-driven tests for variations
- Mock external dependencies

### 2. Integration Tests
- Test command execution end-to-end
- Use dry-run mode to avoid actual docker calls
- Verify correct command assembly

### 3. Edge Case Tests
- Empty inputs
- Nil pointers
- Invalid paths
- Malformed YAML
- Unicode and special characters

### 4. Error Path Tests
- File not found
- Permission denied
- Invalid configuration
- Parse errors

## Execution Order

1. Add new fixtures
2. Add missing unit tests for project package (quick wins)
3. Add cmd package tests
4. Expand compose package tests
5. Expand config package tests
6. Verify 100% coverage
7. Add race detection
8. Add benchmarks for any slow tests
