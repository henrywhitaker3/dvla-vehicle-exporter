package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/henrywhitaker3/dvla-vehicle-exporter/cmd/root"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/app"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/config"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/logger"
)

var (
	version string
)

func main() {
	conf, err := config.Load(configPath(os.Args))
	if err != nil {
		die(err)
	}

	app, err := app.New(conf)
	if err != nil {
		die(err)
	}
	app.Version = version

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Wrap(ctx, conf.LogLevel.Level())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("recieved interrupt, cancelling...")
		cancel()

		<-sigs
		fmt.Println("Recieved second iterrupt, exiting")
		os.Exit(1)
	}()

	cmd := root.New(app)
	cmd.SetContext(ctx)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func die(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func configPath(args []string) string {
	for i, a := range args {
		if a == "--config" {
			return args[i+1]
		}
	}
	return "dvla-vehicle-exporter.yaml"
}
