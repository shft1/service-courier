package kafka

import (
	"context"
	"net"
	"time"

	"github.com/IBM/sarama"

	"service-courier/internal/config/consumercfg"
	"service-courier/observability/logger"
)

type Arguments struct {
	Log     logger.Logger
	Env     *consumercfg.ConsumerEnv
	Handler consumerHandler
	Topics  []string
}

// kafkaClient - клиент Kafka
type kafkaClient struct {
	log     logger.Logger
	client  sarama.ConsumerGroup
	handler consumerHandler
	topics  []string
}

// NewKafkaClient - конструктор Kafka клиент
func NewKafkaClient(args Arguments) (*kafkaClient, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = args.Env.CommitInterval

	broker := net.JoinHostPort(args.Env.KafkaHost, args.Env.KafkaPort)
	groupID := args.Env.ConsumerGroup

	client, err := sarama.NewConsumerGroup([]string{broker}, groupID, config)
	if err != nil {
		return nil, err
	}
	args.Log.Info("kafka client successfully created")
	return &kafkaClient{log: args.Log, client: client, handler: args.Handler, topics: args.Topics}, nil
}

// Consume - запуск consuming Kafka
func (kc *kafkaClient) Consume(ctx context.Context) {
	for {
		if err := kc.client.Consume(ctx, kc.topics, kc.handler); err != nil {
			kc.log.Error("kafka consume error", logger.NewField("error", err))
		}
		if ctx.Err() != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// Close - закрытие клиента Kafka
func (kc *kafkaClient) Close() {
	kc.log.Info("closing kafka client...")
	if err := kc.client.Close(); err != nil {
		kc.log.Warn("failed to close connection with kafka gracefully", logger.NewField("error", err))
	}
	kc.log.Info("kafka client closed")
}
