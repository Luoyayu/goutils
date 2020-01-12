package biliAPI

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const (
	danmakuOpsSendHeartBeat = 2
	danmakuOpsPopularity    = 3
	danmakuOpsCommand       = 5
	danmakuOpsAuth          = 7
	danmakuOpsRecvHeartBeat = 8
)

type DanmakuClient struct {
	Cid   int64
	Conn  *websocket.Conn
	Debug bool
}

func (r *DanmakuClient) logDanmaku(v ...interface{}) {
	if r.Debug {
		log.Println(v...)
	}
}
func (r *DanmakuClient) logDanmakuf(f string, v ...interface{}) {
	if r.Debug {
		log.Printf(f, v...)
	}
}

func (r *DanmakuClient) NewDanmakuClient(ctx context.Context, roomId interface{}) *DanmakuClient {
	if room, err := RoomInit(roomId); err == nil {
		r.Cid = room.Data.RoomId
		log.Println("连接到弹幕服务器: ", Config.API.DanmakuHost)
		ack := make(chan bool)
		r.Conn, _, err = (&websocket.Dialer{}).Dial(Config.API.DanmakuHost, nil)

		go r.SendAuth()
		go r.RecvMsg(ctx, ack)

		if <-ack == true {
			fmt.Printf("进入房间 https://live.bilibili.com/%v 成功!\n", r.Cid)
			go r.SendHeadBeat(ctx)
		}

		select {
		case <-ctx.Done():
			//log.Println("关闭弹幕子协程!")
			return r
		}
	} else {
		log.Fatalln(err)
	}
	return r
}

func (r *DanmakuClient) SendAuth() bool {
	authParams := map[string]interface{}{
		"uid":      int(rand.Float64()*200000000000000.0 + 100000000000000.0),
		"roomid":   r.Cid,
		"protover": 1,
		"platform": "web",
		//"clientver": "1.4.0",
	}
	jsonBody, _ := json.Marshal(authParams)
	body := string(jsonBody)
	handshake := fmt.Sprintf("%08x00100001%08x00000001", len(body)+16, danmakuOpsAuth)

	buf := make([]byte, len(handshake)>>1)
	_, _ = hex.Decode(buf, []byte(handshake))
	r.logDanmaku("发送确认信息")
	_ = r.Conn.WriteMessage(websocket.BinaryMessage, append(buf, []byte(body)...))
	return true
}

func (r *DanmakuClient) RecvMsg(ctx context.Context, ack chan bool) {
	for {
		select {
		case <-ctx.Done():
			//log.Println("停止接受弹幕!")
			return
		default:
			_, msg, err := r.Conn.ReadMessage()
			if err != nil {
				log.Println("read from wss err: ", err)
				return
			}
			recvMsgOp := msg[11]
			recvMsgBody := msg[16:]

			switch recvMsgOp {
			case danmakuOpsCommand:
				r.logDanmaku("收到弹幕")
				r.logDanmaku("msg: ", string(msg))

				msgs := bytes.SplitAfter(msg, []byte{0, 0, 0, danmakuOpsCommand, 0, 0, 0, 0})
				msgsCnt := len(msgs)
				for i, m := range msgs { // 消息嵌套
					r.logDanmaku("第", i, "条消息")
					r.logDanmaku(string(m))
					if i != msgsCnt-1 {
						r.handleRecMsgType(m[:len(m)-16])
					} else {
						r.handleRecMsgType(m)
					}
				}
			case danmakuOpsRecvHeartBeat:
				r.logDanmaku("收到心跳包确认")
				ack <- true

			case danmakuOpsPopularity:
				r.logDanmaku("收到人气值")
				p := binary.BigEndian.Uint32(recvMsgBody)
				fmt.Println("🔥人气值: ", p)
			}
		}
	}
}

func (r *DanmakuClient) SendHeadBeat(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			//log.Println("停止发送心跳包!")
			return
		default:
			buf := make([]byte, 16)
			_, _ = hex.Decode(buf, []byte("0000001f001000010000000200000001"))
			_ = r.Conn.WriteMessage(websocket.BinaryMessage, buf)
			time.Sleep(30 * time.Second)
		}
	}
}

func (r *DanmakuClient) handleRecMsgType(msg []byte) {
	var jsonMap map[string]interface{}
	_ = json.Unmarshal(msg, &jsonMap)
	cmd := jsonMap["cmd"]

	r.logDanmaku("命令:", cmd)
	r.logDanmakuf("信息: %#v\n", jsonMap)
	switch cmd {
	case "LIVE":
		fmt.Println("✔️直播开始...")
	case "PREPARING":
		fmt.Println("❌房主准备中...")
	case "DANMU_MSG":
		info := jsonMap["info"].([]interface{})

		//r.logDanmakuf("info[0]: %+v\n", info[0].([]interface{}))
		//r.logDanmakuf("info[1]: %+v\n", info[1])
		//r.logDanmakuf("info[2]: %+v\n", info[2].([]interface{}))

		message := info[1]
		postInfo := info[2].([]interface{})
		posterName := postInfo[1]
		isGuard := int64(postInfo[2].(float64))

		if isGuard == 1 {
			fmt.Printf("[🔧\033[32m%s\033[0m]: \033[34m%s\033[0m\n", posterName, message)
		} else {
			fmt.Printf("[\033[32m%s\033[0m]: \033[34m%s\033[0m\n", posterName, message)
		}
	case "SEND_GIFT":
	case "WELCOME":
	case "WELCOME_GUARD":
	case "ENTRY_EFFECT":
	case "ACTIVITY_BANNER_UPDATE_V2":
	case "NOTICE_MSG":
		msgType := int64(jsonMap["msg_type"].(float64))
		switch msgType {
		case 1: // 恭喜 ...

		case 2: // 礼物 ...
			businessId := jsonMap["business_id"].(string)
			linkUrl := jsonMap["link_url"].(string)
			msgCommon := jsonMap["msg_common"].(string)
			_ = int64(jsonMap["real_roomid"].(float64))

			if businessIdMap[businessId] == "" {
				log.Println("未知礼物:", businessId)
			}
			toLink, err := url.Parse(linkUrl)
			if err != nil {
				log.Println(toLink)
				fmt.Printf("🎁\033[31m%s\033[0m: %s %s\n", businessIdMap[businessId], linkUrl, strings.TrimRight(msgCommon, "，点击前往TA的房间去抽奖吧"))
			} else {
				fmt.Printf("🎁\033[31m%s\033[0m: %s %s\n", businessIdMap[businessId], toLink.Host+toLink.Path, strings.TrimRight(msgCommon, "，点击前往TA的房间去抽奖吧"))
			}

		case 3: // 开通 ...

		}

	case "ROOM_RANK":
		data := jsonMap["data"].(map[string]interface{})
		rankDesc := data["rank_desc"]
		fmt.Println("🎤", rankDesc)

	case "ROOM_REAL_TIME_MESSAGE_UPDATE":
	case "USER_TOAST_MSG":
	case "GUARD_BUY":
	case "WEEK_STAR_CLOCK":
	}
}

var businessIdMap = map[string]string{
	"25":    "小电视飞船",
	"20003": "摩天大楼",
}
