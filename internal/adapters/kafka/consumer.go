package kafka

import (
	"encoding/json"
	"github.com/WebCraftersGH/User-service/internal/usecase"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

type Config struct {
	GroupID            string
	BootStrapServers   string
	TimeoutMS          int
	Topic              string
	ReadMessageTimeout time.Duration
	AutoOffsetStore    bool
	AutoCommit         bool
	AutoCommitInterval int
}

type kafkaConsumer struct {
	consumer *kafka.Consumer
	config   *Config
	stop     bool
	userSVC  usecase.UserService
	lg       usecase.Logger
}

var _ usecase.Consumer = (*kafkaConsumer)(nil)

func NewKafkaConsumer(
	config *Config,
	userSVC usecase.UserService,
	lg usecase.Logger,
) (*kafkaConsumer, error) {

	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        config.BootStrapServers,
		"group.id":                 config.GroupID,
		"session.timeout.ms":       config.TimeoutMS,
		"enable.auto.offset.store": config.AutoOffsetStore,
		"enable.auto.commit":       config.AutoCommit,
		"auto.commit.interval.ms":  config.AutoCommitInterval,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		lg.Error("[Kafka-consumer][NewConsumer][ERROR] - Error to create kafka consumer", "kafka_err", err)
		return nil, err
	}

	if err := c.Subscribe(config.Topic, nil); err != nil {
		lg.Error("[Kafka-consumer][SubscribeTopic][ERROR] - Error subscribe to topic", "kafka_err", err)
		return nil, err
	}

	return &kafkaConsumer{
		consumer: c,
		config:   config,
		userSVC:  userSVC,
		lg:       lg,
	}, nil
}

func (c *kafkaConsumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(c.config.ReadMessageTimeout)
		if err != nil {
			if ke, ok := err.(kafka.Error); ok && ke.Code() == kafka.ErrTimedOut {
				continue
			}
			c.lg.Warn("[Kafka-consumer][ReadMessage][ERROR] - Read message error", "kafka_err", err)
			continue
		}

		if kafkaMsg == nil {
			c.lg.Warn("[Kafka-consumer][ReadMessage][WARN] - Kafka msg == nil")
			continue
		}

		if err = c.handleMessage(kafkaMsg.Value); err != nil {
			c.lg.Warn("[Kafka-consumer][HandleMessage][WARN] - Kafka handle msg error", "kafka_err", err)
			continue
		}

		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			c.lg.Warn("[Kafka-consumer][StortMessage][WARN] - Kafka store msg error", "kafka_err", err)
			continue
		}
	}
}

func (c *kafkaConsumer) Stop() error {
	c.stop = true

	if _, err := c.consumer.Commit(); err != nil {
		return err
	}

	return c.consumer.Close()
}

func (c *kafkaConsumer) handleMessage(message []byte) error {
	var u KafkaUser
	if err := json.Unmarshal(message, &u); err != nil {
		return err
	}

	domainUser := toDomainUser(u)

	if _, err := c.userSVC.CreateUser(nil, domainUser); err != nil {
		return err
	}

	return nil
}
