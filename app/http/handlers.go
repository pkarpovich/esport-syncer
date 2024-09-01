package http

import (
	"encoding/json"
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/calendar"
	match "github.com/pkarpovich/esport-syncer/app/store/matches"
	"github.com/pkarpovich/esport-syncer/app/sync"
	"log"
	"net/http"
)

func (c *Client) ServeCalendar(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/calendar")
	headers.Add("Content-Disposition", "attachment; filename=calendar.ics")
	id := r.PathValue("id")

	syncConfig := FirstOrDefault[sync.ConfigItem](c.SyncConfig, func(syncConfig *sync.ConfigItem) bool {
		return syncConfig.Id == id
	})

	if syncConfig == nil {
		log.Printf("[ERROR] config ID is not found")
		http.Error(w, "config ID is not found", http.StatusNotFound)
		return
	}

	matches, err := c.Events.GetByTeamId(syncConfig.TeamId, syncConfig.GameType)
	if err != nil {
		log.Printf("[ERROR] error while querying events: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	iCal := calendar.Create(c.Config.Calendar.Name, c.Config.Calendar.Color, c.Config.Calendar.RefreshInterval)
	for _, m := range matches {
		event := sync.MatchToCalendarEvent(m)
		calendar.AddEvent(iCal, event)
	}

	_, err = fmt.Fprintf(w, iCal.Serialize())
	if err != nil {
		fmt.Println(err)
	}
}

type ErrorResponse struct {
	code    int
	message string
}

type GetEventsResponse struct {
	Data  []match.Match  `json:"data"`
	Error *ErrorResponse `json:"error"`
}

type GetEventsRequest struct {
	Ids []string `json:"ids"`
}

func (c *Client) GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req GetEventsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("[ERROR] error while decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matches := make([]match.Match, 0)
	for _, id := range req.Ids {
		syncConfig := FirstOrDefault[sync.ConfigItem](c.SyncConfig, func(syncConfig *sync.ConfigItem) bool {
			return syncConfig.Id == id
		})

		if syncConfig == nil {
			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(&GetEventsResponse{
				Error: &ErrorResponse{
					code:    http.StatusNotFound,
					message: "Config ID is not found",
				},
				Data: matches,
			})
			if err != nil {
				log.Printf("[ERROR] error while encoding response: %v", err)
			}
			return
		}

		teamMatches, err := c.Events.GetByTeamId(syncConfig.TeamId, syncConfig.GameType)
		if err != nil {
			log.Printf("[ERROR] error while querying events: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		matches = append(matches, teamMatches...)
	}

	err = json.NewEncoder(w).Encode(&GetEventsResponse{
		Data:  matches,
		Error: nil,
	})
	if err != nil {
		log.Printf("[ERROR] error while encoding response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func FirstOrDefault[T any](slice []T, filter func(*T) bool) (element *T) {

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			return &slice[i]
		}
	}

	return nil
}
