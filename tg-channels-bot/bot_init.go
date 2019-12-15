package tg_channels_bot

import (
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/luoyayu/goutils/d4j"
	"github.com/luoyayu/goutils/gcore/gadio"
	myRedis "github.com/luoyayu/goutils/redis"
	"log"
	"strconv"
)

var (
	Config ConfigToml
	DEBUG  bool
	Bots   []*tgbotapi.BotAPI

	GadioRedis  *gadio.Redis
	SanqiuRedis *d4j.Redis
)

func init() {
	// init main config from bot_config.toml
	tomlPath := "bot_config.toml"
	if _, err := toml.DecodeFile(tomlPath, &Config); err != nil {
		panic(err)
	} else {
		DEBUG = Config.Debug
		if DEBUG {
			log.Printf("bot_config.toml: %+v\n", Config)
			log.Printf("sanqiu_config: %+v\n", Config.SanqiuConfig.Sanqiu)
			log.Printf("SanqiuRedis_config: %+v\n", Config.SanqiuConfig.SanqiuRedis)

			log.Printf("Gcore_config: %+v\n", Config.GadioConfig.Gcore)
			log.Printf("Gadio_config: %+v\n", Config.GadioConfig.Gadio)
			log.Printf("GadioRedis_config: %+v\n", Config.GadioConfig.GadioRedis)

		}

		// init services' config
		gadio.Config = Config.GadioConfig
		d4j.Config = Config.SanqiuConfig

		// init services' redis config
		if gadio.Config.GadioRedis.Enable {
			GadioRedis = &gadio.Redis{
				Client: myRedis.InitRedis(&redis.Options{
					Addr: gadio.Config.GadioRedis.Addr,
					DB:   gadio.Config.GadioRedis.Db,
				}),
			}
		}

		if d4j.Config.SanqiuRedis.Enable {
			SanqiuRedis = &d4j.Redis{
				Client: myRedis.InitRedis(&redis.Options{
					Addr: d4j.Config.SanqiuRedis.Addr,
					DB:   d4j.Config.SanqiuRedis.Db,
				}),
			}
		}

	}

	// init tg-bots config
	if len(Config.Bots) == 0 {
		panic("no bot found in bot_config.toml")
	}
	Bots = make([]*tgbotapi.BotAPI, len(Config.Bots))
	for botNumStr, bot := range Config.Bots {
		var err error
		var botNum int64
		if botNum, err = strconv.ParseInt(botNumStr, 10, 64); err != nil {
			panic(err)
		}

		if Bots[botNum], err = tgbotapi.NewBotAPI(bot.Token); err != nil {
			panic(err)
		}
		Bots[botNum].Debug = Config.Bots[botNumStr].Debug
		if DEBUG {
			log.Printf("bots[%d]:%+v\n", botNum, Bots[botNum])
		}
	}
}
