package app

import (
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/collector"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/config"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/metrics"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/pkg/dvla"
)

type App struct {
	Version string

	Config *config.Config

	Client *dvla.Client

	Collectors map[string]*collector.Collector

	Metrics *metrics.Metrics
}

func New(conf *config.Config) (*App, error) {
	client := dvla.NewClient(dvla.ClientOptions{
		Endpoint: conf.Endpoint,
		ApiKey:   conf.ApiKey,
	})

	app := &App{
		Config:     conf,
		Client:     client,
		Collectors: make(map[string]*collector.Collector),
		Metrics:    metrics.New(conf.Port),
	}

	for _, reg := range conf.Vehicles {
		app.Collectors[reg] = collector.New(app.Client, reg, conf.Interval)
	}

	return app, nil
}
