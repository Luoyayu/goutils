package main

import "github.com/AlecAivazis/survey/v2"

func showPageFollowing() {
	options := []string{
		FollowingCmdSelect,
		FollowingCmdAdd,
		CMDHome,
		CMDExit,
	}

	nextRoute := promptSelect("following", options,
		survey.WithPageSize(len(options),
		),
	)

	switch nextRoute {
	case FollowingCmdSelect:
		pageFollowingCmdSelect()
	case FollowingCmdAdd:
		pageFollowingCmdAdd()
	case CMDHome:
		showPageHome()
		return
	case CMDExit:
		exitClear()
	}

	showPageHome()
}
