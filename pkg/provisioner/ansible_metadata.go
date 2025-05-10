package provisioner

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/z0mbix/rolecule/pkg/filesystem"
	"gopkg.in/yaml.v3"
)

// RoleMetadata represents the metadata of an Ansible role
type RoleMetadata struct {
	Dependencies []interface{}   `yaml:"dependencies"`
	resolvedDeps map[string]bool // To keep track of already resolved dependencies
	localPath    string          // Path to the role directory
}

// GetRoleMetadata reads the role metadata from meta/main.yml
func GetRoleMetadata() (*RoleMetadata, error) {
	metadataFile := "meta/main.yml"
	data, err := filesystem.ReadFile(metadataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata RoleMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	cwd, err := filesystem.GetCurrentDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	metadata.localPath = filepath.Dir(cwd)
	metadata.resolvedDeps = make(map[string]bool)

	return &metadata, nil
}

// AllLocalDependencies returns a complete list of all local dependencies (including nested ones)
// and a map of role mount paths for use in container setup
func (rm *RoleMetadata) AllLocalDependencies() ([]string, map[string]string, error) {
	allDeps := []string{}
	roleMounts := make(map[string]string)

	// Get direct and nested local dependencies
	localDeps := rm.LocalDependencies()

	// Generate mount paths for all dependencies
	for _, dep := range localDeps {
		srcPath := filepath.Join(rm.localPath, dep)
		dstPath := filepath.Join("/etc/ansible/roles", dep)
		roleMounts[srcPath] = dstPath
		allDeps = append(allDeps, dep)

		log.Debugf("adding mount for dependency %s: %s -> %s", dep, srcPath, dstPath)
	}

	return allDeps, roleMounts, nil
}

// LocalDependencies returns the list of local dependencies (both direct and nested)
func (rm *RoleMetadata) LocalDependencies() []string {
	var localDeps []string
	depSet := make(map[string]bool) // Use a map to avoid duplicates

	// First, collect direct dependencies
	for _, dep := range rm.parseDependencies() {
		// Check if it's a local role (not a galaxy role)
		if !strings.Contains(dep, ".") && !strings.Contains(dep, "/") {
			if !depSet[dep] {
				depSet[dep] = true
				localDeps = append(localDeps, dep)

				// Now recursively resolve dependencies of this dependency
				log.Debugf("resolving dependencies for role: %s", dep)
				nestedDeps, err := rm.resolveNestedDependencies(dep)
				if err != nil {
					log.Warnf("failed to resolve dependencies for %s: %v", dep, err)
					continue
				}

				// Add nested dependencies
				for _, nestedDep := range nestedDeps {
					if !depSet[nestedDep] {
						depSet[nestedDep] = true
						localDeps = append(localDeps, nestedDep)
						log.Debugf("added nested dependency: %s", nestedDep)
					}
				}
			}
		}
	}

	log.Debugf("all local dependencies: %v", localDeps)
	return localDeps
}

// GalaxyDependencies returns the list of galaxy dependencies (both direct and nested)
func (rm *RoleMetadata) GalaxyDependencies() []string {
	var galaxyDeps []string
	galaxyDepSet := make(map[string]bool)

	// Process direct dependencies
	for _, dep := range rm.parseDependencies() {
		// Check if it's a galaxy role (contains a dot)
		if strings.Contains(dep, ".") {
			if !galaxyDepSet[dep] {
				galaxyDepSet[dep] = true
				galaxyDeps = append(galaxyDeps, dep)
			}
		} else {
			// For local roles, we need to check their dependencies too
			nestedGalaxyDeps, err := rm.resolveNestedGalaxyDependencies(dep)
			if err != nil {
				log.Warnf("failed to resolve galaxy dependencies for %s: %v", dep, err)
				continue
			}

			// Add nested galaxy dependencies
			for _, nestedDep := range nestedGalaxyDeps {
				if !galaxyDepSet[nestedDep] {
					galaxyDepSet[nestedDep] = true
					galaxyDeps = append(galaxyDeps, nestedDep)
					log.Debugf("added nested galaxy dependency: %s", nestedDep)
				}
			}
		}
	}

	log.Debugf("all galaxy dependencies: %v", galaxyDeps)
	return galaxyDeps
}

// resolveNestedDependencies reads the dependencies of a dependency and returns all nested local dependencies
func (rm *RoleMetadata) resolveNestedDependencies(roleName string) ([]string, error) {
	var nestedDeps []string
	depSet := make(map[string]bool) // Track dependencies to avoid duplicates

	// Construct path to the dependency's metadata file
	depMetaPath := filepath.Join(rm.localPath, roleName, "meta", "main.yml")

	// Check if metadata file exists
	if !filesystem.FileExists(depMetaPath) {
		log.Debugf("no metadata file found for role %s at %s", roleName, depMetaPath)
		return nestedDeps, nil // Return empty slice, not an error
	}

	// Read and parse the dependency's metadata
	data, err := filesystem.ReadFile(depMetaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata for %s: %w", roleName, err)
	}

	var depMetadata RoleMetadata
	if err := yaml.Unmarshal(data, &depMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata for %s: %w", roleName, err)
	}

	// Set the local path for the dependency
	depMetadata.localPath = rm.localPath

	// Get dependencies of this dependency
	for _, dep := range depMetadata.parseDependencies() {
		// Only process local dependencies
		if !strings.Contains(dep, ".") && !strings.Contains(dep, "/") {
			if !depSet[dep] {
				depSet[dep] = true
				nestedDeps = append(nestedDeps, dep)
				log.Debugf("found nested dependency: %s from role: %s", dep, roleName)

				// Recursively get dependencies of this dependency
				subDeps, err := rm.resolveNestedDependencies(dep)
				if err != nil {
					log.Warnf("failed to resolve sub-dependencies for %s: %v", dep, err)
					continue
				}

				// Add sub-dependencies
				for _, subDep := range subDeps {
					if !depSet[subDep] {
						depSet[subDep] = true
						nestedDeps = append(nestedDeps, subDep)
						log.Debugf("added sub-dependency: %s from role: %s", subDep, dep)
					}
				}
			}
		}
	}

	return nestedDeps, nil
}

// resolveNestedGalaxyDependencies reads the dependencies of a dependency and returns all nested galaxy dependencies
func (rm *RoleMetadata) resolveNestedGalaxyDependencies(roleName string) ([]string, error) {
	var galaxyDeps []string
	galaxyDepSet := make(map[string]bool) // Track dependencies to avoid duplicates

	// Construct path to the dependency's metadata file
	depMetaPath := filepath.Join(rm.localPath, roleName, "meta", "main.yml")

	// Check if metadata file exists
	if !filesystem.FileExists(depMetaPath) {
		log.Debugf("no metadata file found for role %s at %s", roleName, depMetaPath)
		return galaxyDeps, nil // Return empty slice, not an error
	}

	// Read and parse the dependency's metadata
	data, err := filesystem.ReadFile(depMetaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata for %s: %w", roleName, err)
	}

	var depMetadata RoleMetadata
	if err := yaml.Unmarshal(data, &depMetadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata for %s: %w", roleName, err)
	}

	// Set the local path for the dependency
	depMetadata.localPath = rm.localPath

	// Get dependencies of this dependency
	for _, dep := range depMetadata.parseDependencies() {
		// Process galaxy dependencies
		if strings.Contains(dep, ".") || strings.Contains(dep, "/") {
			if !galaxyDepSet[dep] {
				galaxyDepSet[dep] = true
				galaxyDeps = append(galaxyDeps, dep)
				log.Debugf("found nested galaxy dependency: %s from role: %s", dep, roleName)
			}
		} else {
			// For local roles, check their galaxy dependencies
			nestedGalaxyDeps, err := rm.resolveNestedGalaxyDependencies(dep)
			if err != nil {
				log.Warnf("failed to resolve nested galaxy dependencies for %s: %v", dep, err)
				continue
			}

			// Add nested galaxy dependencies
			for _, nestedDep := range nestedGalaxyDeps {
				if !galaxyDepSet[nestedDep] {
					galaxyDepSet[nestedDep] = true
					galaxyDeps = append(galaxyDeps, nestedDep)
					log.Debugf("added nested galaxy dependency: %s from role: %s", nestedDep, dep)
				}
			}
		}
	}

	return galaxyDeps, nil
}

// parseDependencies parses the dependencies field and returns a list of role names
func (rm *RoleMetadata) parseDependencies() []string {
	var dependencies []string

	for _, dep := range rm.Dependencies {
		switch v := dep.(type) {
		case string:
			dependencies = append(dependencies, v)
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok {
				dependencies = append(dependencies, name)
			} else if role, ok := v["role"].(string); ok {
				dependencies = append(dependencies, role)
			}
		}
	}

	return dependencies
}
