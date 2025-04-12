package filesystem

import (
	"os"
	"os/exec"
)

// CommandExists looks for a command in the PATH
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	if stat, err := os.Stat(path); !os.IsNotExist(err) && !stat.IsDir() {
		return true
	}
	return false
}

// ReadFile reads the entire file at path and returns its contents
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// GetCurrentDir returns the current working directory
func GetCurrentDir() (string, error) {
	return os.Getwd()
}
