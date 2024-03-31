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
	ModifiedAt time.Time
}

type Provider interface {
	GetMatches() ([]Match, error)
}
