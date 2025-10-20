package sorting

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"jcelaya775/go-cli/models"
)

type InsertionSortIterator struct {
	stateMsgCh chan models.SortingStateMsg
	ctx        context.Context
}

type InsertionSortIteratorFactory struct{}

func (_ InsertionSortIteratorFactory) New(ctx context.Context, nums []int) models.SortingAlgorithmIterator {
	it := &InsertionSortIterator{stateMsgCh: make(chan models.SortingStateMsg), ctx: ctx}
	go it.start(append([]int(nil), nums...))
	return it
}

func (it *InsertionSortIterator) start(nums []int) {
	defer close(it.stateMsgCh)

	for i := 1; i < len(nums); i++ {
		if !it.sendState(models.SortingStateMsg{
			Nums: append([]int(nil), nums...),
			I:    i,
			J:    i,
		}) {
			return
		}

		for j := i - 1; j >= 0 && nums[j+1] < nums[j]; j-- {
			nums[j], nums[j+1] = nums[j+1], nums[j]

			if !it.sendState(models.SortingStateMsg{
				Nums: append([]int(nil), nums...),
				I:    i + 1, // include as part of the sorted portion
				J:    j,
			}) {
				return
			}
		}
	}
}

func (it *InsertionSortIterator) sendState(state models.SortingStateMsg) bool {
	select {
	case <-it.ctx.Done():
		return false
	case it.stateMsgCh <- state:
		return true
	}
}

func (it *InsertionSortIterator) NextCmd() tea.Cmd {
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
