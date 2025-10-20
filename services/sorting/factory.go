package sorting

import (
	"context"
	"fmt"
	"jcelaya775/go-cli/models"
)

type SortingAlgorithmIteratorFactory interface {
	New(context.Context, []int) models.SortingAlgorithmIterator
}

func GetIteratorFactory(algorithm models.SortingAlgorithm) (SortingAlgorithmIteratorFactory, error) {
	switch algorithm {
	case models.InsertionSort:
		return InsertionSortingIteratorFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported sorting algorithm: %s", algorithm)
	}
}
