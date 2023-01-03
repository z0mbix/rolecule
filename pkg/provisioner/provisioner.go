package provisioner

import (
	"fmt"

	"github.com/apex/log"
)

type Provisioner interface {
	GetCommand() (map[string]string, string, []string)
	String() string
}

type Config struct {
	Name    string            `mapstructure:"name"`
	Command string            `mapstructure:"command"`
	Args    []string          `mapstructure:"args"`
	Env     map[string]string `mapstructure:"env"`
}

func NewProvisioner(config Config) (Provisioner, error) {
	if config.Name == "ansible" {
		ansibleConfig := defaultAnsibleConfig
		if config.Command != "" {
			log.Debugf("using ansible command from config file: %v", config.Command)
			ansibleConfig.Command = config.Command
		}
		if len(config.Args) > 0 {
			log.Debugf("using ansible args from config file: %v", config.Args)
			ansibleConfig.Args = config.Args
		}
		if len(config.Env) > 0 {
			log.Debugf("using ansible env from config file: %v", config.Env)
			ansibleConfig.Env = config.Env
		}

		return ansibleConfig, nil
	}

	return nil, fmt.Errorf("provisioner '%s' not recognised (only ansible currently supported)", config.Name)
}
