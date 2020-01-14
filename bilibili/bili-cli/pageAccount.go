package main

import (
	"github.com/manifoldco/promptui"
)

func showPageAccount() {
	accountPageSelectedItems := []string{
		AccountCmdSelect,
		AccountCmdAdd,
		AccountCmdEdit,
		AccountCmdDelete,
		AccountCmdBlock,
		CMDHome,
		CMDExit,
	}

	nextPage := promptSelect(
		"Account: "+promptui.Styler(promptui.FGRed)(AccountSelected.NikeName)+" - "+promptui.Styler(promptui.FGRed)(AccountSelected.Uid),
		accountPageSelectedItems,
	)

	switch nextPage {
	case AccountCmdSelect:
		pageAccountSelect()
	case AccountCmdAdd:
		pageAccountCMDAdd()
	case AccountCmdEdit:
		pageAccountCmdEdit()
	case AccountCmdDelete:
		pageAccountCmdDelete()
	case AccountCmdBlock:
		pageAccountCmdBlock()
	case CMDHome:
		showPageHome()
	case CMDExit:
		exitClear()
	default:
		//Logger.Fatal("unknown route:", nextPage)
	}
}
