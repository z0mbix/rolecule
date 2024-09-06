package instance

import (
	"fmt"
	"path/filepath"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/container"
	"github.com/z0mbix/rolecule/pkg/provisioner"
	"github.com/z0mbix/rolecule/pkg/verifier"
)

type Instances []Instance

type Config struct {
	Name     string   `mapstructure:"name"`
	Image    string   `mapstructure:"image"`
	Arch     string   `mapstructure:"arch"`
	Args     []string `mapstructure:"args"`
	Playbook string   `mapstructure:"playbook"`
	TestFile string   `mapstructure:"testfile"`
	SkipTags []string `mapstructure:"skip_tags"`
	Tags     []string `mapstructure:"tags"`
}

type Instance struct {
	Name       string
	Image      string
	Arch       string
	Args       []string
	Playbook   string
	TestFile   string
	SkipTags   []string
	Tags       []string
	WorkDir    string
	RoleName   string
	RoleDir    string
	RoleMounts map[string]string
	container.Engine
	Provisioner provisioner.Provisioner
	Verifier    verifier.Verifier
}

func (i *Instance) Create() (string, error) {
	workDir := filepath.Join("/etc/ansible/roles", i.RoleName)
	instanceArgs := []string{
		"run",
		"--privileged",
		"--rm",
		"--detach",
		"--tmpfs", "/tmp",
		"--tmpfs", "/run",
		"--tmpfs", "/run/lock",
		"--tmpfs", "/var/lib/docker",
		"--cgroupns", "host",
		"--workdir", workDir,
		"--volume", "/sys/fs/cgroup:/sys/fs/cgroup:rw",
		"--volume", fmt.Sprintf("%s:%s", i.WorkDir, workDir),
	}

	for src, dst := range i.RoleMounts {
		instanceArgs = append(instanceArgs, "--volume", fmt.Sprintf("%s:%s", src, dst))
	}

	if i.Arch != "" {
		instanceArgs = append(instanceArgs, "--platform", fmt.Sprintf("linux/%s", i.Arch))
	}

	instanceArgs = append(instanceArgs, "--name", i.Name)

	args := append(instanceArgs, i.Args...)

	log.Debugf("%+v", args)
	output, err := i.Run(i.Image, args)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (i *Instance) Prepare() error {
	env, cmd, args := i.Provisioner.GetInstallDependenciesCommand()
	return i.Exec(i.Name, env, cmd, args)
}

func (i *Instance) Converge() error {
	env, cmd, args := i.Provisioner.GetCommand()
	return i.Exec(i.Name, env, cmd, args)
}

func (i *Instance) Verify() error {
	env, cmd, args := i.Verifier.GetCommand()
	return i.Exec(i.Name, env, cmd, args)
}

func (i *Instance) Shell() error {
	return i.Engine.Shell(i.Name)
}

func (i *Instance) Exists() bool {
	return i.Engine.Exists(i.Name)
}

func (i *Instance) Destroy() error {
	return i.Remove(i.Name)
}
