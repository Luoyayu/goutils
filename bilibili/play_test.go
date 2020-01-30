package biliAPI

import (
	"fmt"
	"log"
	"testing"
)

func TestGetPlayFromCDN(t *testing.T) {
	ret, err := GetPlayUrl(74417313, 127291576, 112, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", ret.Data)
}

func TestDownloadVideo(t *testing.T) {
	DownloadVideo(74417313, 1, 112, "6354df00%2C1582698818%2C4670c911", -1)
}

func TestGetCidByAid(t *testing.T) {
	ret, _ := GetCidByAid(74417313)
	log.Printf("%+v\n", ret.Data.Owner)
}