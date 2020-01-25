package biliAPI

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

/*
https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history?
access_key=&
actionKey=appkey&
appkey=&
build=12300&
device=pad&
from=space&
host_uid=25876945&
mobi_app=ipad&
offset_dynamic_id=0&
page=1&
platform=ios&
qn=32&
sign=&
statistics={"appId":5,"version":"2.4","abtest":"886","platform":2}&
ts=1576605048&
video_meta=qn:32,fnval:16,fnver:0,fourk:1,&
visitor_uid=12312321
*/

func GetLatestDynamicSpaceHistory(hostId int64) (*CardStruct, error) {
	if all, err := GetDynamicSpaceHistory(hostId, 0, true); err == nil {
		return all.Data.Cards[0], nil
	} else {
		return nil, err
	}
}

// 20 entities/page
func GetDynamicSpaceHistory(hostId interface{}, offsetDynamicId int64, latest bool) (data *DynamicRetStruct, err error) {
	params := map[string]interface{}{
		//"actionKey": "appkey",
		//"appkey":    Config.AppKey,
		//"from":      "space",
		//"mobi_app":  "android",
		//"build":     "5500300",
		//"platform":  "android",
		//"ts":        time.Now().Unix(),
		//"page":       2,
		"host_uid": hostId,
		//"access_key":        accessKey,
		"offset_dynamic_id": offsetDynamicId,
		//"visitor_uid":       visitorUid,
	}
	//params["sign"] = signParams(params)

	l := url.Values{}
	c := &http.Client{}
	req := &http.Request{}
	resp := &http.Response{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}
	if req, err = http.NewRequest("GET", Config.API.DynamicSpaceHistory+"?"+l.Encode(), nil); err == nil {
		req.Header.Add("user-agent", "Mozilla/5.0 BiliDroid/5.37.0")
		if resp, err = c.Do(req); err == nil {
			b, _ := ioutil.ReadAll(resp.Body)
			ret := &DynamicRetStruct{}
			err = json.Unmarshal(b, &ret)
			for _, card := range ret.Data.Cards {
				cardBytes := []byte(card.CardString)
				card.CartType = card.Desc.Type

				switch card.CartType {
				case 1: // 图文 + 引用< 动态, 投稿, 专栏, 房间, 投票 >
					card.CardContent.T1 = &Type1{}
					err = json.Unmarshal(cardBytes, &card.CardContent.T1)
					card.CardContent.T1.OriginType = card.Desc.OrigType
					card.CardContent.T1.OriginContent = &CardContentInterface{}
					originBytes := []byte(card.CardContent.T1.OriginString)

					switch card.CardContent.T1.OriginType {
					case 2: // 图文
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T2)
					case 4: // 图文
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T4)
					case 8: // 投稿
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T8)
					case 16: // 小视频
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T16)
					case 32:
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T32)
					case 64: // 专栏
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T64)
					case 256: // 创作
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T256)
					case 512: // 番剧
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T512)
					case 1024:
						// PASS
					case 2048: // sketch 宣传网页
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T2048)
					case 4099: // 番剧
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T4099)
					case 4200: // 房间
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T4200)
					case 4300: // 收藏夹
						err = json.Unmarshal(originBytes, &card.CardContent.T1.OriginContent.T4300)
					}

				case 2: // 图文
					err = json.Unmarshal(cardBytes, &card.CardContent.T2)
				case 4: // 图文
					err = json.Unmarshal(cardBytes, &card.CardContent.T4)
				case 8: // 投稿
					err = json.Unmarshal(cardBytes, &card.CardContent.T8)
				case 16: // 小视频
					err = json.Unmarshal(cardBytes, &card.CardContent.T16)
				case 32: //
					err = json.Unmarshal(cardBytes, &card.CardContent.T32)
				case 64: // 专栏
					err = json.Unmarshal(cardBytes, &card.CardContent.T64)
				case 256: // 创作
					err = json.Unmarshal(cardBytes, &card.CardContent.T256)
				case 512: // 番剧
					err = json.Unmarshal(cardBytes, &card.CardContent.T512)
				case 2048: // sketch 宣传网页
					err = json.Unmarshal(cardBytes, &card.CardContent.T2048)
				case 4099: // 番剧
					err = json.Unmarshal(cardBytes, &card.CardContent.T4099)
				case 4200: // 房间
					err = json.Unmarshal(cardBytes, &card.CardContent.T4200)
				case 4300: // 收藏夹
					err = json.Unmarshal(cardBytes, &card.CardContent.T4300)
				default:
					log.Println(card.Desc.DynamicId)
					return nil, errors.New(fmt.Sprint("unknown type of dynamic: ", card.Desc.Type, card.CardString))
				}

				if latest {
					break
				}
			}
			return ret, nil
		}
	}
	return nil, err
}

type DynamicRetStruct struct {
	Code    int                `json:"code"`
	Msg     string             `json:"msg"`
	Message string             `json:"message"`
	Data    *DynamicDataStruct `json:"data"`
}

type DynamicDataStruct struct {
	Attentions *AttentionsStruct `json:"attentions"`
	Cards      []*CardStruct     `json:"cards"`
	NextOffset int64             `json:"next_offset"`
	HasMore    int               `json:"has_more"`
	Gt         int               `json:"_gt_"`
}

type AttentionsStruct struct {
	UIds []int64 `json:"uids"`
}

type DescStruct struct {
	Uid       int64 `json:"uid"`
	Type      int   `json:"type"`
	Rid       int64 `json:"rid"`
	View      int   `json:"view"`
	RePost    int   `json:"repost"`
	Comment   int   `json:"comment"`
	Like      int   `json:"like"`
	DynamicId int64 `json:"dynamic_id"`
	Timestamp int64 `json:"timestamp"`
	PreDyId   int64 `json:"pre_dy_id"`
	OrigDyId  int64 `json:"orig_dy_id"`
	OrigType  int   `json:"orig_type"`
	UidType   int   `json:"uid_type"`
	Status    int   `json:"status"`
}

type CardStruct struct {
	Desc       *DescStruct `json:"desc"`
	CardString string      `json:"card"`

	CartType    int                  // equal to DescStruct.Type
	CardContent CardContentInterface //
}

type CardContentInterface struct {
	T1    *Type1    // 图文 + 引用<动态, 投稿, 专栏, 房间, 投票>
	T2    *Type2    // 图文 有图片
	T4    *Type4    // 图文 无图片 no origin string in card
	T8    *Type8    // 投稿
	T16   *Type16   // 小视频
	T32   *Type32   // 番剧更新
	T64   *Type64   // 专栏
	T256  *Type256  // 创作
	T512  *Type512  // 番剧
	T1024 *Type1024 // 已失效资源 no origin string in card!
	T2048 *Type2048 // 宣传页
	T4099 *Type4099 // 番剧
	T4200 *Type4200 // 房间
	T4300 *Type4300 // 收藏夹
}

// 图文 + 引用<动态, 投稿, 专栏, 房间, 投票>
type Type1 struct {
	User       *userStruct `json:"user"`
	OriginUser *struct {
		Info *userStruct `json:"info"`
	} `json:"origin_user"`
	Item *struct {
		RpId      int64  `json:"rp_id"`
		Uid       int64  `json:"uid"`
		Content   string `json:"content"`
		OrigDyId  int64  `json:"orig_dy_id"`
		PreDyId   int64  `json:"pre_dy_id"`
		Timestamp int64  `json:"timestamp"`
		Reply     int    `json:"reply"`
		OrigType  int    `json:"orig_type"`
		Tips      string `json:"tips"` // 视频资源已失效
	} `json:"item"` // dynamic content

	OriginString  string                `json:"origin"`
	OriginType    int                   // equal to Item.OrigType
	OriginContent *CardContentInterface // 引用内容
}

// 图文 有图片
type Type2 struct {
	Item *struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Pictures    []*struct {
			ImgSrc    string `json:"img_src"`
			ImgWidth  int    `json:"img_width"`
			ImgHeight int    `json:"img_height"`
			ImgSize   int    `json:"img_size"`
		} `json:"pictures"`
		PicturesCount int   `json:"pictures_count"`
		UploadTime    int64 `json:"upload_time"`
		Reply         int   `json:"reply"`
		IsFav         int   `json:"is_fav"`
	} `json:"item"`
}

// 图文 无图片
type Type4 struct {
	Item *struct {
		RpId      int64  `json:"rp_id"`
		Uid       int64  `json:"uid"`
		Content   string `json:"content"`
		Timestamp int64  `json:"timestamp"`
		Reply     int    `json:"reply"`
	} `json:"item"`
}

// 投稿
type Type8 struct {
	Aid       int64  `json:"aid"`
	Attribute int64  `json:"attribute"`
	Cid       int64  `json:"cid"`
	Copyright int    `json:"copyright"`
	Desc      string `json:"desc"`
	Dimension *struct {
		Height int `json:"height"`
		Rotate int `json:"rotate"`
		Width  int `json:"width"`
	} `json:"dimension"`
	Duration int64  `json:"duration"`
	Dynamic  string `json:"dynamic"`
	JumpUrl  string `json:"jump_url"`
	Owner    *struct {
		Face string `json:"face"`
		Mid  int64  `json:"mid"`
		Name string `json:"name"`
	} `json:"owner"`
	Pic     string `json:"pic"`
	PubDate int64  `json:"pubdate"`
	Tid     int    `json:"tid"`
	TName   string `json:"tname"`
	Videos  int    `json:"videos"`

	C     string `json:"c"`
	CTime int64  `json:"ctime"`
	Title string `json:"title"`
}

// 小视频
type Type16 struct {
	Item *struct {
		Id    int64 `json:"id"`
		Cover *struct {
			Default   string `json:"default"`
			UnClipped string `json:"unclipped"`
		} `json:"cover"`
		Tags             []string `json:"tags"`
		Description      string   `json:"description"`
		VideoTime        int      `json:"video_time"`
		UploadTimeText   string   `json:"upload_time_text"`
		UploadTime       string   `json:"upload_time"`
		VideoPlayUrl     string   `json:"video_playurl"`
		BackupPlayUrl    string   `json:"backup_playurl"`
		VideoSize        string   `json:"video_size"`
		VerifyStatus     string   `json:"verify_status"`
		VerifyStatusText string   `json:"verify_status_text"`
		Reply            int      `json:"reply"`
		Width            int      `json:"width"`
		Height           int      `json:"height"`
		WatchedNum       int      `json:"watched_num"`
	} `json:"item"`
}

// FIXME 猜测： 投票
type Type32 struct {
}

// 专栏
type Type64 struct {
	Id          int64             `json:"id"`
	Category    *categoryStruct   `json:"category"`
	Categories  []*categoryStruct `json:"categories"`
	Title       string            `json:"title"`
	Summary     string            `json:"summary"`
	BannerUrl   string            `json:"banner_url"`
	TemplateId  int               `json:"template_id"`
	State       int               `json:"state"`
	ImageUrls   []string          `json:"image_urls"`
	PublishTime int64             `json:"publish_time"`
	CTime       int64             `json:"ctime"`
	Tags        []*struct {
		Tid  int64  `json:"tid"`
		Name string `json:"name"`
	} `json:"tags"`
	Words           int      `json:"words"`
	Dynamic         string   `json:"dynamic"`
	OriginImageUrls []string `json:"origin_image_urls"`
}

// FIXME: UNKNOWN
type _type128 struct {
}

// 官方创作
type Type256 struct {
	Id          int64  `json:"id"`
	UpId        int64  `json:"up_id"`
	Title       string `json:"title"`
	Upper       string `json:"upper"`
	Cover       string `json:"cover"`
	Author      string `json:"author"`
	CTime       int64  `json:"ctime"`
	ReplyCnt    int    `json:"reply_cnt"`
	PlayCnt     int    `json:"play_cnt"`
	Intro       string `json:"intro"`
	Schema      string `json:"schema"`
	TypeInfo    string `json:"type_info"`
	UpperAvatar string `json:"upper_avatar"`
}

// 番剧
type Type512 struct {
	Type4099
}

// 引用已失效稿件?
type Type1024 struct {
}

// 宣传
type Type2048 struct {
	Rid  int64       `json:"rid"`
	User *userStruct `json:"user"`
	Vest *struct {
		Uid     int64  `json:"uid"`
		Content string `json:"content"`
		Ctrl    string `json:"ctrl"`
	} `json:"vest"`
	Sketch *struct {
		Title     string `json:"title"`
		DescText  string `json:"desc_text"`
		CoverUrl  string `json:"cover_url"`
		TargetUrl string `json:"target_url"`
		SketchId  int64  `json:"sketch_id"`
		BizType   int    `json:"biz_type"`
	} `json:"sketch"`
}

// 番剧
type Type4099 struct {
	Aid           int64 `json:"aid"`
	ApiSeasonInfo *struct {
		BgmType    int    `json:"bgm_type"`
		Cover      string `json:"cover"`
		IsFinish   int    `json:"is_finish"`
		SeasonId   int    `json:"season_id"`
		Title      string `json:"title"`
		TotalCount int    `json:"total_count"`
		Ts         int64  `json:"ts"`
	} `json:"apiSeasonInfo"`
	BulletCount  int    `json:"bullet_count"`
	Cover        string `json:"cover"`
	EpisodeId    int64  `json:"episode_id"`
	Index        string `json:"index"`
	IndexTitle   string `json:"index_title"`
	NewDesc      string `json:"new_desc"`
	OnlineFinish int    `json:"online_finish"`
	PlayCount    int    `json:"play_count"`
	ReplyCount   int    `json:"reply_count"`
	Url          string `json:"url"`
}

// 房间
type Type4200 struct {
	RoomId           int64  `json:"roomid"`
	UId              int64  `json:"uid"`
	UName            string `json:"uname"`
	Verify           string `json:"verify"`
	Virtual          int    `json:"virtual"`
	Cover            string `json:"cover"`
	LiveTime         string `json:"live_time"`
	RoundStatus      int    `json:"round_status"`
	OnFlag           int    `json:"on_flag"`
	Title            string `json:"title"`
	Tags             string `json:"tags"`
	UserCover        string `json:"user_cover"`
	ShortId          int    `json:"short_id"`
	Area             int    `json:"area"`
	AreaV2Name       string `json:"area_v_2_name"`
	AreaV2ParentName string `json:"area_v_2_parent_name"`
	Face             string `json:"face"`
}

// 收藏夹
type Type4300 struct {
	Cover      string `json:"cover"`
	CoverType  int    `json:"cover_type"`
	Fid        int64  `json:"fid"`
	Id         int64  `json:"id"`
	Intro      string `json:"intro"`
	MediaCount int    `json:"media_count"`
	Mid        int64  `json:"mid"`
	Sharable   bool   `json:"sharable"`
	Title      string `json:"title"`
	Type       int    `json:"type"`
	Upper      *struct {
		Face     string `json:"face"`
		Followed int    `json:"followed"`
		Mid      int64  `json:"mid"`
		Name     string `json:"name"`
	} `json:"upper"`
}

type categoryStruct struct {
	Id       int    `json:"id"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
}

type userStruct struct {
	Uid     int64  `json:"uid"`
	UName   string `json:"uname"`
	Face    string `json:"face"`
	Name    string `json:"name"`
	HeadUrl string `json:"head_url"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func DynamicTypeId2String(typeId int) string {
	switch typeId {
	case 1:
		return "图文&引用"
	case 2:
		return "图文&图片"
	case 4:
		return "图文&无图片"
	case 8:
		return "投稿"
	case 16:
		return "小视频"
	case 64:
		return "专栏"
	case 256:
		return "创作"
	case 512:
		return "番剧"
	case 2048:
		return "宣传"
	case 4099:
		return "番剧"
	case 4200:
		return "房间"
	case 4300:
		return "收藏夹"
	default:
		return "未知"
	}
}