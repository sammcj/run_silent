package main

import (
	"testing"
	"time"
)

func TestRunSuccess(t *testing.T) {
	exitCode := run("", time.Second, "true")
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
}

func TestRunFailure(t *testing.T) {
	exitCode := run("", time.Second, "false")
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

func TestRunCustomExitCode(t *testing.T) {
	exitCode := run("", time.Second, "sh", "-c", "exit 42")
	if exitCode != 42 {
		t.Errorf("expected exit code 42, got %d", exitCode)
	}
}

func TestRunTimeout(t *testing.T) {
	exitCode := run("", 50*time.Millisecond, "sleep", "10")
	if exitCode != 124 {
		t.Errorf("expected exit code 124 for timeout, got %d", exitCode)
	}
}

func TestRunCommandNotFound(t *testing.T) {
	exitCode := run("", time.Second, "nonexistent-command-xyz")
	if exitCode == 0 {
		t.Error("expected non-zero exit code for missing command")
	}
}

func TestFormatCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		args     []string
		expected string
	}{
		{
			name:     "command only",
			cmd:      "ls",
			args:     nil,
			expected: "ls",
		},
		{
			name:     "command with args",
			cmd:      "go",
			args:     []string{"build", "-v"},
			expected: "go build -v",
		},
		{
			name:     "arg with spaces",
			cmd:      "echo",
			args:     []string{"hello world"},
			expected: `echo "hello world"`,
		},
		{
			name:     "arg with quotes",
			cmd:      "echo",
			args:     []string{`it's`},
			expected: `echo "it's"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCommand(tt.cmd, tt.args)
			if result != tt.expected {
				t.Errorf("formatCommand(%q, %v) = %q, want %q", tt.cmd, tt.args, result, tt.expected)
			}
		})
	}
}
