package main

import (
	iCal "github.com/arran4/golang-ical"
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
	calendar *iCal.Calendar
}

func (c *Calendar) CreateCalendar(name string) {
	c.calendar = iCal.NewCalendar()
	c.calendar.SetMethod(iCal.MethodPublish)
	c.calendar.SetName(name)
}

func (c *Calendar) CreateEvent(event CalendarEvent) {
	e := c.calendar.AddEvent(event.Id)

	e.SetCreatedTime(time.Now())
	e.SetStartAt(event.StartAt)
	e.SetEndAt(event.EndAt)
	e.SetSummary(event.Summary)
	e.SetDescription(event.Description)
	e.SetLocation(event.Location)
	e.SetURL(event.URL)
}

func (c *Calendar) GetCalendar() *iCal.Calendar {
	return c.calendar
}

func (c *Calendar) PublishCalendar() string {
	return c.calendar.Serialize()
}
