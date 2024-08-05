package rabbitqueue

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Logger тип логгера.
type Logger interface {
	Log(msg string)
	Error(msg string)
}

// RabbitQueue модель очереди.
type RabbitQueue struct {
	q   amqp.Queue
	ch  *amqp.Channel
	log Logger
}

// NewQueue создание очереди.
func NewQueue(address string, log Logger) *RabbitQueue {
	conn, err := amqp.Dial(address)
	failOnError(err, "Failed to connect to RabbitMQ", log)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel", log)

	q, err := ch.QueueDeclare(
		"events_to_notify", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare a queue", log)

	return &RabbitQueue{q: q, ch: ch, log: log}
}

// Send отправить сообщение в очередь.
func (r *RabbitQueue) Send() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err := r.ch.PublishWithContext(ctx,
		"",       // exchange
		r.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message", r.log)
	r.log.Log(fmt.Sprintf(" [x] Sent %s\n", body))
}

// Receive получение сообщений из очереди.
func (r *RabbitQueue) Receive() {
	msgs, err := r.ch.Consume(
		r.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	failOnError(err, "Failed to register a consumer", r.log)

	var forever chan struct{}

	go func() {
		for d := range msgs {
			r.log.Log(fmt.Sprintf("Received a message: %s", d.Body))
		}
	}()

	r.log.Log(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// failOnError запись ошибки в логи.
func failOnError(err error, msg string, log Logger) {
	if err != nil {
		log.Error(fmt.Sprintf("%s: %s", msg, err))
	}
}
