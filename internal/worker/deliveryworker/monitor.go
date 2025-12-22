package deliveryworker

import (
	"context"
	"time"

	"service-courier/observability/logger"
)

type deliveryChecker interface {
	CheckDelivery(ctx context.Context) error
}

type deliveryMonitor struct {
	log logger.Logger
	period time.Duration
	checker deliveryChecker
}

func NewDeliveryMonitor(log logger.Logger, period time.Duration, checker deliveryChecker) *deliveryMonitor {
	return &deliveryMonitor{
		log: log,
		period: period,
		checker:     checker,
	}
}

// Start - запуск фонового воркера проверки доставок
func (dm *deliveryMonitor) Start(ctx context.Context) {
	ticker := time.NewTicker(dm.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dm.execute(ctx)
			dm.log.Info("checked delivery")
		case <-ctx.Done():
			dm.stop()
			return
		}
	}
}

func (dm *deliveryMonitor) execute(ctx context.Context) {
	if err := dm.checker.CheckDelivery(ctx); err != nil {
		dm.log.Error("failed to check delivery", logger.NewField("error", err))
	}
}

func (dm *deliveryMonitor) stop() {
	dm.log.Info("stop delivery monitoring")
}
