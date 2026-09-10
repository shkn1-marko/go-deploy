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
		if err := Make(os.Args[2]); err != nil {
			fail(err)
		}
		fmt.Printf("created /etc/gdep/%s\n", os.Args[2])

	case "remove":
		requireName()
		if err := Remove(os.Args[2]); err != nil {
			fail(err)
		}
		fmt.Printf("removed /etc/gdep/%s\n", os.Args[2])

	case "list":
		names, err := List()
		if err != nil {
			fail(err)
		}
		if len(names) == 0 {
			fmt.Println("no projects registered")
			return
		}
		for _, n := range names {
			fmt.Println(n)
		}

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
