package internal

type MqttAdaptor interface {
	Publish(topic string, message []byte) bool
}
