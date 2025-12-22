package orderhandler

import (
	"context"
	"encoding/json"
	"service-courier/internal/domain/order"
	"service-courier/observability/logger"

	"github.com/IBM/sarama"
)

type consumeHandler struct {
	log logger.Logger
	gateway orderGateway
	factory eventStrategyFactory
}

func NewConsumeHandler(log logger.Logger, gw orderGateway, f eventStrategyFactory) *consumeHandler {
	return &consumeHandler{
		log: log,
		gateway: gw,
		factory: f,
	}
}

func (ch *consumeHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (ch *consumeHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (ch *consumeHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()
	for dtoMsg := range claim.Messages() {
		ch.log.Info(
			"order.changed handler: received message",
			logger.NewField("value", string(dtoMsg.Value)),
			logger.NewField("partition", dtoMsg.Partition),
			logger.NewField("offset", dtoMsg.Offset),
		)
		var msg message
		if err := json.Unmarshal(dtoMsg.Value, &msg); err != nil {
			ch.log.Error(
				"order.changed handler: received bad message",
				logger.NewField("error", err),
			)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		curStatus, err := ch.getOrderStatus(ctx, dtoToDomainOrderID(&msg))
		if err != nil {
			ch.log.Warn(
				"order.changed handler: failed to get actual order status, retry...",
				logger.NewField("error", err),
			)
			continue
		}
		st, err := ch.factory.GetEventStrategy(msg.Status, curStatus)
		if err != nil {
			ch.log.Warn(
				"order.changed handler: failed to get event strategy",
				logger.NewField("error", err),
			)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		if err := st.Execute(ctx, dtoToDomainOrderID(&msg)); err != nil {
			ch.log.Error(
				"order.changed handler: failed to execute order logic",
				logger.NewField("error", err),
			)
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
