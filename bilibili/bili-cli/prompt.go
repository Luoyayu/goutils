package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
)

func promptSelect(question string, options []string, opt ...survey.AskOpt) (ret string) {
	if err := survey.AskOne(&survey.Select{
		Message: question,
		Options: options,
	}, &ret, opt...); err != nil {
		//Logger.Error(err)
	}
	//Logger.Info("your select: %q\n", ret)
	return
}

func promptPassword(question string, opt ...survey.AskOpt) (ret string) {
	if err := survey.AskOne(&survey.Password{
		Message: promptui.Styler(promptui.FGCyan)(question),
	}, &ret, opt...,
	); err != nil {
		//Logger.Error(err)
	}
	//Logger.Info("your password: %q\n", ret)
	return
}

func promptInput(question *survey.Input, opt ...survey.AskOpt) (ret string) {
	if err := survey.AskOne(question, &ret, opt...); err != nil {
		//Logger.Error(err)
	}

	//Logger.Info("your password: %q\n", ret)
	return
}

func promptConfirm(question string, opt ...survey.AskOpt) (ret bool) {
	if err := survey.AskOne(&survey.Confirm{
		Message: question,
	}, &ret, opt...); err != nil {
		//Logger.Error(err)
	}
	return
}
