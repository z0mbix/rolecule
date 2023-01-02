package container

import (
	"fmt"

	"github.com/z0mbix/rolecule/pkg/utils"
)

type Engine interface {
	Run(string, []string) (string, error)
	Exec(string, string, []string) (string, error)
	Exists(string) bool
	Shell(string) error
	Remove(string) error
}

func NewEngine(name string) (Engine, error) {
	if !utils.CommandExists(name) {
		return nil, fmt.Errorf("container engine '%s' not found in PATH", name)
	}

	if name == "docker" {
		return &DockerEngine{
			Name:   "docker",
			Socket: "docker://",
		}, nil
	}

	if name == "podman" {
		return &PodmanEngine{
			Name:   "podman",
			Socket: "podman://",
		}, nil
	}

	return nil, fmt.Errorf("container engine '%s' not recognised (docker and podman currently supported)", name)
}
