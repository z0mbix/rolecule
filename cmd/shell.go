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
	shellInstance := cfg.Instances[0]

	// If we have more than one container configured, find it
	if len(cfg.Instances) > 1 {
		if shellContainerName == "" {
			var instanceNames []string
			for _, instance := range cfg.Instances {
				instanceNames = append(instanceNames, instance.Name)
			}
			return fmt.Errorf("more than one container, you need to specify which container with -n %v", instanceNames)
		}

		for _, instance := range cfg.Instances {
			if instance.Name == shellContainerName {
				shellInstance = instance
			}
		}
	}

	if !shellInstance.Exists() {
		return fmt.Errorf("container does not exist yet, you need to create it first")
	}

	err := shellInstance.Shell()
	if err != nil {
		return err
	}

	return nil
}
