package main

import (
	"github.com/0xAX/notificator"
	"log"
)

func notify(title string, text string) {
	notify := notificator.New(notificator.Options{
		AppName: "Wifilogin",
	})
	err := notify.Push(title, text, "", notificator.UR_NORMAL)
	if err != nil {
		log.Fatal(err)
	}
}
