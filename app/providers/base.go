package providers

import "time"

type Match struct {
	Id         string
	Tournament string
	Team1      string
	Team2      string
	Score      string
	Time       time.Time
	Location   string
	URL        string
	IsLive     bool
}

type Provider interface {
	GetMatches() (error, []Match)
}
