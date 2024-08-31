package providers

import "time"

type Match struct {
	Id         string    `json:"id"`
	Tournament string    `json:"tournament"`
	Team1      string    `json:"team1"`
	Team2      string    `json:"team2"`
	Score      string    `json:"score"`
	Time       time.Time `json:"time"`
	BestOf     int       `json:"bestOf"`
	Location   string    `json:"location"`
	URL        string    `json:"url"`
	IsLive     bool      `json:"isLive"`
	TeamId     int       `json:"targetTeamId"`
	GameType   string    `json:"gameType"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type Provider interface {
	GetMatches(teamID int, discipline string) ([]Match, error)
}
