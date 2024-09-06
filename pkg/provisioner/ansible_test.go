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
				"ANSIBLE_ROLES_PATH": "/etc/ansible/roles",
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

func TestAnsibleLocalProvisioner_WithExtraArgs(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Provisioner
	}{
		{
			name: "MultipleExtraArgs",
			fields: fields{
				Name: "ansible",
			},
			args: args{
				args: []string{
					"--diff",
					"--verbose",
				},
			},
			want: AnsibleLocalProvisioner{
				Name:      "ansible",
				ExtraArgs: []string{"--diff", "--verbose"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name: tt.fields.Name,
			}
			if got := a.WithExtraArgs(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExtraArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_WithTags(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Provisioner
	}{
		{
			name: "MultipleTags",
			fields: fields{
				Name: "ansible",
			},
			args: args{
				args: []string{
					"build",
					"configure",
				},
			},
			want: AnsibleLocalProvisioner{
				Name: "ansible",
				Tags: []string{"build", "configure"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name: tt.fields.Name,
			}
			if got := a.WithTags(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExtraArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_WithSkipTags(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Provisioner
	}{
		{
			name: "SkipTags",
			fields: fields{
				Name: "ansible",
			},
			args: args{
				args: []string{
					"ignore",
				},
			},
			want: AnsibleLocalProvisioner{
				Name:     "ansible",
				SkipTags: []string{"ignore"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name: tt.fields.Name,
			}
			if got := a.WithSkipTags(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExtraArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_WithPlaybook(t *testing.T) {
	type fields struct {
		Name         string
		Command      string
		Args         []string
		ExtraArgs    []string
		SkipTags     []string
		Tags         []string
		EnvVars      map[string]string
		Playbook     string
		Dependencies Dependencies
	}
	type args struct {
		playbook string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Provisioner
	}{
		{
			name: "Playbook",
			fields: fields{
				Name: "ansible",
			},
			args: args{
				playbook: "playbook.yaml",
			},
			want: AnsibleLocalProvisioner{
				Name:     "ansible",
				Playbook: "playbook.yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name:     tt.fields.Name,
				Playbook: tt.fields.Playbook,
			}
			if got := a.WithPlaybook(tt.args.playbook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPlaybook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_WithLocalDependencies(t *testing.T) {
	type fields struct {
		Name         string
		Command      string
		Args         []string
		ExtraArgs    []string
		SkipTags     []string
		Tags         []string
		EnvVars      map[string]string
		Playbook     string
		Dependencies Dependencies
	}
	type args struct {
		dependencies []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Provisioner
	}{
		{
			name: "local",
			fields: fields{
				Name: "ansible",
			},
			args: args{
				dependencies: []string{
					"depone",
					"deptwo",
				},
			},
			want: AnsibleLocalProvisioner{
				Name: "ansible",
				Dependencies: Dependencies{
					LocalRoles: []string{
						"depone",
						"deptwo",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name:         tt.fields.Name,
				Dependencies: tt.fields.Dependencies,
			}
			if got := a.WithLocalDependencies(tt.args.dependencies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLocalDependencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_WithGalaxyDependencies(t *testing.T) {
	type fields struct {
		Name         string
		Dependencies Dependencies
	}
	type args struct {
		dependencies []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Provisioner
	}{
		{
			name: "galaxy",
			fields: fields{
				Name: "ansible",
			},
			args: args{
				dependencies: []string{
					"z0mbix.depone",
					"z0mbix.deptwo",
				},
			},
			want: AnsibleLocalProvisioner{
				Name: "ansible",
				Dependencies: Dependencies{
					GalaxyRoles: []string{
						"z0mbix.depone",
						"z0mbix.deptwo",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name:         tt.fields.Name,
				Dependencies: tt.fields.Dependencies,
			}
			if got := a.WithGalaxyDependencies(tt.args.dependencies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithGalaxyDependencies() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_GetInstallDependenciesCommand(t *testing.T) {
	type fields struct {
		Name         string
		Command      string
		Args         []string
		ExtraArgs    []string
		SkipTags     []string
		Tags         []string
		EnvVars      map[string]string
		Playbook     string
		Dependencies Dependencies
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
		want1  string
		want2  []string
	}{
		{
			name: "basic",
			fields: fields{
				Name:    "ansible",
				Command: "ansible-playbook",
				Dependencies: Dependencies{
					Collections: nil,
					LocalRoles:  nil,
					GalaxyRoles: []string{
						"z0mbix.depone",
						"z0mbix.deptwo",
					},
				},
			},
			want: nil,
			//want: map[string]string{
			//	"ANSIBLE_ROLES_PATH": "/etc/ansible/roles",
			//	"ANSIBLE_NOCOWS":     "True",
			//},
			want1: "ansible-galaxy",
			want2: []string{
				"install",
				"--roles-path",
				"/etc/ansible/roles",
				"z0mbix.depone",
				"z0mbix.deptwo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name:         tt.fields.Name,
				Command:      tt.fields.Command,
				Args:         tt.fields.Args,
				ExtraArgs:    tt.fields.ExtraArgs,
				SkipTags:     tt.fields.SkipTags,
				Tags:         tt.fields.Tags,
				EnvVars:      tt.fields.EnvVars,
				Playbook:     tt.fields.Playbook,
				Dependencies: tt.fields.Dependencies,
			}
			got, got1, got2 := a.GetInstallDependenciesCommand()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstallDependenciesCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetInstallDependenciesCommand() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetInstallDependenciesCommand() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_GetCommand(t *testing.T) {
	type fields struct {
		Name         string
		Command      string
		Args         []string
		ExtraArgs    []string
		SkipTags     []string
		Tags         []string
		EnvVars      map[string]string
		Playbook     string
		Dependencies Dependencies
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
		want1  string
		want2  []string
	}{
		{
			name: "basic",
			fields: fields{
				Name:    "ansible",
				Command: "ansible-playbook",
				EnvVars: map[string]string{
					"ANSIBLE_ROLES_PATH": "/etc/ansible/roles",
					"ANSIBLE_NOCOWS":     "True",
				},
				Args: []string{
					"--connection",
					"local",
					"--inventory",
					"localhost,",
				},
				Playbook: "playbook.yml",
			},
			want: map[string]string{
				"ANSIBLE_ROLES_PATH": "/etc/ansible/roles",
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
			a := AnsibleLocalProvisioner{
				Name:         tt.fields.Name,
				Command:      tt.fields.Command,
				Args:         tt.fields.Args,
				ExtraArgs:    tt.fields.ExtraArgs,
				SkipTags:     tt.fields.SkipTags,
				Tags:         tt.fields.Tags,
				EnvVars:      tt.fields.EnvVars,
				Playbook:     tt.fields.Playbook,
				Dependencies: tt.fields.Dependencies,
			}
			got, got1, got2 := a.GetCommand()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCommand() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetCommand() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getAnsibleConfig(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
		want AnsibleLocalProvisioner
	}{
		{
			name: "foobar",
			args: args{
				config: Config{
					Name: "foo",
				},
			},
			want: AnsibleLocalProvisioner{
				Name:    "ansible",
				Command: "ansible-playbook",
				Args: []string{
					"--connection",
					"local",
					"--inventory",
					"localhost,",
				},
				EnvVars: map[string]string{
					"ANSIBLE_ROLES_PATH": ansibleRoleDir,
					"ANSIBLE_NOCOWS":     "True",
				},
				Playbook: "playbook.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAnsibleConfig(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAnsibleConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnsibleLocalProvisioner_String(t *testing.T) {
	type fields struct {
		Name         string
		Command      string
		Args         []string
		ExtraArgs    []string
		SkipTags     []string
		Tags         []string
		EnvVars      map[string]string
		Playbook     string
		Dependencies Dependencies
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ansible",
			fields: fields{
				Name: "ansible",
			},
			want: "ansible",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnsibleLocalProvisioner{
				Name: tt.fields.Name,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
