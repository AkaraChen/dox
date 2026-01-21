package compose

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Executor executes commands
type Executor struct {
	DryRun bool
	Dir    string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

// NewExecutor creates a new command executor
func NewExecutor(dryRun bool) *Executor {
	return &Executor{
	 DryRun: dryRun,
	 Stdout: os.Stdout,
	 Stderr: os.Stderr,
	}
}

// SetDir sets the working directory for command execution
func (e *Executor) SetDir(dir string) {
	e.Dir = dir
}

// SetEnv sets environment variables for command execution
func (e *Executor) SetEnv(env []string) {
	e.Env = env
}

// RunCommand executes a single command
func (e *Executor) RunCommand(cmd []string) (string, error) {
	if e.DryRun {
	 output := FormatCommand(cmd)
	 fmt.Fprintf(e.Stdout, "%s\n", output)
	 return output, nil
	}

	if len(cmd) == 0 {
	 return "", fmt.Errorf("empty command")
	}

	var stdout, stderr bytes.Buffer

	// Create command
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = &stdout
	c.Stderr = &stderr
	c.Dir = e.Dir

	if len(e.Env) > 0 {
	 c.Env = append(os.Environ(), e.Env...)
	}

	err := c.Run()
	if err != nil {
	 return "", fmt.Errorf("command failed: %s\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// RunCommandWithOutput executes a command with custom output writers
func (e *Executor) RunCommandWithOutput(cmd []string, stdout, stderr io.Writer) error {
	if e.DryRun {
	 output := FormatCommand(cmd)
	 fmt.Fprintf(e.Stdout, "%s\n", output)
	 return nil
	}

	if len(cmd) == 0 {
	 return fmt.Errorf("empty command")
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Dir = e.Dir

	if len(e.Env) > 0 {
	 c.Env = append(os.Environ(), e.Env...)
	}

	return c.Run()
}

// RunCommands executes multiple commands sequentially
func (e *Executor) RunCommands(commands [][]string) (string, error) {
	var output strings.Builder

	for i, cmd := range commands {
	 if e.DryRun {
   line := FormatCommand(cmd)
   fmt.Fprintf(e.Stdout, "%s\n", line)
   output.WriteString(line)
   output.WriteString("\n")
	 } else {
   out, err := e.RunCommand(cmd)
   if err != nil {
    return output.String(), fmt.Errorf("command %d failed: %w", i+1, err)
   }
   output.WriteString(out)
	 }
	}

	return output.String(), nil
}

// RunInteractive executes a command with inherited stdio
func (e *Executor) RunInteractive(cmd []string) error {
	if e.DryRun {
	 output := FormatCommand(cmd)
	 fmt.Fprintf(e.Stdout, "%s\n", output)
	 return nil
	}

	if len(cmd) == 0 {
	 return fmt.Errorf("empty command")
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = e.Stdout
	c.Stderr = e.Stderr
	c.Dir = e.Dir

	if len(e.Env) > 0 {
	 c.Env = append(os.Environ(), e.Env...)
	}

	return c.Run()
}
