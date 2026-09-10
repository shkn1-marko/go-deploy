package main

import (
	"fmt"
	"os"
)

func OutOK(format string, args ...any) {
	fmt.Fprintf(os.Stdout, "[ok] "+format+"\n", args...)
}

func OutEmpty() {
	fmt.Fprintln(os.Stdout, "[empty] no deployments found")
}

func OutList(items []string) {
	for i, item := range items {
		fmt.Fprintf(os.Stdout, "[%d] %s\n", i, item)
	}
}

func OutErr(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[error] "+format+"\n", args...)
}
