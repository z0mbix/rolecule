package container

import (
	"reflect"
	"testing"
)

func TestNewEngine(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Engine
		wantErr bool
	}{
		{
			name: "docker",
			args: args{
				name: "docker",
			},
			want: &DockerEngine{
				Name:   "docker",
				Socket: "docker://",
			},
		},
		{
			name: "podman",
			args: args{
				name: "podman",
			},
			want: &PodmanEngine{
				Name:   "podman",
				Socket: "podman://",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEngine(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEngine() got = %v, want %v", got, tt.want)
			}
		})
	}
}
