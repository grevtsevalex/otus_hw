package memorystorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

// Storage модель хранилища.
type Storage struct {
	mu     sync.RWMutex
	events map[storage.EventID]storage.Event
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

	if !s.eventDateIsFree(event.StartDate) {
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

	if !s.eventDateIsFree(event.StartDate) {
		return storage.ErrDateBusy
	}
	s.events[event.ID] = event
	return nil
}

// Delete удалить событие.
func (s *Storage) Delete(eventID storage.EventID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, eventID)

	return nil
}

// GetAll получить все события.
func (s *Storage) GetAll() ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]storage.Event, 0, len(s.events))
	for _, v := range s.events {
		events = append(events, v)
	}

	return events, nil
}

// Get получить событие.
func (s *Storage) Get(eventID storage.EventID) (storage.Event, error) {
	event, ok := s.events[eventID]
	if !ok {
		return event, storage.ErrNoEvent
	}

	return event, nil
}

// eventDateIsFree свободна ли дата для нового события.
func (s *Storage) eventDateIsFree(date time.Time) bool {
	fmt.Println(date)
	return true
}
