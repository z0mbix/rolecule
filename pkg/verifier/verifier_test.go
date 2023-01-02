package verifier

import (
	"reflect"
	"testing"
)

func TestNewVerifier(t *testing.T) {
	type args struct {
		name string
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
				name: "goss",
			},
			want:    defaultGossConfig,
			wantErr: false,
		},
		{
			name: "testinfra",
			args: args{
				name: "testinfra",
			},
			want:    defaultTestInfraConfig,
			wantErr: false,
		},
		{
			name: "inspec",
			args: args{
				name: "inspec",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVerifier(tt.args.name)
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
