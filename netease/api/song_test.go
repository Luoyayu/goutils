package api

import (
	"fmt"
	"testing"
)

func TestNetEaseClient_GetSongDetail(t *testing.T) {
	song := _client.GetSongDetail(_TestSongFreeID)
	fmt.Printf("%+v\n%+v\n%+v\n%+v\n", song, song.AL, song.AR, song.Privilege)
}
