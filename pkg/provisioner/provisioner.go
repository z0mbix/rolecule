package provisioner

import "fmt"

type Provisioner interface {
	GetCommand() (map[string]string, string, []string)
	String() string
}

func NewProvisioner(name string) (Provisioner, error) {
	if name == "ansible" {
		return defaultAnsibleConfig, nil
	}

	return nil, fmt.Errorf("provisioner '%s' not recognised (only ansible currently supported)", name)
}
