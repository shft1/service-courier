package consumercfg

import (
	"os"
	"service-courier/observability/logger"
	"time"
)

type ConsumerEnv struct {
	OrderHost      string
	OrderPort      string
	KafkaHost      string
	KafkaPort      string
	KafkaTopic     string
	ConsumerGroup  string
	CommitInterval time.Duration
}

// SetupConsumerEnv - парсер env переменных
func SetupConsumerEnv(log logger.Logger) *ConsumerEnv {
	interval, err := time.ParseDuration(os.Getenv("COMMIT_INTERVAL"))
	if err != nil {
		log.Warn("failed to parse kafka commit offset interval", logger.NewField("error", err))
		interval = 1 * time.Second
	}

	return &ConsumerEnv{
		OrderHost:      os.Getenv("ORDER_HOST"),
		OrderPort:      os.Getenv("ORDER_GRPC_PORT"),
		KafkaHost:      os.Getenv("KAFKA_HOST"),
		KafkaPort:      os.Getenv("KAFKA_PORT"),
		KafkaTopic:     os.Getenv("KAFKA_TOPIC"),
		ConsumerGroup:  os.Getenv("CONSUMER_GROUP"),
		CommitInterval: interval,
	}
}
