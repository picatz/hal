package main

import (
	"fmt"
	"os"
)

// CursorBuffer is a buffer that stores the input characters, and a cursor to
// keep track of the current position in the buffer.
type CursorBuffer struct {
	Buffer []byte
	Cursor int
}

type Screen struct {
	// Mode is the current mode of operation for the screen.
	Mode

	// The file handle to the terminal.
	FH *os.File

	// Buffer to store the input characters, and a cursor to
	// keep track of the current position in the buffer.
	Input CursorBuffer

	// Buffer to store the output characters, and a cursor to
	// keep track of the current position in the buffer.
	Output CursorBuffer

	// Buffer to store the history of commands.
	History []string

	// Buffer to store the current command.
	CurrentCommand []byte
}

func NewScreen(fh *os.File, mode Mode) *Screen {
	return &Screen{
		Mode: mode,
		FH:   fh,
		Input: CursorBuffer{
			Buffer: make([]byte, 0, 1024),
			Cursor: 0,
		},
		Output: CursorBuffer{
			Buffer: make([]byte, 0, 1024),
			Cursor: 0,
		},
		History:        make([]string, 0, 1024),
		CurrentCommand: make([]byte, 0, 1024),
	}
}

// ReadBytes reads next byte from the terminal.
func (screen *Screen) ReadByte() (byte, error) {
	// Read the next character from the terminal.
	char := make([]byte, 1)
	_, err := screen.FH.Read(char)
	if err != nil {
		return 0, fmt.Errorf("failed to read byte from terminal screen: %w", err)
	}

	return char[0], nil
}

// ReadKey reads the next key from the terminal, reading multiple bytes if
// necessary to determine the key.
func (screen *Screen) ReadKey() (Key, error) {
	// Read the next character from the terminal.
	char, err := screen.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("failed to read key from terminal screen: %w", err)
	}

	// Handle backspace.
	if char == 0x7f {
		return KeyBackspace, nil
	}

	// Handle CTL+C.
	if char == 0x03 {
		return KeyCtrlC, nil
	}

	// Handle CTL+D.
	if char == 0x04 {
		return KeyCtrlD, nil
	}

	// Handle CTRL+V.
	if char == 0x16 {
		return KeyCtrlV, nil
	}

	// Handle CTL+Z.
	if char == 0x1a {
		return KeyCtrlZ, nil
	}

	// Handle enter key.
	if char == 0x0d {
		return KeyEnter, nil
	}

	// If the character is not a control character, then it is a normal
	// character.
	if char != 0x1b {
		// panic(fmt.Sprintf("unhandled key: %02x", char))
		return Key(char), nil
	}

	// Read the next character from the terminal.
	char, err = screen.ReadByte()
	if err != nil {
		return 0, err
	}

	// If the character is not a control character, then it is a normal
	// character.
	if char != 0x5b {
		return Key(char), nil
	}

	// Read the next character from the terminal.
	char, err = screen.ReadByte()
	if err != nil {
		return 0, err
	}

	switch char {
	case 0x41:
		return KeyUp, nil
	case 0x42:
		return KeyDown, nil
	case 0x43:
		return KeyRight, nil
	case 0x44:
		return KeyLeft, nil
	}

	// panic(fmt.Sprintf("unhandled key: %02x", char))

	return Key(char), nil
}

type Key int

const (
	// KeyEnter is the enter key.
	KeyEnter Key = iota

	// KeyBackspace is the backspace key.
	KeyBackspace

	// KeyDelete is the delete key.
	KeyDelete

	// KeyLeft is the left arrow key.
	KeyLeft

	// KeyRight is the right arrow key.
	KeyRight

	// KeyUp is the up arrow key.
	KeyUp

	// KeyDown is the down arrow key.
	KeyDown

	// KeyTab is the tab key.
	KeyTab

	// KeyEscape is the escape key.
	KeyEscape

	// KeyCtrlA is the control key and the a key.
	KeyCtrlA

	// KeyCtrlB is the control key and the b key.
	KeyCtrlB

	// KeyCtrlC is the control key and the c key.
	KeyCtrlC

	// KeyCtrlD is the control key and the d key.
	KeyCtrlD

	// KeyCtrlV is the control key and the v key.
	KeyCtrlV

	// KeyCtrlZ is the control key and the z key.
	KeyCtrlZ

	// TODO: Add more keys.
)

// String returns the string representation of the key that can be printed to the
// a terminal.
func (k Key) String() string {
	switch k {
	case KeyEnter:
		return "\r\n"
	case KeyBackspace:
		return "\b \b"
	case KeyDelete:
		return "\x1b[3~"
	case KeyLeft:
		return "\x1b[D"
	case KeyRight:
		return "\x1b[C"
	case KeyUp:
		return "\x1b[A"
	case KeyDown:
		return "\x1b[B"
	case KeyTab:
		return "\t"
	case KeyEscape:
		return "\x1b"
	case KeyCtrlA:
		return "\x01"
	case KeyCtrlB:
		return "\x02"
	case KeyCtrlC:
		return "\x03"
	case KeyCtrlD:
		return "\x04"
	case KeyCtrlV:
		return "\x16"
	case KeyCtrlZ:
		return "\x1a"
	default:
		// If the key is a normal character, then return the character.
		if k >= 0x20 && k <= 0x7e {
			return string(rune(k))
		}
	}

	return "<unknown key:" + string(rune(k)) + ">"
}
