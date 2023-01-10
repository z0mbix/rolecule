package provisioner

import (
	"reflect"
	"testing"
)

func TestAnsibleProvisioner_String(t *testing.T) {
	tests := []struct {
		name string
		a    AnsibleLocalProvisioner
		want string
	}{
		{
			name: "Ansible",
			a:    defaultAnsibleConfig,
			want: "ansible",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("AnsibleProvisioner.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleProvisioner_GetCommand(t *testing.T) {
	tests := []struct {
		name  string
		a     AnsibleLocalProvisioner
		want  map[string]string
		want1 string
		want2 []string
	}{
		{
			name: "command",
			a:    defaultAnsibleConfig,
			want: map[string]string{
				"ANSIBLE_ROLES_PATH": ".",
				"ANSIBLE_NOCOWS":     "True",
			},
			want1: "ansible-playbook",
			want2: []string{
				"--connection",
				"local",
				"--inventory",
				"localhost,",
				"tests/playbook.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.a.GetCommand()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnsibleProvisioner.GetCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AnsibleProvisioner.GetCommand() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("AnsibleProvisioner.GetCommand() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
