package tg_channels_bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/luoyayu/goutils/common"
	"github.com/luoyayu/goutils/d4j"
	"github.com/luoyayu/goutils/date"
	"github.com/luoyayu/goutils/gcore/gadio"
	"github.com/luoyayu/goutils/logger"
	"github.com/pkg/errors"
	"log"
	"strings"
)

var Logger = logger.NewDebugLogger()

func (r *ChannelStruct) Notify(msgText string) (resp *tgbotapi.Message, err error) {
	defer func() {
		err = common.RecoverRuntimeError(recover(), errors.Wrap(err, "Notify ->"))
	}()

	if r.Name != "" {
		resp, err = sendMsgToChannelByName(r.Owner, r.Name, msgText)
	} else if r.ChatId < 0 {
		resp, err = sendMsgToChannelById(r.Owner, r.ChatId, msgText)
	} else {
		panic("can't Notify: channel's name and chat_id all illegal!")
	}
	return
}

// DelMsg :Delete Message by ChatId and MessageId
func (r *ChannelStruct) DelMsg(msgId int, chatId int64) (*tgbotapi.Message, error) {
	var err error
	var resp tgbotapi.Message
	defer func() {
		err = common.RecoverRuntimeError(recover(), errors.Wrap(err, "DelMesFromChannel ->"))
	}()
	msg := tgbotapi.NewDeleteMessage(chatId, msgId)
	resp, err = Bots[r.Owner].Send(msg)
	return &resp, err
}

func (r *ChannelStruct) QueryGadio(number int) (err error) {
	var radioQuery *gadio.RadiosQuery
	var exists bool
	if radioQuery, err = gadio.GetLatestN(number); err == nil {
		for i := len(radioQuery.Data) - 1; i >= 0; i-- {
			radio := radioQuery.Data[i]
			key := strings.Join([]string{"gcore", "gadio", radio.ID}, ":")
			if Config.Gcore.Debug {
				log.Println(key, radio.Attrs.Title)
			}
			if gadio.Config.GadioRedis.Enable == false {
				panic("please enable redis to check whether the radio already pushed to the channel!")
			}
			if exists, err = GadioRedis.Exists(key); err == nil && exists == false {
				msgText := date.ParseDate("", radio.Attrs.PublishedAt) + "\nhttps://www.gcores.com/radios/" + radio.ID
				var resp *tgbotapi.Message
				if resp, err = r.Notify(msgText); err == nil && resp != nil && resp.Chat != nil {
					if err = GadioRedis.Stores(key, &radio, r.Owner, r.Name, resp.MessageID, resp.Chat.ID); err == nil {
						if Config.Gcore.Debug {
							log.Printf("store %s successfully!", key)
						}
					}
				}
			}
		}
	}
	return errors.Wrap(err, "QueryGadio ->")
}

func (r *ChannelStruct) QuerySanQiu() (err error) {
	var exists bool
	if books, _, err := d4j.GetBookInfoFromRss(false); err == nil {
		for i := len(books) - 1; i >= 0; i-- {
			if Config.Sanqiu.Debug {
				log.Println(books[i].Name, books[i].Author, books[i].ShareLink.Url+"/#"+books[i].ShareLink.Key)
			}
			if d4j.Config.SanqiuRedis.Enable == false {
				panic("please enable redis to check whether the books[i] already pushed to the channel!")
			}
			key := "d4j:" + books[i].ID
			if exists, err = SanqiuRedis.Exists(key); exists == false && err == nil {
				msgText := fmt.Sprintf("[%s](%s)\n作者: %s\n标签: %v\n[云盘分享](%s)",
					books[i].Name, Config.Sanqiu.BaseUrl+"/"+books[i].ID+".html",
					books[i].Author,
					books[i].Tags,
					books[i].ShareLink.Url+"/#"+books[i].ShareLink.Key)
				var resp *tgbotapi.Message
				if resp, err = r.Notify(msgText); err == nil && resp != nil && resp.Chat != nil {
					if err = SanqiuRedis.Stores(key, books[i]); err == nil {
						if Config.Sanqiu.Debug {
							log.Printf("store %s successfully!", books[i].ID)
						}
					}
				}
			}
		}
	}
	return
}
