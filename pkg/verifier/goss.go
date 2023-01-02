package verifier

type GossVerifier struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string
}

func (v *GossVerifier) String() string {
	return v.Name
}

func (v *GossVerifier) GetCommand() (map[string]string, string, []string) {
	return v.Env, v.Command, v.Args
}

var defaultGossConfig = &GossVerifier{
	Name:    "goss",
	Command: "goss",
	Args: []string{
		"--gossfile",
		"tests/goss.yaml",
		"validate",
		"--format",
		"tap",
	},
}
