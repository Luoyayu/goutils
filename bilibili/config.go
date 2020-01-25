package biliAPI

import (
	"github.com/BurntSushi/toml"
	"log"
)

var Config = &ConfigToml{}

type ConfigToml struct {
	Debug     bool           `toml:"debug"`
	AppKey    string         `toml:"app_key"`
	AppSecret string         `toml:"app_secret"`
	MagicKey  string         `toml:"magic_key"`
	API       *ConfigBiliAPI `toml:"api"`
}

type ConfigBiliAPI struct {
	Login    string `toml:"login"`
	UserInfo string `toml:"user_info"`

	RelationStat       string `toml:"relation_stat"`
	RelationFollowings string `toml:"relation_followings"`

	SearchMobi        string `toml:"search_mobi"`
	SearchMobiSuggest string `toml:"search_mobi_suggest"`
	SearchWeb         string `toml:"search_web"`
	SearchWebSuggest  string `toml:"search_web_suggest"`

	SpaceArcSearch      string `toml:"space_arc_search"`
	DynamicSpaceHistory string `toml:"dynamic_space_history"`

	RoomInit            string `toml:"room_init"`
	RoomPlayUrl         string `toml:"room_play_url"`
	RoomNewsGet         string `toml:"room_news_get"`
	RoomGetInfoByRoomID string `toml:"room_get_info_by_roomID"`

	LiveGetUserRecommend string `toml:"live_get_user_recommend"`
	LiveMyFollowing      string `toml:"live_my_following"`
	DanmakuHost          string `toml:"danmaku_host"`

	GetCIdsByAId string `toml:"get_cids_by_aid"`
	GetPlayUrl   string `toml:"get_play_url"`
}

// implement by caller ! Config

func init() {
	//tomlPath := "config.toml"
	tomlContent := `
# android
app_key    = "1d8b6e7d45233436"
app_secret = "560c52ccd288fed045859ed18bffd973"

# host aes key
magic_key = "7a840o62v39c41b8"

debug = false

[api]
    login = "https://passport.bilibili.com"
    user_info = "https://api.bilibili.com/x/space/acc/info"

    relation_stat = "https://api.bilibili.com/x/relation/stat"
    relation_followings = "https://api.bilibili.com/x/relation/followings" # can be replaced by dynamic url

    search_mobi = "https://app.bilibili.com/x/v2/search"
    search_mobi_suggest = "https://grpc.biliapi.net/bilibili.app.interface.v1.Search/Suggest3"

    search_web = "https://api.bilibili.com/x/web-interface/search/type"
    search_web_suggest = "https://s.search.bilibili.com/main/suggest"
	
	space_arc_search = "https://api.bilibili.com/x/space/arc/search" # 投稿查询

    dynamic_space_history = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history"

    room_init = "https://api.live.bilibili.com/room/v1/Room/room_init"
    room_play_url = "https://api.live.bilibili.com/room/v1/Room/playUrl"
    room_news_get = "https://api.live.bilibili.com/room_ex/v1/RoomNews/get"
	room_get_info_by_roomID = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom"

    live_get_user_recommend = "https://api.live.bilibili.com/room/v1/room/get_user_recommend"
    live_my_following = "https://api.live.bilibili.com/xlive/web-ucenter/user/following" # need SESSDATA
	
	get_cids_by_aid = "https://api.bilibili.com/x/web-interface/view"
	get_play_url = "https://api.bilibili.com/x/player/playurl"

#    danmaku_host = "wss://broadcastlv.chat.bilibili.com:2245/sub"
    danmaku_host = "wss://ks-live-dmcmt-sh2-pm-03.chat.bilibili.com/sub"
    danmaku_host_2 = "wss://tx-gz3-live-comet-03.chat.bilibili.com/sub"
    danmaku_host_3 = "wss://tx-sh3-live-comet-05.chat.bilibili.com/sub"

`

	if _, err := toml.Decode(tomlContent, &Config); err != nil {
	} else if Config.Debug {
		log.Printf("read from bilibili api config %+v\n", Config)
		log.Printf("%#v\n", Config.API)
	}
}
