package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/SisyphusSQ/go-web-starter/vars"
)

var initOnce sync.Once

var rootCmd = &cobra.Command{
	Use:   vars.AppName,
	Short: "Generate Go web starter projects",
	Long:  "go-web-starter generates a ready-to-use Go web project template.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func initAll() {
	initVersion()
	initInit()
	initNew()
}

func ensureInitialized() {
	initOnce.Do(initAll)
}

func Execute() {
	ensureInitialized()
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
