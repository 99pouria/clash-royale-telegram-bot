package data

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Start               = 0
	LogIn               = 1
	AboutUs             = 2
	HomePage            = 4
	ViewingProfile      = 5
	ViewingAccountStats = 6
)

type MyBotMembers map[int64]*Member

type Member struct {
	State         int
	Tags          []string
	CurrentTag    string
	ChatId        int64
	LastMessageID int
	*tgbotapi.User
}

func (bm MyBotMembers) Init() MyBotMembers {
	return make(map[int64]*Member)
}

func (bm MyBotMembers) AddNewMember(user *tgbotapi.User, chatId int64) error {

	if _, ok := bm[chatId]; ok {
		return fmt.Errorf("can't initialize new member. user is not a new member")
	}
	newMember := &Member{
		State:  Start,
		Tags:   nil,
		ChatId: chatId,
		User:   user,
	}
	bm[chatId] = newMember
	return nil
}

func (bm MyBotMembers) ChangeState(destState int, chatId int64) error {
	if _, ok := bm[chatId]; !ok {
		return fmt.Errorf("can't change member's state. user not found")
	}
	bm[chatId].State = destState
	return nil
}

func (bm MyBotMembers) ChangeCurrentTag(newTag string, chatId int64) error {
	if bm[chatId].CurrentTag == newTag {
		return fmt.Errorf("can't change tag. new tag is current tag")
	}
	bm[chatId].CurrentTag = newTag

	// if new tag doesn't exist in tags, we add that
	isExist := false
	for _, tag := range bm[chatId].Tags {
		if tag == newTag {
			isExist = true
			break
		}
	}
	if !isExist {
		bm[chatId].Tags = append(bm[chatId].Tags, newTag)
	}
	return nil
}

func (bm MyBotMembers) FindChatIDByUsername(un string) int64 {
	var res int64
	for chatId, member := range bm {
		if member.UserName == un {
			res = chatId
			break
		}
	}
	return res
}

func (bm MyBotMembers) IsMemberNewUser(chatId int64) bool {
	_, ok := bm[chatId]
	return !ok
}

func (bm MyBotMembers) GetState(chatId int64) int {
	return bm[chatId].State
}

func (bm MyBotMembers) GetTags(chatId int64) []string {
	return bm[chatId].Tags
}

func (bm MyBotMembers) GetCurrentTag(chatId int64) string {
	return bm[chatId].CurrentTag
}

func (bm MyBotMembers) GetUser(chatId int64) *tgbotapi.User {
	return bm[chatId].User
}
