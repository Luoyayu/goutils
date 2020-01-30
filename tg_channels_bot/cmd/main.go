package main

import (
	"github.com/luoyayu/goutils/tg_channels_bot"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func push(c *bot.ChannelStruct, f func(*bot.ChannelStruct) error, hangOn int) {
	cnt := 0
	for {
		err := f(c)
		if err != nil {
			cnt += 1
			log.Println(err)
		} else {
			cnt = 0
		}
		if cnt >= hangOn {
			_, _ = c.Notify(c.Desc + " no response for 1h!")
			wg.Done()
			return
		}
		time.Sleep(time.Minute * 5)
	}
}

func main() {
	wg.Add(2)
	go push(bot.Config.Channels["sanqiu"], func(c *bot.ChannelStruct) error {
		return c.QuerySanQiu()
	}, 12)
	go push(bot.Config.Channels["gadio"], func(c *bot.ChannelStruct) error {
		return c.QueryGadio(5)
	}, 12)
	wg.Wait()
}
