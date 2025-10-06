package services

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

type InsertionSortIterator struct {
	stateMsgCh chan InsertionSortStateMsg
	cancelled  bool
}

type InsertionSortStateMsg struct {
	Nums    []int
	I       int
	J       int
	Done    bool
	Restart bool
}

func NewInsertionSortIterator(ctx context.Context, nums []int) *InsertionSortIterator {
	it := &InsertionSortIterator{stateMsgCh: make(chan InsertionSortStateMsg), cancelled: false}
	go it.awaitSortSteps(ctx, append([]int(nil), nums...))
	return it
}

func (it *InsertionSortIterator) awaitSortSteps(ctx context.Context, nums []int) {
	n := len(nums)
	for i := 1; i < n; i++ {
		select {
		case <-ctx.Done():
			it.cancelled = true
			close(it.stateMsgCh)
			return
		default:
			it.stateMsgCh <- InsertionSortStateMsg{
				Nums: append([]int(nil), nums...),
				I:    i,
				J:    i, // mark the current index being compared
			}
		}

		for j := i - 1; j >= 0 && nums[j+1] < nums[j]; j-- {
			nums[j], nums[j+1] = nums[j+1], nums[j]

			select {
			case <-ctx.Done():
				it.cancelled = true
				close(it.stateMsgCh)
				return
			default:
				it.stateMsgCh <- InsertionSortStateMsg{
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
			//fmt.Println("NextCmd Step:", state.Nums, "i:", state.I, "j:", state.J, "done:", state.Done, "cancelled:", state.Cancelled)
			return InsertionSortStateMsg{
				Nums: state.Nums,
				I:    state.I,
				J:    state.J,
				Done: false,
			}
		} else {
			if it.cancelled {
				return InsertionSortStateMsg{Done: false, Restart: true}
			} else {
				return InsertionSortStateMsg{Done: true}
			}
		}
	}
}

func (it *InsertionSortIterator) Cancel() {
	it.cancelled = true
}
