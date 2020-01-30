package tg_channels_bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func sendMsgToChannelByName(botNum int64, channelName string, msgText string) (*tgbotapi.Message, error) {
	var err error
	var resp tgbotapi.Message
	msg := tgbotapi.NewMessageToChannel(channelName, msgText)
	msg.ParseMode = "markdown"
	resp, err = Bots[botNum].Send(msg)
	return &resp, errors.Wrap(err, "sendMsgToChannelByName ->")
}

func sendMsgToChannelById(botNum int64, channelId int64, msgText string) (*tgbotapi.Message, error) {
	var err error
	var resp tgbotapi.Message
	msg := tgbotapi.NewMessage(channelId, msgText)
	msg.ParseMode = "markdown"
	resp, err = Bots[botNum].Send(msg)
	return &resp, errors.Wrap(err, "sendMsgToChannelById ->")
}
