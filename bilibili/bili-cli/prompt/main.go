package main

import (
	"github.com/AlecAivazis/survey/v2"
)

func PromptSelect(question string, options []string, opt ...survey.AskOpt) (ret string) {
	if err := survey.AskOne(&survey.Select{
		Message: question,
		Options: options,
	}, &ret, opt...); err != nil {
		//Logger.Error(err)
	}
	//Logger.Info("your select: %q\n", ret)
	return
}

func main() {
	PromptSelect("test", []string{"1. 123\n  456\n  789","2. aaa\n  bb\n  ccc"})
}
