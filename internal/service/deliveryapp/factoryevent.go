package deliveryapp

import (
	"context"
	"fmt"

	"github.com/shft1/service-courier/internal/domain/order"
)

const (
	Created   = "created"
	Deleted   = "deleted"
	Completed = "completed"
)

type factoryEventStrategy struct {
	delExec deliveryExecutor
}

func NewFactoryEventStrategy(delExec deliveryExecutor) factoryEventStrategy {
	return factoryEventStrategy{delExec: delExec}
}

func (f factoryEventStrategy) GetEventStrategy(statusMsg, statusNow string) (Executor, error) {
	switch {
	case statusMsg == Created && statusNow != Deleted && statusNow != Completed:
		return AssignStrategy{f.delExec}, nil
	case statusMsg == Deleted && statusNow == Deleted:
		return UnassignStrategy{f.delExec}, nil
	case statusMsg == Completed && statusNow == Completed:
		return CompleteStrategy{f.delExec}, nil
	default:
		return nil, fmt.Errorf("the type of status (%s) is not do to processing", statusMsg)
	}
}

type AssignStrategy struct {
	DeliveryAssign
}

func (as AssignStrategy) Execute(ctx context.Context, orderID order.OrderID) error {
	if _, err := as.Assign(ctx, orderID); err != nil {
		return err
	}
	return nil
}

type UnassignStrategy struct {
	DeliveryUnassign
}

func (us UnassignStrategy) Execute(ctx context.Context, orderID order.OrderID) error {
	if _, err := us.Unassign(ctx, orderID); err != nil {
		return err
	}
	return nil
}

type CompleteStrategy struct {
	DeliveryComplete
}

func (cmp CompleteStrategy) Execute(ctx context.Context, orderID order.OrderID) error {
	if _, err := cmp.Complete(ctx, orderID); err != nil {
		return err
	}
	return nil
}
