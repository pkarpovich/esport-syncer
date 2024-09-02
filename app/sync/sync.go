package sync

import (
	"encoding/json"
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/calendar"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"github.com/pkarpovich/esport-syncer/app/store/matches"
	"log"
	"os"
	"time"
)

type ConfigItem struct {
	Id       string `json:"id"`
	TeamId   int    `json:"teamId"`
	GameType string `json:"gameType"`
}

func GetSyncConfig(localPath string) ([]ConfigItem, error) {
	bytes, err := os.ReadFile(localPath)
	if err != nil {
		log.Printf("[ERROR] error while reading file: %v", err)
		return nil, err
	}

	var syncConfig []ConfigItem
	err = json.Unmarshal(bytes, &syncConfig)
	if err != nil {
		log.Printf("[ERROR] error while unmarshalling JSON: %v", err)
		return nil, err
	}

	return syncConfig, nil
}

func Start(provider providers.Provider, events *match.Repository, syncConfig []ConfigItem) error {
	for _, item := range syncConfig {
		matches, err := provider.GetMatches(item.TeamId, item.GameType)
		if err != nil {
			log.Printf("[ERROR] error while fetching matches: %v", err)
			return err
		}
		log.Printf("[INFO] Fetched %d matches for game type '%s' and team ID %d", len(matches), item.GameType, item.TeamId)

		for _, m := range matches {
			err := events.CreateOrReplace(m)
			if err != nil {
				log.Printf("[ERROR] error while saving match: %v", err)
				continue
			}

			summary := fmt.Sprintf("%s vs %s", m.Team1.Name, m.Team2.Name)
			startAt := m.Time.Local().Format("2006-01-02 15:04:05")
			log.Printf("[INFO] create or replace event: %s at %s", summary, startAt)
		}
	}

	return nil
}

func MatchToCalendarEvent(match match.Match) calendar.Event {
	duration := time.Duration(match.BestOf) * time.Hour

	return calendar.Event{
		Id:          match.Id,
		Summary:     fmt.Sprintf("%s vs %s", match.Team1.Name, match.Team2.Name),
		Description: fmt.Sprintf("Tournament: %s | Result: %s", match.Tournament, match.Score),
		Location:    match.URL,
		StartAt:     match.Time,
		EndAt:       match.Time.Add(duration),
	}
}
