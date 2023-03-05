package main

import "fmt"

// handleNormalMode handles the normal mode of operation for a screen.
//
// In this mode, the shell is waiting for a command to be entered. The command
// can be a shell command or a command to the editor.
func handleNormalMode(screen *Screen) error {
	// Read the next key from the terminal.
	k, err := screen.ReadKey()
	if err != nil {
		return fmt.Errorf("failed to read key from screen: %w", err)
	}

	if k == KeyCtrlC {
		return fmt.Errorf("user interrupted")
	}

	if k == KeyCtrlV {
		// TODO: Handle the insert/visual mode better?
		// screen.Mode = Insert

		// Write ANSI escape sequence for highlighting the current line.
		_, err := screen.WriteString("\x1b[7m")
		if err != nil {
			return fmt.Errorf("failed to write ANSI escape sequence to screen: %w", err)
		}

		return nil
	}

	_, err = screen.WriteString(k.String())
	if err != nil {
		return fmt.Errorf("failed to write key to screen: %w", err)
	}

	return nil
}
