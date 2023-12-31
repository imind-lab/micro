/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package broker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog"
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
}

func (k *kBroker) Connect(ctx context.Context) error {
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

	k.scMutex.Lock()
	k.producer = producer
	k.consumer = consumerGroup
	k.connected = true
	k.scMutex.Unlock()

	go func() {
		var enqueued, errs int
		for {
			select {
			case <-ctx.Done():
				return
			case <-k.producer.Successes():
				enqueued++
			case err := <-k.producer.Errors():
				if err != nil && err.Msg != nil {
					topic := err.Msg.Topic
					data, _ := err.Msg.Value.Encode()
					k.opts.Logger.Error().Err(err).Str("topic", topic).Str("data", string(data)).Msg("Failed to produce message")
				}
				errs++
			case err := <-k.consumer.Errors():
				k.opts.Logger.Error().Err(err).Msg("Failed to Consume message")
			}
		}
	}()

	return nil
}

func (k *kBroker) Close(cancel context.CancelFunc) error {
	if !k.isConnected() {
		return nil
	}

	k.scMutex.Lock()
	defer k.scMutex.Unlock()

	_ = k.producer.Close()
	_ = k.consumer.Close()

	k.connected = false
	cancel()
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

func (k *kBroker) Publish(_ context.Context, msg *Message) error {
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

func (k *kBroker) Subscribe(ctx context.Context, procs ...Processor) error {

	var ts []string
	ps := make(map[string]Processor, len(procs))

	for _, p := range procs {
		ts = append(ts, p.Topic)
		ps[p.Topic] = p
	}

	handler := &consumerGroupHandler{logger: k.opts.Logger, processors: ps}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := k.consumer.Consume(ctx, ts, handler)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					continue
				default:
					k.opts.Logger.Error().Err(err).Msg("Consume error")
					return
				}
			}
		}
	}()
	return nil
}

func NewKafkaBroker(ctx context.Context, name string, opt ...Option) Broker {
	opts := NewOptions(name)

	for _, o := range opt {
		o(&opts)
	}

	return &kBroker{
		opts: opts,
	}
}

type consumerGroupHandler struct {
	logger     *zerolog.Logger
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
			h.logger.Warn().Str("topic", msg.Topic).Msg("processor not exist")
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
			h.logger.Error().Err(err).Str("topic", msg.Topic).Str("content", string(msg.Value)).Int("retry", i).Msg("retry process error")
		}
		h.logger.Error().Err(err).Str("topic", msg.Topic).Str("content", string(msg.Value)).Msg("process error")
	}
	return nil
}
