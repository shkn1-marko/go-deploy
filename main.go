package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const baseDir = "/etc/gdep"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: gdep <command> <name>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "make":
		makeProject(os.Args[2])
	default:
		fmt.Println("unknown command:", os.Args[1])
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
