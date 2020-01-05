package api

import (
	"testing"
)

func TestNetEaseClient_GetLyric(t *testing.T) {
	_client.GetLyric(_TestSongNoLyric)
	_client.GetLyric(_TestSongPayedID)
}
