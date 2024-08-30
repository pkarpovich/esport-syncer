package scheduler

import (
	"github.com/go-co-op/gocron"
	"time"
)

type Client struct {
	Scheduler *gocron.Scheduler
}

func NewClient() *Client {
	return &Client{
		Scheduler: gocron.NewScheduler(time.Local),
	}
}

func (c *Client) Start(jobFun interface{}, params ...interface{}) error {
	_, err := c.Scheduler.Cron("0 */8 * * *").StartImmediately().Do(jobFun, params...)
	if err != nil {
		return err
	}
	c.Scheduler.StartAsync()

	return nil
}

func (c *Client) NextRun() (*gocron.Job, time.Time) {
	return c.Scheduler.NextRun()
}
