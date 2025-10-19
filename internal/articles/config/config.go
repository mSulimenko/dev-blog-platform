package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

const (
	standardConfigPath = "./configs/articles.yaml"
)

type Config struct {
	Env string `yaml:"env" env:"ENV" env-default:"local"`
	Db  DB     `yaml:"db" env-required:"true"`
}

type DB struct {
	MigrationsDir string `yaml:"migrations_dir" env:"MIGRATIONS_DIR" env-required:"true"`
	Dsn           string `yaml:"dsn" env:"MIGRATIONS_DIR" env-required:"true"`
}

func Load() *Config {

	configPath := os.Getenv("ARTICLES_CONFIG_PATH")
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
