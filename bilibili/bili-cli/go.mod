module bili-cli

go 1.13

require (
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/BurntSushi/toml v0.3.1
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/luoyayu/goutils/bilibili v0.0.0-20200112160343-05ea778e84c8
	github.com/luoyayu/goutils/date v0.0.0-00010101000000-000000000000
	github.com/luoyayu/goutils/logger v0.0.0-20200112075104-e395eab4880e
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-sqlite3 v2.0.2+incompatible
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/pkg/errors v0.8.1
	github.com/shirou/gopsutil v2.19.12+incompatible
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76 // indirect
	gopkg.in/toast.v1 v1.0.0-20180812000517-0a84660828b2
)

replace (
	github.com/luoyayu/goutils/bilibili => /Volumes/MacMisc/github/goutils/bilibili
	github.com/luoyayu/goutils/date => /Volumes/MacMisc/github/goutils/date
	github.com/luoyayu/goutils/enc => /Volumes/MacMisc/github/goutils/enc
	github.com/luoyayu/goutils/logger => /Volumes/MacMisc/github/goutils/logger
	github.com/luoyayu/goutils/net => /Volumes/MacMisc/github/goutils/net
)
