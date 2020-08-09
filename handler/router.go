package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pooria1/clash-royale-telegram-bot/data"
	"log"
	"net/http"
	"strconv"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1037921974:AAH7XCyAy-eVIUwTlcmLPX_WKCUD069qGzg")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		tag := update.Message.Text

		log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, getPlayerStats(tag))
		msg.ReplyToMessageID = update.Message.MessageID
		fmt.Println(msg)

	}
}

func getPlayerStats(playerTag string) string {
	//Q2LL8LYUG

	url := "https://api.clashroyale.com/v1/players/%23" + playerTag
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6IjU0YWM1MGJiLTZhMjEtNDg0Mi04NTkyLWVhNGU0NDgwYWRmOSIsImlhdCI6MTU5Njk3MzcxNSwic3ViIjoiZGV2ZWxvcGVyL2E0MmNjYzI3LWIzYWItZGYxYS05YTFiLWY1YTk3ZjI0ZWFjZiIsInNjb3BlcyI6WyJyb3lhbGUiXSwibGltaXRzIjpbeyJ0aWVyIjoiZGV2ZWxvcGVyL3NpbHZlciIsInR5cGUiOiJ0aHJvdHRsaW5nIn0seyJjaWRycyI6WyI1NC4xNzQuMTU1Ljk5Il0sInR5cGUiOiJjbGllbnQifV19.2pILzRxmvSs4AFVnbc8OYeKLEEQzEvUM080LHm7Ck4Gc6v8LKUEhVEjwDiPkJqh5Ot86tA6kkxBbb60MpJRmKw")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("response status:", res.Status)
	if err != nil {
		fmt.Println(err)
	} else if res.StatusCode != 200 {
		fmt.Printf("none 200 response from [%s]\n", url)
		return "error :("
	}
	p := &data.Player{}
	if err = json.NewDecoder(res.Body).Decode(p); err != nil {
		fmt.Println(err)
	}
	message := "name: " + p.Name + "\nlevel: " + strconv.FormatInt(int64(p.ExpLevel), 10)
	return message
}
