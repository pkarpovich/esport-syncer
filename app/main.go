package main

import (
	"github.com/pkarpovich/esport-syncer/app/config"
	"github.com/pkarpovich/esport-syncer/app/database"
	"github.com/pkarpovich/esport-syncer/app/events"
	"github.com/pkarpovich/esport-syncer/app/http"
	"github.com/pkarpovich/esport-syncer/app/providers"
	"github.com/pkarpovich/esport-syncer/app/scheduler"
	"log"
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

	err = ctx.Scheduler.Start(ctx.Sync)
	if err != nil {
		log.Fatalf("[ERROR] error while starting scheduler: %v", err)
	}

	httpClient := http.NewClient(ctx.Config, ctx.Events)
	err = httpClient.Listen()
	if err != nil {
		log.Fatalf("[ERROR] error while starting http server: %v", err)
	}
}
