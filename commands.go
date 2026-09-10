package main

import (
	"os"
	"path/filepath"
)

const gdepRoot = "/etc/gdep"

func Make(name string) {
	dir := filepath.Join(gdepRoot, name)

	if _, err := os.Stat(dir); err == nil {
		OutErr("'%s' already exists", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		OutErr("failed to create '%s': '%v'", name, err)
		os.Exit(1)
	}

	content := "#!/bin/sh\n"

	for _, file := range []string{"build.sh", "deploy.sh"} {
		if err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0755); err != nil {
			OutErr("failed to write %s: %v", file, err)
			os.Exit(1)
		}
	}

	OutOK("'%s' created", name)
}

func Remove(name string) {
	dir := filepath.Join(gdepRoot, name)

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			OutErr("'%s' does not exist", name)
		} else {
			OutErr("failed to check '%s': '%v'", name, err)
		}
		os.Exit(1)
	}

	if err := os.RemoveAll(dir); err != nil {
		OutErr("failed to remove '%s': %v", name, err)
		os.Exit(1)
	}

	OutOK("'%s' removed", name)
}

func List() {
	entries, err := os.ReadDir(gdepRoot)
	if err != nil {
		if os.IsNotExist(err) {
			OutEmpty()
			return
		}
		OutErr("failed to read '%s': %v", gdepRoot, err)
		os.Exit(1)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	if len(names) == 0 {
		OutEmpty()
		return
	}

	OutList(names)
}
