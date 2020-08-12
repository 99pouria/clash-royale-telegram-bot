package data

func (b *BattleLog) TrophyList() []int {
	var a []int
	for _, battle := range *b {
		if battle.Type == "PvP" {
			a = append(a, battle.Team[0].StartingTrophies, battle.Team[0].TrophyChange)
		}
	}
	return a
}

func (b *BattleLog) RecentWinPercentage() float64 {
	wins := 0.0
	losses := 0.0
	draws := 0.0
	for _, battle := range *b {
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
