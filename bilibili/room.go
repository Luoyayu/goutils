package biliAPI

import (
	"fmt"
	"net/http"
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

func RoomInit(roomID interface{}) (rett *Room, err error) {
	if ret, err := GetDefault(Config.API.RoomInit, map[string]interface{}{"id": roomID}, &Room{}); err == nil {
		rett = ret.(*Room)
	}
	return
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

func GetRoomPlayUrl(roomID interface{}) (rett *RoomPlayUrl, err error) {
	if ret, err := GetDefault(Config.API.RoomPlayUrl, map[string]interface{}{
		"cid": roomID,
	}, &RoomPlayUrl{}); err == nil {
		rett = ret.(*RoomPlayUrl)
	}
	return
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

func GetRoomNews(roomId interface{}, uid interface{}) (rett *RoomNews, err error) {
	if ret, err := GetDefault(Config.API.RoomNewsGet, map[string]interface{}{
		"roomid": roomId, "uid": uid,
	}, &RoomNews{}); err == nil {
		rett = ret.(*RoomNews)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type LiveMyFollowingRet struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    *LiveMyFollowingData `json:"data"`
}

type LiveMyFollowingData struct {
	Title     string `json:"title"`
	PageSize  int    `json:"pageSize"`
	TotalPage int    `json:"totalPage"`
	List      []*struct {
		RoomId     int64  `json:"roomid"`
		Uid        int64  `json:"uid"`
		UName      string `json:"uname"`
		Title      string `json:"title"`
		Face       string `json:"face"`
		LiveStatus int    `json:"live_status"`
	} `json:"list"`
}

// 10 entities per page
func GetLiveMyFollowing(uid int64, SESSDATA string, pn, ps int) (rett *LiveMyFollowingRet, err error) {
	if ps <= 0 {
		ps = 20
	}

	req := &http.Request{}
	req.Header.Set("Cookie", fmt.Sprintf("INTVER=1;DedeUserID=%v;SESSDATA=%v", uid, SESSDATA))

	ret, err := GetWithReq(Config.API.LiveMyFollowing, map[string]interface{}{"page": pn, "page_size": ps}, req, &LiveMyFollowingRet{})
	if err == nil {
		rett = ret.(*LiveMyFollowingRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type LiveUserRecommendRet struct {
	Code int `json:"code"` // 0 / 404
	//Msg     string               `json:"msg"`     // always "ok"
	//Message string               `json:"message"` // always "ok"
	Data []*struct {
		Face        string `json:"face"`
		Link        string `json:"link"`
		Online      int    `json:"online"`
		RoomId      int64  `json:""`
		SystemCover string `json:"system_cover"`
		Title       string `json:"title"`
		Uid         int64  `json:"uid"`
		UName       string `json:"uname"`
		UserCover   string `json:"user_cover"`
	} `json:"data"`
}

// ps = 30
func GetLiveUserRecommend(uid interface{}, SESSDATA interface{}, pn int) (rett *LiveUserRecommendRet, err error) { // need login
	req := &http.Request{}
	req.Header.Set("Cookie", fmt.Sprintf("INTVER=1;DedeUserID=%v;SESSDATA=%v", uid, SESSDATA))

	ret, err := GetWithReq(Config.API.LiveGetUserRecommend, map[string]interface{}{"page": pn}, req, &LiveUserRecommendRet{})
	if err == nil {
		rett = ret.(*LiveUserRecommendRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type GetRoomInfoByRoomIdRet struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RoomInfo   RoomInfo   `json:"room_info"`
		AnchorInfo AnchorInfo `json:"anchor_info"`
	} `json:"data"`
}

type RoomInfo struct {
	Uid            int64  `json:"uid"`
	RoomId         int64  `json:"room_id"`
	ShortId        int    `json:"short_id"`
	Title          string `json:"title"`
	Cover          string `json:"cover"`
	Tags           string `json:"tags"`
	LiveStatus     int    `json:"live_status"`
	LiveStartTime  int64  `json:"live_start_time"`
	AreaName       string `json:"area_name"`
	ParentAreaName string `json:"parent_area_name"`
	Keyframe       string `json:"keyframe"`
	Online         int    `json:"online"`
}

type AnchorInfo struct {
	BaseInfo *struct {
		UName  string `json:"uname"`
		Face   string `json:"face"`
		Gender string `json:"gender"`
	} `json:"base_info"`
}

func GetRoomInfoByRoomId(roomId interface{}) (rett *GetRoomInfoByRoomIdRet, err error) {
	if ret, err := GetDefault(Config.API.RoomGetInfoByRoomID, map[string]interface{}{
		"room_id": roomId,
	}, &GetRoomInfoByRoomIdRet{}); err == nil {
		rett = ret.(*GetRoomInfoByRoomIdRet)
	}
	return
}
