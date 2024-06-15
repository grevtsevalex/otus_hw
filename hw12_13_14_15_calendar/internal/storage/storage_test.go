package storage_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	DBName := os.Getenv("TEST_DB_NAME")
	DBPass := os.Getenv("TEST_DB_PASS")
	DBUser := os.Getenv("TEST_DB_USER")

	var st storage.EventStorage

	if DBName == "" || DBPass == "" || DBUser == "" {
		st = memorystorage.New()
		os.Stdout.WriteString("\n!!!Used memory storage!!!\n")
	} else {
		var err error
		config := sqlstorage.Config{
			DBName: DBName,
			Pass:   DBPass,
			User:   DBUser,
		}

		st, err = sqlstorage.New(config)
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
		os.Stdout.WriteString("\n!!!Used sql storage!!!\n")
	}
	t.Run("basic add", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		eventFromStorage, err := st.Get(event.ID)
		require.NoError(t, err)
		require.Equal(t, event, eventFromStorage)
	})

	t.Run("basic update", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		newTitle := "Second title"
		event.Title = newTitle
		event.StartDate = event.StartDate.Add(time.Hour * 48)
		event.EndDate = event.EndDate.Add(time.Hour * 48)
		err = st.Update(event)
		require.NoError(t, err)

		eventFromStorage, err := st.Get(event.ID)
		require.NoError(t, err)

		require.Equal(t, newTitle, eventFromStorage.Title)
	})

	t.Run("basic delete", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		require.NoError(t, err)

		err = st.Delete(event.ID)
		require.NoError(t, err)

		_, err = st.Get(event.ID)
		require.ErrorIs(t, err, storage.ErrNoEvent)
	})

	t.Run("key already exist add", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		err = st.Add(event)
		require.ErrorIs(t, err, storage.ErrEventIDIsAlreadyExists)
	})

	t.Run("no event on update", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		newEventID, err := uuid.DefaultGenerator.NewV4()
		require.NoError(t, err)

		newEvent := storage.Event{ID: storage.EventID(newEventID.String()), Title: "First Title"}
		err = st.Update(newEvent)
		require.ErrorIs(t, err, storage.ErrNoEvent)
	})

	t.Run("get events by this day", func(t *testing.T) {
		now := time.Now()
		startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

		todayEvent, err := generateEvent()
		require.NoError(t, err)
		todayEvent.StartDate = time.Now()
		todayEvent.EndDate = todayEvent.StartDate.Add(time.Hour * 3)
		err = st.Add(todayEvent)
		defer st.Delete(todayEvent.ID)
		require.NoError(t, err)

		endDayEvent, err := generateEvent()
		require.NoError(t, err)
		endDayEvent.StartDate = time.Now()
		endDayEvent.EndDate = endOfToday
		err = st.Add(endDayEvent)
		defer st.Delete(endDayEvent.ID)
		require.NoError(t, err)

		veryLongEvent, err := generateEvent()
		require.NoError(t, err)
		veryLongEvent.StartDate = time.Now()
		veryLongEvent.EndDate = veryLongEvent.StartDate.AddDate(0, 0, 1)
		err = st.Add(veryLongEvent)
		defer st.Delete(veryLongEvent.ID)
		require.NoError(t, err)

		todayEvents, err := st.GetEventsByDateRange(startOfToday, endOfToday)
		require.NoError(t, err)
		expectedEventsIDs := []storage.EventID{todayEvent.ID, endDayEvent.ID}

		if len(todayEvents) != len(expectedEventsIDs) {
			t.Error("wrong events")
		}

		for _, eventID := range expectedEventsIDs {
			if _, ok := todayEvents[eventID]; !ok {
				t.Error("not all event")
			}
		}
	})
}

// generateEvent сгенерировать событие для теста.
func generateEvent() (storage.Event, error) {
	eventID, err := uuid.DefaultGenerator.NewV4()
	authorID, err := uuid.DefaultGenerator.NewV4()
	startDate := time.Date(2024, 0o5, 28, 15, 0, 0, 0, time.UTC)
	if err != nil {
		return storage.Event{}, fmt.Errorf("generate event: %w", err)
	}
	return storage.Event{
		ID:                  storage.EventID(eventID.String()),
		Title:               "First Title",
		StartDate:           startDate,
		EndDate:             startDate.AddDate(0, 0, 1),
		AuthorID:            authorID.String(),
		HoursBeforeToNotify: 5,
	}, nil
}
