package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pooria1/clash-royale-telegram-bot/data"
	"log"
	"strconv"
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
			c := NewClient(auth)
			_, err := c.Player(tag)
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
			msg := CreateHomePageMessage(chatID)
			_, _ = Bot.Send(msg)
			continue
		case data.HomePage:
			if update.CallbackQuery.Data == ProfileInfoQuery {
			}
		}
		if update.Message == nil || update.CallbackQuery == nil { // ignore any non-Message Updates

			// I need to show an invalid message and return last state

			continue
		}

	}
}

func getPlayerStats(playerTag string) string {
	c := NewClient(auth)
	p, err := c.Player(playerTag)
	if err != nil {
		return err.Error()
	}
	b, err := c.BattleLog(playerTag)
	message := "name: " + p.Name + "\nking level: " + strconv.FormatInt(int64(p.ExpLevel), 10)
	message += "\nwins: " + strconv.FormatInt(int64(p.Wins), 10)
	message += "\nlosses: " + strconv.FormatInt(int64(p.Losses), 10)
	message += "\nrecent wins percentage: " + fmt.Sprint(b.TrophyList())
	message += "\nrecent wins percentage: " + fmt.Sprint(b.RecentWinPercentage(false))
	message += "\nrecent rank games wins percentage: " + fmt.Sprint(b.RecentWinPercentage(true))
	message += "\ntotal wins percentage: " + fmt.Sprint(p.TotalWinPercentage())
	return message
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

func CreateStartMessage(chatId int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "Hi!\nChoose one these options below:")
	msg.ParseMode = tgbotapi.ModeHTML
	aboutButton := tgbotapi.NewKeyboardButton(AboutMessage)
	loginButton := tgbotapi.NewKeyboardButton(LogInMessage)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(loginButton, aboutButton),
	)
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard
	return msg
}

func CreateAboutUsMessage(chatId int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, "blah blah blah blah\nblah blah blah blah\n@pooria2 is creator :D")
	msg.ParseMode = tgbotapi.ModeHTML
	backButton := tgbotapi.NewKeyboardButton(BackMessage)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(backButton),
	)
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard
	return msg
}

func CreateLoginMessage(chatId int64) tgbotapi.MessageConfig {
	fmt.Println("tags:", members.GetTags(chatId))
	if len(members.GetTags(chatId)) == 0 {
		msg := tgbotapi.NewMessage(
			chatId,
			"Send your account tag.\nYou can copy your tag from your clash royale profile, under your username",
		)
		msg.ParseMode = tgbotapi.ModeHTML
		return msg
	}
	var buttons [][]tgbotapi.KeyboardButton
	for _, tag := range members.GetTags(chatId) {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(tag)))
	}
	var msg = tgbotapi.NewMessage(chatId,
		"Choose a tag from list below or send a new tag:",
	)
	msg.ParseMode = tgbotapi.ModeHTML
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	keyboard.OneTimeKeyboard = true
	msg.ReplyMarkup = keyboard
	return msg
}

func CreateHomePageMessage(chatId int64) tgbotapi.MessageConfig {
	tag := members.GetCurrentTag(chatId)
	c := NewClient(auth)
	p, _ := c.Player(tag)
	textMessage := "Dear " + members.GetUser(chatId).FirstName + "!\n" + "You are logged in as " + p.Name
	textMessage += ". Choose one of options below to see result;"
	textMessage += "\nAnd don't forget!! If you like this bot please share "
	textMessage += "for your clan mates and send me your suggestions and feedback with /about_us command"
	msg := tgbotapi.NewMessage(chatId, textMessage)
	newButton1 := tgbotapi.NewInlineKeyboardButtonData("Profile info", ProfileInfoQuery)
	newButton2 := tgbotapi.NewInlineKeyboardButtonData("Game data stats", PlayerStatsQuery)
	newButton3 := tgbotapi.NewInlineKeyboardButtonData("Change Account", ChangeAccountQuery)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(newButton1, newButton2),
		tgbotapi.NewInlineKeyboardRow(newButton3),
	)
	msg.ReplyMarkup = keyboard
	return msg
}
