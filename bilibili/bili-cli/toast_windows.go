// +build windows

package main

import (
	"gopkg.in/toast.v1"
	"log"
)

func _toast_demo_() {
	notification := toast.Notification{
		AppID:   "Example App",
		Title:   "My notification",
		Message: "Some message about how important something is...",
		//Icon:    "go.png", // This file must exist (remove this line if it doesn't)
		Actions: []toast.Action{
			{"protocol", "I'm a button", ""},
			{"protocol", "Me too!", ""},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
