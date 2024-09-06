package container

import (
	"fmt"
)

type Engine interface {
	Exec(string, map[string]string, string, []string) error
	Exists(string) bool
	List(string) (string, error)
	Remove(string) error
	Run(string, []string) (string, error)
	Shell(string) error
	String() string
}

type EngineConfig struct {
	Name string `mapstructure:"name"`
}

func NewEngine(name string) (Engine, error) {
	switch name {
	case "docker":
		return &DockerEngine{
			Name:   "docker",
			Socket: "docker://",
		}, nil
	case "podman":
		return &PodmanEngine{
			Name:   "podman",
			Socket: "podman://",
		}, nil
	default:
		return nil, fmt.Errorf("container engine '%s' not recognised (docker and podman currently supported)", name)
	}
}
