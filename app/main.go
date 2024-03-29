package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	. "github.com/pkarpovich/esport-syncer/app/providers"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	TeamSpiritId           = "6224"
	DefaultCalendarName    = "Esport matches"
	DefaultCalendarColor   = "red"
	DefaultRefreshInterval = "P1D"
	DefaultPort            = "1710"
)

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

type Context struct {
	Provider Provider
	Matches  []Match
	Calendar Calendar
}

func (ctx *Context) UpdateMatches() {
	log.Printf("try to fetch matches")

	err, matches := ctx.Provider.GetMatches()
	if err != nil {
		log.Fatalf("error while getting matches: %v", err)
	}

	ctx.Matches = matches
	log.Printf("matches fetched: %d", len(matches))
}

func (ctx *Context) InitCalendarEvents() {
	ctx.Calendar.CreateCalendar()

	for _, match := range ctx.Matches {
		ctx.Calendar.CreateEvent(matchToCalendarEvent(match))
	}
}

func (ctx *Context) ServeCalendar(w http.ResponseWriter, _ *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/calendar")
	headers.Add("Content-Disposition", "attachment; filename=calendar.ics")

	_, err := fmt.Fprintf(w, ctx.Calendar.PublishCalendar())
	if err != nil {
		fmt.Println(err)
	}
}

func matchToCalendarEvent(match Match) CalendarEvent {
	return CalendarEvent{
		Id:          match.Id,
		Summary:     fmt.Sprintf("%s vs %s", match.Team1, match.Team2),
		Description: fmt.Sprintf("Tournament: %s | Result: %s", match.Tournament, match.Score),
		Location:    match.URL,
		StartAt:     match.Time,
		EndAt:       match.Time.Add(2 * time.Hour),
	}
}

func main() {
	s := gocron.NewScheduler(time.UTC)

	ctx := Context{
		Provider: &PandaScoreProvider{
			TeamID: GetEnvOrDefault("PANDASCORE_TEAM_ID", TeamSpiritId),
			ApiKey: GetEnvOrDefault("PANDASCORE_API_KEY", ""),
		},
		Calendar: Calendar{
			Name:            GetEnvOrDefault("CALENDAR_NAME", DefaultCalendarName),
			Color:           GetEnvOrDefault("CALENDAR_COLOR", DefaultCalendarColor),
			RefreshInterval: GetEnvOrDefault("CALENDAR_REFRESH_INTERVAL", DefaultRefreshInterval),
		},
	}

	_, err := s.Cron("0 0 * * *").StartImmediately().Do(func() {
		log.Printf("cron job started")
		ctx.UpdateMatches()
		ctx.InitCalendarEvents()

		_, nextRun := s.NextRun()
		log.Printf("cron job finished, next run at %s", nextRun.UTC())
	})
	if err != nil {
		log.Fatalf("error while creating cron job: %v", err)
	}
	s.StartAsync()

	http.HandleFunc("/calendar.ics", ctx.ServeCalendar)

	port := GetEnvOrDefault("PORT", DefaultPort)

	log.Printf("Calendar published at http://localhost:%s/calendar.ics\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("error while starting server : %v", err)
	}
}
