package verifier

import (
	"fmt"

	"github.com/apex/log"
)

type GossVerifier struct {
	Name      string
	Command   string
	Args      []string
	ExtraArgs []string
	TestFile  string
}

func (v GossVerifier) String() string {
	return v.Name
}

func (v GossVerifier) WithTestFile(file string) Verifier {
	v.TestFile = file
	return v
}

func (v GossVerifier) GetCommand() (map[string]string, string, []string) {
	gossfilePath := fmt.Sprintf("tests/%s", v.TestFile)
	args := []string{"--gossfile", gossfilePath, "validate"}
	args = append(args, v.ExtraArgs...)
	return nil, v.Command, args
}

func (v GossVerifier) GetTestFile() string {
	return v.TestFile
}

var defaultGossConfig = GossVerifier{
	Name:     "goss",
	Command:  "goss",
	TestFile: "goss.yaml",
}

func getGossConfig(config Config) GossVerifier {
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
