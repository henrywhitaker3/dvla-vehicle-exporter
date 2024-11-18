package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/logger"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	reg                            = &sync.Once{}
	VehicleDetailsCollectionErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "dvla_vehicle_details_collection_errors_count",
		Help: "The number of errors encountered when collecting vehicle data",
	}, []string{"reg"})
	VehicleDetails = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dvla_vehicle_details",
		Help: "The details of the vehicle",
	}, []string{"reg", "co2Emissions", "colour", "engineCapacity", "fuelType", "make", "markedForExport", "monthOfFirstRegistration", "revenueWeight", "typeApproval", "wheelPlan", "yearOfManufacture", "euroStatus", "realDrivingEmissions"})
	TaxExpirySeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dvla_vehicle_details_tax_expiry_seconds",
		Help: "The number of seconds until the tax expires",
	}, []string{"reg"})
	MotExpirySeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dvla_vehicle_details_mot_expiry_seconds",
		Help: "The number of seconds until the tax expires",
	}, []string{"reg"})
	MotStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dvla_vehicle_details_tax_status",
	}, []string{"reg"})
	TaxStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dvla_vehicle_details_mot_status",
	}, []string{"reg"})
)

type Metrics struct {
	e    *echo.Echo
	port int
}

func New(port int) *Metrics {
	reg.Do(func() {
		prometheus.MustRegister(VehicleDetailsCollectionErrors)
		prometheus.MustRegister(VehicleDetails)
		prometheus.MustRegister(TaxExpirySeconds)
		prometheus.MustRegister(TaxStatus)
		prometheus.MustRegister(MotStatus)
		prometheus.MustRegister(MotExpirySeconds)
	})

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware(""))

	e.GET("/metrics", echoprometheus.NewHandler())

	return &Metrics{e: e, port: port}
}

func (m *Metrics) Start(ctx context.Context) error {
	logger.Logger(ctx).Infow("starting metrics server", "port", m.port)
	if err := m.e.Start(fmt.Sprintf(":%d", m.port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (m *Metrics) Shutdown(ctx context.Context) error {
	logger.Logger(ctx).Info("stopping metrics server")
	return m.e.Shutdown(context.Background())
}
