package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/z0mbix/rolecule/pkg/config"
)

func init() {
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"v"},
	Short:   "Verify your containers are configured how you expect",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Get()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = verify(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func verify(cfg *config.Config) error {
	for _, instance := range cfg.Instances {
		log.Infof("verifying container %s with %s (%s)", instance.Name, instance.Verifier, instance.Verifier.GetTestFile())
		if err := instance.Verify(); err != nil {
			return err
		}
	}

	return nil
}
