package verifier

import (
	"fmt"

	"github.com/apex/log"
)

type Config struct {
	Name      string   `mapstructure:"name"`
	TestFile  string   `mapstructure:"testfile"`
	ExtraArgs []string `mapstructure:"extra_args"`
}

type Verifier interface {
	GetCommand() (map[string]string, string, []string)
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

func getGossConfig(config Config) *GossVerifier {
	gossConfig := defaultGossConfig
	if config.TestFile != "" {
		log.Debugf("using gossfile from config file: %v", config.TestFile)
		gossConfig.TestFile = config.TestFile
	}
	if len(config.ExtraArgs) > 0 {
		log.Debugf("using goss extra args from config file: %v", config.ExtraArgs)
		gossConfig.ExtraArgs = config.ExtraArgs
	}

	return gossConfig
}
