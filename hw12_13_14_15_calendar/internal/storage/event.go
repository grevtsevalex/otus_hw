package storage

import "time"

type EventID string

type Event struct {
	ID                  EventID
	Title               string
	StartDate           time.Time
	EndDate             time.Time
	Description         string
	AuthorID            string
	HoursBeforeToNotify int
}

type Events map[EventID]Event
