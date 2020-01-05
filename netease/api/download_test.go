package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _client = &NetEaseClient{}

func init() {
	_client.NewDefaultHttpClient()
}

const (
	_TestSongFreeID  = 537854740 //free to download DownloadRAW br
	_TestSongPayedID = 1406448632
	_TestSongNoLyric = 34528129
)

func TestGetFreeSongRAWDownloadUrl(t *testing.T) {
	SINGED = false
	DEBUG = true
	song := SongStruct{Id: _TestSongFreeID}
	ret := _client.GetSongDownloadUrl(fmt.Sprint(song.Id), DownloadRAW)
	if ret == nil {
		panic("TestGetFreeSongRAWDownloadUrl failed!")
	}
	if ret.Data.Url == "" {
		panic("GetSongDownloadUrl Failed!, url is null")
	}
	fmt.Printf("free song download url: %v\n", ret.Data.Url)
}

func TestGetPayedSongRAWDownloadUrl(t *testing.T) {
	assert.Equal(t, SINGED, true, "must singed!")
	//DEBUG = true
	song := SongStruct{Id: _TestSongPayedID}
	ret := _client.GetSongDownloadUrl(fmt.Sprint(song.Id), DownloadRAW)
	if ret == nil {
		panic("TestGetFreeSongRAWDownloadUrl failed!")
	}
	if ret.Data.Url == "" {
		panic("GetSongDownloadUrl Failed!, url is null")
	}
	fmt.Printf("free song download url: %v\n", ret.Data.Url)
}

func TestDownloadSong(t *testing.T) {
	song := SongStruct{Id: 560079}
	if err := song.Download(_client, "", -1); err != nil {
		panic(err)
	}
}

func TestDownloadSongs(t *testing.T) {
	songs := SongsStruct{
		Songs: []SongStruct{
			{Id: 560079},
			{Id: 1387099649},
			{Id: 514774419},
		},
	}

	if err := songs.Download(_client, "", nil); err != nil {
		panic(err)
	}
}
