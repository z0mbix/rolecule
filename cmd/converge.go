/*
Copyright © 2022 David Wooldridge <zombie@zombix.org>
*/
package cmd

import (
	"fmt"

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
	Short:   "Run your configuration management tool to converge the container(s)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = converge(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func converge(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		if !instance.Engine.Exists(instance.Name) {
			log.Errorf("container does not exist, creating...")
			err := create(cfg)
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}

		log.Infof("converging container %s with %s", instance.Name, cfg.Provisioner)
		output, err := instance.Converge()
		if err != nil {
			log.Error(err.Error())
		}
		fmt.Println(output)
	}

	return nil
}
