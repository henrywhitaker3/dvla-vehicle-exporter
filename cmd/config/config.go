package config

import (
	"fmt"

	"github.com/henrywhitaker3/dvla-vehicle-exporter/internal/app"
	"github.com/spf13/cobra"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

func New(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print out the current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			by, err := yaml.Marshal(app.Config)
			if err != nil {
				return err
			}

			fmt.Println(string(by))

			return nil
		},
	}
}
