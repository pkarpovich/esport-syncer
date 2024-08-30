package http

import (
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/calendar"
	"github.com/pkarpovich/esport-syncer/app/events"
	"github.com/pkarpovich/esport-syncer/app/sync"
	"log"
	"net/http"
)

func (c *Client) ServeCalendar(w http.ResponseWriter, _ *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/calendar")
	headers.Add("Content-Disposition", "attachment; filename=calendar.ics")

	matches, err := c.Events.GetAll()
	if err != nil {
		log.Printf("[ERROR] error while querying events: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	iCal := calendar.Create(c.Config.Calendar.Name, c.Config.Calendar.Color, c.Config.Calendar.RefreshInterval)
	for _, match := range matches {
		event := events.MatchToCalendarEvent(match)
		calendar.AddEvent(iCal, event)
	}

	_, err = fmt.Fprintf(w, iCal.Serialize())
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "application/json")

	matches, err := c.Events.GetAll()
	if err != nil {
		log.Printf("[ERROR] error while querying events: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	_, err = fmt.Fprintf(w, `{"msg": "OK", "events": %d}`, len(matches))
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) RefreshEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] refresh events request")
	if r.Header.Get("App-Token") != c.Config.Secret {
		log.Printf("[WARN] unauthorized request")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	headers := w.Header()
	headers.Add("Content-Type", "application/json")

	err := sync.Start(c.Provider, c.Events, c.SyncConfig)
	if err != nil {
		log.Printf("[ERROR] error while refreshing events: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(w, `{"msg": "OK"}`)
	if err != nil {
		log.Printf("[ERROR] error while writing response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
