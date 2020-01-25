package biliAPI

import (
	"errors"
	"time"
)

// search use: <"keyword", "bili_user", "">
func GetSearchUser(keyword string, pageId int) (*SearchRet, error) {
	return GetSearch(keyword, "bili_user", "", pageId)
}

type SearchRet struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *searchData `json:"data"`
}

type searchData struct {
	Result []*searchResult `json:"result"`
}

type searchResult struct {
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

func GetSearch(keyword, searchType, order string, pageId int) (rett *SearchRet, err error) {
	if ret, err := GetDefault(Config.API.SearchWeb, map[string]interface{}{
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
	}, &SearchRet{}); err == nil {
		rett = ret.(*SearchRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SearchSuggestRet struct {
	Code   int                  `json:"code"`
	Result *searchSuggestResult `json:"result"`
}

type searchSuggestResult struct {
	Tags []*searchSuggestResultTag `json:"tag"`
}

type searchSuggestResultTag struct {
	Value string `json:"value"`
	Ref   int    `json:"ref"`
	Name  string `json:"name"`
	SpId  int    `json:"spid"`
}

// NOTICE: some term got no suggestion without error return!
func GetSearchSuggest(term string) (rett *SearchSuggestRet, err error) {
	if ret, err := GetDefault(Config.API.SearchWebSuggest, map[string]interface{}{
		"func":         "suggest",
		"suggest_type": "accurate",
		"sub_type":     "tag",
		"main_ver":     "v1",
		"term":         term,
		"_":            time.Now().UnixNano(),
		//"highlight":    "", // uncomment for <em class="suggest_high_light">term</em>
	}, &SearchSuggestRet{}); err == nil {
		rett = ret.(*SearchSuggestRet)
		if len(rett.Result.Tags) == 0 {
			err = errors.New("no result")
		}
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpaceArcSearchRet struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		List *struct {
			TList map[string]SpaceArcSearchTList `json:"tlist"`
			VList []*SpaceArcSearchVList         `json:"vlist"`
		} `json:"list"`
		Page *struct {
			Count int `json:"count"`
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
		} `json:"page"`
	} `json:"data"`
}

type SpaceArcSearchTList struct {
	Tid   int64  `json:"tid"`
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type SpaceArcSearchVList struct {
	Comment      int    `json:"comment"`
	TypeId       int    `json:"typeid"`
	Play         int    `json:"play"`
	Pic          string `json:"pic"`
	Description  string `json:"description"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	MId          int64  `json:"mid"`
	Created      int64  `json:"created"`
	Length       string `json:"length"`
	VideoReview  int    `json:"video_review"`
	Aid          int64  `json:"aid"`
	isPay        int
	isUnionVideo int
}

func GetSpaceArcSearch(mid interface{}, ps, pn, tid int) (rett *SpaceArcSearchRet, err error) {
	if ret, err := GetDefault(Config.API.SpaceArcSearch, map[string]interface{}{
		"mid":     mid,
		"ps":      ps,
		"pn":      pn,
		"order":   "pubdate",
		"jsonp":   "jsonp",
		"tid":     tid,
		"keyword": "",
	}, &SpaceArcSearchRet{}); err == nil {
		rett = ret.(*SpaceArcSearchRet)
	}
	return
}
