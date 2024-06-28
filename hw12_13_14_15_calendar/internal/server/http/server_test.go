package internalhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type stubLogger struct{}

func (l *stubLogger) Log(_ string)   {}
func (l *stubLogger) Error(_ string) {}
func (l *stubLogger) Warn(_ string)  {}
func (l *stubLogger) Debug(_ string) {}
func (l *stubLogger) Info(_ string)  {}

func TestApi(t *testing.T) {
	calendar := app.New(&stubLogger{}, memorystorage.New())

	handler := &EventServiceHandler{logger: &stubLogger{}, app: calendar}

	t.Run("basic add", func(t *testing.T) {
		bodyString := `{
			"authorId": "sunt irure",
			"description": "et Excepteur",
			"endDate": "2024-07-19T21:54:42.123Z",
			"hoursBeforeToNotify": 3,
			"id": "1",
			"startDate":"2024-07-18T21:54:42.123Z",
			"title": "title"
		}`

		body := strings.NewReader(bodyString)
		req := httptest.NewRequest(http.MethodGet, "/add", body)
		w := httptest.NewRecorder()

		handler.Add(w, req)
		defer calendar.DeleteEvent("1")

		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, http.StatusCreated, res.StatusCode)

		eventID, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		if string(eventID) != "1" {
			t.Errorf("expected id: 1 got %v", string(eventID))
		}

		event, err := calendar.GetEventByID(storage.EventID(eventID))
		require.NoError(t, err)

		eventSerialized, err := json.Marshal(event)
		require.NoError(t, err)

		require.JSONEq(t, bodyString, string(eventSerialized))
	})
}
