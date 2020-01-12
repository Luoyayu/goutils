package biliAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Room struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    *struct {
		RoomId     int64 `json:"room_id"`
		ShortId    int64 `json:"short_id"`
		Uid        int64 `json:"uid"`
		LiveStatus int   `json:"live_status"`
		LiveTime   int64 `json:"live_time"`
	} `json:"data"`
}

func RoomInit(roomID interface{}) (*Room, error) {
	params := map[string]interface{}{
		"id": roomID,
	}
	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	if resp, err := http.Get(Config.API.RoomInit + "?" + l.Encode()); err == nil {
		var ret = &Room{}
		err = json.NewDecoder(resp.Body).Decode(&ret)
		defer resp.Body.Close()
		return ret, nil
	} else {
		return nil, err
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RoomPlayUrl struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Durl []*struct {
			Url string `json:"url"`
		} `json:"durl"`
	} `json:"data"`
}

func GetRoomPlayUrl(roomID interface{}) (*RoomPlayUrl, error) {
	params := map[string]interface{}{
		"cid": roomID,
	}

	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	if resp, err := http.Get(Config.API.RoomPlayUrl + "?" + l.Encode()); err == nil {
		var ret = &RoomPlayUrl{}
		err = json.NewDecoder(resp.Body).Decode(&ret)
		defer resp.Body.Close()
		return ret, nil
	} else {
		return nil, err
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RoomNews struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"` //map[string]interface{} -> map[string]string

	/*Data    *struct {
		Roomid interface{}  `json:"roomid"`
		Notice string `json:"content"`
		Uname  string `json:"uname"`
	} `json:"data"`*/
}

func GetRoomNews(roomId interface{}, uid interface{}) (*RoomNews, error) {
	if fmt.Sprint(roomId) == "" && fmt.Sprint(uid) == "" {
		log.Fatalln("roomId and uid all null!")
	}
	params := map[string]interface{}{
		"roomid": roomId,
		"uid":    uid,
	}
	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	if resp, err := http.Get(Config.API.RoomNewsGet + "?" + l.Encode()); err == nil {
		var ret = &RoomNews{}
		err = json.NewDecoder(resp.Body).Decode(&ret)
		if err != nil {
			return nil, err
		}

		//log.Printf("%#v\n", ret.Data)

		if ret.Code == 0 {
			//log.Printf("%q\n", ret.Data.(map[string]interface{})["roomid"])
		}
		defer resp.Body.Close()
		return ret, nil
	} else {
		return nil, err
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type LiveMyFollowingRet struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    *LiveMyFollowingData `json:"data"`
}

type LiveMyFollowingData struct {
	Title     string                   `json:"title"`
	PageSize  int                      `json:"pageSize"`
	TotalPage int                      `json:"totalPage"`
	List      []*LiveMyFollowingStruct `json:"list"`
}

type LiveMyFollowingStruct struct {
	RoomId     int64  `json:"roomid"`
	Uid        int64  `json:"uid"`
	UName      string `json:"uname"`
	Title      string `json:"title"`
	Face       string `json:"face"`
	LiveStatus int    `json:"live_status"`
}

// 10 entities per page
func GetLiveMyFollowing(uid int64, SESSDATA string, pn, ps int) (ret *LiveMyFollowingRet, err error) {
	if ps <= 0 {
		ps = 20
	}
	params := map[string]interface{}{
		"page":      pn,
		"page_size": ps,
	}
	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	req, _ := http.NewRequest("GET", Config.API.LiveMyFollowing+"?"+l.Encode(), nil)
	req.Header.Set("Cookie",
		fmt.Sprint("INTVER=1;DedeUserID=", uid, ";SESSDATA=", SESSDATA, ";"))

	var resp *http.Response
	if resp, err = (&http.Client{}).Do(req); err == nil {
		defer resp.Body.Close()
		ret = &LiveMyFollowingRet{}
		err = json.NewDecoder(resp.Body).Decode(&ret)
		if err != nil {
			return
		}
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type LiveUserRecommendRet struct {
	Code    int                  `json:"code"`    // 0 / 404
	//Msg     string               `json:"msg"`     // always "ok"
	//Message string               `json:"message"` // always "ok"
	Data    []*LiveUserRecommend `json:"data"`
}

type LiveUserRecommend struct {
	Face        string `json:"face"`
	Link        string `json:"link"`
	Online      int    `json:"online"`
	RoomId      int64  `json:""`
	SystemCover string `json:"system_cover"`
	Title       string `json:"title"`
	Uid         int64  `json:"uid"`
	UName       string `json:"uname"`
	UserCover   string `json:"user_cover"`
}

// ps = 30
func GetLiveUserRecommend(uid interface{}, SESSDATA interface{}, pn int) (ret *LiveUserRecommendRet, err error) { // need login
	params := map[string]interface{}{
		"page": pn,
	}
	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	req, _ := http.NewRequest("GET", Config.API.LiveGetUserRecommend+"?"+l.Encode(), nil)
	req.Header.Set("Cookie",
		fmt.Sprint("INTVER=1;DedeUserID=", uid, ";SESSDATA=", SESSDATA, ";"))

	var resp *http.Response
	if resp, err = (&http.Client{}).Do(req); err == nil {
		defer resp.Body.Close()
		ret = &LiveUserRecommendRet{}
		//b, _ := ioutil.ReadAll(resp.Body)
		//log.Println(string(b))
		err = json.NewDecoder(resp.Body).Decode(&ret)
		if err != nil {
			return
		}
	}
	return
}
