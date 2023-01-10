package verifier

import (
	"fmt"
)

type Config struct {
	Name      string   `mapstructure:"name"`
	TestFile  string   `mapstructure:"testfile"`
	ExtraArgs []string `mapstructure:"extra_args"`
}

type Verifier interface {
	GetCommand() (map[string]string, string, []string)
	GetTestFile() string
	WithTestFile(file string) Verifier
	String() string
}

func NewVerifier(config Config) (Verifier, error) {
	if config.Name == "goss" {
		return getGossConfig(config), nil
	}

	if config.Name == "testinfra" {
		return defaultTestInfraConfig, nil
	}

	return nil, fmt.Errorf("Verifier '%s' not recognised (only goss is currently supported)", config.Name)
}
