package provisioner

import (
	"reflect"
	"testing"
)

func TestNewProvisioner(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    Provisioner
		wantErr bool
	}{
		{
			name: "ansible",
			args: args{
				name: "ansible",
			},
			want:    defaultAnsibleConfig,
			wantErr: false,
		},
		{
			name: "chef",
			args: args{
				name: "chef",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "puppet",
			args: args{
				name: "puppet",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProvisioner(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvisioner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvisioner() = %v, want %v", got, tt.want)
			}
		})
	}
}
