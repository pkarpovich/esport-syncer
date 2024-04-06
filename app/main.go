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

	eventsRepository, err := events.NewRepository(db)
	if err != nil {
		log.Fatalf("[ERROR] error while creating events repository: %v", err)
	}

	ctx := Context{
		Provider: &providers.PandaScoreProvider{
			TeamID: cfg.PandaScore.TeamId,
			ApiKey: cfg.PandaScore.ApiKey,
		},
		Events:    eventsRepository,
		Scheduler: scheduler.NewClient(),
		Config:    cfg,
	}

	err = ctx.Scheduler.Start(ctx.handleSync)
	if err != nil {
		log.Fatalf("[ERROR] error while starting scheduler: %v", err)
	}

	httpClient := http.NewClient(ctx.Config, ctx.Events, ctx.Provider)
	err = httpClient.Listen()
	if err != nil {
		log.Fatalf("[ERROR] error while starting http server: %v", err)
	}
}

func (ctx *Context) handleSync() {
	log.Printf("[INFO] cron job started at %s", time.Now().Format("2006-01-02 15:04:05"))
	err := sync.Start(ctx.Provider, ctx.Events)
	if err != nil {
		log.Printf("[ERROR] error while syncing events: %v", err)
	}

	_, nextRun := ctx.Scheduler.NextRun()
	log.Printf("[INFO] cron job finished, next run at %s", nextRun.Local())
}
