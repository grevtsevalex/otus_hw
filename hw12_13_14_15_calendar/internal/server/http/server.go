package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server"
)

// Server модель сервера.
type Server struct {
	logger server.Logger
	config Config
	app    *app.App
}

// Config конфиг сервера.
type Config struct {
	Port            int
	ReadTimeoutMS   int
	WriteTimeoutMS  int
	HandlerTimeoutS int
}

// NewServer конструктор сервера.
func NewServer(logger server.Logger, app *app.App, config Config) *Server {
	return &Server{logger: logger, config: config, app: app}
}

// Start старт сервера.
func (s *Server) Start() error {
	handler := &EventServiceHandler{logger: s.logger, app: s.app}
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/event", loggingMiddleware(handler.Add, s.logger)).Methods("POST")
	r.HandleFunc("/event", loggingMiddleware(handler.Update, s.logger)).Methods("PUT")
	r.HandleFunc("/event/{id}", loggingMiddleware(handler.Delete, s.logger)).Methods("DELETE")
	r.HandleFunc("/get-all", loggingMiddleware(handler.GetAll, s.logger)).Methods("GET")
	r.HandleFunc("/get-today-events", loggingMiddleware(handler.GetTodayEvents, s.logger)).Methods("GET")
	r.HandleFunc("/get-week-events", loggingMiddleware(handler.GetWeekEvents, s.logger)).Methods("GET")
	r.HandleFunc("/get-month-events", loggingMiddleware(handler.GetMonthEvents, s.logger)).Methods("GET")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      http.TimeoutHandler(r, time.Duration(s.config.HandlerTimeoutS)*time.Second, "handler timeout"),
		ReadTimeout:  time.Duration(s.config.ReadTimeoutMS) * time.Millisecond,
		WriteTimeout: time.Duration(s.config.WriteTimeoutMS) * time.Millisecond,
	}

	s.logger.Log(fmt.Sprintf("starting http server on %d", s.config.Port))

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("serve http connections: %w", err)
	}

	return nil
}

// Stop стоп сервера.
func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	s.logger.Log("Stopping http server...")
	return nil
}
