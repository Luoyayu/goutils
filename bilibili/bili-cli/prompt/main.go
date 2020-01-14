package main

import (
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"github.com/luoyayu/goutils/logger"
)

func main() {

	r, _ := biliAPI.RoomInit("")
	logger.NewDebugLogger().Info(r.Data.LiveStatus)
	logger.NewDebugLogger().Info(r.Data)
}
