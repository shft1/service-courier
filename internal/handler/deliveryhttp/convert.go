package deliveryhttp

import (
	"github.com/shft1/service-courier/internal/domain/delivery"
	"github.com/shft1/service-courier/internal/domain/order"
)

func toDomainOrderID(req DeliveryOrderRequest) order.OrderID {
	return order.OrderID{OrderID: req.OrderID}
}

func domainToDTOAssign(del *delivery.AssignResult) DeliveryAssignResponse {
	return DeliveryAssignResponse{
		CourierID:     del.CourierID,
		OrderID:       del.OrderID,
		TransportType: del.TransportType,
		Deadline:      del.Deadline,
	}
}

func domainToDTOUnassign(del *delivery.UnassignResult) DeliveryUnassignResponse {
	return DeliveryUnassignResponse{
		OrderID:   del.OrderID,
		Status:    del.Status,
		CourierID: del.CourierID,
	}
}
