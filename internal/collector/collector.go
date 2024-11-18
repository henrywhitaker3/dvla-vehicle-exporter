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
	logger.Infow("starting vehicle details collector", "reg", c.reg, "interval", c.interval)
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

	logger.Debugw("getting vehicle details", "reg", c.reg)

	vehicle, err := c.client.GetVehicle(ctx, c.reg)
	if err != nil {
		metrics.VehicleDetailsCollectionErrors.With(c.vehicleLabels(vehicle)).Add(1)
		logger.Errorw("failed to collect vehicle details", "error", err)
		return
	}

	logger.Debugw("got vehicle details", "reg", c.reg, "vehicle", vehicle)

	taxExpiry := time.Until(time.Time(vehicle.TaxDueDate))
	metrics.TaxExpirySeconds.With(c.vehicleLabels(vehicle)).Set(taxExpiry.Seconds())
	taxed := 0
	if vehicle.TaxStatus == "Taxed" {
		taxed = 1
	}
	metrics.TaxStatus.With(c.vehicleLabels(vehicle)).Set(float64(taxed))
	moted := 0
	if vehicle.MotStatus == "Valid" {
		moted = 1
	}
	metrics.MotStatus.With(c.vehicleLabels(vehicle)).Set(float64(moted))
	metrics.YearOfManufacture.With(c.vehicleLabels(vehicle)).Set(float64(vehicle.YearOfManufacture))
	metrics.CO2Emissions.With(c.vehicleLabels(vehicle)).Set(float64(vehicle.Co2Emissions))
	metrics.EngineCapactiy.With(c.vehicleLabels(vehicle)).Set(float64(vehicle.EngineCapacity))
}

func (c *Collector) vehicleLabels(v *dvla.Vehicle) prometheus.Labels {
	return prometheus.Labels{"reg": c.reg}
}
