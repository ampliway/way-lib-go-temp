package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/ampliway/way-lib-go/msg"
)

var _ msg.SubscriberV1[any] = (*Subscriber[any])(nil)

type Subscriber[T any] struct {
	producer *Producer
	id       id.ID
	cfg      *Config
	client   sarama.ConsumerGroup
}

func NewSub[T any](cfg *Config, producer *Producer, id id.ID) (*Subscriber[T], error) {
	return &Subscriber[T]{
		producer: producer,
		id:       id,
		cfg:      cfg,
	}, nil
}

func (s *Subscriber[T]) Publish(key string, msg interface{}) error {
	return s.producer.Publish(key, msg)
}

func (s *Subscriber[T]) PublishT(topicName, key string, msg interface{}) error {
	return s.producer.PublishT(topicName, key, msg)
}

func (s *Subscriber[T]) Subscribe(queueGroup string, execution func(m *msg.Message[T]) bool) error {
	topicName := topicName(new(T))

	return s.SubscribeT(topicName, queueGroup, execution)
}

func (s *Subscriber[T]) SubscribeT(topicName, queueGroup string, execution func(msg *msg.Message[T]) bool) error {
	config := defaultConfig(s.cfg)

	err := s.producer.CreateTopicIfNotExist(topicName, 3, 3)
	if err != nil {
		return err
	}

	queueGroup = fmt.Sprintf("%s.%s", reflection.AppNamePkg(), queueGroup)

	client, err := sarama.NewConsumerGroup(strings.Split(s.cfg.KafkaServers, ","), queueGroup, config)
	if err != nil {
		return err
	}

	s.client = client

	consumer := Consumer[T]{
		ready:     make(chan bool),
		execution: execution,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx := context.Background()

		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(topicName, ","), &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready

	return nil
}

func (s *Subscriber[T]) Shutdown() {
	if err := s.client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

	defer s.producer.Shutdown()
}

// Consumer represents a Sarama consumer group consumer
type Consumer[T any] struct {
	ready     chan bool
	execution func(m *msg.Message[T]) bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer[T]) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer[T]) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			var finalValue T

			err := json.Unmarshal(message.Value, &finalValue)
			if err != nil {
				log.Printf("Cannot unmarshal value: %s - %s - %s", string(message.Value), message.Timestamp, message.Topic)

				continue
			}

			result := consumer.execution(&msg.Message[T]{
				MessageID: "",
				TraceID:   "",
				Timestamp: message.Timestamp.Unix(),
				Body:      finalValue,
			})

			if result {
				session.MarkMessage(message, "")
			}

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
