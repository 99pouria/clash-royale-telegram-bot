package main

import (
	"log"

	"github.com/99pouria/clash-royale-telegram-bot/handler"
	tbotapi "github.com/99pouria/clash-royale-telegram-bot/pkg/tbot-api"
)

func main() {

	members := tbotapi.Init()
	updates, err := handler.InitBot("")

	if err != nil {
		log.Println(err)
	}

	handler.HandleMessages(updates, members)
}
