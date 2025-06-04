package rabbitmq

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ServicesHeath struct {
	name      string
	updatedAt time.Time
}

func (s *ServicesHeath) getConnectionStatus() bool {

	// check whether service disconnected
	currentAt := time.Now()

	if s.updatedAt.IsZero() || currentAt.Sub(s.updatedAt) > 10*time.Second {
		return false
	}

	return true
}

type ServicesHealthManager struct {
	OngoingHealth *[]ServicesHeath
}

func (shm *ServicesHealthManager) registerService(name string) {
	service := ServicesHeath{
		name:      name,
		updatedAt: time.Now(),
	}

	if shm.OngoingHealth == nil {
		shm.OngoingHealth = &[]ServicesHeath{}
	}

	*shm.OngoingHealth = append(*shm.OngoingHealth, service)
}

func (shm *ServicesHealthManager) updateService(name string) {

	for i, service := range *shm.OngoingHealth {
		if service.name == name {
			(*shm.OngoingHealth)[i].updatedAt = time.Now()
			return
		}
	}
}

func ConsumeHealthMessages() {

	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		panic("Could not connect to rabbitmq server")
	}

	defer conn.Close()

	// define exchange
	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to open a channel")
	}

	defer ch.Close()

	// declare exchange
	err = ch.ExchangeDeclare(
		"health", // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		panic("Failed to declare health exchange")
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic("Failed to declare empty queue")
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"health", // exchange
		false,
		nil,
	)
	if err != nil {
		panic("Failed to bind queue to exchange")
	}

	// Start consuming messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic("Failed to register health consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Process the message
			serviceName := string(d.Body)
			if serviceName == "" {
				continue
			}
		}
	}()

	<-forever
}
