package deliveryapp

import (
	"fmt"
	"time"
)

const (
	deadlineFoot    = time.Minute * 30
	deadlineScooter = time.Minute * 15
	deadlineCar     = time.Minute * 5
)

type factoryTimeCalculator struct{}

func NewFactoryTimeCalculator() factoryTimeCalculator {
	return factoryTimeCalculator{}
}

// GetDeliveryCalculator - фабрика для работы со временем доставки
func (f factoryTimeCalculator) GetDeliveryCalculator(transportType string) (TimeCalculator, error) {
	switch transportType {
	case "on_foot":
		return OnFootCalculator{}, nil
	case "scooter":
		return OnScooterCalculator{}, nil
	case "car":
		return OnCarCalculator{}, nil
	default:
		return nil, fmt.Errorf("the type of transport (%s) was not found", transportType)
	}
}

type OnFootCalculator struct{}

func (c OnFootCalculator) Calculate() time.Time {
	return time.Now().UTC().Add(deadlineFoot)
}

type OnScooterCalculator struct{}

func (c OnScooterCalculator) Calculate() time.Time {
	return time.Now().UTC().Add(deadlineScooter)
}

type OnCarCalculator struct{}

func (c OnCarCalculator) Calculate() time.Time {
	return time.Now().UTC().Add(deadlineCar)
}
