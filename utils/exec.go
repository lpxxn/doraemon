package utils

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// ExecCommand function to execute a given command.
func ExecCommand(command string, options []string) error {
	//
	if command == "" {
		return errors.New("No command to execute!")
	}

	// Create buffer for stderr.
	stderr := &bytes.Buffer{}

	// Collect command line
	cmd := exec.Command(command, options...) // #nosec G204

	// Set buffer for stderr from cmd
	cmd.Stderr = stderr

	// Create a new reader
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Start executing command.
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create a new scanner and run goroutine func with output.
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			SendMsg(false, "*", scanner.Text(), Cyan, false)
		}
	}()

	// Wait for executing command.
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func RunCmd(cmdStr string, cmdDir ...string) error {
	var err error
	//fmt.Println("begin run command")
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
	c := exec.Command("/bin/sh", "-c", cmd)
	if len(cmdDir) > 0 && cmdDir[0] != "" {
		c.Dir = cmdDir[0]
	}
	f, err := pty.Start(c)
	return c, f, f, err
}
