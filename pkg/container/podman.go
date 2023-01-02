package container

import (
	"fmt"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/command"
)

type PodmanEngine struct {
	Name   string
	Socket string
}

func (p *PodmanEngine) Run(image string, args []string) (string, error) {
	containerArgs := append(args, image)
	_, output, err := command.Execute(p.Name, containerArgs...)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (p *PodmanEngine) Exec(containerName string, envVars map[string]string, cmd string, args []string) (string, error) {
	log.Debug("executing command in container")

	execArgs := []string{
		"exec",
		"--interactive",
		"--tty",
	}

	if len(envVars) > 0 {
		for k, v := range envVars {
			execArgs = append(execArgs, "--env")
			execArgs = append(execArgs, fmt.Sprintf("%s=%s", k, v))
		}
	}

	execArgs = append(execArgs, containerName)
	execArgs = append(execArgs, cmd)
	execArgs = append(execArgs, args...)

	_, output, err := command.Execute(p.Name, execArgs...)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (p *PodmanEngine) Shell(containerName string) error {
	log.Debug("executing command in container")

	shell := "bash"

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

func (p *PodmanEngine) Remove(name string) error {
	log.Debug("removing container")
	return nil
}

func (p *PodmanEngine) Exists(name string) bool {
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

	if output == name {
		return true
	}

	return false
}
