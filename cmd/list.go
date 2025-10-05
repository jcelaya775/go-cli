package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"jcelaya775/go-cli/cmd/ui/list"
	"jcelaya775/go-cli/cmd/ui/textinput"
)

func init() {
	rootCmd.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "Choose items from a list",
	Args:    cobra.NoArgs,
	Aliases: []string{"h"},
	Run: func(cmd *cobra.Command, args []string) {
		listModel := list.InitialModel()
		p := tea.NewProgram(listModel)
		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Println("\nHere's what you selected:")
		for i := range listModel.Selected {
			fmt.Printf("  - %s\n", listModel.Choices[i])
		}

		textInputModel := textinput.InitialModel()
		p = tea.NewProgram(textInputModel)
		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}
		for _, input := range textInputModel.Inputs {
			fmt.Printf("\nHere's what you typed for %s:\n", input.Placeholder)
			fmt.Printf("  - %s\n", input.Value())
		}
	},
}
