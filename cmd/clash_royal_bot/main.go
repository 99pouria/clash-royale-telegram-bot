package main

import (
	"log"

	"github.com/99pouria/clash-royale-telegram-bot/internal/config"
	"github.com/99pouria/clash-royale-telegram-bot/internal/handler"
	tbotapi "github.com/99pouria/clash-royale-telegram-bot/pkg/tbot-api"
)

func main() {
	// loading config
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	members := tbotapi.Init()
	updates, err := handler.InitBot(config.GetBotToken())

	if err != nil {
		log.Println(err)
	}

	handler.HandleMessages(updates, members)
}
