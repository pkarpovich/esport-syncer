package providers

import "time"

type Match struct {
	Id         string
	Tournament string
	Team1      string
	Team2      string
	Score      string
	Time       time.Time
	BestOf     int
	Location   string
	URL        string
	IsLive     bool
	TeamId     int
	GameType   string
	ModifiedAt time.Time
}

type Provider interface {
	GetMatches(teamID int, discipline string) ([]Match, error)
}
