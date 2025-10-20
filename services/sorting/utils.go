package sorting

import (
	"jcelaya775/go-cli/models"
	"slices"
)

func IsValidSortingAlgorithm(algorithm string) bool {
	return slices.Contains(models.SortingAlgorithms, models.SortingAlgorithm(algorithm))
}
