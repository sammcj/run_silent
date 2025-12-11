# run_silent

Executes commands silently, showing output only on failure. Useful for reducing noise when working with AI coding agents.

## Install

```bash
go install github.com/sammcj/run_silent@latest
```

Or build from source:

```bash
make build
```

## Usage

```bash
run_silent [-d description] [-t timeout] <command> [args...]
```

### Options

- `-d` - Custom description (default: the command itself)
- `-t` - Timeout duration (default: 5m)

### Behaviour

- **Success**: Prints `✓ <description>`, exits 0
- **Failure**: Prints captured output then `✗ <description>`, preserves exit code
- **Timeout**: Prints `⏱ <description> (timed out)`, exits 124

### Examples

```bash
# Run tests silently
run_silent -d "Running tests" go test ./...

# Build with timeout
run_silent -t 2m make build

# Chain multiple commands
run_silent -d "Lint" golangci-lint run && run_silent -d "Test" go test ./...
```

## Agent Rule

```markdown
The `run_silent` command wrapper reduces token usage by only providing the exit status and any stderr. You MUST use run_silent to wrap any CLI command that you do not truly need to see all the output from such as installs, builds, tests, linting etc... example: `run_silent pnpm install`.
```

## Why

From [Context-Efficient Backpressure for Coding Agents](https://www.humanlayer.dev/blog/context-efficient-backpressure):

> Every line of `PASS src/utils/helper.test.ts` is waste. When all tests pass, we waste context conveying what could take under 10 tokens.
