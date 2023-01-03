package provisioner

import (
	"reflect"
	"testing"
)

func TestNewProvisioner(t *testing.T) {
	type args struct {
		config Config
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
				config: Config{
					Name: "ansible",
				},
			},
			want:    defaultAnsibleConfig,
			wantErr: false,
		},
		{
			name: "puppet",
			args: args{
				config: Config{
					Name: "puppet",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProvisioner(tt.args.config)
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
