package main

import (
	"clash-royale-telegram-bot/data"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	clanURL   = "https://api.clashofclans.com/v1/clans/%s"
	playerURL = "https://api.clashofclans.com/v1/players/%s"
)

// Client that fetches model objects from Clash of Clan's API
type Client interface {
	Clan(tag string) (*data.Clan, error)
	Player(tag string) (*data.Player, error)
	//FetchAllPlayers(clan *data.Clan) error
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
	if err := t.unmarshalURL(fmt.Sprintf(clanURL, url.PathEscape(tag)), &clan); err != nil {
		return nil, err
	}
	return clan, nil
}

func (t *client) Player(tag string) (*data.Player, error) {
	player := &data.Player{}
	player.Tag = tag
	if err := t.hydratePlayer(player); err != nil {
		return nil, err
	}
	return player, nil
}

/*func (t *client) FetchAllPlayers(clan *data.Clan) error {
	var wg sync.WaitGroup
	for _, member := range clan.MemberList {
		wg.Add(1)
		go func(player *data.Player) {
			if err := t.hydratePlayer(player); err != nil {
				fmt.Println(err) // What do with errors? maybe error channel...
			}
			wg.Done()
		}(member)
	}
	wg.Wait()

	return nil
}*/

func (t *client) hydratePlayer(player *data.Player) error {
	if err := t.unmarshalURL(fmt.Sprintf(playerURL, url.PathEscape(player.Tag)), player); err != nil {
		return err
	}

	return nil
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
	fmt.Println("request :", req.Header)
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
