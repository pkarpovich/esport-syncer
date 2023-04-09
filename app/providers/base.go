package providers

import "time"

type Match struct {
	Tournament string
	Team1      string
	Team2      string
	Score      string
	Time       time.Time
	IsLive     bool
}

type Provider interface {
	GetMatches() (error, []Match)
}
