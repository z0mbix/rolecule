package verifier

import (
	"reflect"
	"testing"
)

func TestNewVerifier(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		args    args
		want    Verifier
		wantErr bool
	}{
		{
			name: "goss",
			args: args{
				config: Config{
					Name: "goss",
				},
			},
			want:    defaultGossConfig,
			wantErr: false,
		},
		{
			name: "testinfra",
			args: args{
				config: Config{
					Name: "testinfra",
				},
			},
			want:    defaultTestInfraConfig,
			wantErr: false,
		},
		{
			name: "inspec",
			args: args{
				config: Config{
					Name: "inspec",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVerifier(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVerifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVerifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
