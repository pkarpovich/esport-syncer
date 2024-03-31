package events

import (
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/calendar"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"time"
)

func MatchToCalendarEvent(match providers.Match) calendar.Event {
	duration := time.Duration(match.BestOf) * time.Hour

	return calendar.Event{
		Id:          match.Id,
		Summary:     fmt.Sprintf("%s vs %s", match.Team1, match.Team2),
		Description: fmt.Sprintf("Tournament: %s | Result: %s", match.Tournament, match.Score),
		Location:    match.URL,
		StartAt:     match.Time,
		EndAt:       match.Time.Add(duration),
	}
}
