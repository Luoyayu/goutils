package api

import (
	"encoding/json"
	"errors"
)

func (r *NetEaseClient) GetRecommend() (recommend *recommendStruct, err error) {
	if Signed == false {
		err = errors.New("please login firstly")
		return
	}
	params := map[string]interface{}{
		"e_r":    false,
		"offset": 0,
		"limit":  20,
		"total":  true,
	}
	_ = r.NewEnc("/api/v1/discovery/recommend/songs", params)
	var ret []byte
	if ret, err = r.DoPost(); err == nil {
		//log.Println(string(ret))
		err = json.Unmarshal(ret, &recommend)
	}
	return
}

type recommendStruct struct {
	Recommend []recommendSongStruct `json:"recommend"`
}

type recommendSongStruct struct {
	Code       int64           `json:"code"`
	Name       string          `json:"name"`
	Id         int64           `json:"id"`
	Alias      []string        `json:"alias"`
	Fee        int64           `json:"fee"`
	Disc       string          `json:"disc"`
	No         int64           `json:"no"`
	Artists    []artistsStruct `json:"artists"`
	Album      albumStruct     `json:"album"`
	Reason     string          `json:"reason"`
	Privilege  songPrivilege   `json:"privilege"`
	Duration   int64           `json:"duration"`
	Popularity float64         `json:"popularity"`
	Score      float64         `json:"score"`
}

type artistsStruct struct {
	Name      string `json:"name"`
	Id        int64  `json:"id"`
	PicUrl    string `json:"picUrl"`
	Img1v1Url string `json:"img1v1Url"`
}

type albumStruct struct {
	Name        string `json:"name"`
	Id          int64  `json:"id"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	BlurPicUrl  string `json:"blurPicUrl"`
	PicUrl      string `json:"picUrl"`
	PublishTime int64  `json:"publishTime"`
	Description string `json:"description"`
	SubType     string `json:"subType"`
}
