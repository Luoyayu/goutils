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
	if runtime.GOOS == "windows" {
		mpv += ".exe"
	}

	promptui.IconInitial = promptui.Styler(promptui.FGCyan)(">")
	promptui.IconSelect = promptui.Styler(promptui.FGRed)("λ")

	AccountSelected, UserSelected = &Account{}, &User{}

	AccountSelected.NikeName = "offline mode"
	AccountSelected.Uid = 2333

	var dbFile *os.File
	var err error

	if dbFile, err = os.Open("data.db"); err != nil { // firstly use the cli-app
		fmt.Println("initialize database...")
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
	//log.Println("MpvPid:", MpvPid)
	if MpvPid != -1 {
		p, _ := process.NewProcess(MpvPid)

		if err := p.Kill(); err != nil {
			log.Println(err)
		} else {

		}
		/*if err := syscall.Kill(MpvPid, syscall.SIGKILL); err == nil {
		} else {
			log.Println(err)
		}*/
		MpvPid = -1
	}
}

func exitClear() {
	if MpvPid != -1 {
		stopMpvSafely()
	}
	os.Exit(0)
}
