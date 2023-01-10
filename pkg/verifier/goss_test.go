package verifier

import (
	"reflect"
	"testing"
)

func TestGossVerifier_String(t *testing.T) {
	tests := []struct {
		name string
		v    GossVerifier
		want string
	}{
		{
			name: "goss",
			v:    defaultGossConfig,
			want: "goss",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("GossVerifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGossVerifier_GetCommand(t *testing.T) {
	tests := []struct {
		name  string
		v     GossVerifier
		want  map[string]string
		want1 string
		want2 []string
	}{
		{
			name:  "goss",
			v:     defaultGossConfig,
			want:  nil,
			want1: "goss",
			want2: []string{
				"--gossfile",
				"tests/goss.yaml",
				"validate",
			},
		},
		{
			name: "goss-custom",
			v: GossVerifier{
				Name:     "goss",
				Command:  "goss",
				TestFile: "gossfile.yaml",
				ExtraArgs: []string{
					"--format",
					"json",
				},
			},
			want:  nil,
			want1: "goss",
			want2: []string{
				"--gossfile",
				"tests/gossfile.yaml",
				"validate",
				"--format",
				"json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := tt.v.GetCommand()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GossVerifier.GetCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GossVerifier.GetCommand() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GossVerifier.GetCommand() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
