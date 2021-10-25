package tbotapi

import (
	"fmt"
	"strconv"
)

func (b *BattleLog) TrophyList() []int {
	var a []int
	for _, battle := range *b {
		if battle.Type == "PvP" {
			a = append(a, battle.Team[0].StartingTrophies, battle.Team[0].TrophyChange)
		}
	}
	return a
}

func (b *BattleLog) RecentWinPercentage(onlyRankMatches bool) float64 {
	wins := 0.0
	losses := 0.0
	draws := 0.0
	for _, battle := range *b {
		if onlyRankMatches && battle.Type != "PvP" {
			continue
		}
		if battle.Team[0].Crowns > battle.Opponent[0].Crowns {
			wins++
			continue
		}
		if battle.Team[0].Crowns < battle.Opponent[0].Crowns {
			losses++
			continue
		}
		draws++
	}
	return wins / (wins + losses + draws) * 100
}

func (p *Player) TotalWinPercentage() (result float64) {
	result = float64(p.Wins) / float64(p.Wins+p.Losses) * 100
	return
}

func (p *Player) GetPlayerProfileInfo() string {
	message := "name: " + p.Name + "\nking level: " + strconv.FormatInt(int64(p.ExpLevel), 10)
	message += "\nwins: " + strconv.FormatInt(int64(p.Wins), 10)
	message += "\nlosses: " + strconv.FormatInt(int64(p.Losses), 10)
	message += "\ntrophies: " + strconv.FormatInt(int64(p.Trophies), 10)
	message += "\nbest trophies: " + strconv.FormatInt(int64(p.BestTrophies), 10)
	return message
}

func (b *BattleLog) GetPlayerStats(totalWinPercentage float64) string {
	message := "\nrecent wins percentage: " + fmt.Sprint(b.RecentWinPercentage(false))
	message += "\nrecent rank games wins percentage: " + fmt.Sprint(b.RecentWinPercentage(true))
	message += "\ntotal wins percentage: " + fmt.Sprint(totalWinPercentage)
	return message
}
