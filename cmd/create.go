package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/config"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"cr"},
	Short:   "Create a new container(s) to test the role in",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = create(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func create(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		if instance.Engine.Exists(instance.Name) {
			log.Infof("container %s already exists!", instance.Name)
			continue
		}

		log.Infof("creating container %s with %s", instance.Name, instance.Engine)
		output, err := instance.Create()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
		log.Debug(output)
	}

	return nil
}
