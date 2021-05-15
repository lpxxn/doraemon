package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
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

func initConf() error {
	if err := config.ParseConfig(); err != nil {
		return err
	}
	// fmt.Println(*config.LoginConf)
	// config.OpenConfDir()
	setSSHSuggest()
	return nil
}

// 👻 >
//const consolePrefix = "⚡️>>> "
const consolePrefix = "🤪 >> "

func main() {
	if err := initConf(); err != nil {
		panic(err)
	}
	fmt.Println("Please select ssh name.")
	sshName := prompt.Input(consolePrefix, sshCompleter)
	fmt.Println("You selected " + sshName)
	sshConfig, err := config.SSHConfigByName(sshName)
	if err != nil {
		panic(err)
	}
	client, err := utils.NewSSHClient(sshConfig)
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
		fmt.Println(err)
	}
}

func runLoginCmd(cmd *cobra.Command, args []string) {
	utils.SendMsg(true, "go ...", "login ~", utils.Yellow, true)
}

var sshSuggest []prompt.Suggest

func sshCompleter(d prompt.Document) []prompt.Suggest {
	return prompt.FilterHasPrefix(sshSuggest, d.GetWordBeforeCursor(), true)
}

func setSSHSuggest() {
	sshSuggest = sshSuggest[:0]
	for _, item := range config.LoginConf.SSHInfo {
		sshSuggest = append(sshSuggest, prompt.Suggest{
			Text:        item.Name,
			Description: item.Desc,
		})
	}
}
