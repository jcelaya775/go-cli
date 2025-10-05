package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"jcelaya775/go-cli/cmd/ui/request"
)

func init() {
	rootCmd.AddCommand(requestCmd)
}

var requestCmd = &cobra.Command{
	Use:     "request",
	Short:   "Make a GET request to https://charm.sh/",
	Args:    cobra.NoArgs,
	Aliases: []string{"req"},
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(request.Model{})
		_, err := p.Run()
		return err
	},
}
