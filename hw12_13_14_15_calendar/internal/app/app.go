package app

import (
	"fmt"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

// App модель приложения.
type App struct {
	Logger  Logger
	storage Storage
}

// Logger тип логгера приложения.
type Logger interface {
	Log(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

// Storage тип хранилища приложения.
type Storage interface {
	storage.EventStorage
}

// New конструктор приложения.
func New(logger Logger, storage Storage) *App {
	return &App{Logger: logger, storage: storage}
}

// RegisterNewEvent зарегистрировать новое событие.
func (a App) RegisterNewEvent(event storage.Event) error {
	err := a.storage.Add(event)
	if err != nil {
		return fmt.Errorf("создание нового события: %w", err)
	}

	return nil
}

// DeleteEvent удалить событие.
func (a App) DeleteEvent(eventID storage.EventID) error {
	err := a.storage.Delete(eventID)
	if err != nil {
		return fmt.Errorf("удаление события: %w", err)
	}

	return nil
}

// ChangeEvent изменить событие.
func (a App) ChangeEvent(event storage.Event) error {
	err := a.storage.Update(event)
	if err != nil {
		return fmt.Errorf("изменение события: %w", err)
	}

	return nil
}

// GetAllEvents получить список всех событий.
func (a App) GetAllEvents() ([]storage.Event, error) {
	events, err := a.storage.GetAll()
	if err != nil {
		return nil, fmt.Errorf("получение всех событий: %w", err)
	}

	return events, nil
}

// GetEventByID получить событие по идентификатору.
func (a App) GetEventByID(eventID storage.EventID) (storage.Event, error) {
	event, err := a.storage.Get(eventID)
	if err != nil {
		return storage.Event{}, fmt.Errorf("получение события по идентификатору: %w", err)
	}

	return event, nil
}

// GetEventsForDay получить события за день.
func (a App) GetEventsForDay(date time.Time) (storage.Events, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location())

	var events storage.Events
	events, err := a.storage.GetEventsByDateRange(startOfDay, endOfDay)
	if err != nil {
		return events, fmt.Errorf("получение событий за день: %w", err)
	}

	return events, nil
}

// GetEventsForWeek получить события за неделю.
func (a App) GetEventsForWeek(date time.Time) (storage.Events, error) {
	daysToReset := (date.Weekday() - 1)
	if date.Weekday() == time.Sunday {
		daysToReset = 6
	}

	startOfWeek := date.AddDate(0, 0, -int(daysToReset))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var events storage.Events
	events, err := a.storage.GetEventsByDateRange(startOfWeek, endOfWeek)
	if err != nil {
		return events, fmt.Errorf("получение событий за неделю: %w", err)
	}

	return events, nil
}

// GetEventsForMonth получить события за месяц.
func (a App) GetEventsForMonth(date time.Time) (storage.Events, error) {
	firstDayOfMonth := time.Date(date.Year(), date.Month(), 0, 0, 0, 0, 0, date.Location())
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	var events storage.Events
	events, err := a.storage.GetEventsByDateRange(firstDayOfMonth, lastDayOfMonth)
	if err != nil {
		return events, fmt.Errorf("получение событий за неделю: %w", err)
	}

	return events, nil
}
