package verifier

type GossVerifier struct {
	Name    string
	Command string
	Args    []string
}

func (a *GossVerifier) GetCommand() []string {
	cmdSlice := []string{a.Command}
	return append(cmdSlice, a.Args...)
}
