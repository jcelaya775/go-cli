package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"jcelaya775/go-cli/cmd/ui/table"
)

func init() {
	rootCmd.AddCommand(tableCommand)
}

var tableCommand = &cobra.Command{
	Use:     "table",
	Short:   "Show a table",
	Args:    cobra.NoArgs,
	Aliases: []string{"t"},
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(table.InitialModel())
		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}
