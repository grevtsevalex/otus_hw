package storage

import "time"

type EventID string

type Event struct {
	ID                  EventID   `json:"id,omitempty"`
	Title               string    `json:"title,omitempty"`
	StartDate           time.Time `json:"startDate,omitempty"`
	EndDate             time.Time `json:"endDate,omitempty"`
	Description         string    `json:"description,omitempty"`
	AuthorID            string    `json:"authorId,omitempty"`
	HoursBeforeToNotify int       `json:"hoursBeforeToNotify,omitempty"`
}

type Events map[EventID]Event
