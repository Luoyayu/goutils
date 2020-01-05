package api

import (
	"encoding/json"
)

func (r *NetEaseClient) GetLyric(songId int64) *LyricStruct {
	params := map[string]interface{}{
		"e_r": false,
		"tv":  "-1", "lv": "-1", "kv": "-1",
		"id": songId,
	}

	_ = r.NewEnc("/api/song/lyric", params)
	if ret, err := r.DoPost(); err == nil {
		var lyric = &LyricStruct{}
		err = json.Unmarshal(ret, &lyric)
		//log.Printf("%+v\n", lyric)
		return lyric
	}
	return nil
}

type LyricUser struct {
	Sgc      bool   `json:"sgc"`
	Sfy      bool   `json:"sfy"`
	Qfy      bool   `json:"qfy"`
	Id       int64  `json:"id"`
	Userid   int64  `json:"userid"`
	Nickname string `json:"nickname"`
	Uptime   int64  `json:"uptime"`
}

type Lyric struct {
	Version int64  `json:"version"`
	Lyric   string `json:"lyric"`
}

type LyricStruct struct {
	LyricUser LyricUser `json:"lyricUser"`
	TransUser LyricUser `json:"transUser"`
	Lrc       Lyric     `json:"lrc"`
	Klyric    Lyric     `json:"klyric"`
	Tlyric    Lyric     `json:"tlyric"`
}
