package biliAPI

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RelationFollowingsRet struct {
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

// order: [(desc), ]
func GetRelationFollowings(vmid int64, pn int, order string) (rett *RelationFollowingsRet, err error) {
	if order == "" {
		order = "desc"
	}

	if ret, err := GetDefault(Config.API.RelationFollowings, map[string]interface{}{
		"vmid":  vmid,
		"order": order,
		"pn":    pn,
		"ps":    50,
	}, &RelationFollowingsRet{}); err == nil {
		rett = ret.(*RelationFollowingsRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type RelationStatRet struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Mid       int64 `json:"mid"`
		Following int   `json:"following"`
		Whisper   int   `json:"whisper"`
		Black     int   `json:"black"`
		Follower  int   `json:"follower"`
	} `json:"data"`
}

// following total number
func GetRelationStat(vmid int64) (rett *RelationStatRet, err error) {
	if ret, err := GetDefault(Config.API.RelationStat, map[string]interface{}{
		"vmid":  vmid,
		"jsonp": "jsonp",
	}, &RelationStatRet{}); err == nil {
		rett = ret.(*RelationStatRet)
	}
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
