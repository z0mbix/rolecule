package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/viper"
	"github.com/z0mbix/rolecule/pkg/container"
	"github.com/z0mbix/rolecule/pkg/provisioner"
	"github.com/z0mbix/rolecule/pkg/verifier"
)

var AppName = "rolecule"

type engineConfig struct {
	Name string `mapstructure:"name"`
}

type containerConfig struct {
	Name  string   `mapstructure:"name"`
	Image string   `mapstructure:"image"`
	Arch  string   `mapstructure:"arch"`
	Args  []string `mapstructure:"args"`
}

type provisionerConfig struct {
	Name    string            `mapstructure:"name"`
	Command string            `mapstructure:"command"`
	Env     map[string]string `mapstructure:"env"`
}

type verifierConfig struct {
	Name    string `mapstructure:"name"`
	Command string `mapstructure:"command"`
}

type configFile struct {
	Engine      engineConfig      `mapstructure:"engine"`
	Containers  []containerConfig `mapstructure:"containers"`
	Provisioner provisionerConfig `mapstructure:"provisioner"`
	Verifier    verifierConfig    `mapstructure:"verifier"`
}

type Config struct {
	WorkDir     string
	RoleName    string
	Instances   container.Instances
	EngineName  string
	Provisioner provisioner.Provisioner
	Verifier    verifier.Verifier
	Engine      container.Engine
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
		return nil, fmt.Errorf("unable to decode config file into struct, %v", err)
	}

	log.Debugf("config file: %+v", configValues)

	engine, err := container.NewEngine(configValues.Engine.Name)
	if err != nil {
		return nil, err
	}

	provisioner, err := provisioner.NewProvisioner(configValues.Provisioner.Name)
	if err != nil {
		return nil, err
	}

	verifier, err := verifier.NewVerifier(configValues.Verifier.Name)
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

	var instances container.Instances
	for _, i := range configValues.Containers {
		instanceConfig := container.Instance{
			Name:        i.Name,
			Image:       i.Image,
			Arch:        i.Arch,
			Args:        i.Args,
			WorkDir:     cwdNoSymlinks,
			Engine:      engine,
			Provisioner: provisioner,
			Verifier:    verifier,
		}

		instances = append(instances, instanceConfig)
	}

	return &Config{
		RoleName:    filepath.Base(cwd),
		WorkDir:     cwd,
		Provisioner: provisioner,
		Verifier:    verifier,
		Engine:      engine,
		Instances:   instances,
	}, nil
}

// Create creates a rolecule.yml file in the current directory
func Create(engine, provisioner, verifier string) error {
	log.Debugf("creating config with: %s/%s/%s", engine, provisioner, verifier)
	return nil
}
