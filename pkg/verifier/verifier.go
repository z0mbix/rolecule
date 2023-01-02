package verifier

import "fmt"

type Verifier interface {
	GetCommand() []string
}

func NewVerifier(name string) (Verifier, error) {
	if name == "goss" {
		return &GossVerifier{
			Name:    "goss",
			Command: "goss",
			Args: []string{
				"validate",
			},
		}, nil
	}

	if name == "testinfra" {
		// TODO: how to get socket and container name?
		containerName := "rolecule-rockylinux-systemd-9.1"
		engineName := "podman"

		return &GossVerifier{
			Name:    "testinfra",
			Command: "py.test",
			Args: []string{
				"-vv",
				"--hosts",
				fmt.Sprintf("%s://%s", engineName, containerName),
			},
		}, nil
	}

	return nil, fmt.Errorf("Verifier '%s' not recognised (only godd is currently supported)", name)
}
