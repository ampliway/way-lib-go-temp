package v1

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/ampliway/way-lib-go/msg"
	"github.com/iancoleman/strcase"
	"github.com/nats-io/nats.go"
)

var _ msg.MsgV1 = (*Producer)(nil)

type Producer struct {
	js            nats.JetStreamContext
	nc            *nats.Conn
	subscriptions []*nats.Subscription
	id            id.ID
	topics        map[string]bool
	topicMux      sync.Mutex
}

func New(cfg *Config, id id.ID) (*Producer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%s: %w", msg.MODULE_NAME, errConfigNull)
	}

	if cfg.NatsServers == "" {
		return nil, fmt.Errorf("%s: %w", msg.MODULE_NAME, errConfigServersEmpty)
	}

	nc, err := nats.Connect(cfg.NatsServers)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", msg.MODULE_NAME, errNatsConnect, err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", msg.MODULE_NAME, errNatsConnect, err)
	}

	return &Producer{
		nc:            nc,
		js:            js,
		id:            id,
		topics:        map[string]bool{},
		topicMux:      sync.Mutex{},
		subscriptions: []*nats.Subscription{},
	}, nil
}

func (p *Producer) Publish(m interface{}) error {
	topicName := topicName(m)

	return p.PublishT(topicName, m)
}

func (p *Producer) PublishT(topicName string, m interface{}) error {
	err := p.CreateTopicIfNotExist(strings.ToLower(reflection.AppName(m)))
	if err != nil {
		return err
	}

	value, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errUnmarshal, topicName)
	}

	_, err = p.js.Publish(topicName, value)
	if err != nil {
		return fmt.Errorf("%s: %w: %s: %+v", msg.MODULE_NAME, errPublish, topicName, err)
	}

	return nil
}

func (p *Producer) Subscribe(m interface{}, queueGroup string, exec func(data []byte) bool) error {
	err := p.CreateTopicIfNotExist(strings.ToLower(reflection.AppName(m)))
	if err != nil {
		return err
	}

	topicName := topicName(m)

	return p.SubscribeT(m, topicName, queueGroup, exec)
}

func (p *Producer) SubscribeT(m interface{}, topicName string, queueGroup string, exec func(data []byte) bool) error {
	subscription, err := p.js.QueueSubscribe(topicName, queueGroup, func(msg *nats.Msg) {
		var err error
		if exec(msg.Data) {
			err = msg.Ack()
		} else {
			err = msg.Nak()
		}

		if err != nil {
			fmt.Println(err)
		}
	})

	p.subscriptions = append(p.subscriptions, subscription)

	return err
}

func (p *Producer) Shutdown() {
	defer p.nc.Drain()
}

func topicName(msg interface{}) string {
	msgName := strcase.ToKebab(reflection.TypeName(msg))
	packageName := strings.ToLower(reflection.AppName(msg))
	return fmt.Sprintf("%s.%s", packageName, msgName)
}

func (p *Producer) CreateTopicIfNotExist(topicName string) error {
	p.topicMux.Lock()
	defer p.topicMux.Unlock()
	if item, ok := p.topics[topicName]; ok && item {
		return nil
	}

	_, err := p.js.AddStream(&nats.StreamConfig{
		Name:      topicName,
		Subjects:  []string{topicName + ".>"},
		MaxAge:    time.Hour * 24 * 365,
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", msg.MODULE_NAME, err)
	}

	p.topics[topicName] = true

	return nil
}
