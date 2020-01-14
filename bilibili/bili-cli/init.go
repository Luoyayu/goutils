package main

import (
	"database/sql"
	"fmt"
	"github.com/luoyayu/goutils/logger"
	"github.com/manifoldco/promptui"
	"github.com/shirou/gopsutil/process"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	AccountsInDB  []*Account
	FollowingInDB []*Following
	LiveInDB      []*Live

	UserSelected    *User
	AccountSelected *Account

	IsSingleMode int

	Logger = logger.NewDebugLogger()

	db *sql.DB
)

func init() {
	initConfig()
	fmt.Println(promptui.Styler(promptui.FGRed)("DEMO VERSION! any issue post to 'https://github.com/Luoyayu/goutils/tree/master/bilibili/bili-cli'"))
	if runtime.GOOS == "windows" {
		mpv += ".exe"
	}

	AccountSelected, UserSelected = &Account{}, &User{}

	AccountSelected.NikeName = "本地模式"
	AccountSelected.Uid = 2333

	var dbFile *os.File
	var err error

	if dbFile, err = os.Open("data.db"); err != nil { // firstly use the cli-app
		Logger.Info("initialize database...")
		Logger.Info("欢迎使用，本软件支持多账户，推荐登录使用：Account -> add -> select")
		//Logger.Info("本地使用，请使用 Sync 同步关注")
		if dbFile, err = os.Create(databaseName); err == nil {
			defer dbFile.Close()

			db, _ = sql.Open("sqlite3", databaseName)

			if err = initDB(); err != nil {
				log.Fatalln("initialize database error:", err)
			} else {
				if dbFile.Sync() != nil {
					log.Println("sync database error:", err)
				}
			}
		}
	} else {
		db, _ = sql.Open("sqlite3", databaseName)
		if err = db.QueryRow(`SELECT singleMode FROM Preference`).Scan(&IsSingleMode); err != nil {
			_, err = db.Exec(`INSERT INTO Preference VALUES (0);`) // default multi-user mod
		}
		checkIfAccountExpired()
	}
}

func checkIfAccountExpired() {
	if findAccountLastUsed() == nil {
		if AccountSelected.Expire <= time.Now().Unix() {
			Logger.Warnf("the account:%v[%v]'s token is expired, \n"+
				"please add it again, or edit manually\n",
				AccountSelected.NikeName, AccountSelected.Uid)
		}
	}
}

func printLoading(done chan bool, d time.Duration) {
	if d == 0 {
		d = time.Millisecond * 200
	}
	fmt.Print("loading")
	for {
		tm := time.NewTimer(d)
		select {
		case <-tm.C:
			fmt.Print(".")
		case <-done:
			return
		}
	}
}

func stopMpvSafely() {
	if MpvPid != -1 {
		p, _ := process.NewProcess(MpvPid)

		if err := p.Kill(); err != nil {
			Logger.Error(err)
		} else {

		}
		MpvPid = -1
	}
}

func exitClear() {
	if MpvPid != -1 {
		stopMpvSafely()
	}
	os.Exit(0)
}
