/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package broker

import "context"

type Broker interface {
	Init(...Option) error
	Options() Options
	Connect(context.Context) error
	Close(context.CancelFunc) error
	Publish(context.Context, *Message) error
	Subscribe(context.Context, ...Processor) error
	String() string
}

// Handler is used to process messages via a subscription of a topic.
type Handler func(*Message) error

type Processor struct {
	Topic   string
	Handler Handler
	Retry   int
}

type ErrorHandler func(*Message, error)

type Message struct {
	Topic     string
	Key       string
	Partition int32
	Body      []byte
}

func NewMessage(topic string, content []byte, opt ...MessageOption) *Message {
	msg := &Message{Key: "", Partition: -1}
	for _, o := range opt {
		o(msg)
	}
	msg.Topic = topic
	msg.Body = content
	return msg
}

type MessageOption func(*Message)

func MessageKey(key string) MessageOption {
	return func(msg *Message) {
		msg.Key = key
	}
}

func MessagePartition(partition int32) MessageOption {
	return func(msg *Message) {
		msg.Partition = partition
	}
}

func NewBroker(ctx context.Context, opt ...Option) (Broker, error) {
	return NewBrokerWithName(ctx, "default", opt...)
}

func NewBrokerWithName(ctx context.Context, name string, opt ...Option) (Broker, error) {
	ep := NewKafkaBroker(ctx, name, opt...)
	err := ep.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return ep, nil
}
