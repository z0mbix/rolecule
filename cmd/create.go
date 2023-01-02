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
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr"},
	Short:   "Create a new container(s) to test the role in",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = create(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func create(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		log.Infof("creating container: %s", instance.GetContainerName())
		if instance.Engine.Exists(instance.GetContainerName()) {
			log.Errorf("container already exists!")
			continue
		}

		output, err := instance.Create()
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(output)
	}

	return nil
}
