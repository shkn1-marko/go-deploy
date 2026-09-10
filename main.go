package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fail(fmt.Errorf("usage: gdep <command> [name]"))
	}

	switch os.Args[1] {
	case "make":
		requireName()
		name := os.Args[2]
		err := Make(name)
		OutputStatus("make", err, fmt.Sprintf("created /etc/gdep/%s", name), "build.sh", "deploy.sh")
		if err != nil {
			os.Exit(1)
		}

	case "remove":
		requireName()
		name := os.Args[2]
		err := Remove(name)
		OutputStatus("remove", err, fmt.Sprintf("removed /etc/gdep/%s", name))
		if err != nil {
			os.Exit(1)
		}

	case "list":
		names, err := List()
		if err != nil {
			fail(err)
		}
		OutputList("list", names)

	default:
		fail(fmt.Errorf("unknown command: %s", os.Args[1]))
	}
}

func requireName() {
	if len(os.Args) < 3 {
		fail(fmt.Errorf("usage: gdep %s <name>", os.Args[1]))
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
