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

var containerName string

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&containerName, "container-name", "n", "", "Login to a specific instance")
}

var shellCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"sh", "login"},
	Short:   "get a shell in a container",
	// Long: `to quickly create a Cobra application.`,

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

// TODO: how to support using the containerName when passed as a flag?
func shell(cfg *config.Config) error {
	if len(cfg.Instances) > 1 {
		return fmt.Errorf("more than one container, you need to specify which container with -n [container name]")
	}

	err := cfg.Instances[0].Shell()
	if err != nil {
		return err
	}

	return nil
}
