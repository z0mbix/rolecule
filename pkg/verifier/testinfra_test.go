package verifier

import (
	"reflect"
	"testing"
)

func TestTestInfraVerifier_String(t *testing.T) {
	tests := []struct {
		name string
		v    *TestInfraVerifier
		want string
	}{
		{
			name: "testinfra",
			v:    defaultTestInfraConfig,
			want: "testinfra",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("TestInfraVerifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestInfraVerifier_GetCommand(t *testing.T) {
	tests := []struct {
		name  string
		v     *TestInfraVerifier
		want  map[string]string
		want1 string
		want2 []string
	}{
		{
			name:  "testinfra",
			v:     defaultTestInfraConfig,
			want:  nil,
			want1: "py.test",
			want2: []string{
				"-vv",
				"--hosts",
				"podman://foobar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.v.GetCommand()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestInfraVerifier.GetCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("TestInfraVerifier.GetCommand() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("TestInfraVerifier.GetCommand() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
