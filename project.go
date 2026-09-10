package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const baseDir = "/etc/gdep"

func Make(name string) error {
	dir := filepath.Join(baseDir, name)

	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("project %q already exists", name)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	for _, script := range []string{"build.sh", "deploy.sh"} {
		content := "#!/bin/bash\n# exit 0 on success, nonzero on failure\n"
		if err := os.WriteFile(filepath.Join(dir, script), []byte(content), 0o755); err != nil {
			return err
		}
	}

	return nil
}

func Remove(name string) error {
	dir := filepath.Join(baseDir, name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("project %q does not exist", name)
	}
	return os.RemoveAll(dir)
}

func List() ([]string, error) {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
