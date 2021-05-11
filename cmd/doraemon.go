package main

import (
	"fmt"
	"os"

	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/ssh_utils"
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
	//config.OpenConfDir()
	client, err := ssh_utils.CreateSSHClient(config.LoginConf.ConfigByName("sandbox1").ToSSHConfig(),
		ssh_utils.ProxyConfig(config.LoginConf.ConfigByName("proxy").ToSSHConfig()))
	if err != nil {
		panic(err)
	}

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
