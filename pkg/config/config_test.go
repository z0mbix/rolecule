package config

import "testing"

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
