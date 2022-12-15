/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package broker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var (
	defaultKafkaConfig = sarama.NewConfig()
)

type kBroker struct {
	connected bool
	scMutex   sync.RWMutex

	producer sarama.AsyncProducer
	consumer sarama.ConsumerGroup
	opts     Options

	ctx    context.Context
	cancel context.CancelFunc
}

func (k *kBroker) Connect() error {
	if k.isConnected() {
		return nil
	}

	config := defaultKafkaConfig
	config.Version = sarama.V0_10_2_0
	producer, err := sarama.NewAsyncProducer(k.opts.ProducerAddr, config)
	if err != nil {
		return err
	}

	config.Consumer.Return.Errors = true
	consumerGroup, err := sarama.NewConsumerGroup(k.opts.ConsumerAddr, k.opts.GroupId, config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(k.opts.Context)

	k.scMutex.Lock()
	k.producer = producer
	k.consumer = consumerGroup
	k.connected = true
	k.ctx = ctx
	k.cancel = cancel
	k.scMutex.Unlock()

	go func() {
		var enqueued, errs int
		for {
			select {
			case <-k.ctx.Done():
				return
			case <-k.producer.Successes():
				enqueued++
			case err := <-k.producer.Errors():
				if err != nil && err.Msg != nil {
					topic := err.Msg.Topic
					data, _ := err.Msg.Value.Encode()
					ctxzap.Error(k.ctx, "Failed to produce message", zap.String("topic", topic), zap.String("data", string(data)), zap.Error(err))
				}
				errs++
			case err := <-k.consumer.Errors():
				ctxzap.Error(k.ctx, "consumer error", zap.Error(err))
			}
		}
	}()

	return nil
}

func (k *kBroker) Close() error {
	if !k.isConnected() {
		return nil
	}

	k.scMutex.Lock()
	defer k.scMutex.Unlock()

	_ = k.producer.Close()
	_ = k.consumer.Close()

	k.connected = false
	k.cancel()
	return nil
}

func (k *kBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&k.opts)
	}
	return nil
}

func (k *kBroker) isConnected() bool {
	k.scMutex.RLock()
	defer k.scMutex.RUnlock()
	return k.connected
}

func (k *kBroker) Options() Options {
	return k.opts
}

func (k *kBroker) String() string {
	return "kafka"
}

func (k *kBroker) Publish(msg *Message) error {
	if !k.isConnected() {
		return errors.New("[kafka] broker not connected")
	}

	message := &sarama.ProducerMessage{}
	message.Topic = msg.Topic
	message.Partition = msg.Partition
	message.Timestamp = time.Now()
	message.Value = sarama.ByteEncoder(msg.Body)

	k.producer.Input() <- message

	return nil
}

func (k *kBroker) Subscribe(procs ...Processor) error {

	var ts []string
	ps := make(map[string]Processor, len(procs))

	for _, p := range procs {
		ts = append(ts, p.Topic)
		ps[p.Topic] = p
	}

	handler := &consumerGroupHandler{cxt: k.ctx, processors: ps}

	go func() {
		for {
			select {
			case <-k.ctx.Done():
				return
			default:
				err := k.consumer.Consume(k.ctx, ts, handler)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					continue
				default:
					ctxzap.Error(k.ctx, "Consume error", zap.Error(err))
					return
				}
			}
		}
	}()
	return nil
}

func NewKafkaBroker(name string, opt ...Option) Broker {
	opts := NewOptions(name)

	for _, o := range opt {
		o(&opts)
	}

	return &kBroker{
		opts: opts,
	}
}

type consumerGroupHandler struct {
	cxt        context.Context
	processors map[string]Processor
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		m := &Message{
			Topic:     msg.Topic,
			Key:       string(msg.Key),
			Partition: msg.Partition,
			Body:      msg.Value,
		}

		processor, ok := h.processors[msg.Topic]
		if !ok {
			ctxzap.Warn(h.cxt, "processor not exist", zap.String("topic", msg.Topic))
			continue
		}
		err := processor.Handler(m)
		if err == nil {
			sess.MarkMessage(msg, "")
			continue
		}
		retry := processor.Retry
		if retry <= 0 {
			continue
		}
		for i := 0; i < retry; i++ {
			err := processor.Handler(m)
			if err == nil {
				sess.MarkMessage(msg, "")
				break
			}
			ctxzap.Error(h.cxt, "retry process error", zap.String("topic", msg.Topic), zap.String("content", string(msg.Value)), zap.Int("retry", i), zap.Error(err))

		}
		ctxzap.Error(h.cxt, "process error", zap.String("topic", msg.Topic), zap.String("content", string(msg.Value)), zap.Error(err))
	}
	return nil
}
