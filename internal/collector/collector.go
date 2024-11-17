package collector

import (
	"context"
	"time"

	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/logger"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/metrics"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/pkg/dvla"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	client   *dvla.Client
	interval time.Duration
	reg      string
}

func New(c *dvla.Client, reg string, interval time.Duration) *Collector {
	return &Collector{
		client:   c,
		interval: interval,
		reg:      reg,
	}
}

func (c *Collector) CollectVehicleDetails(ctx context.Context) {
	logger := logger.Logger(ctx)
	logger.Infow("starting vehicle details collector", "reg", c.reg)
	tick := time.NewTicker(c.interval)
	defer tick.Stop()

	c.collectVehicleDetails(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.Infow("stopping collecting vehicle details", "reg", c.reg)
			return
		case <-tick.C:
			c.collectVehicleDetails(ctx)
		}
	}
}

func (c *Collector) collectVehicleDetails(ctx context.Context) {
	logger := logger.Logger(ctx)

	vehicle, err := c.client.GetVehicle(ctx, c.reg)
	if err != nil {
		metrics.VehicleDetailsCollectionErrors.With(prometheus.Labels{"reg": c.reg}).Add(1)
		logger.Errorw("failed to collect vehicle details", "error", err)
		return
	}
	logger.Info(vehicle)
}
