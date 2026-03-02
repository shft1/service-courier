package deliverydb

import "github.com/shft1/service-courier/internal/domain/delivery"

func domainToRowCreate(del *delivery.AssignCreate) *deliveryCreateRow {
	return &deliveryCreateRow{
		CourierID: del.CourierID,
		OrderID:   del.OrderID,
		Deadline:  del.Deadline,
	}
}

func rowToDomainDelivery(delRow *deliveryRow) *delivery.Delivery {
	return &delivery.Delivery{
		DeliveryID: delRow.DeliveryID,
		CourierID:  delRow.CourierID,
		OrderID:    delRow.OrderID,
		AssignedAt: delRow.AssignedAt,
		Deadline:   delRow.Deadline,
	}
}
