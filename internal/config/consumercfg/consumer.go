package consumercfg

import (
	"fmt"
	"os"
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
func SetupConsumerEnv() *ConsumerEnv {
	interval, err := time.ParseDuration(os.Getenv("COMMIT_INTERVAL"))
	if err != nil {
		fmt.Println("error parsing kafka commit interval:", err)
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
