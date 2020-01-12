package biliAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// order: [(desc), ]
// 50 entities / page
func GetRelationFollowings(vmid int64, pn int, order string) (*relationFollowingsRetStruct, error) {
	if order == "" {
		order = "desc"
	}
	params := map[string]interface{}{
		"vmid":  vmid,
		"order": order,
		"pn":    pn,
		"ps":    50,
	}

	var err error
	l := url.Values{}
	resp := &http.Response{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}
	if resp, err = http.Get(Config.API.RelationFollowings + "?" + l.Encode()); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		ret := &relationFollowingsRetStruct{}
		err = json.Unmarshal(b, &ret)
		return ret, err
	}
	return nil, err
}

type relationFollowingsRetStruct struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    *relationFollowingsDataStruct `json:"data"`
}
type relationFollowingsDataStruct struct {
	List []*relationFollowingsUserStruct `json:"list"`
}

type relationFollowingsUserStruct struct {
	Mid            int64                 `json:"mid"`
	Attribute      int                   `json:"attribute"`
	MTime          int64                 `json:"mtime"`
	Special        int                   `json:"special"`
	Uname          string                `json:"uname"`
	Face           string                `json:"face"`
	Sign           string                `json:"sign"`
	OfficialVerify *officialVerifyStruct `json:"official_verify"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// following total number
func GetRelationStat(vmid int64) (*RelationStatRetStruct, error) {
	params := map[string]interface{}{
		"vmid":  vmid,
		"jsonp": "jsonp",
	}
	l := url.Values{}
	var err error
	resp := &http.Response{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}
	if resp, err = http.Get(Config.API.RelationStat + "?" + l.Encode()); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		ret := &RelationStatRetStruct{}
		err = json.Unmarshal(b, &ret)
		return ret, err
	}
	return nil, err

}

type RelationStatRetStruct struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    *RelationStatDataStruct `json:"data"`
}

type RelationStatDataStruct struct {
	Mid       int64 `json:"mid"`
	Following int   `json:"following"`
	Whisper   int   `json:"whisper"`
	Black     int   `json:"black"`
	Follower  int   `json:"follower"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
