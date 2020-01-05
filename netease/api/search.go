package api

import (
	"encoding/json"
)

type SearchSongRet struct {
	Result struct {
		Songs []*recommendSongStruct `json:"songs"`
	} `json:"result"`
}

func (r *NetEaseClient) SearchSongByName(songName string) *SearchSongRet {
	params := map[string]interface{}{
		"e_r":   false,
		"limit": 1,
		"s":     songName,
	}
	_ = r.NewEnc("/api/search/suggest/web", params)
	result := &SearchSongRet{}
	if ret, err := r.DoPost(); err == nil {
		//log.Println(string(ret))
		err = json.Unmarshal(ret, &result)
	}
	return result
}

func (r *NetEaseClient) SearchSongByNameAndAlbum(songName, Album string) *recommendSongStruct {
	results := r.SearchSongByName(songName)
	for _, song := range results.Result.Songs {
		if song.Album.Name == Album {
			return song
		}
	}
	return nil
}
