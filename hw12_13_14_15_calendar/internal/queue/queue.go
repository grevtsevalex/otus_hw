package queue

import "github.com/grevtsevalex/otus_hw/hw12_13_14_15_calendar/internal/notify"

// Queue тип очереди.
type Queue interface {
	Send(notify notify.Notify)
	Receive() <-chan notify.Notify
}
