package calendar

import (
	iCal "github.com/arran4/golang-ical"
	"time"
)

type Event struct {
	Id          string
	Summary     string
	Description string
	Location    string
	URL         string
	StartAt     time.Time
	EndAt       time.Time
}

func Create(name, color, refreshInterval string) *iCal.Calendar {
	calendar := iCal.NewCalendar()
	calendar.SetMethod(iCal.MethodPublish)
	calendar.SetName(name)
	calendar.SetColor(color)
	calendar.SetRefreshInterval(refreshInterval)

	return calendar
}

func AddEvent(calendar *iCal.Calendar, event Event) {
	e := calendar.AddEvent(event.Id)
	e.SetCreatedTime(time.Now())
	e.SetDtStampTime(time.Now())
	e.SetStartAt(event.StartAt)
	e.SetEndAt(event.EndAt)
	e.SetSummary(event.Summary)
	e.SetDescription(event.Description)
	e.SetLocation(event.Location)
}
