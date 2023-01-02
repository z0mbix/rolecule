package verifier

type TestInfraVerifier struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string
}

func (v *TestInfraVerifier) String() string {
	return v.Name
}

func (v *TestInfraVerifier) GetCommand() (map[string]string, string, []string) {
	return v.Env, v.Command, v.Args
}

// TODO: how to get socket and container name?
var defaultTestInfraConfig = &TestInfraVerifier{
	Name:    "testinfra",
	Command: "py.test",
	Args: []string{
		"-vv",
		"--hosts",
		"podman://foobar",
	},
}
