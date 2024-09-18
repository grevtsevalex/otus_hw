package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/notify"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

// scheduler модель планировщика.
type scheduler struct {
	logger       Logger
	queue        queue.Queue
	storage      storage.EventStorage
	deletePeriod time.Duration
}

// Scheduler тип планировщика.
type Scheduler interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

// NewScheduler конструктор.
func NewScheduler(storage storage.EventStorage, l Logger, q queue.Queue, period time.Duration) Scheduler {
	return &scheduler{storage: storage, logger: l, queue: q, deletePeriod: period}
}

// Start запустить планировщик.
func (s *scheduler) Start(ctx context.Context) {
	go s.startOldEventsCrawler(ctx)
L:
	for {
		select {
		case <-ctx.Done():
			s.logger.Log("stopping scheduler by context")
			break L

		default:
			events, err := s.storage.GetAll() // лучше метод который сразу выбирает нужные события
			if err != nil {
				s.logger.Error(fmt.Sprintf("selecting all events: %s", err.Error()))
				break L
			}

			for _, event := range events {
				if event.HoursBeforeToNotify == 0 {
					continue
				}

				notifyDate := event.StartDate.Add(time.Hour * time.Duration(-event.HoursBeforeToNotify))

				event.HoursBeforeToNotify = 0
				err = s.storage.Update(event)
				if err != nil {
					s.logger.Error(fmt.Sprintf("set 0 to HoursBeforeToNotify: %s", err.Error()))
				}

				s.queue.Send(notify.Notify{
					ID:          "id",
					Title:       "notify message",
					Date:        notifyDate,
					RecipientID: event.AuthorID,
				})
			}
		}
	}
}

func (s *scheduler) startOldEventsCrawler(ctx context.Context) {
	s.logger.Log("starting old events crawler")
L:
	for {
		select {
		case <-ctx.Done():
			s.logger.Log("stopping remove crawler by context")
			break L

		default:
			events, err := s.storage.GetAll() // лучше метод который сразу выбирает нужные события
			if err != nil {
				s.logger.Error(fmt.Sprintf("selecting all events: %s", err.Error()))
				break L
			}

			for _, event := range events {
				lowBorder := time.Now().AddDate(-1, 0, 0)
				if event.EndDate.After(lowBorder) {
					continue
				}

				err := s.storage.Delete(event.ID)
				if err != nil {
					s.logger.Error(fmt.Sprintf("delete old event: %s eventId: %s", err.Error(), event.ID))
				}
			}
		}

		time.Sleep(s.deletePeriod)
	}
}

// Stop остановить планировщик.
func (s *scheduler) Stop(ctx context.Context) {
	<-ctx.Done()
	s.logger.Log("Stopping scheduler...")
}
