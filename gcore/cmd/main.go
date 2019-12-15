package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	bot := &tgbotapi.BotAPI{}
	bot.Debug = true
	var err error
	bot, err = tgbotapi.NewBotAPI("1058178782:AAHIN3JG50FhRKUPAyJoVHO42ZP1CFwWX9Q")
	if err != nil {
		panic(err)
	} else {
		log.Println(bot.Self.UserName)
	}
}
