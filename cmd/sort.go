package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"jcelaya775/go-cli/cmd/ui/sort"
	"time"
)

var tickDuration int

func init() {
	sortCmd.Flags().IntVarP(&tickDuration, "tickDuration", "t", 1000, "Tick duration for sorting steps in milliseconds")
	rootCmd.AddCommand(sortCmd)
}

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort an array of numbers",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		sortModel := sort.InitialModel([]int{5, 2, 9, 1, 5, 6, 3, 23, 12, 5}, time.Duration(tickDuration)*time.Millisecond)
		teaProgram := tea.NewProgram(sortModel)
		if _, err := teaProgram.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}
