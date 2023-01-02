/*
Copyright © 2023 David Wooldridge <zombie@zombix.org>
*/
package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/config"
)

var (
	engineName      string
	provisionerName string
	verifierName    string
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&engineName, "engine", "e", "podman", "Specify the container engine to use (podman or docker)")
	initCmd.Flags().StringVarP(&provisionerName, "provisioner", "p", "ansible", "Specify the provisioner to use")
	initCmd.Flags().StringVarP(&verifierName, "verifier", "v", "goss", "Specify the verifier to use")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise the project with a nice new rolecule.yml file",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Create(engineName, provisionerName, verifierName)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}
