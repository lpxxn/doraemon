package main

import (
	"fmt"

	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/internal"
)

func runGlobalCmd(cmdName string) (ran bool, needExist bool) {
	if _, ok := existCommand[cmdName]; ok {
		return true, true
	}
	if openConfigDir == cmdName {
		if err := config.OpenConfDir(); err != nil {
			fmt.Println(err)
		}
		return true, false
	}
	return false, false
}

func runCustomCmd(cmdName string) error {
	item, err := config.CustomConfigByName(cmdName)
	if err != nil {
		return err
	}
	return internal.RunCmd(item.Cmd)
}

func startSSHShell(sshName string) error {
	sshConfig, err := config.SSHConfigByName(sshName)
	if err != nil {
		return err
	}
	client, err := internal.NewSSHClient(sshConfig)
	if err != nil {
		return err
	}
	session, err := client.CreateSession()
	if err != nil {
		return err
	}
	return client.Shell(session)
}
