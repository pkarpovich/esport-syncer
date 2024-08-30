package sync

import (
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/events"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"log"
)

type ConfigItem struct {
	Id       string `json:"id"`
	TeamId   int    `json:"teamId"`
	GameType string `json:"gameType"`
}

func Start(provider providers.Provider, events *events.Repository, syncConfig []ConfigItem) error {
	for _, item := range syncConfig {
		matches, err := provider.GetMatches(item.TeamId, item.GameType)
		if err != nil {
			log.Printf("[ERROR] error while fetching matches: %v", err)
			return err
		}
		log.Printf("[INFO] Fetched %d matches for game type '%s' and team ID %d", len(matches), item.GameType, item.TeamId)

		for _, match := range matches {
			err := events.CreateOrReplace(match)
			if err != nil {
				log.Printf("[ERROR] error while saving match: %v", err)
				continue
			}

			summary := fmt.Sprintf("%s vs %s", match.Team1, match.Team2)
			startAt := match.Time.Local().Format("2006-01-02 15:04:05")
			log.Printf("[INFO] create or replace event: %s at %s", summary, startAt)
		}
	}

	return nil
}
