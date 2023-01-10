package verifier

type TestInfraVerifier struct {
	Name     string
	Command  string
	Args     []string
	EnvVars  map[string]string
	TestFile string
}

func (v TestInfraVerifier) String() string {
	return v.Name
}

func (v TestInfraVerifier) WithTestFile(file string) Verifier {
	v.TestFile = file
	return v
}

func (v TestInfraVerifier) GetCommand() (map[string]string, string, []string) {
	return v.EnvVars, v.Command, v.Args
}

func (v TestInfraVerifier) GetTestFile() string {
	return v.TestFile
}

// TODO: how to get socket and container name?
var defaultTestInfraConfig = TestInfraVerifier{
	Name:    "testinfra",
	Command: "py.test",
	Args: []string{
		"-vv",
		"--hosts",
		"podman://foobar",
	},
	TestFile: "test.py",
}
