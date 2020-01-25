package biliAPI

func GetUserInfo(mid int64) (rett *UserInfoRet, err error) {
	if ret, err := GetDefault(Config.API.UserInfo,
		map[string]interface{}{"mid": mid, "jsonp": "jsonp"}, &UserInfoRet{}); err == nil {
		rett = ret.(*UserInfoRet)
	}
	return
}

type UserInfoRet struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *UserInfoData `json:"data"`
}

type UserInfoData struct {
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
