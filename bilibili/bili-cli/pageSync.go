package main

import (
	"github.com/AlecAivazis/survey/v2"
)

func showPageSync() {

	Logger.Error("本地模式测试中...")
	showPageHome()
	return

	pageSyncSelectItems := []string{
		SyncCmdSyncFollowing,
		SyncCmdSyncLive,
		SyncCmdSyncBoth,
		CMDHome,
		CMDExit,
	}

	nextRoute := PromptSelect(
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
		//Logger.Fatal("unknown route:", nextRoute)
	}
}
