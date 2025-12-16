# run_silent

Executes commands silently, showing output only on failure. Useful for AI coding agents to reduce token usage from noisy commands such as package managers, test and build scripts.

## Usage

```bash
run_silent [-d description] [-t timeout] <command> [args...]
```

### Agent Rule

Here is the rule I add to my AI coding agents (e.g. CLAUDE.md) to get them to use `run_silent`:

```
<CLI_COMMANDS>
- IMPORTANT: Use the `run_silent` command wrapper (a command you run prefixed before the command you want to run) to reduce token usage by buffering stdout/stderr and only showing them on non-zero exit.
  - You MUST use run_silent to wrap any bash / CLI commands unless you need to see all the stdout.
  - Good commands to prefix with run_silent include package installs, builds, tests, linting etc...
  - Examples:
    - run_silent pnpm install
    - run_silent cargo check
    - run_silent make lint
</CLI_COMMANDS>
```

The result being that the agent only sees output on failure (non-zero exit), saving valuable context tokens.

_Tip: Set your linters to exit non-zero on warnings to catch those too!_

## Install

```bash
go install github.com/sammcj/run_silent@HEAD
```

Or build from source:

```bash
make build
make install
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

## Why

Inspiration taken from [Context-Efficient Backpressure for Coding Agents](https://www.humanlayer.dev/blog/context-efficient-backpressure):

> Every line of `PASS src/utils/helper.test.ts` is waste. When all tests pass, we waste context conveying what could take under 10 tokens.
