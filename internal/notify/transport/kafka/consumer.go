package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/mSulimenko/dev-blog-platform/internal/notify/service"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/events/dto"
	"go.uber.org/zap"
)

const (
	userRegisteredTopic = "user-registered"
)

type Consumer struct {
	consumer     sarama.Consumer
	log          *zap.SugaredLogger
	emailService *service.EmailService
}

func NewConsumer(brokers []string, emailService *service.EmailService, log *zap.SugaredLogger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{
		consumer:     consumer,
		emailService: emailService,
		log:          log,
	}, nil
}

func (c *Consumer) Start() {
	c.log.Info("Starting Kafka consumer with Sarama...")

	partitionConsumer, err := c.consumer.ConsumePartition(userRegisteredTopic, 0, sarama.OffsetNewest)
	if err != nil {
		c.log.Errorf("failed to start consumer: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			c.handleMessage(msg)
		case err = <-partitionConsumer.Errors():
			c.log.Errorf("Kafka consumer error: %v", err)
		}
	}
}

func (c *Consumer) handleMessage(msg *sarama.ConsumerMessage) {
	var event dto.UserRegisteredEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		c.log.Errorf("Parse event error: %v", err)
		return
	}

	if err := c.handleUserRegistered(event); err != nil {
		c.log.Errorf("Handle event error: %v", err)
	}
}

func (c *Consumer) handleUserRegistered(event dto.UserRegisteredEvent) error {
	c.log.Infof("Processing user registration: %s", event.Email)

	err := c.emailService.SendVerificationEmail(context.Background(), event.Email, event.Username, event.Token)
	if err != nil {
		c.log.Errorf("Failed to send verification email to %s: %v", event.Email, err)
		return err
	}

	c.log.Infof("Verification email sent to %s", event.Email)
	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
