package tg_channels_bot

import (
	"github.com/luoyayu/goutils/d4j"
	"github.com/luoyayu/goutils/gcore/gadio"
)

type ConfigToml struct {
	Bots     map[string]*BotStruct     `toml:"bot"`
	Channels map[string]*ChannelStruct `toml:"channel"`
	Debug    bool                      `toml:"debug"`

	// pay attention to name shadow ! recommend start with package name
	*gadio.GadioConfig
	*d4j.SanqiuConfig
}

type BotStruct struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
	Id    int64  `toml:"id"`
	Debug bool   `toml:"debug"`
}

type ChannelStruct struct {
	PubUrl     string `toml:"pub_url"`
	PrivateUrl string `toml:"private_url"`
	Name       string `toml:"name"`
	ChatId     int64  `toml:"chat_id"`
	Desc       string `toml:"desc"`
	Owner      int64  `toml:"owner"`
}
