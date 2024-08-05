package queue

// Queue тип очереди.
type Queue interface {
	Send()
	Receive()
}
