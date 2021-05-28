package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/lpxxn/doraemon/utils"
	"go.uber.org/dig"
)

const consolePrefix = "🤪 >> "

func main() {

}

type ICommand interface {
	Desc() string
	CmdName() string
	Exec() error
}

type app struct {
	currentCmd ICommand
	cmds       map[string]ICommand
}

func (a *app) AddCmd(c ICommand) {
	a.cmds[c.CmdName()] = c
}

func (a *app) ChangeCurrentCmd(name string) {
	a.currentCmd = a.cmds[name]
}

type sshPromptInfo struct {
	cmdName   string
	Completer prompt.Completer `name:"sshCompleter"`
}

var sshPrompt *sshPromptInfo

type sshCmdParam struct {
	dig.In
	Completer prompt.Completer `name:"sshCompleter"`
}

func NewSSHPrompt(param sshCmdParam) {
	sshPrompt = &sshPromptInfo{Completer: param.Completer, cmdName: "servCMD"}
}

func (s *sshPromptInfo) CmdName() string {
	return s.cmdName
}

func (s *sshPromptInfo) Desc() string {
	return "ssh services manager"
}

func (s *sshPromptInfo) Exec() error {
	utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
	//fmt.Println("Please select a command.")
	cmdName := prompt.Input(consolePrefix, s.Completer, prompt.OptionAddKeyBind(prompt.KeyBind{
		Key: prompt.ControlC,
		Fn: func(buffer *prompt.Buffer) {
			fmt.Println("👋👋👋 bye ~")
			os.Exit(0)
		},
	}))
	fmt.Println(cmdName)
	return nil
}
