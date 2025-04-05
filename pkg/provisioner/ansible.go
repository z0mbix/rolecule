package provisioner

import (
	"strings"

	"github.com/apex/log"
	"golang.org/x/exp/maps"
)

type Dependencies struct {
	Collections []string
	LocalRoles  []string
	GalaxyRoles []string
}

type AnsibleLocalProvisioner struct {
	Name         string
	Command      string
	Args         []string
	ExtraArgs    []string
	SkipTags     []string
	Tags         []string
	EnvVars      map[string]string
	Playbook     string
	Dependencies Dependencies
}

var ansibleRoleDir = "/etc/ansible/roles"

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
		"ANSIBLE_ROLES_PATH": ansibleRoleDir,
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
	a.ExtraArgs = append(a.ExtraArgs, args...)
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

func (a AnsibleLocalProvisioner) WithLocalDependencies(dependencies []string) Provisioner {
	a.Dependencies.LocalRoles = dependencies
	return a
}

func (a AnsibleLocalProvisioner) WithGalaxyDependencies(dependencies []string) Provisioner {
	a.Dependencies.GalaxyRoles = dependencies
	return a
}

func (a AnsibleLocalProvisioner) GetDependencies() Dependencies {
	return a.Dependencies
}

func (a AnsibleLocalProvisioner) GetInstallDependenciesCommand() (map[string]string, string, []string) {
	log.Debugf("installing galaxy role(s):")
	args := []string{"install", "--roles-path", ansibleRoleDir}
	args = append(args, a.Dependencies.GalaxyRoles...)

	return a.EnvVars, "ansible-galaxy", args
}

func (a AnsibleLocalProvisioner) GetCommand() (map[string]string, string, []string) {
	testPlaybook := []string{testDirectory, a.Playbook}
	playbookPath := strings.Join(testPlaybook, "/")
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

func (a AnsibleLocalProvisioner) String() string {
	return a.Name
}
