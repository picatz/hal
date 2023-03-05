package main

// Mode is the current mode of operation HAL is in.
type Mode int

const (
	// Normal mode is the default mode. In this mode, the shell is waiting for a
	// command to be entered. The command can be a shell command or a command to
	// the editor.
	Normal Mode = iota

	// Insert mode is used place text at the current cursor position.
	Insert

	// Visual mode is used to select a range of text. In this mode, the shell is
	// waiting for a command to be entered. The command can be a shell command or
	// a command to the editor. The difference between this mode and normal mode
	// is that the cursor is not moved. Instead, the cursor is used to select a
	// range of text.
	Visual
)

// String returns the string representation of the mode.
func (m Mode) String() string {
	switch m {
	case Normal:
		return "Normal"
	case Insert:
		return "Insert"
	case Visual:
		return "Visual"
	default:
		return "Unknown"
	}
}
