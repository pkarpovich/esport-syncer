package http

import (
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/config"
	"github.com/pkarpovich/esport-syncer/app/events"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"log"
	"net/http"
)

type Client struct {
	Provider providers.Provider
	Events   *events.Repository
	Config   *config.Config
}

func NewClient(cfg *config.Config, events *events.Repository, provider providers.Provider) *Client {
	return &Client{
		Provider: provider,
		Events:   events,
		Config:   cfg,
	}
}

func (c *Client) Listen() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /calendar.ics", c.ServeCalendar)
	mux.HandleFunc("POST /refresh", c.RefreshEvents)
	mux.HandleFunc("GET /health", c.HealthCheck)

	log.Printf("[INFO] Calendar published at http://localhost:%s/calendar.ics\n", c.Config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", c.Config.Port), mux)
}
