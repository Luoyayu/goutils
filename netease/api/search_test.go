package api

import (
	"log"
	"testing"
)

func TestNetEaseClient_SearchSongByName(t *testing.T) {
	ret := _client.SearchSongByName(`Don't say "lazy"`)
	log.Printf("%+v\n", ret.Result.Songs[0])
}

func TestNetEaseClient_SearchSongByNameAndAlbum(t *testing.T) {
	ret := _client.SearchSongByNameAndAlbum(`Don't say "lazy"`, "K-ON! Music History's Box")
	log.Printf("%+v\n", ret)
}
