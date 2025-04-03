package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/actions"
	"github.com/z0mbix/rolecule/pkg/config"
)

var execContainerName string

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&execContainerName, "name", "n", "", "Execute command in a specific container")
}

var execCmd = &cobra.Command{
	Use:     "exec [command]",
	Aliases: []string{"e", "execute"},
	Short:   "Execute a command in a container",
	Long: `Execute a command in a running container without opening an interactive shell.

Example:
  rolecule exec systemctl status
  rolecule exec -n rolecule-sshd-ubuntu-22.04 apt-get -y update`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = actions.Exec(cfg.Instances, execContainerName, args[0], args[1:])
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}
