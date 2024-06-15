package storage

import "time"

type EventID string

type Event struct {
	ID                  EventID   `json:"id"`
	Title               string    `json:"title"`
	StartDate           time.Time `json:"startDate"`
	EndDate             time.Time `json:"endDate"`
	Description         string    `json:"description"`
	AuthorID            string    `json:"authorId"`
	HoursBeforeToNotify int       `json:"hoursBeforeToNotify"`
}

type Events map[EventID]Event
