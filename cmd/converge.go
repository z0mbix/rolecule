/*
Copyright © 2022 David Wooldridge <zombie@zombix.org>
*/
package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/config"
)

func init() {
	rootCmd.AddCommand(convergeCmd)
}

var convergeCmd = &cobra.Command{
	Use:     "converge",
	Aliases: []string{"co"},
	Short:   "Run your configuration management tool to converge the configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Get()
		if err != nil {
			return err
		}

		return converge(cfg)
	},
}

func converge(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		if !instance.Engine.Exists(instance.GetContainerName()) {
			log.Errorf("container does not exist, creating...")
			err := create(cfg)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}

		log.Infof("converging container: %s", instance.GetContainerName())
		output, err := instance.Converge()
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(output)
	}

	return nil
}
