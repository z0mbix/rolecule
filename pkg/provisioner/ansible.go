package provisioner

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"golang.org/x/exp/maps"
)

type AnsibleProvisioner struct {
	Name      string
	Command   string
	Args      []string
	ExtraArgs []string
	Env       map[string]string
	Playbook  string
}

func (a *AnsibleProvisioner) String() string {
	return a.Name
}

func (a *AnsibleProvisioner) GetCommand() (map[string]string, string, []string) {
	// TODO: how to handle getting the playbook path better, and support scenarios?
	playbookPath := fmt.Sprintf("tests/%s", a.Playbook)
	args := append(a.Args, playbookPath)
	return a.Env, a.Command, args
}

var defaultAnsibleConfig = &AnsibleProvisioner{
	Name:    "ansible",
	Command: "ansible-playbook",
	Args: []string{
		"--connection",
		"local",
		"--inventory",
		"localhost,",
	},
	Env: map[string]string{
		"ANSIBLE_ROLES_PATH": ".",
		"ANSIBLE_NOCOWS":     "True",
	},
	Playbook: "playbook.yml",
}

func getAnsibleConfig(config Config) *AnsibleProvisioner {
	ansibleConfig := defaultAnsibleConfig
	if config.Command != "" {
		log.Debugf("using ansible command from config file: %v", config.Command)
		ansibleConfig.Command = config.Command
	}

	if len(config.Args) > 0 {
		log.Debugf("using ansible args from config file: %v", config.Args)
		ansibleConfig.Args = config.Args
	}

	if len(config.ExtraArgs) > 0 {
		log.Debugf("using ansible extra args from config file: %v", config.ExtraArgs)
		ansibleConfig.Args = append(ansibleConfig.Args, config.ExtraArgs...)
	}

	if len(config.Env) > 0 {
		// Work around viper lowercasing map keys: https://github.com/spf13/viper/issues/373
		uppercaseEnv := make(map[string]string)
		for k, v := range config.Env {
			uppercaseEnv[strings.ToUpper(k)] = v
		}
		log.Debugf("using ansible env from config file: %v", uppercaseEnv)
		maps.Copy(ansibleConfig.Env, uppercaseEnv)
	}

	if config.Playbook != "" {
		log.Debugf("using ansible playbook from config file: %v", config.Playbook)
		ansibleConfig.Playbook = config.Playbook
	}

	return ansibleConfig
}
