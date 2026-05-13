package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-collector"`
}

type Services struct {
	Subscribe string `yaml:"subscribe" env:"SUBSCRIBE_ADDRESS" env-default:"localhost:8084"`
}

type Kafka struct {
	Brokers       []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	TaskTopic     string   `yaml:"task_topic" env:"KAFKA_TASK_TOPIC" env-default:"repository.collect.tasks"`
	ResultTopic   string   `yaml:"result_topic" env:"KAFKA_RESULT_TOPIC" env-default:"repository.collect.results"`
	ConsumerGroup string   `yaml:"consumer_group" env:"KAFKA_TASK_CONSUMER_GROUP" env-default:"collector-tasks"`
}

type Config struct {
	App      App           `yaml:"app"`
	Services Services      `yaml:"services"`
	Kafka    Kafka         `yaml:"kafka"`
	Logger   logger.Config `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
