package main

import (
	"os"
	"strings"
)

// This is an interactive shell and text editor that can be used
// to control and build systems with natural language.
//
// At its core, it is a graph-based modal shell and text editor.
//
// It is inspired by ChatGPT, BASH, FISH, ZSH, and VIM.

func main() {
	err := startPrompt(os.Stdin)
	if err != nil {
		if strings.Contains(err.Error(), "user interrupted") {
			os.Exit(1)
		}
		panic(err)
	}
}
