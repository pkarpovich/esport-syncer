package http

import (
	"fmt"
	"github.com/pkarpovich/esport-syncer/app/config"
	"github.com/pkarpovich/esport-syncer/app/providers"
	match "github.com/pkarpovich/esport-syncer/app/store/matches"
	"github.com/pkarpovich/esport-syncer/app/sync"
	"log"
	"net/http"
)

type Client struct {
	Provider   providers.Provider
	Events     *match.Repository
	Config     *config.Config
	SyncConfig []sync.ConfigItem
}

type ClientOptions struct {
	Provider   providers.Provider
	Events     *match.Repository
	SyncConfig []sync.ConfigItem
	Config     *config.Config
}

func NewClient(opt ClientOptions) *Client {
	return &Client{
		SyncConfig: opt.SyncConfig,
		Provider:   opt.Provider,
		Events:     opt.Events,
		Config:     opt.Config,
	}
}

func (c *Client) Listen() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /events/{id}/calendar.ics", c.ServeCalendar)
	mux.HandleFunc("GET /events/{id}", c.GetEvents)
	mux.HandleFunc("POST /refresh", c.RefreshEvents)
	mux.HandleFunc("GET /health", c.HealthCheck)

	log.Printf("[INFO] Calendars published at http://localhost:%s\n", c.Config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%s", c.Config.Port), mux)
}
