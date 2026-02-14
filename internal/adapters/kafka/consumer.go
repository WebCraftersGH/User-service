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
	}

	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, nil
	}

	if err := c.Subscribe(config.Topic, nil); err != nil {
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
			// логирование ошибки
			continue
		}

		if kafkaMsg == nil {
			// логирование пустого сообщения
			continue
		}

		if err = c.handleMessage(kafkaMsg.Value); err != nil {
			// логирование ошибки
			continue
		}

		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			// логирование ошибки
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
