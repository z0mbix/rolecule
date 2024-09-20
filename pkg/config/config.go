package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/z0mbix/rolecule/pkg/filesystem"

	"github.com/apex/log"
	"github.com/spf13/viper"
	"github.com/z0mbix/rolecule/pkg/container"
	"github.com/z0mbix/rolecule/pkg/instance"
	"github.com/z0mbix/rolecule/pkg/provisioner"
	"github.com/z0mbix/rolecule/pkg/verifier"
)

var (
	AppName       = "rolecule"
	defaultEngine = "docker"
)

type configFile struct {
	Engine      container.EngineConfig `mapstructure:"engine"`
	Provisioner provisioner.Config     `mapstructure:"provisioner"`
	Verifier    verifier.Config        `mapstructure:"verifier"`
	Instances   []instance.Config      `mapstructure:"instances"`
}

type Config struct {
	RoleName  string
	Instances instance.Instances
	Engine    container.Engine
}

func Get() (*Config, error) {
	// config file is 'rolecule.yml|rolecule.yaml' in the current directory
	viper.SetEnvPrefix(strings.ToUpper(AppName))
	viper.SetConfigName(AppName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("tests")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Fatalf("config file not found: %s.yml", AppName)
		}
	}

	var configValues configFile
	err := viper.Unmarshal(&configValues)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config file: %v", err)
	}

	log.Debugf("config file: %+v", configValues)

	if configValues.Engine.Name == "" {
		log.Debugf("engine not specified, using default engine: %s", defaultEngine)
		configValues.Engine.Name = defaultEngine
	}
	engine, err := container.NewEngine(configValues.Engine.Name)
	if err != nil {
		return nil, err
	}

	if !filesystem.CommandExists(configValues.Engine.Name) {
		return nil, fmt.Errorf("container engine '%s' not found in PATH", configValues.Engine.Name)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// resolve any symlinks in the current working directory
	cwdNoSymlinks, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return nil, err
	}

	roleName := filepath.Base(cwd)
	roleDir := filepath.Dir(cwd)
	log.Debugf("role name: %s", roleName)
	log.Debugf("role dir: %s", roleDir)

	prov, err := provisioner.NewProvisioner(configValues.Provisioner)
	if err != nil {
		return nil, err
	}

	verif, err := verifier.NewVerifier(configValues.Verifier)
	if err != nil {
		return nil, err
	}

	// Check if the role has a meta/main.yml file to determine if it has dependencies
	roleMounts := make(map[string]string)
	if filesystem.FileExists("meta/main.yml") {
		roleMetadata, err := provisioner.GetRoleMetadata()
		if err != nil {
			return nil, err
		}

		for _, dep := range roleMetadata.LocalDependencies() {
			log.Debugf("found local role dependency: %s", dep)
			roleMounts[filepath.Join(roleDir, dep)] = filepath.Join("/etc/ansible/roles", dep)
		}

		for _, dep := range roleMetadata.GalaxyDependencies() {
			log.Debugf("found galaxy role dependency: %s", dep)
		}

		prov = prov.WithLocalDependencies(roleMetadata.LocalDependencies())
		prov = prov.WithGalaxyDependencies(roleMetadata.GalaxyDependencies())
	}

	var localRoleDependencies []string
	for _, v := range roleMounts {
		log.Debugf("adding local dependency: %s", v)
		localRoleDependencies = append(localRoleDependencies, v)
	}

	var instances instance.Instances
	for _, i := range configValues.Instances {
		iProvisioner := prov.WithTags(i.Tags).WithSkipTags(i.SkipTags)

		if i.Playbook != "" {
			iProvisioner = iProvisioner.WithPlaybook(i.Playbook)
		}

		iVerifier := verif
		if len(i.TestFile) > 0 {
			iVerifier = verif.WithTestFile(i.TestFile)
		}

		instanceConfig := instance.Instance{
			Name:        generateContainerName(i.Name, roleName),
			Image:       i.Image,
			Arch:        i.Arch,
			Args:        i.Args,
			Playbook:    i.Playbook,
			WorkDir:     cwdNoSymlinks,
			RoleName:    roleName,
			RoleDir:     roleDir,
			Engine:      engine,
			Provisioner: iProvisioner,
			Verifier:    iVerifier,
			RoleMounts:  roleMounts,
			Volumes:     i.Volumes,
		}

		instances = append(instances, instanceConfig)
	}

	cfg := &Config{
		RoleName:  roleName,
		Engine:    engine,
		Instances: instances,
	}

	return cfg, nil
}

func generateContainerName(name, roleName string) string {
	replacer := strings.NewReplacer("_", "-", " ", "-", ":", "-")
	return fmt.Sprintf("%s-%s-%s", AppName, replacer.Replace(roleName), replacer.Replace(name))
}

// Create creates a rolecule.yml file in the current directory
func Create(engine, provisioner, verifier string) error {
	// TODO: yeah, actually implement this
	log.Debugf("creating config with: %s/%s/%s", engine, provisioner, verifier)

	return nil
}

var roleculeFileTemplate = `engine:
  name: {{.Engine}}

provisioner:
  name: ansible

verifier:
  name: goss

instances:
  - name: ubuntu-24.04
    image: ubuntu-systemd:24.04
`
