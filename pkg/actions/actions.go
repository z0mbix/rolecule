package actions

import (
	"fmt"

	"github.com/z0mbix/rolecule/pkg/instance"
)

// Shell opens a shell in the specified container
func Shell(instances instance.Instances, containerName string) error {
	targetInstance, err := findInstance(instances, containerName)
	if err != nil {
		return err
	}

	if !targetInstance.Exists() {
		return fmt.Errorf("container does not exist yet, you need to create it first")
	}

	return targetInstance.Shell()
}

// Exec executes a command in the specified container
func Exec(instances instance.Instances, containerName string, cmd string, args []string) error {
	targetInstance, err := findInstance(instances, containerName)
	if err != nil {
		return err
	}

	if !targetInstance.Exists() {
		return fmt.Errorf("container does not exist yet, you need to create it first")
	}

	// Execute the command
	envVars := map[string]string{}                                             // Empty map or add default environment variables
	return targetInstance.Engine.Exec(targetInstance.Name, envVars, cmd, args) // Assuming you added the interactive parameter
}

// findInstance finds an instance by name, or returns the first instance if no name is specified
// or an error if multiple instances exist and no name is specified
func findInstance(instances instance.Instances, containerName string) (instance.Instance, error) {
	if len(instances) == 0 {
		return instance.Instance{}, fmt.Errorf("no containers configured")
	}

	// If only one instance, use it
	if len(instances) == 1 {
		return instances[0], nil
	}

	// Multiple instances but no name specified
	if containerName == "" {
		var instanceNames []string
		for _, instance := range instances {
			instanceNames = append(instanceNames, instance.Name)
		}
		return instance.Instance{}, fmt.Errorf("more than one container, you need to specify which container with -n %v", instanceNames)
	}

	// Find the specified instance
	for _, instance := range instances {
		if instance.Name == containerName {
			return instance, nil
		}
	}

	return instance.Instance{}, fmt.Errorf("container %s not found", containerName)
}
