// run_silent executes commands silently, showing output only on failure.
// Designed to reduce noise when working with AI coding agents.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)

const defaultTimeout = 5 * time.Minute

func main() {
	var (
		description    string
		timeout        time.Duration
		showVersion    bool
	)

	flag.StringVar(&description, "d", "", "description to show instead of command")
	flag.DurationVar(&timeout, "t", defaultTimeout, "command timeout")
	flag.BoolVar(&showVersion, "v", false, "show version and exit")
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: run_silent [-d description] [-t timeout] [-v] <command> [args...]")
		os.Exit(1)
	}

	args := flag.Args()
	exitCode := run(description, timeout, args[0], args[1:]...)
	os.Exit(exitCode)
}

func run(description string, timeout time.Duration, name string, args ...string) int {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()
	defer signal.Stop(sigChan)

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()

	if description == "" {
		description = formatCommand(name, args)
	}

	// Handle timeout
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Print(output.String())
		fmt.Printf("  ⏱ %s (timed out)\n", description)
		return 124
	}

	// Success: just print checkmark
	if err == nil {
		fmt.Printf("  ✓ %s\n", description)
		return 0
	}

	// Failure: print output then status
	if output.Len() > 0 {
		fmt.Print(output.String())
	} else if _, isExitErr := err.(*exec.ExitError); !isExitErr {
		// No output and not a simple exit error - show the underlying error (e.g. command not found)
		fmt.Println(err)
	}
	fmt.Printf("  ✗ %s\n", description)

	// Extract exit code
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}

	return 1
}

func formatCommand(name string, args []string) string {
	if len(args) == 0 {
		return name
	}

	parts := make([]string, 0, len(args)+1)
	parts = append(parts, name)

	for _, arg := range args {
		if strings.ContainsAny(arg, " \t\n\"'") {
			parts = append(parts, fmt.Sprintf("%q", arg))
		} else {
			parts = append(parts, arg)
		}
	}

	return strings.Join(parts, " ")
}

func printVersion() {
	fmt.Printf("run_silent %s", Version)
	if Commit != "" {
		fmt.Printf(" (%s)", Commit)
	}
	if BuildDate != "" {
		fmt.Printf(" built %s", BuildDate)
	}
	fmt.Println()
}
