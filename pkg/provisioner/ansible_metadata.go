package provisioner

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type RoleMetadata struct {
	Collections  []string `yaml:"collections"`
	Dependencies []struct {
		Role string `yaml:"role"`
	} `yaml:"dependencies"`
}

func (md *RoleMetadata) LocalDependencies() []string {
	var roles []string

	for _, dep := range md.Dependencies {
		// if role name does not contain a dot it is a local role
		if !strings.Contains(dep.Role, ".") {
			roles = append(roles, dep.Role)
		}
	}

	return roles
}

func (md *RoleMetadata) GalaxyDependencies() []string {
	var roles []string

	for _, dep := range md.Dependencies {
		// if role name contains a dot it is a galaxy role
		if strings.Contains(dep.Role, ".") {
			roles = append(roles, dep.Role)
		}
	}

	return roles
}

// GetRoleMetadata parses the role meta/main.yml file
// and returns a RoleMetadata struct
func GetRoleMetadata() (*RoleMetadata, error) {
	metaFile := "meta/main.yml"
	if _, err := os.Stat(metaFile); os.IsNotExist(err) {
		return &RoleMetadata{}, nil
	}

	// Open the meta/main.yml file
	file, err := os.Open(metaFile)
	if err != nil {
		return &RoleMetadata{}, err
	}
	defer file.Close()

	// Create a RoleMetadata struct to hold the parsed data
	var metadata RoleMetadata

	// Decode the YAML file into the RoleMetadata struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&metadata); err != nil {
		return &RoleMetadata{}, err
	}

	return &metadata, nil
}
