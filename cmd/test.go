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
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Create the container(s), converge them, test them, then clean up",
	Long: `"test" will create the containers defined, run the provisioner of choice
against them, test them with your verifier of choice, then destroy everything.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Get()
		if err != nil {
			return err
		}
		log.Debugf("config: %+v", cfg)

		log.Info("creating containers...")
		err = create(cfg)
		if err != nil {
			return err
		}

		log.Info("converging containers...")
		err = converge(cfg)
		if err != nil {
			return err
		}

		log.Info("verifing containers...")
		err = verify(cfg)
		if err != nil {
			return err
		}

		log.Info("destroying containers...")
		err = destroy(cfg)
		if err != nil {
			return err
		}

		return nil
	},
}
