package sorting

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"jcelaya775/go-cli/models"
)

type MergeSortIterator struct {
	stateMsgCh chan models.SortingStateMsg
	ctx        context.Context
}

type MergeSortIteratorFactory struct{}

func (_ MergeSortIteratorFactory) New(ctx context.Context, nums []int) models.SortingAlgorithmIterator {
	it := &MergeSortIterator{stateMsgCh: make(chan models.SortingStateMsg), ctx: ctx}
	go it.start(append([]int(nil), nums...))
	return it
}

func (it *MergeSortIterator) start(nums []int) {
	defer close(it.stateMsgCh)

	merge := func(left, right []int) []int {
		merged := make([]int, 0, len(left)+len(right))
		i, j := 0, 0

		if !it.sendState(models.SortingStateMsg{
			Nums: append([]int(nil), append(merged, append(left[i:], right[j:]...)...)...),
			I:    len(merged),
			J:    len(merged),
		}) {
			return nil
		}

		for i < len(left) && j < len(right) {
			if left[i] < right[j] {
				//if !it.sendState(models.SortingStateMsg{
				//	Nums: append([]int(nil), append(merged, append(left[i:], right[j:]...)...)...),
				//	I:    len(merged) + 1,
				//	J:    len(merged),
				//}) {
				//	return nil
				//}

				merged = append(merged, left[i])
				i++
			} else {
				//if !it.sendState(models.SortingStateMsg{
				//	Nums: append([]int(nil), append(merged, append(left[i:], right[j:]...)...)...),
				//	I:    len(merged),
				//	J:    len(merged) + 1,
				//}) {
				//	return nil
				//}

				merged = append(merged, right[j])
				j++
			}

			if !it.sendState(models.SortingStateMsg{
				Nums: append([]int(nil), append(merged, append(left[i:], right[j:]...)...)...),
				I:    len(merged),
				J:    len(merged),
			}) {
				return nil
			}
		}

		merged = append(merged, left[i:]...)
		merged = append(merged, right[j:]...)
		return merged
	}

	var mergeSort func([]int) []int
	mergeSort = func(subNums []int) []int {
		if len(subNums) <= 1 {
			return subNums
		}

		mid := len(subNums) / 2
		leftSorted := mergeSort(subNums[:mid])
		rightSorted := mergeSort(subNums[mid:])

		return merge(leftSorted, rightSorted)
	}

	nums = mergeSort(nums)
	if !it.sendState(models.SortingStateMsg{
		Nums: append([]int(nil), nums...),
		I:    len(nums),
		J:    0,
	}) {
		return
	}
}

func (it *MergeSortIterator) sendState(state models.SortingStateMsg) bool {
	select {
	case <-it.ctx.Done():
		return false
	case it.stateMsgCh <- state:
		return true
	}
}

func (it *MergeSortIterator) NextCmd() tea.Cmd {
	return func() tea.Msg {
		if state, ok := <-it.stateMsgCh; ok {
			return models.SortingStateMsg{
				Nums: state.Nums,
				I:    state.I,
				J:    state.J,
				Done: false,
			}
		}
		return models.SortingStateMsg{Done: true}
	}
}
