package consumercfg

import (
	"os"
	"strconv"
	"time"

	"service-courier/observability/logger"
)

type ConsumerEnv struct {
	Host           string
	Port           string
	OrderHost      string
	OrderPort      string
	KafkaHost      string
	KafkaPort      string
	KafkaTopic     string
	ConsumerGroup  string
	CommitInterval time.Duration
	MaxAttempts    int
	Multiplier     float64
	Jitter         float64
	InitDelay      time.Duration
	MaxDelay       time.Duration
}

// SetupConsumerEnv - парсер env переменных
func SetupConsumerEnv(log logger.Logger) *ConsumerEnv {
	port := os.Getenv("CONSUMER_LOCALPORT")
	if port == "" {
		log.Info("CONSUMER_LOCALPORT is 8080")
		port = "8080"
	}
	interval, err := time.ParseDuration(os.Getenv("COMMIT_INTERVAL"))
	if err != nil {
		log.Warn("failed to parse kafka COMMIT_INTERVAL", logger.NewField("error", err))
		interval = 1 * time.Second
	}
	maxAttempts, err := strconv.Atoi(os.Getenv("MAX_ATTEMPTS"))
	if err != nil {
		log.Warn("failed to get MAX_ATTEMPTS as integer", logger.NewField("error", err))
		maxAttempts = 3
	}
	multiplier, err := strconv.ParseFloat(os.Getenv("MULTIPLIER"), 64)
	if err != nil {
		log.Warn("failed to get MULTIPLIER as integer", logger.NewField("error", err))
		multiplier = 2.0
	}
	jitter, err := strconv.ParseFloat(os.Getenv("JITTER"), 64)
	if err != nil {
		log.Warn("failed to get JITTER as integer", logger.NewField("error", err))
		jitter = 0.1
	}
	initDelay, err := time.ParseDuration(os.Getenv("INIT_DELAY"))
	if err != nil {
		log.Warn("failed to get INIT_DELAY as time", logger.NewField("error", err))
		initDelay = 200 * time.Millisecond
	}
	maxDelay, err := time.ParseDuration(os.Getenv("MAX_DELAY"))
	if err != nil {
		log.Warn("failed to get MAX_DELAY as time", logger.NewField("error", err))
		maxDelay = 5 * time.Second
	}

	return &ConsumerEnv{
		Port:           port,
		OrderHost:      os.Getenv("ORDER_HOST"),
		OrderPort:      os.Getenv("ORDER_GRPC_PORT"),
		KafkaHost:      os.Getenv("KAFKA_HOST"),
		KafkaPort:      os.Getenv("KAFKA_PORT"),
		KafkaTopic:     os.Getenv("KAFKA_TOPIC"),
		ConsumerGroup:  os.Getenv("CONSUMER_GROUP"),
		CommitInterval: interval,
		MaxAttempts:    maxAttempts,
		Multiplier:     multiplier,
		Jitter:         jitter,
		InitDelay:      initDelay,
		MaxDelay:       maxDelay,
	}
}
