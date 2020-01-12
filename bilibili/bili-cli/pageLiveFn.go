package main

import (
	"context"
	"errors"
	"fmt"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"github.com/manifoldco/promptui"
	"strconv"
	"strings"
)

var NowLive []*Live

func updateNowLive() {
	NowLive = pageSyncCmdSyncLive(false, false, false)
}

// if needUpdateForAccountSelected is false, you need to update NowLive []*Live
func subPageSelectLive(needUpdateForAccountSelected bool) int64 {
	if needUpdateForAccountSelected {
		updateNowLive()
	}

	liveSelectItems := make([]string, 0, len(NowLive))
	for _, ll := range NowLive {
		if ll.State == 1 && ll.Blocked == 0 {
			liveSelectItems = append(liveSelectItems, fmt.Sprintf("[%d] %s: %s", ll.Cid, ll.NikeName, ll.Title))
		}
	}
	_, result, _ := (&promptui.Select{
		Label: "online live",
		Items: liveSelectItems,
		Size:  10,
	}).Run()

	tmp_ := strings.Split(result, " ")

	var selectedCid int64
	if len(tmp_) > 0 && len(tmp_[0]) > 3 {
		selectedCid, _ = strconv.ParseInt(tmp_[0][1:len(tmp_[0])-1], 10, 64)
	} else {
		Logger.Error("")
	}
	return selectedCid
}

func subPagePlaySelectedLive(cid interface{}, playControlBackToSelf bool) {
	r, _ := biliAPI.RoomInit(cid)
	u, _ := biliAPI.GetUserInfo(r.Data.Uid)

	if r.Code != 0 {
		Logger.Error(r.Message)
		showPageLive()
	} else if r.Data.LiveStatus == 0 {
		if u.Code == 0 {
			Logger.Error("the live room for ", u.Data.Name, " now is closed!")
		} else {
			Logger.Error("the live room for ", r.Data.Uid, " now is closed!")
		}
		showPageLive()
	}

	uName := ""
	if u.Code == 0 {
		uName = u.Data.Name
	} else {
		uName = fmt.Sprint(r.Data.Uid)
	}

	playSelectedLiveItems := []string{
		"play: video,sound,danmaku",
		"play: sound",
		"play: video,sound",
		"play: sound,danmaku",
		"play: danmaku",
		"play: costumed mpv args",
		CMDBack, // back to subPageSelectLive
		CMDHome,
		CMDExit,
	}
	pageSelectPrompt := promptui.Select{
		Label: "play live for " + uName,
		Items: playSelectedLiveItems,
		Size:  len(playSelectedLiveItems),
	}

	_, playCmd, _ := pageSelectPrompt.Run()
	if strings.HasPrefix(playCmd, "play:") {

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		paramsMap := map[string]int{
			"danmaku": -1,
			"video":   -1,
			"sound":   -1,
		} // 1:on, -1:off, 0:no-set
		mpvOptions := ""

		if playCmd == "play: costumed mpv args" {
			mpvOptions, _ = (&promptui.Prompt{
				Label: "mpv options, set --danmaku=<yes|no> default: yes",
			}).Run()

			if strings.Contains(mpvOptions, "--danmaku=yes") {
				paramsMap["danmaku"] = 1
				mpvOptions = strings.Replace(mpvOptions, "--danmaku=yes", "", -1)
			}
		} else {
			params := strings.Split(strings.TrimPrefix(playCmd, "play: "), ",")
			for i := range params {
				paramsMap[params[i]] = 1
			}
		}

		if paramsMap["danmaku"] == -1 && paramsMap["video"] == -1 && paramsMap["sound"] == -1 && mpvOptions == "" {
			subPagePlaySelectedLive(cid, true)
			return
		}
		// end stupid check

		go playLive(ctx, cid, paramsMap, mpvOptions)

		_, controlCallback, _ := (&promptui.Select{
			Label: "control live " + playCmd,
			Items: []string{
				CMDBack,
				CMDHome,
				CMDExit,
			},
		}).Run()

		switch controlCallback {
		case CMDBack:
			cancel()
			if playControlBackToSelf {
				subPagePlaySelectedLive(cid, true)
				return
			} else {
				showPageLive()
			}
		case CMDHome:
			cancel()
			showPageHome()
		case CMDExit:
			cancel()
			exitClear()
		}
	} else {
		switch playCmd {
		case CMDBack:
			pageLiveCmdSelect()
		case CMDHome:
			showPageHome()
		case CMDExit:
			exitClear()
		}
	}
}

func pageLiveCmdSelect() {

	subPagePlaySelectedLive(subPageSelectLive(true), true)
}

func pageLiveCmdBlock() {
	Logger.Warn("unimplemented method!")
	showPageLive()
}

func pageLiveCmdAdd() {
	input, _ := (&promptui.Prompt{
		Label: "rid:XXX / uid:XXX",
		Validate: func(s string) error {
			if strings.HasPrefix(s, "rid") == false || strings.HasPrefix(s, "uid") {
				return errors.New("input must start with `rid:` or `uid:`")
			} else {
				sp := strings.Split(s, ":")
				if len(sp) != 2 {
					return errors.New("wrong format, input must like `rid:XXX` or `uid:YYY`")
				}
				if _, err := strconv.ParseInt(sp[1], 10, 64); err != nil {
					return errors.New("input must end with a number")
				}
			}
			return nil
		},
	}).Run()

	sp := strings.Split(input, ":")
	if strings.TrimSpace(sp[0]) == "rid" {
		subPagePlaySelectedLive(strings.TrimSpace(sp[1]), false)
	} else {
		r, _ := biliAPI.GetRoomNews("", strings.TrimSpace(sp[1]))
		if r.Code != 0 {

		} else {
			subPagePlaySelectedLive(r.Data.(map[string]interface{})["roomid"], false)
		}
	}
}

func pageLiveCmdRecommend() {
	if AccountSelected.SESSDATA == "" || AccountSelected.Uid == 0 {
		Logger.Error("you need to add and select account firstly")
		showPageLive()
	}
	ret, err := biliAPI.GetLiveUserRecommend(AccountSelected.Uid, AccountSelected.SESSDATA, 1)
	if err != nil {
		Logger.Error(err)
		showPageLive()
	}
	if ret.Code != 0 {
		Logger.Error("Not found your")
		showPageLive()
	}

	NowLive = make([]*Live, 0, len(ret.Data))
	for _, ll := range ret.Data {
		//log.Printf("%+v\n", ll)
		NowLive = append(NowLive, &Live{
			NikeName: ll.UName,
			Cid:      ll.RoomId,
			Title:    ll.Title,
			State:    1,
		})
	}

	subPagePlaySelectedLive(subPageSelectLive(false), false)
}

func pageLiveCmdEdit() {
	showPageLive()
}

func pageLiveCmdDelete() {
	Logger.Warn("unimplemented method!")
	showPageLive()
}
