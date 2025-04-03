package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/actions"
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
	Example: `  rolecule shell
  rolecule shell -n rolecule-sshd-ubuntu-22.04`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = actions.Shell(cfg.Instances, shellContainerName)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}
