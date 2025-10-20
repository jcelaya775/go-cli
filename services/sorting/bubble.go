package sorting

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"jcelaya775/go-cli/models"
)

type BubbleSortIterator struct {
	stateMsgCh chan models.SortingStateMsg
	ctx        context.Context
}

type BubbleSortIteratorFactory struct{}

func (_ BubbleSortIteratorFactory) New(ctx context.Context, nums []int) models.SortingAlgorithmIterator {
	it := &BubbleSortIterator{stateMsgCh: make(chan models.SortingStateMsg), ctx: ctx}
	go it.start(append([]int(nil), nums...))
	return it
}

func (it *BubbleSortIterator) start(nums []int) {
	defer close(it.stateMsgCh)

	for i := len(nums) - 1; i >= 0; i-- {
		if !it.sendState(models.SortingStateMsg{
			Nums: append([]int(nil), nums...),
			I:    i + 1, // include as part of the unsorted portion
			J:    0,
		}) {
			return
		}

		for j := 0; j < i; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}

			if !it.sendState(models.SortingStateMsg{
				Nums: append([]int(nil), nums...),
				I:    i + 1,
				J:    j + 1,
			}) {
				return
			}
		}
	}
}

func (it *BubbleSortIterator) sendState(state models.SortingStateMsg) bool {
	select {
	case <-it.ctx.Done():
		return false
	case it.stateMsgCh <- state:
		return true
	}
}

func (it *BubbleSortIterator) NextCmd() tea.Cmd {
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
