package container

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/command"
)

type DockerEngine struct {
	Name   string
	Socket string
}

func (p *DockerEngine) String() string {
	return p.Name
}

func (p *DockerEngine) Run(image string, args []string) (string, error) {
	containerArgs := append(args, image)
	_, output, err := command.Execute(p.Name, containerArgs...)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (p *DockerEngine) Exec(containerName string, envVars map[string]string, cmd string, args []string) error {
	execArgs := []string{
		"exec",
	}

	for k, v := range envVars {
		execArgs = append(execArgs, "--env")
		execArgs = append(execArgs, fmt.Sprintf("%s=%s", k, v))
	}

	execArgs = append(execArgs, containerName)
	execArgs = append(execArgs, cmd)
	execArgs = append(execArgs, args...)

	_, err := command.Interactive(p.Name, execArgs...)
	return err
}

func (p *DockerEngine) Shell(containerName string) error {
	shell := "bash"
	log.Debugf("opening %s shell in container", shell)

	args := []string{
		"exec",
		"--interactive",
		"--tty",
		containerName,
		shell,
	}

	_, err := command.Interactive(p.Name, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p *DockerEngine) Remove(name string) error {
	log.Debug("removing container")

	args := []string{
		"rm",
		"--force",
		name,
	}

	exitCode, _, err := command.Execute(p.Name, args...)
	if err != nil || exitCode != 0 {
		return err
	}

	return nil
}

func (p *DockerEngine) Exists(name string) bool {
	log.Debug("checking if container already exists")

	args := []string{
		"container",
		"inspect",
		"--format",
		"{{.Name}}",
		name,
	}

	exitCode, output, err := command.Execute(p.Name, args...)
	if err != nil || exitCode != 0 {
		return false
	}

	// docker returns the container name with a forward slash prefix :(
	trimmedOutput := strings.TrimPrefix(output, "/")
	return trimmedOutput == name
}

func (p *DockerEngine) List(name string) (string, error) {
	args := []string{
		"ps",
		"--filter",
		fmt.Sprintf("name=%s", name),
	}

	exitCode, output, err := command.Execute(p.Name, args...)
	if err != nil || exitCode != 0 {
		return "", err
	}

	return output, nil
}
