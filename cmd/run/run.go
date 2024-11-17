package run

import (
	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/app"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			go app.Metrics.Start(cmd.Context())

			for _, c := range app.Collectors {
				c.CollectVehicleDetails(cmd.Context())
			}

			<-cmd.Context().Done()

			return app.Metrics.Shutdown(cmd.Context())
		},
	}
}
