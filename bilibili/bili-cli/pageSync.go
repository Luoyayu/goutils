package main

import (
	"github.com/AlecAivazis/survey/v2"
)

func showPageSync() {
	pageSyncSelectItems := []string{
		SyncCmdSyncFollowing,
		SyncCmdSyncLive,
		SyncCmdSyncBoth,
		CMDHome,
		CMDExit,
	}

	/*	pageSync := promptui.Select{
			Label: "sync followings or live or both from cloud",
			Items: pageSyncSelectItems,
			Size:  len(pageSyncSelectItems),
		}

		_, nextRoute, _ := pageSync.Run()
	*/
	nextRoute := promptSelect(
		"sync followings or live or both from cloud",
		pageSyncSelectItems,
		survey.WithPageSize(len(pageSyncSelectItems)),
	)

	switch nextRoute {
	case SyncCmdSyncFollowing:
		pageSyncCmdSyncFollowing(true)
		showPageSync()
	case SyncCmdSyncLive:
		pageSyncCmdSyncLive(true, false, true)
		showPageSync()
	case SyncCmdSyncBoth:
		pageSyncCmdSyncBoth()
		showPageSync()
	case CMDHome:
		showPageHome()
	case CMDExit:
		exitClear()
	default:
		Logger.Fatal("unknown route:", nextRoute)
	}
}
