package main

import (
	"github.com/BurntSushi/toml"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"log"
)

const (
	databaseName = "data.db"
)

var (
	mpv = "mpv"
)

var MpvPid int32 = -1

func initConfig() {

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
    relation_followings = "https://api.bilibili.com/x/relation/followings" # can be replced by dynamic url

    search_mobi = "https://app.bilibili.com/x/v2/search"
    search_mobi_suggest = "https://grpc.biliapi.net/bilibili.app.interface.v1.Search/Suggest3"

    search_web = "https://api.bilibili.com/x/web-interface/search/type"
    search_web_suggest = "https://s.search.bilibili.com/main/suggest"

    dynamic_space_history = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/space_history"

    room_init = "https://api.live.bilibili.com/room/v1/Room/room_init"
    room_play_url = "https://api.live.bilibili.com/room/v1/Room/playUrl"
    room_news_get = "https://api.live.bilibili.com/room_ex/v1/RoomNews/get"

    live_get_user_recommend = "https://api.live.bilibili.com/room/v1/room/get_user_recommend"
    live_my_following = "https://api.live.bilibili.com/xlive/web-ucenter/user/following" # need SESSDATA

#    danmaku_host = "wss://broadcastlv.chat.bilibili.com:2245/sub"
    danmaku_host = "wss://ks-live-dmcmt-sh2-pm-03.chat.bilibili.com/sub"
    danmaku_host_2 = "wss://tx-gz3-live-comet-03.chat.bilibili.com/sub"
    danmaku_host_3 = "wss://tx-sh3-live-comet-05.chat.bilibili.com/sub"

`
	//tomlPath := "config.toml"

	if _, err := toml.Decode(tomlContent, &biliAPI.Config); err != nil {
		panic(err)
	} else if biliAPI.Config.Debug {
		log.Printf("read from bilibili api config %+v\n", biliAPI.Config)
		log.Printf("%#v\n", biliAPI.Config.API)
	}
}
