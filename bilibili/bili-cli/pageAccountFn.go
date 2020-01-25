package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	biliAPI "github.com/luoyayu/goutils/bilibili"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

func pageAccountSelect() {
	findAllAccounts()
	if len(AccountsInDB) == 0 {
		Logger.Info("no any account in DB")
		showPageAccount()
	}
	pageAccountSelectItems := make([]string, len(AccountsInDB))

	for i, account := range AccountsInDB {
		pageAccountSelectItems[i] = fmt.Sprint(i, ". ", account.NikeName, " ", account.Uid)
	}

	selectedUserStr := PromptSelect(
		"select account to activate",
		pageAccountSelectItems,

	)
	if selectedUserIndex, err := strconv.ParseInt(strings.Split(selectedUserStr, ". ")[0], 10, 32); err == nil {
		AccountSelected = AccountsInDB[selectedUserIndex]
		//log.Println(AccountSelected.Uid)
		_, err := db.Exec(`UPDATE Account SET lastUsedTimestamp=? WHERE uid IS ?;`, time.Now().Unix(), AccountSelected.Uid)
		if err != nil {
			panic(err)
		}
		showPageAccount()
	}
}

func pageAccountCMDAdd() {
	validate := func(input interface{}) error {

		if len(fmt.Sprint(input)) == 0 {
			return errors.New("Username must not null")
		}
		return nil
	}

	loginUserName := PromptInput(&survey.Input{
		Message: "phone/email",
		Default: "",
		Help:    "",
	}, survey.WithValidator(validate))
	password := PromptPassword("password")

	loginDone := make(chan bool)
	go printLoading(loginDone, 0)

	if u := biliAPI.DoLogin(loginUserName, password); u != nil && u.Code == 0 {
		ret_, _ := biliAPI.GetUserInfo(u.Data.TokenInfo.Mid)
		AccountSelected.Info = ret_.Data
		fmt.Print("\n")
		loginDone <- true

		Logger.Info("Your token is:", u.Data.TokenInfo.AccessToken)
		Logger.Info("Your SESSDATA is:", u.Data.CookieInfo.CookiesMap["SESSDATA"])
		Logger.Info("Expire Date:", time.Unix(u.Data.CookieInfo.Cookies[0].Expires, 0))

		if _, err := db.Exec(`
			REPLACE INTO Account(uid,nikeName,loginUserName,accessToken,expire,SESSDATA,sid,DedeUserID__ckMd5,lastUsedTimestamp,blocked) 
			VALUES (?,?,?,?,?,?,?,?,?,?);`,
			u.Data.TokenInfo.Mid, AccountSelected.Info.Name, loginUserName,
			u.Data.TokenInfo.AccessToken,
			time.Now().Unix()+u.Data.TokenInfo.ExpiresIn,
			u.Data.CookieInfo.CookiesMap["SESSDATA"], u.Data.CookieInfo.CookiesMap["sid"], u.Data.CookieInfo.CookiesMap["DedeUserID__ckMd5"],
			time.Now().Unix(), 0,
		); err != nil {
			Logger.Error("update dab failed: ", err)
		}
	} else if u != nil {
		loginDone <- false
		Logger.Warn(u.Message)
		showPageAccount()
	} else {
		loginDone <- false
		Logger.Error("unknown error happened!")
	}
	Logger.Info("remember to select the account")
	showPageHome()
}

func pageAccountCmdEdit() {
	Logger.Warn("unimplemented method!")
	showPageAccount()
}

func pageAccountCmdDelete() {
	Logger.Warn("unimplemented method!")
	showPageAccount()
}

func pageAccountCmdBlock() {
	Logger.Warn("unimplemented method!")
	showPageAccount()
}
