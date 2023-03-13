package music

import (
	"fmt"
	"os/exec"
	"strings"
)

type Track struct {
	Name   string
	Artist string
}

// SpotifyCurrentlyPlaying returns the currently playing artist and track name from Spotify.
func SpotifyCurrentlyPlaying() (*Track, error) {
	// TODO: consider a cache of the last track, so we don't have to
	// call out to Spotify every time to get the artist.

	cmd := exec.Command("osascript", "-e", "tell application \"Spotify\" to artist of current track as string")
	artist, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("music: spotify: failed to get current artist: %w", err)
	}

	cmd = exec.Command("osascript", "-e", "tell application \"Spotify\" to name of current track as string")
	track, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("music: spotify: failed to get current track name: %w", err)
	}

	return &Track{
		Name:   strings.TrimSpace(string(track)),
		Artist: strings.TrimSpace(string(artist)),
	}, nil
}
