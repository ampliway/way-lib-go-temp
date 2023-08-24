package v1

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/IBM/sarama"
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

	config := defaultConfig(cfg)

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

func (p *Producer) Publish(key string, m interface{}) error {
	topicName := topicName(m)

	return p.PublishT(topicName, key, m)
}

func (p *Producer) PublishT(topicName, key string, m interface{}) error {
	err := p.CreateTopicIfNotExist(topicName, 3, 3)
	if err != nil {
		return err
	}

	value, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("%s: %w: %s", msg.MODULE_NAME, errUnmarshal, topicName)
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Key:   sarama.StringEncoder(key),
		Topic: topicName,
		Value: sarama.StringEncoder(value),
	})
	if err != nil {
		return fmt.Errorf("%s: %w: %s: %+v", msg.MODULE_NAME, errPublish, topicName, err)
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

func (p *Producer) CreateTopicIfNotExist(topicName string, numPartitions int32, replicationFactor int16) error {
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
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
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

func defaultConfig(cfg *Config) *sarama.Config {
	clientID, _ := os.Hostname()

	config := sarama.NewConfig()
	config.Version = sarama.V3_3_2_0
	config.ClientID = clientID
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = true

	if cfg.KafkaUsername != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.KafkaUsername
		config.Net.SASL.Password = cfg.KafkaPassword
		config.Net.SASL.Handshake = true
		if cfg.KafkaAlgorithm == "sha512" {
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		} else if cfg.KafkaAlgorithm == "sha256" {
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		} else {
			log.Fatalf("invalid SHA algorithm \"%s\": can be either \"sha256\" or \"sha512\"", cfg.KafkaAlgorithm)
		}
	}

	if cfg.KafkaCAFile != "" {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = createTLSConfiguration(cfg)
	}

	return config
}

func createTLSConfiguration(cfg *Config) (t *tls.Config) {
	cert, err := tls.LoadX509KeyPair(cfg.KafkaCertFile, cfg.KafkaKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := os.ReadFile(cfg.KafkaCAFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	return t
}
