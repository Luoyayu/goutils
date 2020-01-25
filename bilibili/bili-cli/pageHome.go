package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func showPageHome() {

	homePageSelectItems := []string{
		Home2Account,
		Home2Live,
		Home2Following,
		Home2Sync,
		Home2Setting,
		CMDExit,
	}

	nextRoute := PromptSelect(
		promptui.Styler(promptui.FGRed)("hallo! "+AccountSelected.NikeName+": "+fmt.Sprint(AccountSelected.Uid)),
		homePageSelectItems,
	)

	switch nextRoute {
	case Home2Account:
		showPageAccount()
	case Home2Live:
		showPageLive()
	case Home2Following:
		showPageFollowing()
	case Home2Sync:
		showPageSync()
	case Home2Setting:
		showPageSetting()
	case CMDExit:
		exitClear()

	default:
		//Logger.Fatal("unknown route:", nextRoute)
	}
}
