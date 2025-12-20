package orderhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"service-courier/internal/domain/order"

	"github.com/IBM/sarama"
)

type consumeHandler struct {
	gateway orderGateway
	factory eventStrategyFactory
}

func NewConsumeHandler(gw orderGateway, f eventStrategyFactory) *consumeHandler {
	return &consumeHandler{
		gateway: gw,
		factory: f,
	}
}

func (ch *consumeHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (ch *consumeHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (ch *consumeHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()
	for dtoMsg := range claim.Messages() {
		fmt.Printf(
			"order.changed handler: received message: key=%s, value=%s, partition=%d, offset=%d\n",
			string(dtoMsg.Key), string(dtoMsg.Value), dtoMsg.Partition, dtoMsg.Offset,
		)
		var msg message
		if err := json.Unmarshal(dtoMsg.Value, &msg); err != nil {
			fmt.Printf("order.changed handler: received bad message: %v\n", err)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		curStatus, err := ch.getOrderStatus(ctx, dtoToDomainOrderID(&msg))
		if err != nil {
			fmt.Printf("order.changed handler: failed to get actual order status, retry...: %v\n", err)
			continue
		}
		st, err := ch.factory.GetEventStrategy(msg.Status, curStatus)
		if err != nil {
			fmt.Printf("order.changed handler: %v\n", err)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		if err := st.Execute(ctx, dtoToDomainOrderID(&msg)); err != nil {
			fmt.Printf("order.changed handler: failed execute order: %v\n", err)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		sess.MarkMessage(dtoMsg, "")
	}
	return nil
}

func (ch *consumeHandler) getOrderStatus(ctx context.Context, orderID order.OrderID) (string, error) {
	order, err := ch.gateway.GetOrderByID(ctx, orderID)
	if err != nil {
		return "", err
	}
	return order.Status, nil
}
