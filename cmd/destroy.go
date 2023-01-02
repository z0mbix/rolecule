/*
Copyright © 2023 David Wooldridge <zombie@zombix.org>
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
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = destroy(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func destroy(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		log.Infof("destroying container %s", instance.Name)
		err := instance.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}
