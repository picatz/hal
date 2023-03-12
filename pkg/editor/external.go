package editor

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type ExternalFinishedMsg struct {
	Err    error
	Buffer []byte
}

func ConfiguredExternalCommand() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return editor
}

func OpenExternal(buffer string) tea.Cmd {
	editor := ConfiguredExternalCommand()

	// Write to a temp file and open it
	f, err := os.CreateTemp(os.TempDir(), "hal-editor-*")
	if err != nil {
		panic(err)
	}

	f.WriteString(buffer)
	f.Sync()

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

		return ExternalFinishedMsg{
			Err:    nil,
			Buffer: buf[:n],
		}
	})
}
