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
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"v"},
	Short:   "verify your container",
	// Long: `to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = verify(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func verify(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		output, err := instance.Verify()
		if err != nil {
			return fmt.Errorf("%w - %s", err, output)
		}
		fmt.Println(output)
	}

	return nil
}
