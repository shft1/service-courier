package courierdb

import "service-courier/internal/domain/courier"

func domainToRowCreate(cour *courier.CourierCreate) *courierCreateRow {
	return &courierCreateRow{
		Name:          cour.Name,
		Phone:         cour.Phone,
		Status:        cour.Status,
		TransportType: cour.TransportType,
	}
}

func domainToRowUpdate(cour *courier.CourierUpdate) *courierUpdateRow {
	return &courierUpdateRow{
		ID:            cour.ID,
		Name:          cour.Name,
		Phone:         cour.Phone,
		Status:        cour.Status,
		TransportType: cour.TransportType,
	}
}

func rowToDomainCourier(courRow *courierRow) *courier.Courier {
	return &courier.Courier{
		ID:            courRow.ID,
		Name:          courRow.Name,
		Phone:         courRow.Phone,
		Status:        courRow.Status,
		TransportType: courRow.TransportType,
		CreatedAt:     courRow.CreatedAt,
		UpdatedAt:     courRow.UpdatedAt,
	}
}
