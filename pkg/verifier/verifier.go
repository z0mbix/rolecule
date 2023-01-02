package verifier

import "fmt"

type Verifier interface {
	GetCommand() (map[string]string, string, []string)
	String() string
}

func NewVerifier(name string) (Verifier, error) {
	if name == "goss" {
		return defaultGossConfig, nil
	}

	if name == "testinfra" {
		return defaultTestInfraConfig, nil
	}

	return nil, fmt.Errorf("Verifier '%s' not recognised (only goss is currently supported)", name)
}
