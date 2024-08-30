package main

import (
	"github.com/pkarpovich/esport-syncer/app/config"
	"github.com/pkarpovich/esport-syncer/app/database"
	"github.com/pkarpovich/esport-syncer/app/events"
	"github.com/pkarpovich/esport-syncer/app/http"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"github.com/pkarpovich/esport-syncer/app/scheduler"
	"github.com/pkarpovich/esport-syncer/app/sync"
	"log"
	"time"
)

type Context struct {
	Config    *config.Config
	Provider  providers.Provider
	Events    *events.Repository
	Scheduler *scheduler.Client
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("[ERROR] error while reading config: %v", err)
	}

	db, err := database.NewClient("events.db")
	if err != nil {
		log.Fatalf("[ERROR] error while creating database client: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("[ERROR] error while closing database client: %v", err)
		}
	}()

	syncConfig := []sync.ConfigItem{
		{
			Id:       "4c4d2148-8a6d-42bb-94ea-8725e4b96f26",
			TeamId:   1669,
			GameType: "dota2",
		},
		{
			Id:       "7ccc873c-c0a9-4a72-a48e-e88d6e436e26",
			TeamId:   124523,
			GameType: "csgo",
		},
		{
			Id:       "e1575b45-06eb-4f0d-8a4c-fe0496e57189",
			TeamId:   135142,
			GameType: "mlbb",
		},
	}

	eventsRepository, err := events.NewRepository(db)
	if err != nil {
		log.Fatalf("[ERROR] error while creating events repository: %v", err)
	}

	ctx := Context{
		Provider: &providers.PandaScoreProvider{
			ApiKey: cfg.PandaScore.ApiKey,
		},
		Events:    eventsRepository,
		Scheduler: scheduler.NewClient(),
		Config:    cfg,
	}

	err = ctx.Scheduler.Start(ctx.handleSync, syncConfig)
	if err != nil {
		log.Fatalf("[ERROR] error while starting scheduler: %v", err)
	}

	httpClient := http.NewClient(http.ClientOptions{
		Provider:   ctx.Provider,
		Events:     ctx.Events,
		SyncConfig: syncConfig,
		Config:     cfg,
	})
	err = httpClient.Listen()
	if err != nil {
		log.Fatalf("[ERROR] error while starting http server: %v", err)
	}
}

func (ctx *Context) handleSync(syncConfig []sync.ConfigItem) {
	log.Printf("[INFO] cron job started at %s", time.Now().Format("2006-01-02 15:04:05"))
	err := sync.Start(ctx.Provider, ctx.Events, syncConfig)
	if err != nil {
		log.Printf("[ERROR] error while syncing events: %v", err)
	}

	_, nextRun := ctx.Scheduler.NextRun()
	log.Printf("[INFO] cron job finished, next run at %s", nextRun.Local())
}
