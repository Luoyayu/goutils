package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func showPageLive() {
	options := []string{
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

	nextRoute := promptSelect(
		AccountSelected.NikeName+": "+fmt.Sprint(AccountSelected.Uid),
		options,
		survey.WithPageSize(len(options)),
	)

	switch nextRoute {
	case LiveCmdSelect:
		pageLivePlaySelected2Select = true
		pageLiveSelect(true)
	case LiveCmdStop:
		stopMpvSafely()
		showPageLive()
	case LiveCmdRecommend:

		pageLiveRecommend()
	case LiveCmdBlock:
		pageLiveBlock()
	case LiveCmdAdd:
		pageLiveAdd()
	case LiveCmdEdit:
		pageLiveEdit()
	case LiveCmdDelete:
		pageLiveDelete()
	case CMDHome:
		showPageHome()
	case CMDExit:
		exitClear()
	default:
		//Logger.Fatal("unknown route:", nextRoute)
	}
}
