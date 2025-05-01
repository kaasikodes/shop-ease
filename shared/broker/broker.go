package broker

type MessageBroker interface {
	Publish(topic string, messsage []byte) error
	Subscribe(topic string, handler func(msg []byte)) error
	Close() error
}
