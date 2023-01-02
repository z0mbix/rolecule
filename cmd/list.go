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
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"v"},
	Short:   "list the containers",
	// Long: `to quickly create a Cobra application.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Debugf("config: %+v", cfg)

		return list(cfg)
	},
}

func list(cfg *config.Config) error {
	// for _, instance := range cfg.Instances {
	// 	instance.Engine = cfg.Engine
	// 	err := instance.List()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
