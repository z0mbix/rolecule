package verifier

type TestInfraVerifier struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string
}

func (a *TestInfraVerifier) GetCommand() (map[string]string, string, []string) {
	return a.Env, a.Command, a.Args
}
