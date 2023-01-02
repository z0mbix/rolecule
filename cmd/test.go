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

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = create(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = converge(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = verify(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = destroy(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}
