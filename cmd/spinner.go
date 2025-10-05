package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"jcelaya775/go-cli/cmd/ui/spinner"
)

func init() {
	rootCmd.AddCommand(showSpinnerCommand)
}

var showSpinnerCommand = &cobra.Command{
	Use:     "spinner",
	Short:   "Show a spinner",
	Args:    cobra.NoArgs,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(spinner.InitialModel())
		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}
