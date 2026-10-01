package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const scriptTimeout = 1 * time.Minute

func RunDeploy(device, name string) error {
	dir := filepath.Join(deploymentsRoot, name)

	if _, err := os.Stat(dir); err != nil {
		return nil
	}

	for _, script := range []string{"build.sh", "deploy.sh"} {
		output, err := runScript(dir, script)
		if err != nil {
			if device != "" {
				EmailERR(name, script, output, err)
			}
			return fmt.Errorf("%s failed: %w", script, err)
		}
	}

	if device != "" {
		EmailOK(name)
	}

	return nil
}

func runScript(dir, script string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), scriptTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, filepath.Join(dir, script))
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	return string(output), err
}
