package handler

import (
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
	var chatID int64

	for update := range updates {
		if update.CallbackQuery != nil {
			log.Printf("[%s] %s", update.CallbackQuery.From.FirstName, update.CallbackQuery.Data)
			chatID = members.FindChatIDByUsername(update.CallbackQuery.From.UserName)
		} else {

			log.Printf("[%s] %s", update.Message.From.FirstName, update.Message.Text)
			chatID = update.Message.Chat.ID

			if members.IsMemberNewUser(chatID) {
				err := members.AddNewMember(update.Message.From, chatID)
				if err != nil {
					log.Println(err)
				}
				_ = members.ChangeState(data.Start, chatID)
			}

			if strings.HasPrefix(update.Message.Text, "new_authorization|") && update.Message.From.UserName == "Pooria2" {
				auth = update.Message.Text[18:]
				log.Println("new authorization sets!")
				msg := tgbotapi.NewMessage(chatID, "New authorization sets!")
				_, _ = Bot.Send(msg)
				continue
			}

			// Todo: when we receive GetIP message we have to send IP to Pooria

			if update.Message.IsCommand() {
				if members.GetState(chatID) == data.HomePage ||
					members.GetState(chatID) == data.ViewingAccountStats ||
					members.GetState(chatID) == data.ViewingProfile {
					_, err := Bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, members[chatID].LastMessageID))
					if err != nil {
						log.Println(err)
					}
				}
				if update.Message.Command() == "start" {
					err := members.ChangeState(data.Start, chatID)
					if err != nil {
						log.Println(err)
					}
					msg := CreateStartMessage(chatID)
					_, _ = Bot.Send(msg)
					continue
				}
				if update.Message.Command() == "about_us" {
					msg := CreateAboutUsMessage(chatID)
					err := members.ChangeState(data.AboutUs, chatID)
					if err != nil {
						log.Println(err)
					}
					_, _ = Bot.Send(msg)
					continue
				}
			}
		}

		// Todo: All of keyboards should be inline. We can't have loading msg with button keyboard.

		switch members.GetState(chatID) {
		case data.Start:
			if update.Message.Text == LogInMessage {
				msg := CreateLoginMessage(chatID)
				err := members.ChangeState(data.LogIn, chatID)
				if err != nil {
					log.Println(err)
				}
				sentMessage, _ := Bot.Send(msg)
				members[chatID].LastMessageID = sentMessage.MessageID
				continue
			} else if update.Message.Text == AboutMessage {
				msg := CreateAboutUsMessage(chatID)
				err := members.ChangeState(data.AboutUs, chatID)
				if err != nil {
					log.Println(err)
				}
				sentMessage, _ := Bot.Send(msg)
				members[chatID].LastMessageID = sentMessage.MessageID
				continue
			}
		case data.AboutUs:
			if update.Message.Text == BackMessage {
				_ = members.ChangeState(data.Start, chatID)
				msg := CreateStartMessage(chatID)
				sentMessage, _ := Bot.Send(msg)
				members[chatID].LastMessageID = sentMessage.MessageID
				continue
			}
		case data.LogIn:
			m, err := SendLoadingMessage(chatID)
			if err != nil {
				log.Println(err)
			}
			members[chatID].LastMessageID = m.MessageID
			tag := update.Message.Text
			royaleClient = NewClient(auth)
			_, err = royaleClient.Player(tag)
			if err != nil {
				log.Println(err)
				_, _ = Bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, members[chatID].LastMessageID))
				msg := tgbotapi.NewMessage(chatID, "Invalid Tag.\nMake sure the tag starts with '#'\nLike: #ABCDE")
				sentMessage, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
				members[chatID].LastMessageID = sentMessage.MessageID
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
			_, _ = Bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, members[chatID].LastMessageID))
			msg := CreateHomePageMessage(chatID, -1)
			sentMessage, err := Bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			members[chatID].LastMessageID = sentMessage.MessageID
			continue

		case data.HomePage:
			if update.CallbackQuery == nil {
				continue
			}
			if update.CallbackQuery.Data == ProfileInfoQuery {
				msg1 := tgbotapi.NewMessage(chatID, "")
				msg1.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true, Selective: false}
				_, _ = Bot.Send(msg1)

				msg := CreateProfileInfoMessage(chatID, update.CallbackQuery.Message.MessageID)
				_ = members.ChangeState(data.ViewingProfile, chatID)
				sentMessage, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
					_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "something went wrong!"))
					continue
				}
				members[chatID].LastMessageID = sentMessage.MessageID
				_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done!"))
				continue

			} else if update.CallbackQuery.Data == PlayerStatsQuery {

				msg := CreateProfileStatsMessage(chatID, update.CallbackQuery.Message.MessageID)
				_ = members.ChangeState(data.ViewingAccountStats, chatID)
				sentMessage, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
					_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "something went wrong!"))
					continue
				}
				members[chatID].LastMessageID = sentMessage.MessageID
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
				_, _ = Bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, members[chatID].LastMessageID))
				sentMessage, _ := Bot.Send(msg)
				members[chatID].LastMessageID = sentMessage.MessageID
				continue
			}
			continue
		case data.ViewingAccountStats, data.ViewingProfile:
			if update.CallbackQuery == nil {
				continue
			}
			if update.CallbackQuery.Data == BackMessage {
				msg := CreateHomePageMessage(chatID, update.CallbackQuery.Message.MessageID)
				_ = members.ChangeState(data.HomePage, chatID)
				sentMessage, err := Bot.Send(msg)
				if err != nil {
					log.Println(err)
					_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "something went wrong!"))
					continue
				}
				members[chatID].LastMessageID = sentMessage.MessageID
				_, _ = Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
				continue
			}
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
