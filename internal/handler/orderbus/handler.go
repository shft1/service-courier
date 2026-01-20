package orderbus

import (
	"context"
	"encoding/json"
	"service-courier/internal/domain/order"
	"service-courier/observability/logger"
	"time"

	"github.com/IBM/sarama"
)

// consumeHandler - обработчик топика Kafka
type consumeHandler struct {
	log     logger.Logger
	gateway orderGateway
	factory eventStrategyFactory
}

// NewConsumeHandler - конструктор обработчика топика Kafka
func NewConsumeHandler(log logger.Logger, gw orderGateway, f eventStrategyFactory) *consumeHandler {
	return &consumeHandler{
		log:     log,
		gateway: gw,
		factory: f,
	}
}

// Setup - подготовка перед consuming
func (ch *consumeHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup - очистка после consuming
func (ch *consumeHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim - consuming топика Kafka
func (ch *consumeHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()
	for dtoMsg := range claim.Messages() {
		var msg message
		if err := json.Unmarshal(dtoMsg.Value, &msg); err != nil {
			ch.log.Error(
				"received bad message",
				logger.NewField("error", err),
			)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		ch.log.Info(
			"received message",
			logger.NewField("order_id", msg.OrderID),
			logger.NewField("status", msg.Status),
			logger.NewField("created_at", msg.CreatedAt),
			logger.NewField("partition", dtoMsg.Partition),
			logger.NewField("offset", dtoMsg.Offset),
		)
		curStatus, err := ch.getOrderStatus(ctx, dtoToDomainOrderID(&msg))
		if err != nil {
			ch.log.Warn(
				"failed to get actual order status, retry after 30 seconds",
				logger.NewField("error", err),
			)
			time.Sleep(30 * time.Second)
			continue
		}
		st, err := ch.factory.GetEventStrategy(msg.Status, curStatus)
		if err != nil {
			ch.log.Warn(
				"failed to get event strategy for order",
				logger.NewField("error", err),
			)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		if err := st.Execute(ctx, dtoToDomainOrderID(&msg)); err != nil {
			ch.log.Error(
				"failed to execute order logic",
				logger.NewField("error", err),
			)
			sess.MarkMessage(dtoMsg, "")
			continue
		}
		sess.MarkMessage(dtoMsg, "")
	}
	return nil
}

// getOrderStatus - получение актуального статуса заказа через gRPC
func (ch *consumeHandler) getOrderStatus(ctx context.Context, orderID order.OrderID) (string, error) {
	order, err := ch.gateway.GetOrderByID(ctx, orderID)
	if err != nil {
		return "", err
	}
	return order.Status, nil
}
