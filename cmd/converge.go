package cmd

import (
	"github.com/spf13/cobra"
	"github.com/z0mbix/cliout"
	"github.com/z0mbix/rolecule/pkg/config"
)

func init() {
	rootCmd.AddCommand(convergeCmd)
}

var convergeCmd = &cobra.Command{
	Use:     "converge",
	Aliases: []string{"co"},
	Short:   "Run your configuration management tool to converge the container(s)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			cliout.Fatal(err.Error())
		}

		err = converge(cfg)
		if err != nil {
			cliout.Fatal(err.Error())
		}
	},
}

func converge(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		if !instance.Engine.Exists(instance.Name) {
			err := create(cfg)
			if err != nil {
				cliout.Error(err.Error())
				continue
			}
		}

		if len(instance.Provisioner.GetDependencies().GalaxyRoles) > 0 {
			cliout.Infof("preparing container %s", instance.Name)
			if err := instance.Prepare(); err != nil {
				cliout.Error(err.Error())
			}
		}

		cliout.Infof("converging container %s with %s", instance.Name, instance.Provisioner)
		if err := instance.Converge(); err != nil {
			cliout.Error(err.Error())
		}
	}

	return nil
}
