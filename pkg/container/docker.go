package container

import (
	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/command"
)

type DockerEngine struct {
	Name   string
	Socket string
}

func (p *DockerEngine) Run(image string, args []string) (string, error) {
	// args := []string{
	// 	"run",
	// 	"--tty",
	// 	"--interactive",
	// 	"--rm",
	// 	"--detach",
	// 	"--volume /sys/fs/cgroup:/sys/fs/cgroup:ro",
	// 	"--name rocky-systemd-9.1",
	// 	"localhost/rockylinux-9.1-systemd:latest",
	// }
	containerArgs := append(args, image)
	_, output, err := command.Execute(p.Name, containerArgs...)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (p *DockerEngine) Exec(containerName string, envVars map[string]string, cmd string, args []string) (string, error) {
	log.Debug("executing command in container")
	return "", nil
}

func (p *DockerEngine) Shell(name string) error {
	log.Debug("logging in to container")
	return nil
}
func (p *DockerEngine) Remove(name string) error {
	log.Debug("removing container")
	return nil
}

func (p *DockerEngine) Exists(name string) bool {
	log.Debug("checking if container already exist")
	return false
}
