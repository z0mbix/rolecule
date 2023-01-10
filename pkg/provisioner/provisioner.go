package provisioner

import (
	"fmt"
)

type Provisioner interface {
	GetCommand() (map[string]string, string, []string)
	WithTags([]string) Provisioner
	WithPlaybook(string) Provisioner
	String() string
}

type Config struct {
	Name      string            `mapstructure:"name"`
	Command   string            `mapstructure:"command"`
	Args      []string          `mapstructure:"args"`
	ExtraArgs []string          `mapstructure:"extra_args"`
	EnvVars   map[string]string `mapstructure:"env"`
	Playbook  string            `mapstructure:"playbook"`
}

func NewProvisioner(config Config) (Provisioner, error) {
	if config.Name == "ansible" {
		return getAnsibleConfig(config), nil
	}

	return nil, fmt.Errorf("provisioner '%s' not recognised", config.Name)
}
