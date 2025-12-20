package deliveryapp

import (
	"context"
	"fmt"
	"service-courier/internal/domain/order"
)

const (
	created   = "created"
	deleted = "deleted"
	completed = "completed"
)

type factoryEventStrategy struct {
	delExec deliveryExecutor
}

func NewFactoryEventStrategy(delExec deliveryExecutor) factoryEventStrategy {
	return factoryEventStrategy{delExec: delExec}
}

func (f factoryEventStrategy) GetEventStrategy(statusMsg string, statusNow string) (Executor, error) {
	switch {
	case statusMsg == created && statusNow != deleted && statusNow != completed:
		return AssignStrategy{f.delExec}, nil
	case statusMsg == deleted && statusNow == deleted:
		return UnassignStrategy{f.delExec}, nil
	case statusMsg == completed && statusNow == completed:
		return CompleteStrategy{f.delExec}, nil
	default:
		return nil, fmt.Errorf("the type of status (%s) is not do to processing", statusMsg)
	}
}

type AssignStrategy struct {
	assign deliveryAssign
}

func (as AssignStrategy) Execute(ctx context.Context, orderID order.OrderID) error {
	if _, err := as.assign.Assign(ctx, orderID); err != nil {
		return err
	}
	return nil
}

type UnassignStrategy struct {
	unassign deliveryUnassign
}

func (us UnassignStrategy) Execute(ctx context.Context, orderID order.OrderID) error {
	if _, err := us.unassign.Unassign(ctx, orderID); err != nil {
		return err
	}
	return nil
}

type CompleteStrategy struct {
	complete deliveryComplete
}

func (cmp CompleteStrategy) Execute(ctx context.Context, orderID order.OrderID) error {
	if _, err := cmp.complete.Complete(ctx, orderID); err != nil {
		return err
	}
	return nil
}
