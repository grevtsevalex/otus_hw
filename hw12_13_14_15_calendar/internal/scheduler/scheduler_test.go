package scheduler

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
	t.Run("basic select", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		eventFromStorage, err := st.Get(event.ID)
		require.NoError(t, err)
		require.Equal(t, event, eventFromStorage)

		// up queue
		// запустить scheduler
		// check in message in queue
		// check event in storage with zero hours
		// stop scheduler
		//

	})

	t.Run("basic remove", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		eventFromStorage, err := st.Get(event.ID)
		require.NoError(t, err)
		require.Equal(t, event, eventFromStorage)
	})
}

// generateEvent сгенерировать событие для теста.
func generateEvent() (storage.Event, error) {
	eventID, err := uuid.DefaultGenerator.NewV4()
	if err != nil {
		return storage.Event{}, fmt.Errorf("generate event: %w", err)
	}
	authorID, err := uuid.DefaultGenerator.NewV4()
	startDate := time.Now()
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
