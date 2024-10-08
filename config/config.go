package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env               string `env:"ENV"`
	LogLevel          string `env:"LOG_LEVEL"`
	Postgres          Postgres
	KafkaNotification KafkaNotification
	Mail              Mail
}

type Postgres struct {
	Host            string `env:"PG_HOST"`
	Port            int    `env:"PG_PORT"`
	DbName          string `env:"PG_DB_NAME"`
	Password        string `env:"PG_PASSWORD"`
	User            string `env:"PG_USER"`
	PoolMax         int    `env:"PG_POOL_MAX"`
	MaxOpenConns    int    `env:"PG_MAX_OPEN_CONNS"`
	ConnMaxLifetime int    `env:"PG_CONN_MAX_LIFETIME"`
	MaxIdleConns    int    `env:"PG_MAX_IDLE_CONNS"`
	ConnMaxIdleTime int    `env:"PG_CONN_MAX_IDLE_TIME"`
}

type KafkaNotification struct {
	ConsumerGroup string   `env:"KAFKA_NOTIFICATION_CONSUMER_GROUP"`
	ConsumerUrl   []string `env:"KAFKA_NOTIFICATION_CONSUMER_URL"`
	Topic         string   `env:"KAFKA_NOTIFICATION_TOPIC"`
}

type Mail struct {
	Host     string `env:"MAIL_HOST"`
	Port     int    `env:"MAIL_PORT"`
	Address  string `env:"MAIL_ADDRESS"`
	Password string `env:"MAIL_PASSWORD"`
}

func MustLoad() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	opts := env.Options{RequiredIfNoDef: true}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		log.Fatalf("parse config error: %s", err)
	}

	return cfg
}
