package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/events/dto"
	"go.uber.org/zap"
)

const (
	userRegisteredTopic = "user-registered"
)

type Dispatcher struct {
	producer sarama.SyncProducer
	log      *zap.SugaredLogger
}

func NewKafkaDispatcher(brokers []string, log *zap.SugaredLogger) (*Dispatcher, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &Dispatcher{
		producer: producer,
		log:      log,
	}, nil
}

func (d *Dispatcher) UserRegistered(ctx context.Context, email, token, username string) error {
	event := dto.UserRegisteredEvent{
		Email:    email,
		Token:    token,
		Username: username,
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshall event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: userRegisteredTopic,
		Key:   sarama.StringEncoder(email),
		Value: sarama.ByteEncoder(jsonEvent),
	}

	partition, offset, err := d.producer.SendMessage(msg)
	if err != nil {
		d.log.Errorw("failed to send user_registered event",
			"error", err, "email", email)
		return err
	}
	d.log.Infow("user_registered event sent",
		"email", email, "partition", partition, "offset", offset)
	return nil
}

func (d *Dispatcher) Close() error {
	return d.producer.Close()
}
