module github.com/luoyayu/goutils/bilibili

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/gorilla/websocket v1.4.1
	github.com/luoyayu/goutils/enc v0.0.0-20200112053613-a6db8386b494
	github.com/manifoldco/promptui v0.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
)

replace (
	github.com/luoyayu/goutils/date => /Volumes/MacMisc/github/goutils/date
	github.com/luoyayu/goutils/enc => /Volumes/MacMisc/github/goutils/enc
	github.com/luoyayu/goutils/logger => /Volumes/MacMisc/github/goutils/logger
	github.com/luoyayu/goutils/net => /Volumes/MacMisc/github/goutils/net
)
