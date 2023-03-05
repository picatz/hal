package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// startPrompt enables a shell and text editor, but with a more natural language
// interface. It handles changing modes, and executing commands.
func startPrompt(fh *os.File) error {
	// Ensure that the terminal is in raw mode so that we can read the input
	// character by character, and not line by line. This is necessary for
	// the editor to work properly, for example, when moving the cursor.
	//
	// It is important to provide a great user experience, even if it
	// means that we have to do some low-level terminal manipulation.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Start a screen in normal mode.
	screen := NewScreen(fh, Normal)

	// Start the main loop.
	for {
		switch screen.Mode {
		case Normal:
			err = handleNormalMode(screen)
			// case Insert:
			// 	err = handleInsertMode(screen)
			// case Visual:
			// 	err = handleVisualMode(screen)
		}

		if err != nil {
			return fmt.Errorf("error handling %s mode: %w", screen.Mode, err)
		}
	}
}
