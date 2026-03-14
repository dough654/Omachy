package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// Result holds the output and exit status of a command.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Run executes a command and returns the combined result.
func Run(name string, args ...string) (Result, error) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Result{
				Stdout:   string(stdout),
				Stderr:   string(exitErr.Stderr),
				ExitCode: exitErr.ExitCode(),
			}, err
		}
		return Result{}, err
	}

	return Result{
		Stdout:   string(stdout),
		ExitCode: 0,
	}, nil
}

// RunStreaming executes a command and calls onLine for each line of combined output.
func RunStreaming(name string, args []string, onLine func(string)) error {
	cmd := exec.Command(name, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	var wg sync.WaitGroup
	scanLines := func(r io.Reader) {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			onLine(scanner.Text())
		}
	}

	wg.Add(2)
	go scanLines(stdoutPipe)
	go scanLines(stderrPipe)
	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// Which checks if a command exists in PATH.
func Which(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}
