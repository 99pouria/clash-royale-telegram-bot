package main

import (
	"log"

	"github.com/pooria1/clash-royale-telegram-bot/handler"
	tbotapi "github.com/pooria1/clash-royale-telegram-bot/pkg/tbot-api"
)

func main() {

	members := tbotapi.Init()
	updates, err := handler.InitBot("1037921974:AAH7XCyAy-eVIUwTlcmLPX_WKCUD069qGzg")

	if err != nil {
		log.Println(err)
	}

	handler.HandleMessages(updates, members)
}
