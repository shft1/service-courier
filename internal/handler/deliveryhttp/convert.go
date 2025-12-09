package deliveryhttp

import "service-courier/internal/domain/delivery"

func toDomainOrderID(req DeliveryOrderRequest) delivery.OrderID {
	return delivery.OrderID{OrderID: req.OrderID}
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
