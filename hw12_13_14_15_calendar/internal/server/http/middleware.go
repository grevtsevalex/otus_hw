package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/server"
)

// CustomResponseWriter кастомная модель ответа.
type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter конструктор кастомной модели ответа.
func NewLoggingResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{w, http.StatusOK}
}

// WriteHeader запись кода ответа.
func (c *CustomResponseWriter) WriteHeader(code int) {
	c.statusCode = code
}

// loggingMiddleware логирование входящих запросов.
func loggingMiddleware(next http.HandlerFunc, logger server.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		_, offset := t.Zone()
		formatted := fmt.Sprintf("%02d/%s/%d:%02d:%02d:%02d +%04d",
			t.Day(), t.Month().String(), t.Year(),
			t.Hour(), t.Minute(), t.Second(), offset)

		customResponseWriter := NewLoggingResponseWriter(w)
		next(customResponseWriter, r)

		statusCode := customResponseWriter.statusCode
		str := fmt.Sprintf(`%s [%v] %s %s %s %v %v "%s"`,
			r.RemoteAddr, formatted, r.Method, r.RequestURI, r.Proto, statusCode, time.Since(t).Milliseconds(), r.UserAgent())
		logger.Log(str)
	})
}
