package main

// home route
const (
	Home2Account   = "Account"
	Home2Live      = "Live"
	Home2Following = "Following"
	Home2Setting   = "Setting"
	Home2Sync      = "Sync"
	//Home2Home      = "Home"
	//Home2Exit      = "Exit"
)

// Account route
const (
	AccountRouteId = "@" + Home2Account

	AccountCmdSelect       = CMDSelect + AccountRouteId
	AccountCmdBlock        = CMDBlock + AccountRouteId
	AccountCmdAdd          = CMDAdd + AccountRouteId
	AccountCmdEdit         = CMDEdit + AccountRouteId
	AccountCmdDelete       = CMDDelete + AccountRouteId
	AccountCmdRefreshToken = CMDRefreshToken + AccountRouteId
)

// Live route
const (
	LiveRouteId = "@" + Home2Live

	LiveCmdSelect    = CMDSelect + LiveRouteId
	LiveCmdStop      = CMDStop + LiveRouteId
	LiveCmdBlock     = CMDBlock + LiveRouteId
	LiveCmdAdd       = CMDAdd + LiveRouteId
	LiveCmdRecommend = CMDRecommend + LiveRouteId
	LiveCmdEdit      = CMDEdit + LiveRouteId
	LiveCmdDelete    = CMDDelete + LiveRouteId
)

// Following route
const (
	FollowingRouteId = "@" + Home2Following

	FollowingCmdSelect = CMDSelect + FollowingRouteId
	FollowingCmdBlock  = CMDBlock + FollowingRouteId
	FollowingCmdAdd    = CMDAdd + FollowingRouteId
	FollowingCmdEdit   = CMDEdit + FollowingRouteId
	FollowingCmdDelete = CMDDelete + FollowingRouteId
)

// Setting route
const (
	SettingRouteId = "@" + Home2Setting

	SettingCmdSingleMode = "single-user" + SettingRouteId
	SettingCmdMultiMode  = "multi-user" + SettingRouteId
	SettingCmdReset      = CMDReset + SettingRouteId
)

// Sync route
const (
	SyncRouteId = "@" + Home2Sync

	SyncCmdSyncFollowing = "sync following" + SyncRouteId
	SyncCmdSyncLive      = "sync live" + SyncRouteId
	SyncCmdSyncBoth      = "sync both" + SyncRouteId

	SyncCmdUnCover = "uncover" + SyncRouteId // inherit `block` property
	SyncCmdCover   = "cover" + SyncRouteId
)

// common route
const (
	CMDBack = "back"
	CMDHome = "home"
	CMDExit = "exit"
)

// common ops
const (
	CMDAdd          = "add"
	CMDDelete       = "delete"
	CMDEdit         = "edit"
	CMDSearch       = "search"
	CMDBlock        = "block"
	CMDShow         = "show"
	CMDSelect       = "select"
	CMDPlay         = "play"
	CMDReset        = "reset"
	CMDRefresh      = "refresh"
	CMDStart        = "start"
	CMDStop         = "stop"
	CMDRestart      = "restart"
	CMDRecommend    = "recommend"
	CMDRefreshToken = "refresh token"
)
