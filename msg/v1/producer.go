package v1

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/ampliway/way-lib-go/helper/id"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/ampliway/way-lib-go/msg"
	"github.com/iancoleman/strcase"
)

var _ msg.ProducerV1 = (*Producer)(nil)

type Producer struct {
	client   sarama.Client
	producer sarama.SyncProducer
	id       id.ID
	topics   map[string]bool
	topicMux sync.Mutex
}

func New(cfg *Config, id id.ID) (*Producer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%s: %w", msg.MODULE_NAME, errConfigNull)
	}

	if cfg.KafkaServers == "" {
		return nil, fmt.Errorf("%s: %w", msg.MODULE_NAME, errConfigServersEmpty)
	}

	config := defaultConfig()

	servers := strings.Split(cfg.KafkaServers, ",")

	client, err := sarama.NewClient(servers, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", msg.MODULE_NAME, errKafkaConnect, err)
	}

	producer, err := sarama.NewSyncProducer(servers, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", msg.MODULE_NAME, errProducerStart, err)
	}

	return &Producer{
		client:   client,
		producer: producer,
		id:       id,
		topics:   map[string]bool{},
		topicMux: sync.Mutex{},
	}, nil
}

func (p *Producer) Publish(m interface{}) error {
	topicName := topicName(m)

	err := p.createTopicIfNotExist(topicName)
	if err != nil {
		return err
	}

	value, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errUnmarshal, topicName)
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder("a"),
		Topic: topicName,
		Value: sarama.StringEncoder(value),
	})
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errPublish, topicName)
	}

	return nil
}

func (p *Producer) Shutdown() {
	defer p.producer.Close()
	defer p.client.Close()
}

func topicName(msg interface{}) string {
	msgName := strcase.ToKebab(reflection.TypeName(msg))
	packageName := strings.ToLower(reflection.AppName(msg))
	return fmt.Sprintf("%s.%s", packageName, msgName)
}

func (p *Producer) createTopicIfNotExist(topicName string) error {
	p.topicMux.Lock()
	defer p.topicMux.Unlock()
	if item, ok := p.topics[topicName]; ok && item {
		return nil
	}

	admin, err := sarama.NewClusterAdminFromClient(p.client)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errAdminClientStart, topicName)
	}

	foundTopic := false
	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("%s: %w", msg.MODULE_NAME, err)
	}

	for topic := range topics {
		if topic == topicName {
			foundTopic = true

			break
		}
	}

	if !foundTopic {
		cleanupPolicy := "compact"

		topicDetail := &sarama.TopicDetail{
			NumPartitions:     3,
			ReplicationFactor: 3,
			ConfigEntries: map[string]*string{
				"cleanup.policy": &cleanupPolicy,
			},
		}

		err = admin.CreateTopic(topicName, topicDetail, false)
		if err != nil {
			return fmt.Errorf("%s: %w", msg.MODULE_NAME, err)
		}
	}

	p.topics[topicName] = true

	return nil
}

func defaultConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_3_1_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	return config
}
