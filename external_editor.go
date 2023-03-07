package main

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct {
	err      error
	buffer   []byte
	viewport bool
}

func configuredEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	return editor
}

func openEditor(buffer string, viewport bool) tea.Cmd {
	editor := configuredEditor()

	// Write to a temp file and open it
	f, err := os.CreateTemp("", "hal-editor-*")
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
			return editorFinishedMsg{err: err}
		}

		// Read the file back in
		f, err := os.Open(f.Name())
		if err != nil {
			return editorFinishedMsg{err: err}
		}

		buf := make([]byte, 4096)
		n, err := f.Read(buf)
		if err != nil {
			return editorFinishedMsg{err: err}
		}

		return editorFinishedMsg{
			err:      nil,
			buffer:   buf[:n],
			viewport: viewport,
		}
	})
}
