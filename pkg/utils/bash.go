package utils

import (
	"context"
	"os/exec"
	"time"
)

func LaunchCommand(command string, timeout time.Duration) ([]byte, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	output, err := cmd.Output()

	// Call cancel after command execution
	cancel()

	return output, err
}
