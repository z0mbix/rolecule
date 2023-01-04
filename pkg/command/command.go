package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/apex/log"
)

func convertEnvVarMapToSlice(envVars map[string]string) []string {
	var envVarSlice []string
	for k, v := range envVars {
		envVarSlice = append(envVarSlice, fmt.Sprintf("%s=%s", k, v))
	}

	return envVarSlice
}

func Execute(name string, args ...string) (int, string, error) {
	cmd := exec.Command(name, args...)
	log.Debugf("executing command: %s", cmd)
	out, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		return exitCode, string(out), fmt.Errorf("command failed: %s", err)
	}
	output := strings.TrimSuffix(string(out), "\n")
	return exitCode, output, nil
}

func ExecuteWithEnvVars(env map[string]string, name string, args ...string) (int, string, error) {
	cmd := exec.Command(name, args...)
	envVars := convertEnvVarMapToSlice(env)
	log.Debugf("executing command: %s with env vars: %+v", cmd, envVars)
	cmd.Env = envVars
	out, err := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		return exitCode, string(out), fmt.Errorf("command failed: %s", err)
	}
	output := strings.TrimSuffix(string(out), "\n")
	return exitCode, output, nil
}

func Interactive(name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	log.Debugf("executing interactive command: %s", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		return exitCode, fmt.Errorf("command failed: %s", err)
	}
	return exitCode, nil
}
