package scheduler

import (
	"context"

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
	logger  Logger
	queue   queue.Queue
	storage storage.EventStorage
}

// Scheduler тип планировщика.
type Scheduler interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

// NewScheduler конструктор.
func NewScheduler(storage storage.EventStorage, l Logger, q queue.Queue) Scheduler {
	return &scheduler{storage: storage, logger: l, queue: q}
}

// Start запустить планировщик.
func (s *scheduler) Start(ctx context.Context) {
L:
	for j := 0; j < 5; j++ {
		select {
		case <-ctx.Done():
			s.logger.Log("stopping scheduler by context")
			break L

		default:
			s.queue.Send()
			// select events where field HoursBeforeToNotify not null
			// loop for selected events
			// // write message for queue and send it to rabbit
			// // set null to HoursBeforeToNotify

			// delete all events that older than 1 year
		}
	}
}

// Stop остановить планировщик.
func (s *scheduler) Stop(ctx context.Context) {
	<-ctx.Done()
	s.logger.Log("Stopping scheduler...")
}
