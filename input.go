package main

import (
	"os"
)

var argCount = map[string]int{
	"make":   1,
	"remove": 1,
	"list":   0,
	"deploy": 1,
}

func Input(args []string) {
	if len(args) < 2 {
		OutErr("usage: gdep <command> [name]")
		os.Exit(1)
	}

	command := args[1]

	n, ok := argCount[command]
	if !ok {
		OutErr("unknown command")
		os.Exit(1)
	}
	if len(args) != 2+n {
		OutErr("'%s' takes %d arguments", command, n)
		os.Exit(1)
	}

	switch command {
	case "make":
		Make(args[2])
	case "remove":
		Remove(args[2])
	case "list":
		List()
	case "deploy":
		Deploy(args[2])
	default:
		OutErr("'%s' is known but undefined", command)
		os.Exit(1)
	}
}
