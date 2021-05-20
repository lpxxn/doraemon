package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/lpxxn/doraemon/config"
	"github.com/lpxxn/doraemon/utils"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"mc"},
	Short:   "login",
	Long:    "\n oauth login.",
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
	fx.New(fx.NopLogger,
		fx.Provide(
			fx.Populate(&sd, &lc),
			config.ParseConfig,
			setSSHSuggest,
			getSSHCompleter),
		fx.Invoke(RunSSHCommand))
}

func RunSSHCommand(sshCompleter prompt.Completer) {
	for {
		utils.SendMsg(true, "Hi!", "Please select a command.", utils.Yellow, true)
		cmdName := prompt.Input(consolePrefix, sshCompleter)
		//fmt.Println("You selected " + sshName)
		if _, ok := existCommand[cmdName]; ok {
			fmt.Println("👋👋👋 bye ~")
			return
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

func setSSHSuggest(conf *config.AppConfig) []prompt.Suggest{
	var sshSuggest = []prompt.Suggest{}
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
