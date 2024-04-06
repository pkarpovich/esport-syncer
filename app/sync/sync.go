package sync

import (
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/events"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"log"
)

func Start(provider providers.Provider, events *events.Repository) error {
	matches, err := provider.GetMatches()
	if err != nil {
		log.Printf("[ERROR] error while fetching matches: %v", err)
		return err
	}
	log.Printf("[INFO] matches fetched: %d", len(matches))

	for _, match := range matches {
		err := events.CreateOrReplace(match)
		if err != nil {
			log.Printf("[ERROR] error while saving match: %v", err)
			continue
		}

		summary := fmt.Sprintf("%s vs %s", match.Team1, match.Team2)
		startAt := match.Time.Format("2006-01-02 15:04:05")
		log.Printf("[INFO] create or replace event: %s at %s", summary, startAt)
	}

	return nil
}
