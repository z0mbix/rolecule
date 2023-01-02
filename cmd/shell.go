/*
Copyright © 2023 David Wooldridge <zombie@zombix.org>
*/
package cmd

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/config"
)

var shellContainerName string

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&shellContainerName, "name", "n", "", "Login to a specific container")
}

var shellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh", "login"},
	Short:   "Open a shell in a container",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = shell(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func shell(cfg *config.Config) error {
	if len(cfg.Instances) > 1 {
		if shellContainerName == "" {
			return fmt.Errorf("more than one container, you need to specify which container with -n [container name]")
		}
		for _, instance := range cfg.Instances {
			if instance.Name == shellContainerName {
				err := instance.Shell()
				if err != nil {
					return err
				}
			}
		}
	} else {
		err := cfg.Instances[0].Shell()
		if err != nil {
			return err
		}
	}

	return nil
}
