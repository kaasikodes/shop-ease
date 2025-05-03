package broker

type MessageBroker interface {
	Publish(topic string, messsage []byte) error
	Subscribe(topic string, handler func(msg []byte) error) error
	Close() error
}

// Events Definition of Each Service
