package courierhttp

import "service-courier/internal/domain/courier"

func toDomainCreate(req *CourierCreateRequest) *courier.CourierCreate {
	return &courier.CourierCreate{
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}
}

func toDomainUpdate(req *CourierUpdateRequest) *courier.CourierUpdate {
	return &courier.CourierUpdate{
		ID:            req.ID,
		Name:          req.Name,
		Phone:         req.Phone,
		Status:        req.Status,
		TransportType: req.TransportType,
	}
}

func domainToDTO(cour *courier.Courier) CourierResponse {
	return CourierResponse{
		ID:            cour.ID,
		Name:          cour.Name,
		Phone:         cour.Phone,
		Status:        cour.Status,
		TransportType: cour.TransportType,
	}
}

func domainToDTOList(cours []courier.Courier) []CourierResponse {
	coursResp := make([]CourierResponse, 0, len(cours))
	for _, cour := range cours {
		coursResp = append(coursResp, domainToDTO(&cour))
	}
	return coursResp
}
