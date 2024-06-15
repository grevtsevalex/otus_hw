package internalhttp

import (
	"encoding/json"
	"io"
	"net/http"

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

	h.logger.Log(event.Title)

	err = h.app.RegisterNewEvent(event)
	if err != nil {
		h.logger.Error("вызов метода создания события: " + err.Error())
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
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
	w.Write(response)
	w.WriteHeader(200)
}
