package main

import (
	"encoding/json"
	"fmt"
	"github.com/pooria1/clash-royale-telegram-bot/data"
	"net/http"
	"net/url"
)

const (
	clanURL   = "https://api.clashroyale.com/v1/clans/%s"
	playerURL = "https://api.clashroyale.com/v1/players/%s"
)

// Client that fetches model objects from Clash of Clan's API
type Client interface {
	Clan(tag string) (*data.Clan, error)
	Player(tag string) (*data.Player, error)
	BattleLog(tag string) (*data.BattleLog, error)
}

type client struct {
	token string
}

// NewClient creates a new instance of API.
func NewClient(token string) Client {
	return &client{token}
}

func (t *client) Clan(tag string) (*data.Clan, error) {
	clan := &data.Clan{}
	err := t.unmarshalURL(fmt.Sprintf(clanURL, url.PathEscape(tag)), &clan)
	if err != nil {
		return nil, err
	}
	return clan, nil
}

func (t *client) Player(tag string) (*data.Player, error) {
	player := &data.Player{}
	err := t.unmarshalURL(fmt.Sprintf(playerURL, url.PathEscape(tag)), player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (t *client) BattleLog(tag string) (*data.BattleLog, error) {
	battleLog := &data.BattleLog{}
	err := t.unmarshalURL(fmt.Sprintf(playerURL, url.PathEscape(tag))+"/battlelog", battleLog)
	if err != nil {
		return nil, err
	}
	return battleLog, nil
}

func (t *client) unmarshalURL(url string, v interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", t.token))
	r, err := client.Do(req)
	fmt.Println("response status:", r.Status)
	if err != nil {
		return err
	} else if r.StatusCode != 200 {
		return fmt.Errorf("none 200 response from [%s]", url)
	}

	if err = json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
