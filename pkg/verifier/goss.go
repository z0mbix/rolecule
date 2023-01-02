package verifier

type GossVerifier struct {
	Name    string
	Command string
	Args    []string
	Env     map[string]string
}

func (a *GossVerifier) GetCommand() (map[string]string, string, []string) {
	return a.Env, a.Command, a.Args
}
