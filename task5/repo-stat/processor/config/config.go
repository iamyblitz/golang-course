package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-processor"`
}

type Services struct {
	Subscribe string `yaml:"subscribe" env:"SUBSCRIBE_ADDRESS" env-default:"localhost:8084"`
}

type Database struct {
	DSN            string `yaml:"dsn" env:"PROCESSOR_DATABASE_DSN" env-required:"true"`
	MigrationsPath string `yaml:"migrations_path" env:"PROCESSOR_MIGRATIONS_PATH" env-default:"file://processor/migrations"`
}

type Kafka struct {
	Brokers       []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	TaskTopic     string   `yaml:"task_topic" env:"KAFKA_TASK_TOPIC" env-default:"repository.collect.tasks"`
	ResultTopic   string   `yaml:"result_topic" env:"KAFKA_RESULT_TOPIC" env-default:"repository.collect.results"`
	ConsumerGroup string   `yaml:"consumer_group" env:"KAFKA_RESULT_CONSUMER_GROUP" env-default:"processor-results"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	Database Database          `yaml:"database"`
	Kafka    Kafka             `yaml:"kafka"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
