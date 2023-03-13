package music

import "testing"

func TestSpotifyCurrentlyPlaying(t *testing.T) {
	track, err := SpotifyCurrentlyPlaying()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("track: %+v", track)
}
