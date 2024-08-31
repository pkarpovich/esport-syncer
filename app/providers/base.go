package providers

import match "github.com/pkarpovich/esport-syncer/app/store/matches"

type Provider interface {
	GetMatches(teamID int, discipline string) ([]match.Match, error)
}
