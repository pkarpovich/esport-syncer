package main

import (
	"fmt"
	"log"
	"time"
)

func (ctx *Context) Sync() {
	log.Printf("[INFO] cron job started at %s", time.Now().Format("2006-01-02 15:04:05"))
	matches, err := ctx.Provider.GetMatches()
	if err != nil {
		log.Printf("[ERROR] error while fetching matches: %v", err)
		return
	}
	log.Printf("[INFO] matches fetched: %d", len(matches))

	for _, match := range matches {
		err := ctx.Events.CreateOrReplace(match)
		if err != nil {
			log.Printf("[ERROR] error while saving match: %v", err)
			continue
		}

		summary := fmt.Sprintf("%s vs %s", match.Team1, match.Team2)
		startAt := match.Time.Format("2006-01-02 15:04:05")
		log.Printf("[INFO] create or replace event: %s at %s", summary, startAt)
	}

	_, nextRun := ctx.Scheduler.NextRun()
	log.Printf("[INFO] cron job finished, next run at %s", nextRun.UTC())
}
