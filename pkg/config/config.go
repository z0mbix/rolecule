package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/viper"
	"github.com/z0mbix/rolecule/pkg/container"
	"github.com/z0mbix/rolecule/pkg/instance"
	"github.com/z0mbix/rolecule/pkg/provisioner"
	"github.com/z0mbix/rolecule/pkg/verifier"
)

var AppName = "rolecule"

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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("config file not found: %s.yml", AppName)
		} else {
			log.Fatalf("config file not valid: %s", err)
		}
	}

	var configValues configFile
	err := viper.Unmarshal(&configValues)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config file: %v", err)
	}

	log.Debugf("config file: %+v", configValues)
	log.Debugf("config file instances: %+v", configValues.Instances)

	engine, err := container.NewEngine(configValues.Engine.Name)
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cwdNoSymlinks, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return nil, err
	}

	roleName := filepath.Base(cwd)

	prov, err := provisioner.NewProvisioner(configValues.Provisioner)
	if err != nil {
		return nil, err
	}

	verif, err := verifier.NewVerifier(configValues.Verifier)
	if err != nil {
		return nil, err
	}

	var instances instance.Instances
	for _, i := range configValues.Instances {
		iProvisioner := prov.WithTags(i.Tags)

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
			Engine:      engine,
			Provisioner: iProvisioner,
			Verifier:    iVerifier,
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
	return fmt.Sprintf("%s-%s-%s", AppName, roleName, replacer.Replace(name))
}

// Create creates a rolecule.yml file in the current directory
func Create(engine, provisioner, verifier string) error {
	// TODO: yeah, actually implement this
	log.Debugf("creating config with: %s/%s/%s", engine, provisioner, verifier)
	return nil
}
