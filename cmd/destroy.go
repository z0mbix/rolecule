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
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:     "destroy",
	Aliases: []string{"rm"},
	Short:   "Destroy everything",
	Long:    `Destroy the containers created for these tests`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		return destroy(cfg)
	},
}

func destroy(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		log.Infof("destroying container: %s", instance.GetContainerName())
		err := instance.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}
