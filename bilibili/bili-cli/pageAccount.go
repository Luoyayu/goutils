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

	/*accountPage := promptui.Select{
		Label: "Account: " + promptui.Styler(promptui.FGRed)(AccountSelected.NikeName) + " - " + promptui.Styler(promptui.FGRed)(AccountSelected.Uid),
		Items: accountPageSelectedItems,
		Size:  len(accountPageSelectedItems),
	}

	nextPage := ""
	err := errors.New("")

	_, nextPage, err = accountPage.Run()
	if err != nil {
		Logger.Fatal(err)
	}*/

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
		Logger.Fatal("unknown route:", nextPage)
	}
}
