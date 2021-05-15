package main

import (
	"fmt"
	"os"

	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/utils"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"mc"},
	Short:   "login",
	Long:    "\n oauth login.",
	Run:     runLoginCmd,
}

func main() {
	if err := config.ParseConfig(); err != nil {
		panic(err)
	}
	fmt.Println(*config.LoginConf)
	config.OpenConfDir()
	sshConfig, err := config.SSHConfigByName("sandbox1")
	if err != nil {
		panic(err)
	}
	client, err := utils.NewSSHClient(sshConfig)
	if err != nil {
		panic(err)
	}
	/*
		sandbox1Conf, err := config.LoginConf.ConfigByName("sandbox1")
		if err != nil {
			panic(err)
		}
		sandbox1SSHConf, err := sandbox1Conf.ToSSHConfig()
		if err != nil {
			panic(err)
		}
		clientOpts := make([]utils.SSHClientOption, 0)
		if sandbox1Conf.HaveProxy() {
			proxyConfig, err := config.LoginConf.ConfigByName(sandbox1Conf.ProxySSHName)
			if err != nil {
				panic(err)
			}
			proxySSHConf, err := proxyConfig.ToSSHConfig()
			if err != nil {
				panic(err)
			}
			clientOpts = append(clientOpts, utils.ProxyConfig(proxySSHConf))
		}
		client, err := utils.CreateSSHClient(sandbox1SSHConf, clientOpts...)
		if err != nil {
			panic(err)
		}

		/*
	*/
	// Create Session
	session, err := client.CreateSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start ssh shell
	if err := client.Shell(session); err != nil {
		panic(err)
	}
}

func runLoginCmd(cmd *cobra.Command, args []string) {
	utils.SendMsg(true, "go ...", "login ~", utils.Yellow, true)
}
