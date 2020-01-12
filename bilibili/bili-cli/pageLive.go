package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func showPageLive() {
	pageLiveSelectItems := []string{
		LiveCmdSelect,
		LiveCmdStop,
		LiveCmdAdd,
		LiveCmdRecommend,
		LiveCmdBlock,
		LiveCmdEdit,
		LiveCmdDelete,
		CMDHome,
		CMDExit,
	}

	/*	livePage := promptui.Select{
			Label: AccountSelected.NikeName + ": " + fmt.Sprint(AccountSelected.Uid),
			Items: pageLiveSelectItems,
			Size:  len(pageLiveSelectItems),
		}

		var nextRoute = ""
		_, nextRoute, _ = livePage.Run()
	*/
	nextRoute := promptSelect(
		AccountSelected.NikeName+": "+fmt.Sprint(AccountSelected.Uid),
		pageLiveSelectItems,
		survey.WithPageSize(len(pageLiveSelectItems)),
	)

	switch nextRoute {
	case LiveCmdSelect:
		pageLiveCmdSelect()
	case LiveCmdStop:
		stopMpvSafely()
		showPageLive()
	case LiveCmdRecommend:
		pageLiveCmdRecommend()
	case LiveCmdBlock:
		pageLiveCmdBlock()
	case LiveCmdAdd:
		pageLiveCmdAdd()
	case LiveCmdEdit:
		pageLiveCmdEdit()
	case LiveCmdDelete:
		pageLiveCmdDelete()
	case CMDHome:

		showPageHome()
	case CMDExit:
		exitClear()
	default:
		Logger.Fatal("unknown route:", nextRoute)
	}
}
