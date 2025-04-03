package config

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func Test_generateContainerName(t *testing.T) {
	type args struct {
		name     string
		roleName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ubuntu-22.04",
			args: args{
				name:     "ubuntu-22.04",
				roleName: "foobar",
			},
			want: "rolecule-foobar-ubuntu-22.04",
		},
		{
			name: "arch",
			args: args{
				name:     "arch",
				roleName: "i_use_arch_btw",
			},
			want: "rolecule-i-use-arch-btw-arch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateContainerName(tt.args.name, tt.args.roleName); got != tt.want {
				t.Errorf("generateContainerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	// Create memory handler for log testing
	handler := memory.New()
	log.SetHandler(handler)

	tests := []struct {
		name              string
		engine            string
		wantEngineSection bool
		setupFs           func(afero.Fs) error
		wantErr           bool
		wantErrMsg        string
	}{
		{
			name:              "with docker engine",
			engine:            "docker",
			wantEngineSection: false,
			setupFs:           func(fs afero.Fs) error { return nil },
			wantErr:           false,
		},
		{
			name:              "with podman engine",
			engine:            "podman",
			wantEngineSection: true,
			setupFs:           func(fs afero.Fs) error { return nil },
			wantErr:           false,
		},
		{
			name:              "with existing file",
			engine:            "docker",
			wantEngineSection: false,
			setupFs: func(fs afero.Fs) error {
				if err := fs.MkdirAll(testsDir, 0755); err != nil {
					return err
				}
				return afero.WriteFile(fs, filepath.Join(testsDir, "rolecule.yml"), []byte("existing content"), 0644)
			},
			wantErr:    true,
			wantErrMsg: "tests/rolecule.yml file already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset log entries
			handler.Entries = nil

			// Save original filesystem and restore after test
			originalFs := appFs
			defer func() { appFs = originalFs }()

			// Create a new in-memory filesystem for each test
			mockFs := afero.NewMemMapFs()
			appFs = mockFs

			// Setup the filesystem if needed
			err := tt.setupFs(mockFs)
			require.NoError(t, err)

			// Call the Create function
			err = Create(tt.engine)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}

				// If testing existing file case, verify file wasn't changed
				if tt.name == "with existing file" {
					content, err := afero.ReadFile(mockFs, filepath.Join(testsDir, "rolecule.yml"))
					require.NoError(t, err)
					assert.Equal(t, "existing content", string(content))
				}
				return
			}

			assert.NoError(t, err)

			// Check if the file was created
			configPath := filepath.Join(testsDir, "rolecule.yml")
			exists, err := afero.Exists(mockFs, configPath)
			require.NoError(t, err)
			assert.True(t, exists)

			// Read the file content
			content, err := afero.ReadFile(mockFs, configPath)
			require.NoError(t, err)

			// Check if the engine section is present or not
			expectedSection := "engine:\n  name: " + tt.engine
			if tt.wantEngineSection {
				assert.Contains(t, string(content), expectedSection)
			} else {
				assert.NotContains(t, string(content), expectedSection)
			}

			// Check the file starts with "provisioner:" when engine is docker
			if !tt.wantEngineSection {
				assert.True(t, len(content) > 11)
				assert.Equal(t, "provisioner:", string(content)[:12])
			}

			// Verify other required sections are present
			assert.Contains(t, string(content), "provisioner:\n  name: ansible")
			assert.Contains(t, string(content), "verifier:\n  name: goss")
			assert.Contains(t, string(content), "instances:\n  - name: ubuntu-24.04\n    image: ubuntu-systemd:24.04")

			// Check log output
			found := false
			for _, entry := range handler.Entries {
				if entry.Level == log.InfoLevel &&
					strings.Contains(entry.Message, "created") &&
					strings.Contains(entry.Message, configPath) {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected log message about file creation not found")
		})
	}
}
