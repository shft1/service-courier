package deliveryapp

import (
	"context"

	"service-courier/internal/domain/delivery"
	"service-courier/internal/domain/order"
)

type Arguments struct {
	DelRepo   deliveryRepository
	CourRepo  courierRepository
	Factory   timeCalculatorFactory
	TxManager txManagerDo
}

// deliveryService - сервис доставки
type deliveryService struct {
	deliveryRepo deliveryRepository
	courierRepo  courierRepository
	factory      timeCalculatorFactory
	txManager    txManagerDo
}

// NewDeliveryService - конструктор сервиса доставок
func NewDeliveryService(args Arguments) *deliveryService {
	return &deliveryService{
		deliveryRepo: args.DelRepo,
		courierRepo:  args.CourRepo,
		factory:      args.Factory,
		txManager:    args.TxManager,
	}
}

// Assign - создает доставку на свободного курьера
func (ds *deliveryService) Assign(ctx context.Context, orderID order.OrderID) (*delivery.AssignResult, error) {
	var assignRes *delivery.AssignResult

	err := ds.txManager.Do(ctx, func(ctx context.Context) error {
		res, err := ds.doAssign(ctx, orderID)
		assignRes = res
		return err
	})
	if err != nil {
		return nil, mapError(err)
	}
	return assignRes, nil
}

// doAssign - вспомогательная функция для Assign
func (ds *deliveryService) doAssign(ctx context.Context, orderID order.OrderID) (*delivery.AssignResult, error) {
	var assignRes delivery.AssignResult
	var assignCreate delivery.AssignCreate

	cour, err := ds.courierRepo.GetAvailable(ctx)
	if err != nil {
		return nil, err
	}
	calc, err := ds.factory.GetDeliveryCalculator(cour.TransportType)
	if err != nil {
		return nil, err
	}
	assignCreate.CourierID = cour.ID
	assignCreate.OrderID = orderID.OrderID
	assignCreate.Deadline = calc.Calculate()

	del, err := ds.deliveryRepo.Create(ctx, &assignCreate)
	if err != nil {
		return nil, err
	}
	if _, err := ds.courierRepo.SetBusy(ctx, del.CourierID); err != nil {
		return nil, err
	}
	assignRes.CourierID = del.CourierID
	assignRes.OrderID = del.OrderID
	assignRes.TransportType = cour.TransportType
	assignRes.Deadline = del.Deadline

	return &assignRes, nil
}

// Unassign - удаляет доставку и освобождает соответствующего курьера
func (ds *deliveryService) Unassign(ctx context.Context, orderID order.OrderID) (*delivery.UnassignResult, error) {
	var unassignRes *delivery.UnassignResult

	err := ds.txManager.Do(ctx, func(ctx context.Context) error {
		res, err := ds.doUnassign(ctx, orderID)
		unassignRes = res
		return err
	})
	if err != nil {
		return nil, mapError(err)
	}
	return unassignRes, nil
}

// doUnassign - вспомогательная функция для Unassign
func (ds *deliveryService) doUnassign(ctx context.Context, orderID order.OrderID) (*delivery.UnassignResult, error) {
	var unassignRes delivery.UnassignResult

	del, err := ds.deliveryRepo.Delete(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if _, err = ds.courierRepo.SetAvailable(ctx, del.CourierID); err != nil {
		return nil, err
	}
	unassignRes.OrderID = del.OrderID
	unassignRes.Status = delivery.UnassignStatus
	unassignRes.CourierID = del.CourierID

	return &unassignRes, nil
}

// Complete - отмечает доставку как выполненную, освобождает соответствующего курьера
func (ds *deliveryService) Complete(ctx context.Context, orderID order.OrderID) (*delivery.CompleteResult, error) {
	var completeRes *delivery.CompleteResult

	err := ds.txManager.Do(ctx, func(ctx context.Context) error {
		res, err := ds.doComplete(ctx, orderID)
		completeRes = res
		return err
	})
	if err != nil {
		return nil, mapError(err)
	}
	return completeRes, nil
}

// doComplete - вспомогательная функция для Complete
func (ds *deliveryService) doComplete(ctx context.Context, orderID order.OrderID) (*delivery.CompleteResult, error) {
	var completeRes delivery.CompleteResult

	del, err := ds.deliveryRepo.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	courID, err := ds.courierRepo.SetAvailable(ctx, del.CourierID)
	if err != nil {
		return nil, err
	}
	completeRes.CourierID = courID
	completeRes.OrderID = del.OrderID
	completeRes.Deadline = del.Deadline

	return &completeRes, nil
}

// CheckDelivery - проверяет состояние доставок
func (ds *deliveryService) CheckDelivery(ctx context.Context) error {
	if err := ds.courierRepo.ReleaseStaleBusy(ctx); err != nil {
		return mapError(err)
	}
	return nil
}
