package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Verbose bool
	rootCmd = &cobra.Command{
		Use:   "cobra",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("hello %t \n", Verbose)
			fmt.Println("root command")
			return nil
		},
	}
)

/*
go run main.go
go run main.go -v
go run main.go try
go run main.go try -v
*/
func init() {
	rootCmd.AddCommand(tryCmd)
}

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "Try and possibly fail at something",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := someFunc(); err != nil {
			return err
		}
		return nil
	},
}

func someFunc() error {
	fmt.Println("hello")
	fmt.Printf("hello %t \n", Verbose)
	return nil
}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
