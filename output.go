package main

import (
	"fmt"
	"strings"
)

const bannerSize = 24

func banner(command string) (header, footer string) {
	label := fmt.Sprintf(" %s ", strings.ToUpper(command))

	if len(label)%2 != 0 {
		label += " "
	}

	pad := bannerSize - len(label)
	left := pad / 2
	right := pad - left

	header = strings.Repeat("=", left) + label + strings.Repeat("=", right)
	footer = strings.Repeat("=", bannerSize)
	return
}

func OutputStatus(command string, err error, message string, details ...string) {
	header, footer := banner(command)
	fmt.Println(header)
	if err != nil {
		fmt.Printf("[error] %s\n", err)
	} else {
		fmt.Printf("[ok] %s\n", message)
		for _, d := range details {
			fmt.Println("    " + d)
		}
	}
	fmt.Println(footer)
}

func OutputList(command string, items []string) {
	header, footer := banner(command)
	fmt.Println(header)
	if len(items) == 0 {
		fmt.Println("[empty] 0 projects found.")
	}
	for _, item := range items {
		fmt.Println(" " + item)
	}
	fmt.Println(footer)
}
