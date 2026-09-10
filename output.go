package main

import (
	"fmt"
	"os"
)

func OutErr(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[error] "+format+"\n", args...)
}
