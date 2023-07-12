package v1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ampliway/way-lib-go/ctx"
	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/ampliway/way-lib-go/msg"
	"github.com/nats-io/nats.go"
)

var _ msg.SubscriberV1[any] = (*NatsSubscriber[any])(nil)

type NatsSubscriber[T any] struct {
	producer *NatsProducer
	id       id.ID
}

func NewSub[T any](config *Config, producer *NatsProducer, id id.ID) (*NatsSubscriber[T], error) {
	return &NatsSubscriber[T]{
		producer: producer,
		id:       id,
	}, nil
}

func (n *NatsSubscriber[T]) Publish(ctx ctx.V1, msg interface{}) error {
	return n.producer.Publish(ctx, msg)
}

func (n *NatsSubscriber[T]) Subscribe(queueGroup string, execution func(m *msg.Message[T]) bool) error {
	subjectName := subjectName(new(T))
	packageName := strings.ToLower(reflection.AppNamePkg())

	fmt.Printf("Subscribe stream %s with subject %s\n", packageName, subjectName)

	_, err := n.producer.js.AddConsumer(packageName, &nats.ConsumerConfig{
		Durable:        queueGroup,
		DeliverSubject: "consumer." + subjectName,
		DeliverGroup:   queueGroup,
		AckPolicy:      nats.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = n.producer.js.QueueSubscribe(subjectName, queueGroup, func(m *nats.Msg) {
		var finalValue T

		err := json.Unmarshal(m.Data, &finalValue)
		if err != nil {
			m.Nak()

			panic(err)
		}

		metadata, err := m.Metadata()
		if err != nil {
			panic(err)
		}

		result := execution(&msg.Message[T]{
			MessageID: m.Header.Get(msg.HEADER_X_MSG_ID),
			TraceID:   m.Header.Get(msg.HEADER_X_TRACE_ID),
			Timestamp: metadata.Timestamp.Unix(),
			Body:      finalValue,
		})

		if result {
			err = m.Ack()
			if err != nil {
				panic(err)
			}
		}
	})

	return err
}

func (n *NatsSubscriber[T]) Shutdown() {
	n.producer.Shutdown()
}
