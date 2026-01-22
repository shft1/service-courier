package orderworker

import (
	"context"
	"fmt"
	"time"

	"service-courier/internal/domain/delivery"
	"service-courier/internal/domain/order"
	"service-courier/observability/logger"
)

type gateway interface {
	GetOrders(ctx context.Context, cursor time.Time) ([]*order.Order, error)
}

type assigner interface {
	Assign(ctx context.Context, orderID order.OrderID) (*delivery.AssignResult, error)
}

// orderPoller - фоновый воркер получения и назначения заказов
type orderPoller struct {
	log      logger.Logger
	cursor   time.Time
	period   time.Duration
	gateway  gateway
	assigner assigner
}

// NewOrderPoller - конструктор фонового воркера получения и назначения заказов
func NewOrderPoller(log logger.Logger, p time.Duration, gw gateway, as assigner) *orderPoller {
	return &orderPoller{
		log:      log,
		cursor:   time.Now().Add(-5 * time.Second),
		period:   p,
		gateway:  gw,
		assigner: as,
	}
}

// Start - запуск фонового воркера получения и назначения заказов
func (op *orderPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(op.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			op.processTick(ctx)
			op.log.Info("order polled")
		case <-ctx.Done():
			op.stop()
			return
		}
	}
}

// processTick - получение заказов и назначение доставок
func (op *orderPoller) processTick(ctx context.Context) {
	orders, err := op.fetch(ctx)
	if err != nil {
		op.log.Error("failed to fetch orders", logger.NewField("error", err))
		return
	}
	if err := op.assignDeliveries(ctx, orders); err != nil {
		op.log.Error("failed to assign deliveries", logger.NewField("error", err))
		return
	}
}

// fetch - получение заказов
func (op *orderPoller) fetch(ctx context.Context) ([]*order.Order, error) {
	out, err := op.gateway.GetOrders(ctx, op.cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return out, nil
}

// assignDeliveries - назначение доставок
func (op *orderPoller) assignDeliveries(ctx context.Context, ords []*order.Order) error {
	for _, ord := range ords {
		op.log.Info("planning delivery on order...", logger.NewField("order_id", ord.OrderID))
		del, err := op.assigner.Assign(ctx, order.OrderID{OrderID: ord.OrderID})
		if err != nil {
			op.cursor = ord.CreatedAt
			op.log.Debug("order cursor was updated", logger.NewField("cursor", op.cursor))
			return fmt.Errorf("failed to plan delivery on order '%s': %w", ord.OrderID, err)
		}
		op.log.Info("successfully plan delivery on order", logger.NewField("order_id", del.OrderID))
	}
	op.cursor = time.Now()
	op.log.Debug("order cursor was updated", logger.NewField("cursor", op.cursor))
	return nil
}

func (op *orderPoller) stop() {
	op.log.Info("stop order polling")
}
