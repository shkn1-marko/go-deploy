package main

import (
	"fmt"
	"os"
)

func OutOK(format string, args ...any) {
	fmt.Fprintf(os.Stdout, "[ok] "+format+"\n", args...)
}

func OutErr(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[error] "+format+"\n", args...)
}
