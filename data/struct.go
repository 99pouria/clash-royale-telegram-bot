package data

type Player struct {
	Tag                   string `json:"tag"`
	Name                  string `json:"name"`
	ExpLevel              int    `json:"expLevel"`
	Trophies              int    `json:"trophies"`
	BestTrophies          int    `json:"bestTrophies"`
	Wins                  int    `json:"wins"`
	Losses                int    `json:"losses"`
	BattleCount           int    `json:"battleCount"`
	ThreeCrownWins        int    `json:"threeCrownWins"`
	ChallengeCardsWon     int    `json:"challengeCardsWon"`
	ChallengeMaxWins      int    `json:"challengeMaxWins"`
	TournamentCardsWon    int    `json:"tournamentCardsWon"`
	TournamentBattleCount int    `json:"tournamentBattleCount"`
	Role                  string `json:"role"`
	Donations             int    `json:"donations"`
	DonationsReceived     int    `json:"donationsReceived"`
	TotalDonations        int    `json:"totalDonations"`
	WarDayWins            int    `json:"warDayWins"`
	ClanCardsCollected    int    `json:"clanCardsCollected"`
	Clan                  struct {
		Tag     string `json:"tag"`
		Name    string `json:"name"`
		BadgeID int    `json:"badgeId"`
	} `json:"clan"`
	Arena struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"arena"`
	LeagueStatistics struct {
		CurrentSeason struct {
			Trophies     int `json:"trophies"`
			BestTrophies int `json:"bestTrophies"`
		} `json:"currentSeason"`
		PreviousSeason struct {
			ID           string `json:"id"`
			Trophies     int    `json:"trophies"`
			BestTrophies int    `json:"bestTrophies"`
		} `json:"previousSeason"`
		BestSeason struct {
			ID       string `json:"id"`
			Trophies int    `json:"trophies"`
		} `json:"bestSeason"`
	} `json:"leagueStatistics"`
	Badges []struct {
		Name     string `json:"name"`
		Progress int    `json:"progress"`
	} `json:"badges"`
	Achievements []struct {
		Name           string      `json:"name"`
		Stars          int         `json:"stars"`
		Value          int         `json:"value"`
		Target         int         `json:"target"`
		Info           string      `json:"info"`
		CompletionInfo interface{} `json:"completionInfo"`
	} `json:"achievements"`
	Cards []struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		Level    int    `json:"level"`
		MaxLevel int    `json:"maxLevel"`
		Count    int    `json:"count"`
		IconUrls struct {
			Medium string `json:"medium"`
		} `json:"iconUrls"`
	} `json:"cards"`
	CurrentDeck []struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		Level    int    `json:"level"`
		MaxLevel int    `json:"maxLevel"`
		Count    int    `json:"count"`
		IconUrls struct {
			Medium string `json:"medium"`
		} `json:"iconUrls"`
	} `json:"currentDeck"`
	CurrentFavouriteCard struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		MaxLevel int    `json:"maxLevel"`
		IconUrls struct {
			Medium string `json:"medium"`
		} `json:"iconUrls"`
	} `json:"currentFavouriteCard"`
}

type Cards struct {
	Items []struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		MaxLevel int    `json:"maxLevel"`
		IconUrls struct {
			Medium string `json:"medium"`
		} `json:"iconUrls"`
	} `json:"items"`
}

type Clan struct {
	MemberList []struct {
		LastSeen        string `json:"lastSeen"`
		ClanChestPoints int    `json:"clanChestPoints"`
		Arena           struct {
			Name struct {
			} `json:"name"`
			ID       int `json:"id"`
			IconUrls struct {
			} `json:"iconUrls"`
		} `json:"arena"`
		Tag               string `json:"tag"`
		Name              string `json:"name"`
		Role              string `json:"role"`
		ExpLevel          int    `json:"expLevel"`
		Trophies          int    `json:"trophies"`
		ClanRank          int    `json:"clanRank"`
		PreviousClanRank  int    `json:"previousClanRank"`
		Donations         int    `json:"donations"`
		DonationsReceived int    `json:"donationsReceived"`
	} `json:"memberList"`
	Tag               string `json:"tag"`
	DonationsPerWeek  int    `json:"donationsPerWeek"`
	ClanChestStatus   string `json:"clanChestStatus"`
	ClanChestLevel    int    `json:"clanChestLevel"`
	BadgeID           int    `json:"badgeId"`
	RequiredTrophies  int    `json:"requiredTrophies"`
	ClanScore         int    `json:"clanScore"`
	ClanChestMaxLevel int    `json:"clanChestMaxLevel"`
	ClanWarTrophies   int    `json:"clanWarTrophies"`
	Name              string `json:"name"`
	Location          struct {
		LocalizedName string `json:"localizedName"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		IsCountry     bool   `json:"isCountry"`
		CountryCode   string `json:"countryCode"`
	} `json:"location"`
	Type            string `json:"type"`
	Members         int    `json:"members"`
	Description     string `json:"description"`
	ClanChestPoints int    `json:"clanChestPoints"`
	BadgeUrls       struct {
	} `json:"badgeUrls"`
}
