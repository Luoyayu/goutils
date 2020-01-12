package biliAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetUserInfo(mid int64) (*UserInfoRetStruct, error) {
	params := map[string]interface{}{
		"mid":   mid,
		"jsonp": "jsonp",
	}

	var err error
	l := url.Values{}
	resp := &http.Response{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}
	if resp, err = http.Get(Config.API.UserInfo + "?" + l.Encode()); err == nil {
		b, _ := ioutil.ReadAll(resp.Body)
		ret := &UserInfoRetStruct{}
		err = json.Unmarshal(b, &ret)
		return ret, err
	}
	return nil, err
}

type UserInfoRetStruct struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *UserInfoDataStruct `json:"data"`
}

type UserInfoDataStruct struct {
	Mid      int64               `json:"mid"`
	Name     string              `json:"name"`
	Sex      string              `json:"sex"`
	Face     string              `json:"face"`
	Sign     string              `json:"sign"`
	Level    int                 `json:"level"`
	JoinTime int64               `json:"jointime"`
	Moral    int                 `json:"moral"`
	Birthday string              `json:"birthday"`
	Official *UserOfficialStruct `json:"official"`
	TopPhoto string              `json:"top_photo"`
}

type UserOfficialStruct struct {
	Role  int    `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  int    `json:"type"`
}
