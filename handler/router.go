package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pooria1/clash-royale-telegram-bot/data"
	"log"
	"strings"
)

const (
	PlayerStatsQuery   = "player_stats_query📊📈"
	ProfileInfoQuery   = "account_info_query🎫"
	ChangeAccountQuery = "change_account_query"
	AboutMessage       = "About usℹ️"
	LogInMessage       = "Log in"
	BackMessage        = "Back"
)

var Bot *tgbotapi.BotAPI
var auth string
var members data.MyBotMembers
var royaleClient Client

func HandleMessages(updates tgbotapi.UpdatesChannel, m data.MyBotMembers) {
	members = m
	for update := range updates {
		var chatID int64
		if update.CallbackQuery != nil {
			log.Printf("[%s] %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.Data)
			chatID = members.FindChatIDByUsername(update.CallbackQuery.From.UserName)
		} else {

			log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)
			chatID = update.Message.Chat.ID

			if members.IsMemberNewUser(chatID) {
				err := members.AddNewMember(update.Message.From, chatID)
				log.Println(err)
			}
			if strings.HasPrefix(update.Message.Text, "new_authorization|") && update.Message.From.UserName == "Pooria2" {
				auth = update.Message.Text[18:]
				fmt.Println("new authorization sets!")
				msg := tgbotapi.NewMessage(chatID, "Ok shod dash Pooria")
				_, _ = Bot.Send(msg)
				continue
			}

			if update.Message.IsCommand() {
				fmt.Println("________command detected:" + update.Message.Command())
				if update.Message.Command() == "start" {
					err := members.ChangeState(data.Start, chatID)
					if err != nil {
						log.Println(err)
					}
					msg := CreateStartMessage(chatID)
					_, _ = Bot.Send(msg)
				}
				if update.Message.Command() == "about_us" {
					msg := CreateAboutUsMessage(chatID)
					err := members.ChangeState(data.AboutUs, chatID)
					if err != nil {
						log.Println(err)
					}
					_, _ = Bot.Send(msg)
				}
			}
		}
		switch members.GetState(chatID) {
		case data.Start:
			if update.Message.Text == LogInMessage {
				msg := CreateLoginMessage(chatID)
				err := members.ChangeState(data.LogIn, chatID)
				if err != nil {
					log.Println(err)
				}
				_, _ = Bot.Send(msg)
				continue
			} else if update.Message.Text == AboutMessage {
				msg := CreateAboutUsMessage(chatID)
				err := members.ChangeState(data.AboutUs, chatID)
				if err != nil {
					log.Println(err)
				}
				_, _ = Bot.Send(msg)
				continue
			}
		case data.AboutUs:
			if update.Message.Text == BackMessage {
				members.ChangeState(data.Start, chatID)
				msg := CreateStartMessage(chatID)
				_, _ = Bot.Send(msg)
				continue
			}
		case data.LogIn:
			tag := update.Message.Text
			royaleClient = NewClient(auth)
			_, err := royaleClient.Player(tag)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Invalid Tag.\nMake sure the tag starts with '#'\nLike: #ABCDE")
				msg.ParseMode = tgbotapi.ModeHTML
				log.Println(err)
				_, _ = Bot.Send(msg)
				continue
			}
			err = members.ChangeCurrentTag(tag, chatID)
			if err != nil {
				log.Println(err)
			}
			err = members.ChangeState(data.HomePage, chatID)
			if err != nil {
				log.Println(err)
			}
			msg := CreateHomePageMessage(chatID, -1)
			_, _ = Bot.Send(msg)
			continue
		case data.HomePage:
			if update.CallbackQuery.Data == ProfileInfoQuery {

				msg := CreateProfileInfoMessage(chatID, update.CallbackQuery.Message.MessageID)
				_ = members.ChangeState(data.ViewingProfile, chatID)
				_, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
				continue

			} else if update.CallbackQuery.Data == PlayerStatsQuery {

				msg := CreateProfileStatsMessage(chatID, update.CallbackQuery.Message.MessageID)
				_ = members.ChangeState(data.ViewingAccountStats, chatID)
				_, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
				continue

			} else if update.CallbackQuery.Data == ChangeAccountQuery {
				lastPm := tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID)
				_, _ = Bot.DeleteMessage(lastPm)
				msg := CreateLoginMessage(chatID)
				err := members.ChangeState(data.LogIn, chatID)
				if err != nil {
					log.Println(err)
				}
				_, _ = Bot.Send(msg)
				continue
			}
			continue
		case data.ViewingAccountStats, data.ViewingProfile:
			if update.CallbackQuery.Data == BackMessage {
				msg := CreateHomePageMessage(chatID, update.CallbackQuery.Message.MessageID)
				_ = members.ChangeState(data.HomePage, chatID)
				_, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
				continue
			}
		}
		if update.Message == nil || update.CallbackQuery == nil { // ignore any non-Message Updates

			// I need to show an invalid message and return last state

			continue
		}

	}
}

func InitBot(token string) (tgbotapi.UpdatesChannel, error) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	Bot.Debug = false

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return updates, nil
}
