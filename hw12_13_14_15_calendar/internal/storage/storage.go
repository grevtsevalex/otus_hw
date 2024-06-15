package storage

import (
	"errors"
	"time"
)

var (
	ErrDateBusy               = errors.New("this date already used by another event")
	ErrEventIDIsAlreadyExists = errors.New("this eventID already exists")
	ErrNoEvent                = errors.New("no event with this ID")
)

// EventStorage тип - хранилище событий.
type EventStorage interface {
	Add(event Event) error
	Update(event Event) error
	Delete(eventID EventID) error
	GetAll() ([]Event, error)
	Get(eventID EventID) (Event, error)
	GetEventsByDateRange(startDate time.Time, endDate time.Time) (Events, error)
}
