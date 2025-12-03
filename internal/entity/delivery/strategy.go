package delivery

import "time"

const (
	tFoot    = time.Minute * 30
	tScooter = time.Minute * 15
	tCar     = time.Minute * 5
)

type DeliveryTimeStrategy interface {
	TimeCalculate() time.Time
}

type FootTimeStrategy struct{}

func (fs FootTimeStrategy) TimeCalculate() time.Time {
	return time.Now().UTC().Add(tFoot)
}

type ScooterTimeStrategy struct{}

func (ss ScooterTimeStrategy) TimeCalculate() time.Time {
	return time.Now().UTC().Add(tScooter)
}

type CarTimeStrategy struct{}

func (cs CarTimeStrategy) TimeCalculate() time.Time {
	return time.Now().UTC().Add(tCar)
}
