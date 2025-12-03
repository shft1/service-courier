package delivery

import (
	"context"
	"service-courier/internal/entity/delivery"
)

type deliveryService struct {
	deliveryRepo deliveryRepository
	courierRepo  courierRepository
	txManager    TxManagerDo
}

func NewDeliveryService(dr deliveryRepository, cr courierRepository, txManager TxManagerDo) *deliveryService {
	return &deliveryService{
		deliveryRepo: dr,
		courierRepo:  cr,
		txManager:    txManager,
	}
}

func (ds *deliveryService) DeliveryAssign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryAssign, error) {
	var deliveryAssign *delivery.DeliveryAssign
	err := ds.txManager.Do(ctx, func(ctx context.Context) error {
		res, err := ds.deliveryAssignInternal(ctx, orderID)
		deliveryAssign = res
		return err
	})
	if err != nil {
		return nil, deliveryServiceMapError(err)
	}
	return deliveryAssign, nil
}

func (ds *deliveryService) deliveryAssignInternal(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryAssign, error) {
	var deliveryAssign delivery.DeliveryAssign
	var deliveryCreate delivery.DeliveryCreate

	courier, err := ds.courierRepo.GetAvailable(ctx)
	if err != nil {
		return nil, err
	}
	deadline := delivery.DeliveryTimeFactory(courier.TransportType).TimeCalculate()

	deliveryCreate.CourierID = courier.ID
	deliveryCreate.OrderID = orderID.OrderID
	deliveryCreate.DeliveryDeadline = deadline

	delivery, err := ds.deliveryRepo.CreateDelivery(ctx, &deliveryCreate)
	if err != nil {
		return nil, err
	}
	if err := ds.courierRepo.SetBusy(ctx, courier.ID); err != nil {
		return nil, err
	}

	deliveryAssign.CourierID = delivery.CourierID
	deliveryAssign.OrderID = delivery.OrderID
	deliveryAssign.TransportType = courier.TransportType
	deliveryAssign.DeliveryDeadline = delivery.DeliveryDeadline

	return &deliveryAssign, nil
}

func (ds *deliveryService) DeliveryUnassign(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryUnassign, error) {
	var deliveryUnassign *delivery.DeliveryUnassign
	err := ds.txManager.Do(ctx, func(ctx context.Context) error {
		res, err := ds.deliveryUnassignInternal(ctx, orderID)
		deliveryUnassign = res
		return err
	})
	if err != nil {
		return nil, deliveryServiceMapError(err)
	}
	return deliveryUnassign, nil
}

func (ds *deliveryService) deliveryUnassignInternal(ctx context.Context, orderID *delivery.DeliveryOrderID) (*delivery.DeliveryUnassign, error) {
	var deliveryUnassign delivery.DeliveryUnassign
	const status = "unassigned"

	delivery, err := ds.deliveryRepo.DeleteDelivery(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if err = ds.courierRepo.SetAvailable(ctx, delivery.CourierID); err != nil {
		return nil, err
	}

	deliveryUnassign.OrderID = delivery.OrderID
	deliveryUnassign.Status = status
	deliveryUnassign.CourierID = delivery.CourierID

	return &deliveryUnassign, nil

}

func (ds *deliveryService) DeliveryCheck(ctx context.Context) error {
	if err := ds.deliveryRepo.RecheckDelivery(ctx); err != nil {
		return deliveryServiceMapError(err)
	}
	return nil
}
