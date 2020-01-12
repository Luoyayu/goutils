package biliAPI

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

	DynamicSpaceHistory string `toml:"dynamic_space_history"`

	RoomInit    string `toml:"room_init"`
	RoomPlayUrl string `toml:"room_play_url"`
	RoomNewsGet string `toml:"room_news_get"`

	LiveGetUserRecommend string `toml:"live_get_user_recommend"`
	LiveMyFollowing      string `toml:"live_my_following"`
	DanmakuHost          string `toml:"danmaku_host"`
}

// implement by caller ! Config

/*func initConfig() {
	tomlPath := "config.toml"

	if _, err := toml.DecodeFile(tomlPath, &Config); err != nil {
		panic(err)
	} else if Config.Debug {
		log.Printf("read from bilibili api config %+v\n", Config)
		log.Printf("%#v\n", Config.API)
	}
}*/
