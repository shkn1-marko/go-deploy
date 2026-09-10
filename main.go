package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const baseDir = "/etc/gdep"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: gdep <command> [name]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "make":
		requireName()
		makeProject(os.Args[2])
	case "remove":
		requireName()
		removeProject(os.Args[2])
	case "list":
		listProjects()
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}
}

func requireName() {
	if len(os.Args) < 3 {
		fmt.Println("usage: gdep", os.Args[1], "<name>")
		os.Exit(1)
	}
}

func makeProject(name string) {
	dir := filepath.Join(baseDir, name)

	if _, err := os.Stat(dir); err == nil {
		fmt.Printf("project %q already exists\n", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	for _, script := range []string{"build.sh", "deploy.sh"} {
		path := filepath.Join(dir, script)
		content := "#!/bin/bash\n# exit 0 on success, nonzero on failure\n"
		if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}

	fmt.Printf("created %s\n", dir)
}

func removeProject(name string) {
	dir := filepath.Join(baseDir, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("project %q does not exist\n", name)
		os.Exit(1)
	}

	if err := os.RemoveAll(dir); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("removed %s\n", dir)
}

func listProjects() {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("no projects registered")
		return
	}

	for _, e := range entries {
		if e.IsDir() {
			fmt.Println(e.Name())
		}
	}
}
