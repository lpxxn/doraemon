package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/utils"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

const (
	mascot1 = `
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
	consolePrefix = "o((=ﾟェﾟ=))o > "
	openConfigDir = "openConfigDir"
)

var existCommand = map[string]struct{}{"exit": {}, ":q": {}, "\\q": {}}

var (
	sd      fx.Shutdowner
	lc      fx.Lifecycle
	rootCmd *cobra.Command
	loopRun bool
)

func main() {
	fmt.Println(mascot1)
	utils.SendMsg(false, "type exit or :q or \\q to exit app", " ", utils.Yellow, true)
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
		fx.Invoke(Lifecycle),
		fx.Invoke(customCmd))
	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("start err: %#v", err)
	}
	rootCmd.PersistentFlags().BoolVarP(&loopRun, "loopRun", "l", false, "not exist until type :q or \\q")
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("start err: %#v", err)
		return
	}
	if err := app.Stop(context.Background()); err != nil {
		fmt.Printf("stop err: %#v", err)
	}
	utils.SendMsg(false, "bye ~", "👋👋👋 ", utils.Yellow, true)
}

func Lifecycle(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			utils.SetSttySane()
			return nil
		},
	})
}

func getSSHCompleter(conf *config.AppConfig) prompt.Completer {
	return getCompleter(conf.SSHInfo)
}

func getCustomCMDCompleter(conf *config.AppConfig) prompt.Completer {
	return getCompleter(conf.CmdInfo)
}

func getCompleter(c config.InfoCollection) prompt.Completer {
	var sshSuggest []prompt.Suggest

	for iterator := c.GetIterator(); iterator.HasNext(); {
		item := iterator.Next()
		sshSuggest = append(sshSuggest, prompt.Suggest{
			Text:        item.GetName(),
			Description: item.GetDesc(),
		})
	}
	sshSuggest = append(sshSuggest, prompt.Suggest{
		Text:        "openConfigDir",
		Description: "open config directory",
	})
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(sshSuggest, d.GetWordBeforeCursor(), true)
	}
}

type cmdParam struct {
	dig.In
	SSHCompleter prompt.Completer `name:"sshCompleter"`
	CmdCompleter prompt.Completer `name:"customCmdCompleter"`
}

func RootCMD(param cmdParam) *cobra.Command {
	rootCmd = &cobra.Command{
		Use:   "doraemon",
		Short: "doraemon tools",
		Long:  `ssh manager and run custom cmd`,
		RunE: func(cmd *cobra.Command, args []string) error {
		exitCmd:
			for {
				utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
				cmdName := prompt.Input(consolePrefix, param.SSHCompleter)
				if strings.Trim(cmdName, " ") == "" {
					continue
				}
				ran, needExist := runGlobalCmd(cmdName)
				if ran && needExist {
					break exitCmd
				}
				if ran {
					continue
				}
				if err := startSSHShell(cmdName); err != nil {
					fmt.Println(err)
				} else if !loopRun {
					break exitCmd
				}
			}
			if err := sd.Shutdown(); err != nil {
				fmt.Println("shutdown error", err)
			}
			return nil
		},
	}
	return rootCmd
}

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

func customCmd(rootCmd *cobra.Command, param cmdParam) {
	cmd := &cobra.Command{
		Use:   "cmd",
		Short: "run custom cmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, false)
			for {
				cmdName := prompt.Input(consolePrefix, param.CmdCompleter)
				if strings.Trim(cmdName, " ") == "" {
					continue
				}
				ran, needExist := runGlobalCmd(cmdName)
				if ran && needExist {
					break
				}
				if ran {
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
