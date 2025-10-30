package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

const (
	standardConfigPath = "./configs/notify.yaml"
)

type Config struct {
	Env   string      `yaml:"env" env:"ENV" env-default:"local"`
	Kafka KafkaConfig `yaml:"kafka"`
	Email EmailConfig `yaml:"email"`
}

type EmailConfig struct {
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     string `yaml:"smtp_port"`
	FromEmail    string `yaml:"from_email"`
	FromPassword string `yaml:"from_password"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS"`
}

func Load() *Config {

	configPath := os.Getenv("AUTH_CONFIG_PATH")
	if configPath == "" {
		configPath = standardConfigPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("has no articles config file, %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file, %s", configPath)
	}

	return &cfg
}
