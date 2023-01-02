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
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"v"},
	Short:   "verify your container",
	// Long: `to quickly create a Cobra application.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Debugf("config: %+v", cfg)

		return verify(cfg)
	},
}

func verify(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		instance.Engine = cfg.Engine
		err := instance.Verify()
		if err != nil {
			return err
		}
	}

	return nil
}
