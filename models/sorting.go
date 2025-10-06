package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"slices"
)

type SortingAlgorithm string

const (
	BubbleSort    SortingAlgorithm = "bubble"
	InsertionSort SortingAlgorithm = "insertion"
	SelectionSort SortingAlgorithm = "selection"
	MergeSort     SortingAlgorithm = "merge"
)

var SortingAlgorithmNames = map[SortingAlgorithm]string{
	BubbleSort:    "Bubble Sort",
	InsertionSort: "Insertion Sort",
	SelectionSort: "Selection Sort",
	MergeSort:     "Merge Sort",
}

var SortingAlgorithms = []SortingAlgorithm{
	BubbleSort,
	InsertionSort,
	SelectionSort,
	MergeSort,
}

type SortingAlgorithmIterator interface {
	NextCmd() tea.Cmd
	Abort()
}

func IsValidSortingAlgorithm(algorithm string) bool {
	return slices.Contains(SortingAlgorithms, SortingAlgorithm(algorithm))
}
