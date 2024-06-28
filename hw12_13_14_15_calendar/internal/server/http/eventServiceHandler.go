package internalhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

// EventServiceHandler модель обработчика запросов.
type EventServiceHandler struct {
	logger server.Logger
	app    *app.App
}

// Add обработчик.
func (h *EventServiceHandler) Add(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("чтение тела запроса: " + err.Error())
		return
	}
	defer r.Body.Close()

	var event storage.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		h.logger.Error("анмаршалинг тела запроса: " + err.Error())
		return
	}

	err = h.app.RegisterNewEvent(event)
	if err != nil {
		h.logger.Error("вызов метода создания события: " + err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(event.ID))
}

// Update обработчик.
func (h *EventServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("чтение тела запроса: " + err.Error())
		return
	}
	defer r.Body.Close()

	var event storage.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		h.logger.Error("анмаршалинг тела запроса: " + err.Error())
		return
	}

	h.logger.Log(event.Title)

	err = h.app.ChangeEvent(event)
	if err != nil {
		h.logger.Error("вызов метода обновления события: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
}

// Delete обработчик.
func (h *EventServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	err := h.app.DeleteEvent(storage.EventID(eventID))
	if err != nil {
		h.logger.Error("вызов метода удаления события: " + err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

// GetAll обработчик.
func (h *EventServiceHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	events, err := h.app.GetAllEvents()
	if err != nil {
		h.logger.Error("вызов метода получения событий: " + err.Error())
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(events)
	if err != nil {
		h.logger.Error("маршалинг событий в тело ответа: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}

// GetTodayEvents обработчик.
func (h *EventServiceHandler) GetTodayEvents(w http.ResponseWriter, r *http.Request) {
	h.getEventsByRange(w, r, h.app.GetEventsForDay)
}

// GetWeekEvents обработчик.
func (h *EventServiceHandler) GetWeekEvents(w http.ResponseWriter, r *http.Request) {
	h.getEventsByRange(w, r, h.app.GetEventsForWeek)
}

// GetMonthEvents обработчик.
func (h *EventServiceHandler) GetMonthEvents(w http.ResponseWriter, r *http.Request) {
	h.getEventsByRange(w, r, h.app.GetEventsForMonth)
}

// getEventsByRange получить список событий за период.
func (h *EventServiceHandler) getEventsByRange(
	w http.ResponseWriter, r *http.Request, method func(date time.Time) (storage.Events, error),
) {
	timestamp := r.URL.Query().Get("from")
	if timestamp == "" {
		h.logger.Error("дата не передана")
		w.WriteHeader(400)
		return
	}

	date, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		h.logger.Error("парсинг даты: " + err.Error())
		w.WriteHeader(400)
		return
	}

	events, err := method(date)
	if err != nil {
		h.logger.Error("вызов метода получения событий за период: " + err.Error())
		w.WriteHeader(500)
		return
	}

	response, err := json.Marshal(events)
	if err != nil {
		h.logger.Error("маршалинг событий в тело ответа: " + err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(response)
}
