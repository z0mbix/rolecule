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

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Short:   "List the running containers for this role/module/recipe",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(list(cfg))
	},
}

func list(cfg *config.Config) string {
	namePrefix := fmt.Sprintf("%s-%s", config.AppName, cfg.RoleName)
	output, err := cfg.Engine.List(namePrefix)
	if err != nil {
		return ""
	}

	return output
}
