package utils

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"testing"
)

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

func TestMultipleCmd(t *testing.T) {
	if err := RunCmd("cd ~/Downloads && ls -l; cd tmp/tmp; ls -l"); err != nil {
		t.Fatal(err)
	}
}

func TestMultipleCmd2(t *testing.T) {
	if err := RunCmd("cd ~/Downloads && ls -l; cd tmp/tmp/a && ls -l; "); err != nil {
		t.Fatal(err)
	}
}

func TestMultipleCmd3(t *testing.T) {
	if err := RunCmd("cd /Downloads && ls -l && pwd; ls -l;"); err != nil {
		t.Fatal(err)
	}
}

func TestCmd(t *testing.T) {
	//create command
	catCmd := exec.Command( "cat", "exec.go" )
	wcCmd := exec.Command( "wc" )

	//make a pipe
	reader, writer := io.Pipe()
	var buf bytes.Buffer

	//set the output of "cat" command to pipe writer
	catCmd.Stdout = writer
	//set the input of the "wc" command pipe reader

	wcCmd.Stdin = reader

	//cache the output of "wc" to memory
	wcCmd.Stdout = &buf

	//start to execute "cat" command
	catCmd.Start()

	//start to execute "wc" command
	wcCmd.Start()

	//waiting for "cat" command complete and close the writer
	catCmd.Wait()
	writer.Close()

	//waiting for the "wc" command complete and close the reader
	wcCmd.Wait()
	reader.Close()
	//copy the buf to the standard output
	io.Copy( os.Stdout, &buf )
}
