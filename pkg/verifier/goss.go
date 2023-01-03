package verifier

import "path/filepath"

type GossVerifier struct {
	Name      string
	Command   string
	Args      []string
	ExtraArgs []string
	TestFile  string
}

func (v *GossVerifier) String() string {
	return v.Name
}

func (v *GossVerifier) GetCommand() (map[string]string, string, []string) {
	// TODO: how to handle getting the gossfile path better, and support scenarios?
	gossfilePath := filepath.Join("tests", v.TestFile)
	args := []string{"--gossfile", gossfilePath, "validate"}
	args = append(args, v.ExtraArgs...)
	return nil, v.Command, args
}

var defaultGossConfig = &GossVerifier{
	Name:     "goss",
	Command:  "goss",
	TestFile: "goss.yaml",
}
