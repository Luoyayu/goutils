package biliAPI

import (
	"log"
	"testing"
)

func TestGetBangumiInitialSate(t *testing.T) {
	if ret, err := GetBangumiInitialSate(1); err == nil {
		log.Printf("%+v\n", ret)
	} else {
		panic(err)
	}
}
