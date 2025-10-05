package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "go-cli",
	Short:   "It's a CLI",
	Long:    "It's a CLI, duh.",
	Example: "go-cli hey Jeoff",
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
