package worker

import (
	"context"
	"fmt"
	"service-courier/internal/domain/delivery"
	"service-courier/internal/gateway/order"
	"time"
)

type gateway interface {
	GetOrders(ctx context.Context, cursor time.Time) ([]*order.OrderResponse, error)
}

type assigner interface {
	Assign(ctx context.Context, orderID delivery.OrderID) (*delivery.AssignResult, error)
}

// orderPoller - фоновый воркер получения и назначения заказов
type orderPoller struct {
	cursor   time.Time
	period   time.Duration
	gateway  gateway
	assigner assigner
}

func NewOrderPoller(p time.Duration, gw gateway, as assigner) *orderPoller {
	return &orderPoller{
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
			fmt.Println("order polled")
		case <-ctx.Done():
			op.stop()
			return
		}
	}
}

func (op *orderPoller) processTick(ctx context.Context) {
	orders, err := op.fetch(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := op.assignDeliveries(ctx, orders); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (op *orderPoller) fetch(ctx context.Context) ([]*order.OrderResponse, error) {
	out, err := op.gateway.GetOrders(ctx, op.cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return out, nil
}

func (op *orderPoller) assignDeliveries(ctx context.Context, ords []*order.OrderResponse) error {
	for _, order := range ords {
		fmt.Printf("delivery on order '%s' is planning...\n", order.OrderID)
		del, err := op.assigner.Assign(ctx, delivery.OrderID{OrderID: order.OrderID})
		if err != nil {
			op.cursor = order.CreatedAt
			fmt.Printf("cursor updated: %s\n", op.cursor.Format("2006-01-02 15:04:05"))
			return fmt.Errorf("failed to plan delivery on order '%s': %w", order.OrderID, err)
		}
		fmt.Printf("delivery on order '%s' is plan!\n", del.OrderID)
	}
	op.cursor = time.Now()
	fmt.Printf("cursor updated: %s\n", op.cursor.Format("2006-01-02 15:04:05"))
	return nil
}

func (op *orderPoller) stop() {
	fmt.Println("stop order polling")
}
