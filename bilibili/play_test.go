package biliAPI

import (
	"fmt"
	"testing"
)

func TestGetPlayFromCDN(t *testing.T) {
	ret, err := GetPlayUrl(36607048, 64249025, 80)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", ret.Data)
}

func TestDownloadVideo(t *testing.T) {
	DownloadVideo(83199184, 1, 0)
}
