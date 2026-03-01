package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/z0mbix/cliout"
)

var debugLoggingEnabled bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugLoggingEnabled, "debug", "d", false, "enable debug output")
}

var rootCmd = &cobra.Command{
	Use:   "rolecule",
	Short: "rolecule helps you test your ansible roles",
	Long: `rolecule uses docker or podman to test your
ansible roles in a systemd enabled container,
then tests them with a verifier (goss).`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cliout.Default().SetWriter(os.Stderr)

		if debugLoggingEnabled {
			cliout.SetLevel(cliout.LevelDebug)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
