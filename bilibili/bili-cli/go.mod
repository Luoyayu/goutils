module bili-cli

go 1.13

require (
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/BurntSushi/toml v0.3.1
	github.com/luoyayu/goutils/bilibili v0.0.0-20200112053613-a6db8386b494
	github.com/luoyayu/goutils/logger v0.0.0-20200112053613-a6db8386b494
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.4
	github.com/mattn/go-sqlite3 v2.0.2+incompatible
	github.com/pkg/errors v0.8.1
	github.com/shirou/gopsutil v2.19.12+incompatible
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76 // indirect
	gopkg.in/toast.v1 v1.0.0-20180812000517-0a84660828b2
)

replace (
	github.com/luoyayu/goutils => /Volumes/MacMisc/github/goutils/
	github.com/luoyayu/goutils/bilibili => /Volumes/MacMisc/github/goutils/bilibili
	github.com/luoyayu/goutils/enc => /Volumes/MacMisc/github/goutils/enc
	github.com/luoyayu/goutils/logger => /Volumes/MacMisc/github/goutils/logger
)
