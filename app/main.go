package main

import (
	"fmt"
	. "github.com/pkarpovich/esport-syncer/app/providers"
	"log"
	"net/http"
	"time"
)

type Context struct {
	Provider Provider
	Matches  []Match
	Calendar Calendar
}

func (ctx *Context) UpdateMatches() {
	err, matches := ctx.Provider.GetMatches()
	if err != nil {
		log.Fatalf("error while getting matches: %v", err)
	}

	ctx.Matches = matches
}

func (ctx *Context) InitCalendarEvents() {
	ctx.Calendar.CreateCalendar("Esport matches")

	for _, match := range ctx.Matches {
		ctx.Calendar.CreateEvent(matchToCalendarEvent(match))
	}
}

func (ctx *Context) ServeCalendar(w http.ResponseWriter, r *http.Request) {
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
		Location:    "Online",
		URL:         "https://ggscore.com/en/dota-2",
		StartAt:     match.Time,
		EndAt:       match.Time.Add(2 * time.Hour),
	}
}

func main() {
	ctx := Context{
		Provider: &DotaProvider{TeamID: "6224"},
		Calendar: Calendar{},
	}
	ctx.UpdateMatches()
	ctx.InitCalendarEvents()

	http.HandleFunc("/calendar.ics", ctx.ServeCalendar)

	fmt.Printf("Calendar published at http://localhost:1710/calendar.ics\n")
	http.ListenAndServe(":1710", nil)
}
