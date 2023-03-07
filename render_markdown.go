package main

import (
	"bytes"

	"github.com/charmbracelet/glamour"
)

func renderMarkdown(in []byte, wordWrap int) ([]byte, int, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80), // TODO: make this configurable, maybe dynamic?
	)
	if err != nil {
		return nil, 0, err
	}

	b, err := r.RenderBytes(in)
	if err != nil {
		return nil, 0, err
	}

	return b, bytes.Count(b, []byte{'\n'}), nil
}
