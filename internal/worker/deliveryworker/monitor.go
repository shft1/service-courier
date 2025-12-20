package deliveryworker

import (
	"context"
	"fmt"
	"time"
)

type deliveryChecker interface {
	CheckDelivery(ctx context.Context) error
}

type deliveryMonitor struct {
	checkPeriod time.Duration
	checker     deliveryChecker
}

func NewDeliveryMonitor(checkPeriod time.Duration, checker deliveryChecker) *deliveryMonitor {
	return &deliveryMonitor{
		checkPeriod: checkPeriod,
		checker:     checker,
	}
}

// Start - запуск фонового воркера проверки доставок
func (dm *deliveryMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(dm.checkPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dm.execute(ctx)
			fmt.Println("checked delivery")
		case <-ctx.Done():
			dm.stop()
			return
		}
	}
}

func (dm *deliveryMonitor) execute(ctx context.Context) {
	if err := dm.checker.CheckDelivery(ctx); err != nil {
		fmt.Println(err.Error())
	}
}

func (dm *deliveryMonitor) stop() {
	fmt.Println("stop delivery monitoring")
}
