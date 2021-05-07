package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)


func main() {
	for _, env := range os.Environ() {
		// env is
		envPair := strings.SplitN(env, "=", 2)
		key := envPair[0]
		value := envPair[1]

		fmt.Printf("%s : %s\n", key, value)
	}


	cmd := exec.Command("ssh", "cafetest1dev")
	cmd.Env = append(os.Environ())
	if err := cmd.Run(); err != nil {
		fmt.Print(err)
	}
	fmt.Println(cmd.Wait())
	if err := RunCmd("cafe_sandbox_test1"); err != nil {
		fmt.Print(err)
	}
}

func RunCmd(cmdStr string, cmdDir ...string) error {
	var err error
	fmt.Println("begin run command")
	cmd, stdout, stderr, err := startCmd(cmdStr, cmdDir...)
	if err != nil {
		return err
	}
	defer func() {
		stdout.Close()
		stderr.Close()
	}()
	io.Copy(os.Stdout, stdout)
	io.Copy(os.Stderr, stderr)
	// wait for building
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func startCmd(cmd string, cmdDir ...string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	c := exec.Command("/bin/zsh",  cmd)
	c.Env = append(os.Environ())

	if len(cmdDir) > 0 && cmdDir[0] != "" {
		c.Dir = cmdDir[0]
	}
	f, err := pty.Start(c)
	return c, f, f, err
}

func RunSSHCmd(remoteMachine string, cmdStr string) error {
	var err error
	fmt.Println("begin run command")
	cmd, stdout, stderr, err := startSSHCmd(remoteMachine, cmdStr)
	if err != nil {
		return err
	}
	defer func() {
		stdout.Close()
		stderr.Close()
	}()
	io.Copy(os.Stdout, stdout)
	io.Copy(os.Stderr, stderr)
	// wait for building
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func startSSHCmd(remoteMachine string, cmd string) (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	c := exec.Command("ssh", remoteMachine, cmd)
	f, err := pty.Start(c)
	return c, f, f, err
}
