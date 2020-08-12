package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

var auth string

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

		if strings.HasPrefix(update.Message.Text, "new_authorization|") {
			auth = string([]rune(update.Message.Text)[18:])
			fmt.Println("new authorization sets!")
			continue
		}
		tag := update.Message.Text

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, getPlayerStats(tag))

		// Keyboard button code:

		//msg.ParseMode = tgbotapi.ModeHTML
		//keyboard := tgbotapi.NewReplyKeyboard(
		//	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Kang stickers"), tgbotapi.NewKeyboardButton("Get original sticker")),
		//	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Create new pack"), tgbotapi.NewKeyboardButton("List my packs")),
		//	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Synchronize with Telegram Servers")),
		//	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Help")))
		//keyboard.OneTimeKeyboard = true
		//msg.ReplyMarkup = keyboard

		log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)

		msg.ReplyToMessageID = update.Message.MessageID
		fmt.Println(msg.Text)
		bot.Send(msg)

	}
}

func getPlayerStats(playerTag string) string {
	c := NewClient("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiIsImtpZCI6IjI4YTMxOGY3LTAwMDAtYTFlYi03ZmExLTJjNzQzM2M2Y2NhNSJ9.eyJpc3MiOiJzdXBlcmNlbGwiLCJhdWQiOiJzdXBlcmNlbGw6Z2FtZWFwaSIsImp0aSI6ImZkYzk2ZTEzLTk1NzYtNDdjNy1hYzM0LTVlMjJkMTA3NzZlMSIsImlhdCI6MTU5NzI1NTUxNCwic3ViIjoiZGV2ZWxvcGVyL2E0MmNjYzI3LWIzYWItZGYxYS05YTFiLWY1YTk3ZjI0ZWFjZiIsInNjb3BlcyI6WyJyb3lhbGUiXSwibGltaXRzIjpbeyJ0aWVyIjoiZGV2ZWxvcGVyL3NpbHZlciIsInR5cGUiOiJ0aHJvdHRsaW5nIn0seyJjaWRycyI6WyIyMy4xMDUuMTcwLjEzNCJdLCJ0eXBlIjoiY2xpZW50In1dfQ.xs9NgHlJa_3JQZQR3upYYylTjwTxIi_mqpaZ5WSR9Yp72Jl8gRGNFIXEmhxpPqUXwV_IViFxcTaBEN58E7B1Mw")
	p, err := c.Player(playerTag)
	if err != nil {
		return err.Error()
	}
	b, err := c.BattleLog(playerTag)
	message := "name: " + p.Name + "\nlevel: " + strconv.FormatInt(int64(p.ExpLevel), 10)
	message += "\nwins: " + strconv.FormatInt(int64(p.Wins), 10)
	message += "\nlosses: " + strconv.FormatInt(int64(p.Losses), 10)
	message += "\nrecent wins percentage: " + fmt.Sprint(b.TrophyList())
	message += "\nrecent wins percentage: " + fmt.Sprint(b.RecentWinPercentage())
	message += "\ntotal wins percentage: " + fmt.Sprint(p.TotalWinPercentage())
	return message
}
