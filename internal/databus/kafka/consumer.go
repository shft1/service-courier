package kafka

import (
	"context"
	"fmt"
	"service-courier/internal/config/consumercfg"
	"service-courier/observability/logger"
	"time"

	"github.com/IBM/sarama"
)

// kafkaClient - клиент Kafka
type kafkaClient struct {
	log     logger.Logger
	client  sarama.ConsumerGroup
	handler consumerHandler
	topics  []string
}

// NewKafkaClient - конструктор Kafka клиент
func NewKafkaClient(log logger.Logger, env *consumercfg.ConsumerEnv, handler consumerHandler, topics []string) (*kafkaClient, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = env.CommitInterval

	broker := fmt.Sprintf("%v:%v", env.KafkaHost, env.KafkaPort)
	groupID := env.ConsumerGroup

	client, err := sarama.NewConsumerGroup([]string{broker}, groupID, config)
	if err != nil {
		return nil, err
	}
	log.Info("kafka client successfully created")
	return &kafkaClient{log: log, client: client, handler: handler, topics: topics}, nil
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
	kc.client.Close()
	kc.log.Info("kafka client closed")
}
