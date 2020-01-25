package main

import "github.com/AlecAivazis/survey/v2"

func showPageFollowing() {
	options := []string{
		FollowingCmdSelect,
		FollowingCmdPlay,
		CMDHome,
		CMDExit,
	}

	nextRoute := PromptSelect("following", options, survey.WithPageSize(len(options)))

	switch nextRoute {
	case FollowingCmdSelect:
		pageFollowingSelect()
	case FollowingCmdPlay:
		pageFollowingCmdPlay()
	case CMDHome:
		showPageHome()
		return
	case CMDExit:
		exitClear()
	}

	showPageHome()
}
