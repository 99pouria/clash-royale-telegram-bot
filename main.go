package main

import (
	"github.com/pooria1/clash-royale-telegram-bot/data"
	"github.com/pooria1/clash-royale-telegram-bot/handler"
	"log"
)

var members data.MyBotMembers

func main() {

	members = members.Init()

	updates, err := handler.InitBot("1037921974:AAH7XCyAy-eVIUwTlcmLPX_WKCUD069qGzg")

	if err != nil {
		log.Println(err)
	}

	handler.HandleMessages(updates, members)
}
