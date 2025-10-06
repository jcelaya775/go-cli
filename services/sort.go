package services

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
)

type InsertionSortIterator struct {
	ch chan InsertionSortStateMsg
}

type InsertionSortStateMsg struct {
	Nums      []int
	I         int
	J         int
	Done      bool
	Cancelled bool
}

func NewInsertionSortIterator(ctx context.Context, numsInput []int) InsertionSortIterator {
	stateCh := make(chan InsertionSortStateMsg)

	go func() {
		nums := make([]int, len(numsInput))
		copy(nums, numsInput)

		n := len(nums)
		for i := 1; i < n; i++ {
			select {
			case <-ctx.Done():
				stateCh <- InsertionSortStateMsg{
					Nums:      append([]int(nil), nums...),
					I:         i,
					J:         i,
					Done:      false,
					Cancelled: true,
				}
				//close(stateCh)
				return
			case stateCh <- InsertionSortStateMsg{
				Nums: append([]int(nil), nums...),
				I:    i,
				J:    i, // mark the current index being compared
			}:
				// sent successfully
			}

			for j := i - 1; j >= 0 && nums[j+1] < nums[j]; j-- {
				nums[j], nums[j+1] = nums[j+1], nums[j]

				select {
				case <-ctx.Done():
					stateCh <- InsertionSortStateMsg{
						Nums:      append([]int(nil), nums...),
						I:         i,
						J:         j,
						Done:      false,
						Cancelled: true,
					}
					//close(stateCh)
					return
				case stateCh <- InsertionSortStateMsg{
					Nums: append([]int(nil), nums...),
					I:    i + 1, // include the last sorted index as part of the sorted portion
					J:    j,
				}:
					// sent successfully
				}
			}
		}

		close(stateCh)
	}()

	return InsertionSortIterator{stateCh}
}

func (it *InsertionSortIterator) NextCmd() tea.Cmd {
	return func() tea.Msg {
		if state, ok := <-it.ch; ok {
			//fmt.Println("NextCmd Step:", state.Nums, "i:", state.I, "j:", state.J, "done:", state.Done, "cancelled:", state.Cancelled)
			if state.Cancelled {
				//close(it.ch)
				return InsertionSortStateMsg{Done: false, Cancelled: true}
			} else {
				return InsertionSortStateMsg{
					Nums: state.Nums,
					I:    state.I,
					J:    state.J,
					Done: false,
				}
			}
		} else {
			return InsertionSortStateMsg{Done: true}
		}
	}
}

//func (it *InsertionSortIterator) Abort() {
//	for range it.ch {
//		// Drain the channel to stop the goroutine
//	}
//}
