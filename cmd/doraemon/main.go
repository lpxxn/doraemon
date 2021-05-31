package main

import (
	"context"
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/utils"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

var loginCmd = &cobra.Command{
	Use:     "command",
	Aliases: []string{"cmd"},
	Short:   "cmd",
	Long:    "\n command .",
	Run:     runLoginCmd,
}

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
		fx.Provide(
			config.ParseConfig,
			setSSHSuggest),
		fx.Provide(fx.Annotated{
			Name:   "sshCompleter",
			Target: getSSHCompleter,
		}),
		//fx.Provide(NewSSHPrompt),
		fx.Provide(RootCMD),
		fx.Populate(&sd, &lc),
		fx.Invoke(customCmd))
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	<-app.Done()
	fmt.Println("done")
	//app.Stop(context.Background())
}

type cmdParam struct {
	dig.In
	Completer prompt.Completer `name:"sshCompleter"`
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
				//fmt.Println("Please select a command.")
				cmdName := prompt.Input(consolePrefix, param.Completer)
				/*
					, prompt.OptionAddKeyBind(prompt.KeyBind{
							Key: prompt.ControlC,
							Fn: func(buffer *prompt.Buffer) {
								fmt.Println("👋👋👋 bye ~")
								sd.Shutdown()
								//os.Exit(0)
							},
						})
				*/
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
			if err := sd.Shutdown(); err != nil {
				fmt.Println("sd shutdown error", err)
			}
			fmt.Println("-adfasdf")
			return nil
		},
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer rootCmd.Execute()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return rootCmd
}

func customCmd(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "cmd",
		Short: "custom cmd",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("aaaa-----")
			fmt.Println(*cmd)
		},
	}
	rootCmd.AddCommand(cmd)
}

func RunSSHCommand(param cmdParam) {
exitCmd:
	for {
		utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
		//fmt.Println("Please select a command.")
		cmdName := prompt.Input(consolePrefix, param.Completer, prompt.OptionAddKeyBind(prompt.KeyBind{
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

func runLoginCmd(cmd *cobra.Command, args []string) {
	utils.SendMsg(true, "go ...", "login ~", utils.Yellow, true)
}

func getSSHCompleter(sshSuggest []prompt.Suggest) prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(sshSuggest, d.GetWordBeforeCursor(), true)
	}
}

func setSSHSuggest(conf *config.AppConfig) []prompt.Suggest {
	var sshSuggest []prompt.Suggest
	for _, item := range conf.SSHInfo {
		sshSuggest = append(sshSuggest, prompt.Suggest{
			Text:        item.Name,
			Description: item.Desc,
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
