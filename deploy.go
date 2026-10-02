package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const scriptTimeout = 5 * time.Minute

func RunDeploy(name string, store *DeviceStore, sender *FCMSender) error {
	dir := filepath.Join(deploymentsRoot, name)

	if _, err := os.Stat(dir); err != nil {
		return nil
	}

	status := DeployStatus{
		Name:      name,
		Build:     StageFailed,
		Deploy:    StageSkipped,
		Timestamp: time.Now().Unix(),
	}

	buildOutput, buildErr := runScript(dir, "build.sh")
	if buildErr != nil {
		status.Cause = buildErr.Error()
		status.Output = buildOutput
		notify(store, sender, status)
		return fmt.Errorf("build.sh failed: %w", buildErr)
	}
	status.Build = StageOK

	deployOutput, deployErr := runScript(dir, "deploy.sh")
	if deployErr != nil {
		status.Deploy = StageFailed
		status.Cause = deployErr.Error()
		status.Output = deployOutput
		notify(store, sender, status)
		return fmt.Errorf("deploy.sh failed: %w", deployErr)
	}
	status.Deploy = StageOK

	notify(store, sender, status)
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

func notify(store *DeviceStore, sender *FCMSender, status DeployStatus) {
	if store == nil || sender == nil {
		return
	}

	device, ok := store.Current()
	if !ok {
		return
	}

	sender.SendDeployStatus(device, status)
}
