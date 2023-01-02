package container

import (
	"fmt"

	"github.com/z0mbix/rolecule/pkg/utils"
)

type Engine interface {
	Exec(string, map[string]string, string, []string) (string, error)
	Exists(string) bool
	List(string) (string, error)
	Remove(string) error
	Run(string, []string) (string, error)
	Shell(string) error
	String() string
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
