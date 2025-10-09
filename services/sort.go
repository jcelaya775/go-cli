package services

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

type InsertionSortIterator struct {
	stateMsgCh chan InsertionSortStateMsg
	nextCh     chan InsertionSortStateMsg
	cancelCh   chan int
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

// TODO: Look into iter.Pull implementation, it might be a better fit for this use case than channels
func (it *InsertionSortIterator) ListenToStateUpdates() tea.Msg {
	for {
		if state, ok := <-it.stateMsgCh; ok {
			// Allow NextCmd to proceed
			it.nextCh <- state
		} else {
			if it.cancelled {
				// If cancelled, immediately send a Restart message to reset the state
				return InsertionSortStateMsg{
					Restart: true,
				}
			} else {
				it.nextCh <- state
			}
		}
	}
}

// TODO: Implement fan-out pattern to allow sending NextCmd and Cancel messages without blocking
func (it *InsertionSortIterator) NextCmd() tea.Cmd {
	return func() tea.Msg {
		if state, ok := <-it.nextCh; ok {
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
