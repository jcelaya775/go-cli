package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"jcelaya775/go-cli/cmd/ui/sort"
	"jcelaya775/go-cli/models"
	"strings"
	"time"
)

var tickDuration int

const SortingAlgorithmFieldName = "algorithm"

func init() {
	var algorithmOptionsList []string
	for _, alg := range models.SortingAlgorithms {
		algorithmOptionsList = append(algorithmOptionsList, string(alg))
	}
	algorithmHelpText := "Sorting algorithm to use. Options:\n  - " + strings.Join(algorithmOptionsList, "\n  - ")
	sortCmd.Flags().StringP(SortingAlgorithmFieldName, "a", "", algorithmHelpText)

	sortCmd.Flags().IntVarP(&tickDuration, "tickDuration", "t", 1000, "Tick duration for sorting steps in milliseconds")

	rootCmd.AddCommand(sortCmd)
}

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort an array of numbers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sortingAlgorithm, err := userSelectSortingAlgorithm(*cmd)
		if err != nil {
			return err
		}
		sortModel := sort.InitialModel(sortingAlgorithm, []int{5, 2, 9, 1, 5, 6, 3, 23, 12, 5}, time.Duration(tickDuration)*time.Millisecond)
		teaProgram := tea.NewProgram(sortModel)
		if _, err := teaProgram.Run(); err != nil {
			return err
		}
		return nil
	},
}

func userSelectSortingAlgorithm(cmd cobra.Command) (models.SortingAlgorithm, error) {
	sortingAlgorithmInput, err := cmd.Flags().GetString(SortingAlgorithmFieldName)
	if err != nil {
		return "", err
	}

	if sortingAlgorithmInput != "" {
		if !models.IsValidSortingAlgorithm(sortingAlgorithmInput) {
			return "", fmt.Errorf("invalid sorting algorithm: %s", sortingAlgorithmInput)
		}
		return models.SortingAlgorithm(sortingAlgorithmInput), nil
	}

	var sortingAlgorithm models.SortingAlgorithm
	var sortingOptions []huh.Option[models.SortingAlgorithm]
	for _, algorithm := range models.SortingAlgorithms {
		sortingOptions = append(sortingOptions, huh.NewOption(string(algorithm), algorithm))
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[models.SortingAlgorithm]().
				Title("Select a sorting algorithm").
				Options(sortingOptions...).
				Value(&sortingAlgorithm),
		),
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return sortingAlgorithm, nil
}
