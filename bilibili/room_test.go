package biliAPI

import (
	"log"
	"testing"
)

const _testRoomId = 0

func TestRoomInit(t *testing.T) {
	if ret, err := RoomInit(_testRoomId); err != nil {
		panic(err)
	} else if ret.Code != 0 {
		log.Println(ErrorCodeMap[ret.Code])
	} else {
		log.Printf("%+v\n", ret)
		if user, err := GetUserInfo(ret.Data.Uid); err == nil {
			log.Printf("%+v\n", user.Data)
		}
	}
}

func TestGetRoomPlayUrl(t *testing.T) {
	if roomPlayUrl, err := GetRoomPlayUrl(_testRoomId); err == nil {
		log.Printf("%+v\n", roomPlayUrl)
		log.Printf("%+v\n", roomPlayUrl.Data.Durl[0])
	} else {
		panic(err)
	}
}

func TestGetLiveUserRecommend(t *testing.T) {
	ret, err := GetLiveUserRecommend(0, "", 1)
	if err != nil {
		panic(err)
	}
	if ret.Code != 0 {
	} else {
		for _, f := range ret.Data {
			log.Println(f.Uid, f.UName, f.Title)
		}
	}

}

func TestGetRoomInfoByRoomId(t *testing.T) {
	ret, err := GetRoomInfoByRoomId(_testRoomId)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", ret.Data)
}
