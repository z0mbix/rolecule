package container

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/provisioner"
	"github.com/z0mbix/rolecule/pkg/verifier"
)

type Instances []Instance

type Instance struct {
	Name    string
	Image   string
	Args    []string
	Arch    string
	WorkDir string
	Engine
	Provisioner provisioner.Provisioner
	Verifier    verifier.Verifier
}

func (i *Instance) Create() (string, error) {
	workDir := "/src"
	instanceArgs := []string{
		"run",
		"--privileged",
		"--rm",
		"--detach",
		"--volume", "/sys/fs/cgroup:/sys/fs/cgroup:ro",
		"--volume", fmt.Sprintf("%s:%s", i.WorkDir, workDir),
		"--workdir", workDir,
	}

	if i.Arch != "" {
		instanceArgs = append(instanceArgs, "--arch", i.Arch)
	}

	instanceArgs = append(instanceArgs, "--name", i.GetContainerName())

	args := append(instanceArgs, i.Args...)

	log.Debugf("%+v", args)
	output, err := i.Engine.Run(i.Image, args)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (i *Instance) Converge() (string, error) {
	env, cmd, args := i.Provisioner.GetCommand()
	log.Debugf("%s -> %v", cmd, args)
	return i.Engine.Exec(i.GetContainerName(), env, cmd, args)
}

func (i *Instance) Verify() (string, error) {
	env, cmd, args := i.Verifier.GetCommand()
	log.Debugf("%s -> %v", cmd, args)
	return i.Engine.Exec(i.GetContainerName(), env, cmd, args)
}

func (i *Instance) Shell() error {
	return i.Engine.Shell(i.GetContainerName())
}

func (i *Instance) Destroy() error {
	return i.Engine.Remove(i.Name)
}

func (i *Instance) GetContainerName() string {
	replacer := strings.NewReplacer("_", "-", " ", "-")
	containerName := fmt.Sprintf("rolecule-%s", replacer.Replace(i.Name))
	return containerName
}
