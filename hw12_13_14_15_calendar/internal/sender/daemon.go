package sender

import (
	"context"
	"fmt"

	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/notify"
	"github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/queue"
)

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

// Sender тип отправщика уведомлений.
type Sender interface {
	Send(msg notify.Notify)
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

// sender модель отправщика уведомелний.
type sender struct {
	log   Logger
	queue queue.Queue
}

// NewSender конструктор отправщика.
func NewSender(l Logger, q queue.Queue) Sender {
	return &sender{log: l, queue: q}
}

// Send отправить уведомление.
func (s *sender) Send(msg notify.Notify) {
	s.log.Log(fmt.Sprintf("message with id: %s was sended", msg.ID))
}

// Start запуск демона по обработке уведомлений.
func (s *sender) Start(ctx context.Context) {
	notifies := s.queue.Receive()
	select {
	case <-ctx.Done():
		s.log.Log("stopping sender by context")
		break
	case msg := <-notifies:
		s.Send(msg)
	}
}

// Stop остановить планировщик.
func (s *sender) Stop(ctx context.Context) {
	<-ctx.Done()
	s.log.Log("Stopping sender...")
}
