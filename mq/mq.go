package mq

var Default IMessageQueue

type HandlerFunc func(message []byte) error

type IMessageQueue interface {
	Publish(topic string, message []byte) error                 // Publish 发布事件
	Subscribe(topic, channel string, handler HandlerFunc) error // Subscribe 订阅事件
}
