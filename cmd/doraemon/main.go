package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/internal"
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
		fx.Invoke(customCmd), fx.Invoke(httpSrvCmd))
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
	internal.SendMsg(false, "bye ~", "👋👋👋 ", internal.Yellow, true)
}

func Lifecycle(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			internal.SetSttySane()
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
			internal.SendMsg(false, "type exit or :q or \\q to exit", " ", internal.Yellow, true)
		exitCmd:
			for {
				internal.SendMsg(true, "Hi!", "Please select a command.", internal.Yellow, false)
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

func customCmd(rootCmd *cobra.Command, param cmdParam) {
	cmd := &cobra.Command{
		Use:   "cmd",
		Short: "run custom cmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			internal.SendMsg(false, "type exit or :q or \\q to exit", " ", internal.Yellow, true)
			internal.SendMsg(true, "Hi!", "Please select a command.", internal.Yellow, false)
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
				if !loopRun {
					break
				}
			}
			return sd.Shutdown()
		},
	}
	rootCmd.AddCommand(cmd)
}

func httpSrvCmd(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "srv",
		Short: "Simple http service for sharing files",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer sd.Shutdown()
			p := "."
			if len(args) > 0 {
				p = args[0]
			}
			if err := internal.HttpFileServ(p); err != nil {
				internal.SendMsg(false, "Error", err.Error(), internal.Yellow, false)
			}
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
