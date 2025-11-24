package worker

import (
	"context"
	"fmt"
	"time"
)

type checkDelivery interface {
	DeliveryCheck(ctx context.Context) error
}

type deliveryMonitor struct {
	checkPeriod  time.Duration
	checkService checkDelivery
}

func NewDeliveryMonitor(checkPeriod time.Duration, checkService checkDelivery) *deliveryMonitor {
	return &deliveryMonitor{
		checkPeriod:  checkPeriod,
		checkService: checkService,
	}
}

func (dm *deliveryMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(dm.checkPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dm.execute(ctx)
		case <-ctx.Done():
			dm.stop()
			return
		}
	}
}

func (dm *deliveryMonitor) execute(ctx context.Context) {
	if err := dm.checkService.DeliveryCheck(ctx); err != nil {
		fmt.Println(err.Error())
	}
}

func (dm *deliveryMonitor) stop() {
	fmt.Println("stop delivery monitoring")
}
