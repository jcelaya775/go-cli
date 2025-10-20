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

type InsertionSortingIteratorFactory struct{}

func (_ InsertionSortingIteratorFactory) New(ctx context.Context, nums []int) models.SortingAlgorithmIterator {
	it := &InsertionSortIterator{stateMsgCh: make(chan models.SortingStateMsg), ctx: ctx}
	go it.start(append([]int(nil), nums...))
	return it
}

func (it *InsertionSortIterator) start(nums []int) {
	n := len(nums)
	for i := 1; i < n; i++ {
		select {
		case <-it.ctx.Done():
			close(it.stateMsgCh)
			return
		default:
			it.stateMsgCh <- models.SortingStateMsg{
				Nums: append([]int(nil), nums...),
				I:    i,
				J:    i, // mark the current index being compared
			}
		}

		for j := i - 1; j >= 0 && nums[j+1] < nums[j]; j-- {
			nums[j], nums[j+1] = nums[j+1], nums[j]

			select {
			case <-it.ctx.Done():
				close(it.stateMsgCh)
				return
			default:
				it.stateMsgCh <- models.SortingStateMsg{
					Nums: append([]int(nil), nums...),
					I:    i + 1, // include the last sorted index as part of the sorted portion
					J:    j,
				}
			}
		}
	}

	close(it.stateMsgCh)
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
		} else {
			return models.SortingStateMsg{Done: true}
		}
	}
}
