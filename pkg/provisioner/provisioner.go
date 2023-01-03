package provisioner

import (
	"fmt"
)

type Provisioner interface {
	GetCommand() (map[string]string, string, []string)
	String() string
}

type Config struct {
	Name      string            `mapstructure:"name"`
	Command   string            `mapstructure:"command"`
	Args      []string          `mapstructure:"args"`
	ExtraArgs []string          `mapstructure:"extra_args"`
	Env       map[string]string `mapstructure:"env"`
	Playbook  string            `mapstructure:"playbook"`
}

func NewProvisioner(config Config) (Provisioner, error) {
	if config.Name == "ansible" {
		return getAnsibleConfig(config), nil
	}

	return nil, fmt.Errorf("provisioner '%s' not recognised (only ansible currently supported)", config.Name)
}
