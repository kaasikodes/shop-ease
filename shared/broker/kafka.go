package broker

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaHelper struct {
	writer  *kafka.Writer
	readers map[string]*kafka.Reader
	groupID string
	brokers []string
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewKafkaWriter(brokers []string, topic string) *KafkaHelper {
	ctx, cancel := context.WithCancel(context.Background())
	helper := &KafkaHelper{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
		brokers: brokers,
		readers: make(map[string]*kafka.Reader),
		ctx:     ctx,
		cancel:  cancel,
	}
	return helper

}

func (k *KafkaHelper) Subscribe(topic string, handler func(msg []byte)) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.brokers,
		Topic:       topic,
		GroupID:     k.groupID,
		StartOffset: kafka.FirstOffset,
	})

	k.mu.Lock()
	k.readers[topic] = reader
	k.mu.Unlock()

	var errs []error
	go func(errs []error) {
		defer reader.Close()
		for {
			select {
			case <-k.ctx.Done():
				log.Printf("shutting down the consumer/reader for %s topic", topic)
				return
			default:
				msg, err := reader.ReadMessage(context.Background())
				if err != nil {
					log.Printf("error reading message: %v", err)
					errs = append(errs, err)
					continue
				}
				log.Printf("message received: %s", string(msg.Value))
				handler(msg.Value)
			}
		}

	}(errs)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil

}
func (k *KafkaHelper) Publish(topic string, messsage []byte) error {
	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: messsage,
		},
	)
	if err != nil {
		log.Printf("failed to write message: %v", err)
	}
	return err
}

func (k *KafkaHelper) Close() error {
	k.cancel() //Broadcast shutdown to all goroutines
	var errs []error

	k.mu.Lock()
	for topic, reader := range k.readers {
		if err := reader.Close(); err != nil {
			log.Printf("error closing reader for topic - %s: %v", topic, err)
			errs = append(errs, err)
		}
	}
	k.mu.Unlock()

	err := k.writer.Close()
	errs = append(errs, err)

	return errors.Join(errs...)
}
