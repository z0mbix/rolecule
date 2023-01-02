package provisioner

type AnsibleProvisioner struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string
}

func (a *AnsibleProvisioner) String() string {
	return a.Name
}

func (a *AnsibleProvisioner) GetCommand() (map[string]string, string, []string) {
	return a.Env, a.Command, a.Args
}

var defaultAnsibleConfig = &AnsibleProvisioner{
	Name:    "ansible",
	Command: "ansible-playbook",
	Args: []string{
		"--connection",
		"local",
		"--inventory",
		"localhost,",
		"--diff",
		"--verbose",
		"tests/playbook.yml",
	},
	Env: map[string]string{
		"ANSIBLE_ROLES_PATH": ".",
		"ANSIBLE_NOCOWS":     "True",
	},
}
