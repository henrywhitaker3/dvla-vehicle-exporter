package root

import (
	"github.com/henrywhitaker3/dvla-vehicle-exporter/cmd/config"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/cmd/run"
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/app"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dvla-vehicle-exporter",
		Short:   "Prometheus exporter for DVLA vehicle information",
		Version: app.Version,
	}

	cmd.PersistentFlags().StringVar(&configPath, "config", "dvla-vehicle-exporter.yaml", "The path to the config file")

	cmd.AddCommand(run.New(app))
	cmd.AddCommand(config.New(app))

	return cmd
}
