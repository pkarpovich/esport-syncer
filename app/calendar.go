package main

import (
	iCal "github.com/arran4/golang-ical"
	"log"
	"time"
)

type CalendarEvent struct {
	Id          string
	Summary     string
	Description string
	Location    string
	URL         string
	StartAt     time.Time
	EndAt       time.Time
}

type Calendar struct {
	Name            string
	Color           string
	RefreshInterval string
	calendar        *iCal.Calendar
}

func (c *Calendar) CreateCalendar() {
	log.Printf("create calendar: %s", c.Name)

	c.calendar = iCal.NewCalendar()
	c.calendar.SetMethod(iCal.MethodPublish)
	c.calendar.SetName(c.Name)
	c.calendar.SetColor(c.Color)
	c.calendar.SetRefreshInterval(c.RefreshInterval)

	log.Printf("calendar created")
}

func (c *Calendar) CreateEvent(event CalendarEvent) {
	log.Printf("create event: %s", event.Summary)

	e := c.calendar.AddEvent(event.Id)
	e.SetCreatedTime(time.Now())
	e.SetDtStampTime(time.Now())
	e.SetStartAt(event.StartAt)
	e.SetEndAt(event.EndAt)
	e.SetSummary(event.Summary)
	e.SetDescription(event.Description)
	e.SetLocation(event.Location)

}

func (c *Calendar) GetCalendar() *iCal.Calendar {
	return c.calendar
}

func (c *Calendar) PublishCalendar() string {
	return c.calendar.Serialize()
}
