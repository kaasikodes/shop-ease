package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokerAddress, topic, groupID string, handler func([]byte)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}
		log.Printf("message received: %s", string(msg.Value))
		handler(msg.Value)
	}
}
