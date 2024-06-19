package memorystorage

import (
	"sync"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

// Storage модель хранилища.
type Storage struct {
	mu     sync.RWMutex
	events storage.Events
}

// Storage конструктор хранилища.
func New() *Storage {
	return &Storage{events: make(map[storage.EventID]storage.Event)}
}

// Add добавить событие.
func (s *Storage) Add(event storage.Event) error {
	if _, ok := s.events[event.ID]; ok {
		return storage.ErrEventIDIsAlreadyExists
	}

	if !s.eventDateIsFree(event.AuthorID, event.StartDate, event.EndDate) {
		return storage.ErrDateBusy
	}

	s.events[event.ID] = event
	return nil
}

// Update обновить событие.
func (s *Storage) Update(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrNoEvent
	}

	if !s.eventDateIsFree(event.AuthorID, event.StartDate, event.EndDate) {
		return storage.ErrDateBusy
	}
	s.events[event.ID] = event
	return nil
}

// Delete удалить событие.
func (s *Storage) Delete(eventID storage.EventID) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	delete(s.events, eventID)

	return nil
}

// GetAll получить все события.
func (s *Storage) GetAll() ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0, len(s.events))
	for _, v := range s.events {
		events = append(events, v)
	}

	return events, nil
}

// Get получить событие.
func (s *Storage) Get(eventID storage.EventID) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[eventID]
	if !ok {
		return event, storage.ErrNoEvent
	}

	return event, nil
}

// eventDateIsFree свободна ли дата для нового события.
func (s *Storage) eventDateIsFree(authorID string, startDate time.Time, endDate time.Time) bool {
	authorEvents := s.GetEventsByAuthor(authorID)
	for _, event := range authorEvents {
		if startDate.Equal(event.StartDate) || endDate.Equal(event.EndDate) {
			return false
		}

		if startDate.After(event.StartDate) && startDate.Before(event.EndDate) {
			return false
		}

		if startDate.Before(event.StartDate) && endDate.After(event.StartDate) {
			return false
		}
	}
	return true
}

// GetEventsByAuthor получить список событий пользователя.
func (s *Storage) GetEventsByAuthor(authorID string) storage.Events {
	authorEvents := make(storage.Events)
	for eventID, event := range s.events {
		if event.AuthorID == authorID {
			authorEvents[eventID] = event
		}
	}
	return authorEvents
}
