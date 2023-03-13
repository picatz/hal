package editor

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// ExternalFinishedMsg is a message that is sent when the external editor
// finishes (either successfully or with an error). The buffer is the contents
// of the file that was edited.
type ExternalFinishedMsg struct {
	Err    error
	Buffer []byte
}

// ConfiguredExternalCommand returns the external command that should be used
// to open the editor.
//
// This is either the value of the EDITOR environment variable, or "vim"
// if that is not set.
func ConfiguredExternalCommand() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return editor
}

// OpenExternal opens the external editor and returns a command that will
// wait for it to finish.
//
// The buffer is the initial contents of the file.
func OpenExternal(buffer string) tea.Cmd {
	// Get the external editor command.
	editor := ConfiguredExternalCommand()

	// Write to a temp file and open it
	f, err := os.CreateTemp(os.TempDir(), "hal-editor-*")
	if err != nil {
		panic(err)
	}

	f.WriteString(buffer)
	// f.Sync()

	c := exec.Command(editor, f.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()
		if err != nil {
			return ExternalFinishedMsg{Err: err}
		}

		// Read the file back in
		f, err := os.Open(f.Name())
		if err != nil {
			return ExternalFinishedMsg{Err: err}
		}

		// Read the file back in and return it.
		//
		// TODO: this is a bit of a hack, but we don't know how big the file
		// is going to be. We could use bufio.Scanner, but that's a bit
		// annoying to use. We could also use io.ReadAll, but that will
		// allocate a buffer that's the size of the file, which is not
		// ideal.
		buf := make([]byte, 4096)
		n, err := f.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				return ExternalFinishedMsg{
					Err:    nil,
					Buffer: []byte{},
				}
			}

			return ExternalFinishedMsg{Err: err}
		}

		// Return the buffer contents.
		return ExternalFinishedMsg{
			Err:    nil,
			Buffer: buf[:n],
		}
	})
}
