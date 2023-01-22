package provisioner

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"golang.org/x/exp/maps"
)

type AnsibleLocalProvisioner struct {
	Name      string
	Command   string
	Args      []string
	ExtraArgs []string
	SkipTags  []string
	Tags      []string
	EnvVars   map[string]string
	Playbook  string
}

func (a AnsibleLocalProvisioner) String() string {
	return a.Name
}

var defaultAnsibleConfig = AnsibleLocalProvisioner{
	Name:    "ansible",
	Command: "ansible-playbook",
	Args: []string{
		"--connection",
		"local",
		"--inventory",
		"localhost,",
	},
	EnvVars: map[string]string{
		"ANSIBLE_ROLES_PATH": ".",
		"ANSIBLE_NOCOWS":     "True",
	},
	Playbook: "playbook.yml",
}

func getAnsibleConfig(config Config) AnsibleLocalProvisioner {
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

	if len(config.EnvVars) > 0 {
		// Work around viper lowercasing map keys: https://github.com/spf13/viper/issues/373
		uppercaseEnv := make(map[string]string)
		for k, v := range config.EnvVars {
			uppercaseEnv[strings.ToUpper(k)] = v
		}
		log.Debugf("using extra ansible env from config file: %v", uppercaseEnv)
		maps.Copy(ansibleConfig.EnvVars, uppercaseEnv)
	}

	if config.Playbook != "" {
		log.Debugf("using ansible playbook from config file: %v", config.Playbook)
		ansibleConfig.Playbook = config.Playbook
	}

	return ansibleConfig
}

func (a AnsibleLocalProvisioner) WithExtraArgs(args []string) Provisioner {
	a.Tags = append(a.ExtraArgs, args...)
	return a
}

func (a AnsibleLocalProvisioner) WithSkipTags(tags []string) Provisioner {
	a.SkipTags = tags
	return a
}

func (a AnsibleLocalProvisioner) WithTags(tags []string) Provisioner {
	a.Tags = tags
	return a
}

func (a AnsibleLocalProvisioner) WithPlaybook(playbook string) Provisioner {
	a.Playbook = playbook
	return a
}

func (a AnsibleLocalProvisioner) GetCommand() (map[string]string, string, []string) {
	playbookPath := fmt.Sprintf("tests/%s", a.Playbook)
	args := a.Args

	for _, tag := range a.Tags {
		args = append(args, "--tags")
		args = append(args, tag)
	}

	for _, tag := range a.SkipTags {
		args = append(args, "--skip-tags")
		args = append(args, tag)
	}

	args = append(args, playbookPath)
	return a.EnvVars, a.Command, args
}
