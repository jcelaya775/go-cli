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
		return InsertionSortingIteratorFactory{}, nil
	case models.BubbleSort:
		return BubbleSortingIteratorFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported sorting algorithm: %s", algorithm)
	}
}
