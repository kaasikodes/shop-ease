package broker

import (
	"github.com/streadway/amqp"
)

// TODO: Revisit code
type RabbitMQHelper struct {
	conn                                     *amqp.Connection
	channel                                  *amqp.Channel
	RequeueMsgsWithErrorsWhileBeingProcessed bool //this feature is unique to the rabbitmq implementation consider adding to kafka implemention
}

func NewRabbitMQHelper(url string) (*RabbitMQHelper, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQHelper{conn: conn, channel: ch, RequeueMsgsWithErrorsWhileBeingProcessed: false}, nil
}

func (r *RabbitMQHelper) Publish(queue string, message []byte) error {
	_, err := r.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return r.channel.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
}

func (r *RabbitMQHelper) Subscribe(queue string, handler func(msg []byte) error) error {
	msgs, err := r.channel.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			err := handler(d.Body)
			if err != nil {
				d.Ack(true)

			} else {
				d.Nack(false, r.RequeueMsgsWithErrorsWhileBeingProcessed)
			}

		}
	}()
	return nil
}

func (r *RabbitMQHelper) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
