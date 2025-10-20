package sorting

import (
	"context"
	"fmt"
	"jcelaya775/go-cli/models"
)

type IteratorFactory interface {
	New(context.Context, []int) models.SortingAlgorithmIterator
}

func GetIteratorFactory(algorithm models.SortingAlgorithm) (IteratorFactory, error) {
	switch algorithm {
	case models.InsertionSort:
		return InsertionSortIteratorFactory{}, nil
	case models.BubbleSort:
		return BubbleSortIteratorFactory{}, nil
	case models.SelectionSort:
		return SelectionSortIteratorFactory{}, nil
	case models.MergeSort:
		return MergeSortIteratorFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported sorting algorithm: %s", algorithm)
	}
}
