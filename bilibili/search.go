package biliAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// search use: <"keyword", "bili_user", "">
func GetSearchUser(keyword string, pageId int) (*searchRetStruct, error) {
	return GetSearch(keyword, "bili_user", "", pageId)
}

func GetSearch(keyword, searchType, order string, pageId int) (*searchRetStruct, error) {
	params := map[string]interface{}{
		"context":       "",
		"search_type":   searchType,
		"page":          pageId,
		"order":         order,
		"keyword":       keyword,
		"category_id":   "",
		"user_type":     "",
		"order_sort":    "",
		"changing":      "mid",
		"highlight":     0,
		"single_column": 0,
		"jsonp":         "jsonp",
	}

	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}
	resp := &http.Response{}
	var err error
	if resp, err = http.Get(Config.API.SearchWeb + "?" + l.Encode()); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		ret := &searchRetStruct{}
		err = json.Unmarshal(b, &ret)
		return ret, err
	}
	return nil, err
}

type searchRetStruct struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *searchDataStruct `json:"data"`
}

type searchDataStruct struct {
	Result []*searchResultStruct `json:"result"`
}

type searchResultStruct struct {
	Mid            int64                 `json:"mid"`
	UName          string                `json:"uname"`
	USign          string                `json:"usign"`
	Fans           int                   `json:"fans"`
	VideosNum      int                   `json:"videos"`
	UPic           string                `json:"upic"`
	VerifyInfo     string                `json:"verify_info"`
	Level          int                   `json:"level"`
	Gender         int                   `json:"gender"`
	IsUpUser       int                   `json:"is_upuser"`
	IsLive         int                   `json:"is_live"`
	RoomId         int                   `json:"room_id"`
	OfficialVerify *officialVerifyStruct `json:"official_verify"`
}

type officialVerifyStruct struct {
	Type int    `json:"type"`
	Desc string `json:"desc"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// NOTICE: some term got no suggestion without error return!
func GetSearchSuggest(term string) (*searchSuggestRetStruct, error) {
	params := map[string]interface{}{
		"func":         "suggest",
		"suggest_type": "accurate",
		"sub_type":     "tag",
		"main_ver":     "v1",
		"term":         term,
		"_":            time.Now().UnixNano(),
		//"highlight":    "", // uncomment for <em class="suggest_high_light">term</em>
	}

	var err error
	l := url.Values{}
	resp := &http.Response{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}
	if resp, err = http.Get(Config.API.SearchWebSuggest + "?" + l.Encode()); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		ret := &searchSuggestRetStruct{}
		err = json.Unmarshal(b, &ret)
		return ret, err
	}
	return nil, err
}

type searchSuggestRetStruct struct {
	Code   int                        `json:"code"`
	Result *searchSuggestResultStruct `json:"result"`
}
type searchSuggestResultStruct struct {
	Tags []*searchSuggestResultTag `json:"tag"`
}
type searchSuggestResultTag struct {
	Value string `json:"value"`
	Ref   int    `json:"ref"`
	Name  string `json:"name"`
	SpId  int    `json:"spid"`
}
