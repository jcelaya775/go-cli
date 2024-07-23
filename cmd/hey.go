package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var upperCase bool

func init() {
	heyCommand.Flags().BoolVarP(&upperCase, "upper", "u", false, "convert name to uppercase")
	rootCmd.AddCommand(heyCommand)
}

var heyCommand = &cobra.Command{
	Use:     "hey [name]",
	Short:   "Say hey to someone",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"h"},
	Run: func(cmd *cobra.Command, args []string) {
		if upperCase {
			args[0] = strings.ToUpper(args[0])
		}
		fmt.Printf("Hey, %s!\n", args[0])
	},
}
