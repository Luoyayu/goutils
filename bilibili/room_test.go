package biliAPI

import (
	"log"
	"testing"
)

const _testRoomId = 14085407

func TestRoomInit(t *testing.T) {
	if room, err := RoomInit(_testRoomId); err != nil {
		panic(err)
	} else if room.Code != 0 {
		log.Println(ErrorCodeMap[room.Code])
	} else {
		log.Printf("%+v\n", room)
		if user, err := GetUserInfo(room.Data.Uid); err == nil {
			log.Printf("%+v\n", user.Data)
		}
	}
}

func TestGetRoomPlayUrl(t *testing.T) {
	if roomPlayUrl, err := GetRoomPlayUrl(_testRoomId); err == nil {
		log.Printf("%+v\n", roomPlayUrl)
	} else {
		panic(err)
	}
}

func TestGetLiveUserRecommend(t *testing.T) {
	ret, err := GetLiveUserRecommend(0, "", 1)
	if err != nil {
		panic(err)
	} else {
		if ret.Code != 0 {
		} else {
			for _, f := range ret.Data {
				log.Println(f.Uid, f.UName, f.Title)
			}
		}
	}
}
