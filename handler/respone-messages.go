package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

func CreateHomePageMessage(chatId int64, messageId int) tgbotapi.Chattable {
	tag := members.GetCurrentTag(chatId)
	p, _ := royaleClient.Player(tag)
	textMessage := "Dear " + members.GetUser(chatId).FirstName + "!\n" + "You are logged in as " + p.Name
	textMessage += ". Choose one of options below to see result;"
	textMessage += "\nAnd don't forget!! If you like this bot please share "
	textMessage += "for your clan mates and send me your suggestions and feedback with /about_us command"
	newButton1 := tgbotapi.NewInlineKeyboardButtonData("Profile info", ProfileInfoQuery)
	newButton2 := tgbotapi.NewInlineKeyboardButtonData("Game data stats", PlayerStatsQuery)
	newButton3 := tgbotapi.NewInlineKeyboardButtonData("Change Account", ChangeAccountQuery)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(newButton1, newButton2),
		tgbotapi.NewInlineKeyboardRow(newButton3),
	)
	if messageId != -1 {
		msg := tgbotapi.NewEditMessageText(chatId, messageId, textMessage)
		msg.ReplyMarkup = &keyboard
		return msg
	}
	msg := tgbotapi.NewMessage(chatId, textMessage)
	msg.ReplyMarkup = keyboard
	return msg
}

func CreateProfileInfoMessage(chatId int64, messageID int) tgbotapi.EditMessageTextConfig {
	newButton := tgbotapi.NewInlineKeyboardButtonData("Back", BackMessage)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(newButton),
	)
	p, _ := royaleClient.Player(members.GetCurrentTag(chatId))
	msg := tgbotapi.NewEditMessageText(chatId, messageID, p.GetPlayerProfileInfo())
	msg.ReplyMarkup = &keyboard
	return msg
}

func CreateProfileStatsMessage(chatId int64, messageID int) tgbotapi.EditMessageTextConfig {
	newButton := tgbotapi.NewInlineKeyboardButtonData("Back", BackMessage)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(newButton),
	)
	p, _ := royaleClient.Player(members.GetCurrentTag(chatId))
	b, _ := royaleClient.BattleLog(members.GetCurrentTag(chatId))
	msg := tgbotapi.NewEditMessageText(chatId, messageID, b.GetPlayerStats(p.TotalWinPercentage()))
	msg.ReplyMarkup = &keyboard
	return msg
}
