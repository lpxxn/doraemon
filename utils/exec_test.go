package utils

import "testing"

func TestExecCommand(t *testing.T) {
	type args struct {
		command string
		options []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"successfully executing command",
			args{
				command: "echo",
				options: []string{"ping"},
			},
			false,
		},
		{
			"failed executing command",
			args{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecCommand(tt.args.command, tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("ExecCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunCmd1(t *testing.T) {
	if err := RunCmd("ls -l; pwd;"); err != nil {
		panic(err)
	}

	if err := RunCmd(`curl -X GET "http://www.baidu.com/$(date +%M)"`); err != nil {
		panic(err)
	}
}
