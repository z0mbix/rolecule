package container

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/command"
)

type PodmanEngine struct {
	Name   string
	Socket string
}

func (p *PodmanEngine) String() string {
	return p.Name
}

func (p *PodmanEngine) Run(image string, args []string) (string, error) {
	containerArgs := append(args, image)
	_, output, err := command.Execute(p.Name, containerArgs...)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (p *PodmanEngine) Exec(containerName string, envVars map[string]string, cmd string, args []string) error {
	execArgs := []string{"exec", "--interactive", "--tty"}

	if len(envVars) > 0 {
		for k, v := range envVars {
			execArgs = append(execArgs, "--env")
			execArgs = append(execArgs, fmt.Sprintf("%s=%s", k, os.ExpandEnv(v)))
		}
	}

	execArgs = append(execArgs, containerName)
	execArgs = append(execArgs, cmd)
	execArgs = append(execArgs, args...)

	_, err := command.Interactive(p.Name, execArgs...)
	return err
}

func (p *PodmanEngine) Shell(containerName string) error {
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

func (p *PodmanEngine) Remove(name string) error {
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

	return output == name
}

func (p *PodmanEngine) List(name string) (string, error) {
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
