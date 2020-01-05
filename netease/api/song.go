package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 请求: /api/v3/song/detail
// 参数: 歌曲数组
type SongsStruct struct {
	Songs      []SongStruct    `json:"Songs"`
	Privileges []songPrivilege `json:"privileges"`
	Code       int64           `json:"code"`
}

type SongStruct struct {
	Name        string  `json:"name"`
	Id          int64   `json:"id"`
	Pop         float64 `json:"pop"`
	PublishTime int64   `json:"publishTime"`
	AR          []ar    `json:"ar"` // artists
	AL          ar      `json:"al"` // album
}

// 歌手/专辑
type ar struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`

	Pic int64 `json:"pic"` // 专辑图
}

// 私有信息
type songPrivilege struct {
	Id    int64 `json:"id"`
	Fee   int64 `json:"fee"`
	Payed int64 `json:"payed"`
	Dl    int64 `json:"dl"`    // 下载最大比特率
	Maxbr int64 `json:"maxbr"` // 最大比特率

}

func (r *SongsStruct) First() *SongStruct {
	if r.Songs != nil {
		return &r.Songs[0]
	}
	return nil
}

func (r *NetEaseClient) GetSongDetail(songId int64) *SongStruct {
	Songs := r.GetSongsDetail([]int64{songId})
	return &Songs.Songs[0]
}

// GetSongsDetail :获取歌曲详情
func (r *NetEaseClient) GetSongsDetail(songIds []int64) *SongsStruct {
	var c []map[string]interface{}
	for _, songId := range songIds {
		c = append(c, map[string]interface{}{
			"id": songId,
		})
	}

	cBytes, _ := json.Marshal(c) // NECESSARY, only if server could recognize the params!

	params := map[string]interface{}{
		"e_r": false,
		"c":   string(cBytes),
	}

	_ = r.NewEnc("/api/v3/song/detail", params)
	if ret, err := r.DoPost(); err == nil {
		songsInfo := SongsStruct{}
		err = json.Unmarshal(ret, &songsInfo)
		if songsInfo.Code == 200 {
			debugln("songs Info", songsInfo)
			return &songsInfo
		}
	} else {
		debugln("GetSongsDetail error", err)
	}
	return nil
}

func (r *SongStruct) Download(c *NetEaseClient, savePath string, br int64) error {
	if r.Id == 0 {
		return errors.New("song id is None")
	}
	r = c.GetSongDetail(r.Id)
	url := GetSongDownloadUrl(c, fmt.Sprint(r.Id), br)
	return DownloadSong(r.Name, url.Data.Type, url.Data.Url, savePath)
}

func (r *SongsStruct) Download(c *NetEaseClient, savePath string, br []int64) error {
	if br == nil {
		br = make([]int64, len(r.Songs))
	}
	ids := make([]int64, len(r.Songs))
	names := make([]string, len(r.Songs))
	types := make([]string, len(r.Songs))
	urls := make([]string, len(r.Songs))
	for i, song := range r.Songs {
		ids[i] = song.Id
	}
	r = c.GetSongsDetail(ids) // cover old SongsStruct
	for i, song := range r.Songs {
		names[i] = song.Name
		d := GetSongDownloadUrl(c, fmt.Sprint(song.Id), br[i])
		types[i] = d.Data.Type
		urls[i] = d.Data.Url
		debugln("song download info", ids[i], names[i], types[i], urls[i])
	}
	return DownloadSongs(names, types, urls, savePath)
}
