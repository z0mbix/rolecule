/*
Copyright © 2023 David Wooldridge <zombie@zombix.org>
*/
package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

var debugLoggingEnabled bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugLoggingEnabled, "debug", "d", false, "enable debug output")
}

var rootCmd = &cobra.Command{
	Use:   "rolecule",
	Short: "rolecule helps you test your ansible roles",
	Long: `rolecule uses docker or podman to test your
configuration management roles/recipes/modules in a systemd enabled container,
then tests them with a verifier (goss/testinfra).`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetHandler(cli.New(os.Stderr))

		if debugLoggingEnabled {
			log.SetLevel(log.DebugLevel)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
