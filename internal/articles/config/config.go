package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

const (
	standardConfigPath = "./configs/articles.yaml"
)

type Config struct {
	Env  string `yaml:"env" env:"ENV" env-default:"local"`
	HTTP HTTP   `yaml:"http"`
	DB   DB     `yaml:"db"`
}

type HTTP struct {
	Host         string        `yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
	Port         string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DB struct {
	MigrationsDir string `yaml:"migrations_dir"`
	Dsn           string `yaml:"dsn" env:"DB_DSN,required"`
	MaxOpenConns  int    `yaml:"max_open_conns" env:"DB_MAX_CONNS"`
	MaxIdleConns  int    `yaml:"max_idle_conns"`
	MaxIdleTime   string `yaml:"max_idle_time"`
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
