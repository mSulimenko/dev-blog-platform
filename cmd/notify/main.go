package main

import (
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/notify/config"
	"github.com/mSulimenko/dev-blog-platform/internal/notify/service"
	"github.com/mSulimenko/dev-blog-platform/internal/notify/transport/kafka"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting notify")
	cfg := config.Load()

	fmt.Println(cfg)

	log := logger.New(cfg.Env)
	defer log.Sync()

	emailService := service.NewEmailService(service.EmailConfig{
		SMTPHost:     cfg.Email.SMTPHost,
		SMTPPort:     cfg.Email.SMTPPort,
		FromEmail:    cfg.Email.FromEmail,
		FromPassword: cfg.Email.FromPassword,
	}, log)

	consumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, emailService, log)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer", "error", err)
	}
	defer consumer.Close()

	go consumer.Start()

	log.Info("Notify service started")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	log.Info("Notify service stopped")

}
