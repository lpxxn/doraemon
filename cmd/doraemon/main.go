package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/utils"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// 👻 >
//const consolePrefix = "⚡️>>> "
const consolePrefix = "o((=ﾟェﾟ=))o > "

//const consolePrefix = "🤪 >> "
const mascot1 = `
  ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣤⣴⣶⣶⣶⣶⣶⣶⣶⣶⣶⠶⣶⣤⣤⣀⠀⠀⠀⠀⠀⠀ 
  ⠀⠀⠀⠀⠀⠀⠀⢀⣤⣾⣿⣿⣿⣿⣿⠁ ⠀⢀ ⠈⢿⢀⣀   ⠹⣿⣿⣦⣄⠀⠀⠀ 
  ⠀⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⣿⣿⣿⠿⠀⠀ ⣟⡇⢘ ⣾⣽⠀  ⠀⡏⠉⠙⢛⣷⡖⠀ 
  ⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⡿⠿⠷⠶⠤⠙⠒⠀⠒⢻⣿⣿⡷⠋⠀ ⠴⠞⠋⠁⢙⣿⣿⣄ 
  ⠀⠀⠀⠀⢸⣿⣿⣿⣿⣯⣤⣤⣤⣤⣤⡄⠀⠀⠀⠀⠉⢹⡄⠀  ⠀ ⠀⠛⠛⠋⠉⠹⡇ 
  ⠀⠀⠀⠀⢸⣿⣿⣿⣿⠀⠀⠀⣀⣠⣤⣤⣤⣤⣤⣤⣤⣼⣇⣀⣀⣀⣀⣀⣀⣀⣛⣛⣒⣲⢾⡷ 
  ⢀⠤⠒⠒⢼⣿⣿⣿⣿⠶⠞⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁ ⣼⠃ 
  ⢮⠀⠀⠀⠀⣿⣿⣿⣿⣆⠀⠀ ⠻⣿⡿⠛⠉⠉⠁⠀⠉⠉⠛⠿⣿⣿⠟⠁  ⠀⣼⠃⠀ 
  ⠈⠓⠶⣶⣾⣿⣿⣿⣿⣿⣧⡀⠀ ⠈⠒⢤⣀⣀⡀⠀⠀⣀⣀⡠⠚⠁ ⠀⢀⡼⠃⠀⠀ 
  ⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⣿⣷⣤⣤⣤⣤⣭⣭⣭⣭⣭⣥⣤⣤⣤⣴⣟⠁
`

const openConfigDir = "openConfigDir"

var existCommand = map[string]struct{}{"exit": {}, ":q": {}, "\\q": {}}

var sd fx.Shutdowner
var lc fx.Lifecycle

func main() {
	fmt.Println(mascot1)
	fmt.Println("type exit or :q or \\q to exit app")
	app := fx.New(fx.NopLogger,
		fx.Provide(config.ParseConfig),
		fx.Provide(fx.Annotated{
			Name:   "sshCompleter",
			Target: getSSHCompleter,
		}),
		fx.Provide(fx.Annotated{
			Name:   "customCmdCompleter",
			Target: getCustomCMDCompleter,
		}),
		fx.Provide(RootCMD),
		fx.Populate(&sd, &lc),
		fx.Invoke(customCmd))
	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("start err: %#v", err)
	}
	if err := app.Stop(context.Background()); err != nil {
		fmt.Errorf("stop err: %#v", err)
	}
}

type cmdParam struct {
	dig.In
	SSHCompleter prompt.Completer `name:"sshCompleter"`
	CmdCompleter prompt.Completer `name:"customCmdCompleter"`
}

func RootCMD(lc fx.Lifecycle, param cmdParam) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "doraemon",
		Short: "doraemon tools",
		Long:  `ssh manager and .....`,
		RunE: func(cmd *cobra.Command, args []string) error {

		exitCmd:
			for {
				utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
				cmdName := prompt.Input(consolePrefix, param.SSHCompleter)
				if strings.Trim(cmdName, " ") == "" {
					continue
				}
				runed, needExist := runGlobalCmd(cmdName)
				if runed && needExist {
					break exitCmd
				}
				if runed {
					continue
				}
				if err := startSSHShell(cmdName); err != nil {
					fmt.Println(err)
				}
			}
			if err := sd.Shutdown(); err != nil {
				fmt.Println("sd shutdown error", err)
			}
			// fmt.Println("stop ssh command")
			return nil
		},
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer func() {
				// todo 不是特别好。
				if err := rootCmd.Execute(); err != nil {
					fmt.Errorf("start err: %#v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			handleStty()
			// fmt.Println("life stop...")
			return nil
		},
	})
	return rootCmd
}

func runGlobalCmd(cmdName string) (runed bool, needExist bool) {
	if _, ok := existCommand[cmdName]; ok {
		fmt.Println("👋👋👋 bye ~")
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

func handleStty() {
	// https://github.com/c-bata/go-prompt/issues/233
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	//rawModeOff := exec.Command("/bin/stty", "sane")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func customCmd(rootCmd *cobra.Command, param cmdParam) {
	cmd := &cobra.Command{
		Use:   "cmd",
		Short: "custom cmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
			for {
				cmdName := prompt.Input(consolePrefix, param.CmdCompleter)
				if strings.Trim(cmdName, " ") == "" {
					continue
				}
				runed, needExist := runGlobalCmd(cmdName)
				if runed && needExist {
					break
				}
				if runed {
					continue
				}
				if err := runCustomCmd(cmdName); err != nil {
					fmt.Println(err)
				}
			}
			return sd.Shutdown()
		},
	}
	rootCmd.AddCommand(cmd)
}

func startSSHShell(sshName string) error {
	sshConfig, err := config.SSHConfigByName(sshName)
	if err != nil {
		return err
	}
	client, err := utils.NewSSHClient(sshConfig)
	if err != nil {
		return err

	}
	// Create Session
	session, err := client.CreateSession()
	if err != nil {
		return err
	}

	// Start ssh shell
	return client.Shell(session)
}

func runCustomCmd(cmdName string) error {
	item, err := config.CustomConfigByName(cmdName)
	if err != nil {
		return err
	}
	return utils.RunCmd(item.Cmd)
}

func getSSHCompleter(conf *config.AppConfig) prompt.Completer {
	sshSuggest := getSuggest(conf.SSHInfo)
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(sshSuggest, d.GetWordBeforeCursor(), true)
	}
}

func getCustomCMDCompleter(conf *config.AppConfig) prompt.Completer {
	sshSuggest := getSuggest(conf.CmdInfo)
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(sshSuggest, d.GetWordBeforeCursor(), true)
	}
}

func getSuggest(c config.InfoCollection) []prompt.Suggest {
	var sshSuggest []prompt.Suggest

	for iterator := c.GetIterator(); iterator.HasNext(); {
		item := iterator.Next()
		sshSuggest = append(sshSuggest, prompt.Suggest{
			Text:        item.GetName(),
			Description: item.GetDesc(),
		})
	}
	addOpenDirSuggest(&sshSuggest)
	return sshSuggest
}

func addOpenDirSuggest(sshSuggest *[]prompt.Suggest) {
	*sshSuggest = append(*sshSuggest, prompt.Suggest{
		Text:        "openConfigDir",
		Description: "open config directory",
	})
}

func RunSSHCommand(param cmdParam) {
exitCmd:
	for {
		utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
		//fmt.Println("Please select a command.")
		cmdName := prompt.Input(consolePrefix, param.SSHCompleter, prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(buffer *prompt.Buffer) {
				fmt.Println("👋👋👋 bye ~")
				os.Exit(0)
			},
		}))
		if _, ok := existCommand[cmdName]; ok {
			fmt.Println("👋👋👋 bye ~")
			break exitCmd
		}
		if openConfigDir == cmdName {
			if err := config.OpenConfDir(); err != nil {
				fmt.Println(err)
			}
			continue
		}
		if err := startSSHShell(cmdName); err != nil {
			fmt.Println(err)
		}
	}
}
