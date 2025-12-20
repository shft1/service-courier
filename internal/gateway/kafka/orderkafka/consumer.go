package orderkafka

import (
	"context"
	"fmt"
	"service-courier/internal/config/consumercfg"
	"time"

	"github.com/IBM/sarama"
)

type kafkaClient struct {
	client  sarama.ConsumerGroup
	handler consumerHandler
	topics  []string
}

func NewKafkaClient(env *consumercfg.ConsumerEnv, handler consumerHandler, topics []string) (*kafkaClient, error) {
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
	fmt.Println("kafka client successfully created")
	return &kafkaClient{client: client, handler: handler, topics: topics}, nil
}

func (kc *kafkaClient) Consume(ctx context.Context) {
	for {
		if err := kc.client.Consume(ctx, kc.topics, kc.handler); err != nil {
			fmt.Printf("consume error: %v\n", err)
		}
		if ctx.Err() != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (kc *kafkaClient) Close() {
	fmt.Println("closing kafka client...")
	kc.client.Close()
	fmt.Println("kafka client closed")
}
