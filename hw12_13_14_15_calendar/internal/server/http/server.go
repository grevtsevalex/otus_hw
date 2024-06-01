package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server модель сервера.
type Server struct {
	logger Logger
	config Config
	app    Application
}

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

type Application interface{}

// Config конфиг сервера.
type Config struct {
	Port            int
	ReadTimeoutMS   int
	WriteTimeoutMS  int
	HandlerTimeoutS int
}

// NewServer конструктор сервера.
func NewServer(logger Logger, app Application, config Config) *Server {
	return &Server{logger: logger, config: config, app: app}
}

// Start старт сервера.
func (s *Server) Start(ctx context.Context) error {
	handler := &Handler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(handler.Hello, s.logger))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      http.TimeoutHandler(mux, time.Duration(s.config.HandlerTimeoutS)*time.Second, "handler timeout"),
		ReadTimeout:  time.Duration(s.config.ReadTimeoutMS) * time.Millisecond,
		WriteTimeout: time.Duration(s.config.WriteTimeoutMS) * time.Millisecond,
	}

	err := server.ListenAndServe()
	if err != nil {
		s.logger.Error(err.Error())
	}
	<-ctx.Done()
	return nil
}

// Stop стоп сервера.
func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	s.logger.Log("Stopping server...")
	return nil
}

// Handler модель обработчика.
type Handler struct{}

// Hello обработчик.
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Println(r.URL)
}
