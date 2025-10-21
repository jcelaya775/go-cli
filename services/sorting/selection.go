package sorting

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"jcelaya775/go-cli/models"
)

type SelectionSortIterator struct {
	stateMsgCh chan models.SortingStateMsg
	ctx        context.Context
}

type SelectionSortIteratorFactory struct{}

// main test
func (_ SelectionSortIteratorFactory) New(ctx context.Context, nums []int) models.SortingAlgorithmIterator {
	it := &SelectionSortIterator{stateMsgCh: make(chan models.SortingStateMsg), ctx: ctx}
	go it.start(append([]int(nil), nums...))
	return it
}

func (it *SelectionSortIterator) start(nums []int) {
	defer close(it.stateMsgCh)

	for i := 0; i < len(nums)-1; i++ {
		if !it.sendState(models.SortingStateMsg{
			Nums: append([]int(nil), nums...),
			I:    i,
			J:    i,
		}) {
			return
		}

		minIdx := i
		for j := i + 1; j < len(nums); j++ {
			if nums[j] < nums[minIdx] {
				minIdx = j
			}

			if !it.sendState(models.SortingStateMsg{
				Nums: append([]int(nil), nums...),
				I:    i,
				J:    j,
			}) {
				return
			}
		}

		if !it.sendState(models.SortingStateMsg{
			Nums: append([]int(nil), nums...),
			I:    i,
			J:    minIdx,
		}) {
			return
		}

		nums[i], nums[minIdx] = nums[minIdx], nums[i]

		if !it.sendState(models.SortingStateMsg{
			Nums: append([]int(nil), nums...),
			I:    i + 1,
			J:    i,
		}) {
			return
		}
	}
}

func (it *SelectionSortIterator) sendState(state models.SortingStateMsg) bool {
	select {
	case <-it.ctx.Done():
		return false
	case it.stateMsgCh <- state:
		return true
	}
}

func (it *SelectionSortIterator) NextCmd() tea.Cmd {
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
