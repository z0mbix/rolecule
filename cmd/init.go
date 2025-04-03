package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/config"
)

var (
	engineName string
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&engineName, "engine", "e", "docker", "Specify the container engine to use (docker or podman)")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise the project with a tests directory and a rolecule.yml file",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Create(engineName)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}
