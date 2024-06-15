package memorystorage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	st := New()

	t.Run("concurenlty get", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		event.Description = "modified"
		event.StartDate = event.StartDate.Add(time.Hour * 48)
		event.EndDate = event.EndDate.Add(time.Hour * 48)
		err = st.Update(event)
		require.NoError(t, err)
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			event, err := st.Get(event.ID)
			require.NoError(t, err)
			require.Equal(t, "modified", event.Description)
		}()

		wg.Wait()
	})

	t.Run("date is busy", func(t *testing.T) {
		event, err := generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)
		require.NoError(t, err)

		event, err = generateEvent()
		require.NoError(t, err)

		err = st.Add(event)
		defer st.Delete(event.ID)

		require.ErrorIs(t, err, storage.ErrDateBusy)

		event, err = generateEvent()
		require.NoError(t, err)

		event.StartDate = event.StartDate.Add(time.Hour * 1)
		event.EndDate = event.StartDate.Add(time.Hour * 6)
		defer st.Delete(event.ID)

		require.ErrorIs(t, st.Add(event), storage.ErrDateBusy)
	})
}

// generateEvent сгенерировать событие для теста.
func generateEvent() (storage.Event, error) {
	eventID, err := uuid.DefaultGenerator.NewV4()
	authorID := "34a8e2ab-e345-4c0c-9e0d-a2af9abc873a"
	startDate := time.Date(2024, 0o5, 28, 15, 0, 0, 0, time.UTC)
	if err != nil {
		return storage.Event{}, fmt.Errorf("generate event: %w", err)
	}
	return storage.Event{
		ID:                  storage.EventID(eventID.String()),
		Title:               "First Title",
		StartDate:           startDate,
		EndDate:             startDate.AddDate(0, 0, 1),
		AuthorID:            authorID,
		HoursBeforeToNotify: 0,
	}, nil
}
