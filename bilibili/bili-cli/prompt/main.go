package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/luoyayu/goutils/logger"
	"github.com/manifoldco/promptui"
	"time"
)

var Logger = logger.NewDebugLogger()

func promptSelect(question string, options []string) (ret string) {
	if err := survey.AskOne(&survey.Select{
		Message: question,
		Options: options,
	}, &ret); err != nil {
		Logger.Error(err)
	}
	//Logger.Info("your select: %q\n", ret)
	return
}

func promptPassword(question string) (ret string) {
	if err := survey.AskOne(&survey.Password{
		Message: promptui.Styler(promptui.FGCyan)("Please type your password"),
	}, &ret); err != nil {
		Logger.Error(err)
	}
	//Logger.Info("your password: %q\n", ret)
	return
}

func testSelect() {
	color := ""
	prompt := &survey.Select{
		Message: promptui.Styler(promptui.FGYellow)("Choose a color:"),
		Options: []string{
			promptui.Styler(promptui.FGRed, promptui.FGUnderline)("red"),
			promptui.Styler(promptui.FGBlue, promptui.FGUnderline)("blue"),
			promptui.Styler(promptui.FGGreen, promptui.FGUnderline)("green"),
		},
	}
	if err := survey.AskOne(prompt, &color); err != nil {
		panic(err)
	}
	Logger.Infof("you choose: %q\n", color)
}

func testPassword() {
	password := ""
	prompt := &survey.Password{
		Message: promptui.Styler(promptui.FGCyan)("Please type your password"),
	}
	if err := survey.AskOne(prompt, &password); err != nil {
		panic(err)
	}
	Logger.Errorf("your password: %q\n", password)
}

func testMultiSelect() {

	var days []string
	prompt := &survey.MultiSelect{
		Message: promptui.Styler(promptui.FGRed)("What days do you prefer:"),
		Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Gameday"},
		Default: []string{"Sunday", "Saturday"},
		Help:    "使用空格选择或取消选择, 回车确认选择",
	}
	go func() {
		tm := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-tm.C:
				return
			default:
				Logger.Info("log test")
				time.Sleep(time.Second * 1)
			}
		}
	}()

	if err := survey.AskOne(prompt, &days); err != nil {
		panic(err)
	}
	Logger.Infof("your select: %q\n", days)
}
func testConfirm() {
	yn := false
	prompt := &survey.Confirm{
		Message: "Do you like Macintosh operating systems?",
	}
	if err := survey.AskOne(prompt, &yn); err != nil {
		panic(err)
	}
	Logger.Warnf("your choose: %v\n", yn)
}
func main() {
	testSelect()
	testPassword()
	testMultiSelect()
	testConfirm()
}
