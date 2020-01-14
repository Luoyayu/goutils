package main

import "github.com/manifoldco/promptui"

// home route
var (
	Home2Account   = "Account"
	Home2Live      = "Live"
	Home2Following = "Following"
	Home2Setting   = "Setting" + unImplemented
	Home2Sync      = "Sync"
	//Home2Home      = "Home"
	//Home2Exit      = "Exit"
)

var (
	implemented   = promptui.Styler(promptui.FGGreen)(" √")
	unImplemented = promptui.Styler(promptui.FGRed)(" X")
)

// Account route
var (
	AccountRouteId = "@" + Home2Account

	AccountCmdSelect       = CMDSelect + AccountRouteId
	AccountCmdBlock        = CMDBlock + AccountRouteId + unImplemented
	AccountCmdAdd          = CMDAdd + AccountRouteId
	AccountCmdEdit         = CMDEdit + AccountRouteId + unImplemented
	AccountCmdDelete       = CMDDelete + AccountRouteId + unImplemented
	AccountCmdRefreshToken = CMDRefreshToken + AccountRouteId + unImplemented
)

// Live route
var (
	LiveRouteId = "@" + Home2Live

	LiveCmdSelect    = CMDSelect + LiveRouteId
	LiveCmdStop      = CMDStop + LiveRouteId
	LiveCmdBlock     = CMDBlock + LiveRouteId + unImplemented
	LiveCmdAdd       = CMDAdd + LiveRouteId
	LiveCmdRecommend = CMDRecommend + LiveRouteId
	LiveCmdEdit      = CMDEdit + LiveRouteId + unImplemented
	LiveCmdDelete    = CMDDelete + LiveRouteId + unImplemented
)

// Following route
var (
	FollowingRouteId = "@" + Home2Following

	FollowingCmdSelect = CMDSelect +"Video"+ FollowingRouteId
	FollowingCmdBlock  = CMDBlock + FollowingRouteId + unImplemented
	FollowingCmdAdd    = CMDAdd + FollowingRouteId
	FollowingCmdEdit   = CMDEdit + FollowingRouteId + unImplemented
	FollowingCmdDelete = CMDDelete + FollowingRouteId + unImplemented
)

// Setting route
var (
	SettingRouteId = "@" + Home2Setting

	SettingCmdSingleMode = "single-user" + SettingRouteId + unImplemented
	SettingCmdMultiMode  = "multi-user" + SettingRouteId + unImplemented
	SettingCmdReset      = CMDReset + SettingRouteId + unImplemented
)

// Sync route
var (
	SyncRouteId = "@" + Home2Sync

	SyncCmdSyncFollowing = "sync following" + SyncRouteId
	SyncCmdSyncLive      = "sync live" + SyncRouteId
	SyncCmdSyncBoth      = "sync both" + SyncRouteId

	SyncCmdUnCover = "uncover" + SyncRouteId + unImplemented // inherit `block` property
	SyncCmdCover   = "cover" + SyncRouteId + unImplemented
)

// common route
var (
	CMDBack = "back"
	CMDHome = "home"
	CMDExit = "exit"
)

// common ops
var (
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
