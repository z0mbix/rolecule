package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/z0mbix/rolecule/pkg/instance"
)

// Mock implementation of Engine for testing
type mockEngine struct{}

func (m *mockEngine) Run(image string, args []string) (string, error) { return "", nil }
func (m *mockEngine) Exec(containerName string, envVars map[string]string, cmd string, args []string) error {
	return nil
}
func (m *mockEngine) Shell(containerName string) error { return nil }
func (m *mockEngine) Remove(name string) error         { return nil }
func (m *mockEngine) Exists(name string) bool          { return true }
func (m *mockEngine) List(name string) (string, error) { return "", nil }
func (m *mockEngine) String() string                   { return "mock" }

func TestFindInstance(t *testing.T) {
	// Setup mock engine
	mockEng := &mockEngine{}

	// Test cases
	tests := []struct {
		name           string
		instances      instance.Instances
		containerName  string
		expectedError  bool
		expectedErrMsg string
		expectedName   string
	}{
		{
			name:           "no instances",
			instances:      instance.Instances{},
			containerName:  "",
			expectedError:  true,
			expectedErrMsg: "no containers configured",
		},
		{
			name: "single instance, no name specified",
			instances: instance.Instances{
				{Name: "rolecule-test-ubuntu", Engine: mockEng},
			},
			containerName: "",
			expectedError: false,
			expectedName:  "rolecule-test-ubuntu",
		},
		{
			name: "single instance, name matches",
			instances: instance.Instances{
				{Name: "rolecule-test-ubuntu", Engine: mockEng},
			},
			containerName: "rolecule-test-ubuntu",
			expectedError: false,
			expectedName:  "rolecule-test-ubuntu",
		},
		{
			name: "multiple instances, no name specified",
			instances: instance.Instances{
				{Name: "rolecule-test-ubuntu", Engine: mockEng},
				{Name: "rolecule-test-centos", Engine: mockEng},
			},
			containerName:  "",
			expectedError:  true,
			expectedErrMsg: "more than one container",
		},
		{
			name: "multiple instances, name matches first",
			instances: instance.Instances{
				{Name: "rolecule-test-ubuntu", Engine: mockEng},
				{Name: "rolecule-test-centos", Engine: mockEng},
			},
			containerName: "rolecule-test-ubuntu",
			expectedError: false,
			expectedName:  "rolecule-test-ubuntu",
		},
		{
			name: "multiple instances, name matches second",
			instances: instance.Instances{
				{Name: "rolecule-test-ubuntu", Engine: mockEng},
				{Name: "rolecule-test-centos", Engine: mockEng},
			},
			containerName: "rolecule-test-centos",
			expectedError: false,
			expectedName:  "rolecule-test-centos",
		},
		{
			name: "multiple instances, name doesn't match any",
			instances: instance.Instances{
				{Name: "rolecule-test-ubuntu", Engine: mockEng},
				{Name: "rolecule-test-centos", Engine: mockEng},
			},
			containerName:  "non-existent",
			expectedError:  true,
			expectedErrMsg: "container non-existent not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, err := findInstance(tt.instances, tt.containerName)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				// Don't access instance if an error was expected
				return
			}

			// Only check instance properties if no error was expected
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedName, instance.Name)
		})
	}
}
