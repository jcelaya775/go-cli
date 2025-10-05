package services

import (
	tea "github.com/charmbracelet/bubbletea"
)

type InsertionSortIterator struct {
	ch chan stateChannel
}

type stateChannel struct {
	nums []int
	i    int
	j    int
}

func NewInsertionSortIterator(nums []int) InsertionSortIterator {
	numsCh := make(chan stateChannel)

	go func() {
		numsCopy := make([]int, len(nums))
		copy(numsCopy, nums)

		n := len(numsCopy)
		for i := 1; i < n; i++ {
			// Send state when selecting new key to compare
			currentState := make([]int, len(numsCopy))
			copy(currentState, numsCopy)
			numsCh <- stateChannel{
				nums: currentState,
				i:    i,
				j:    i,
			}

			for j := i - 1; j >= 0 && numsCopy[j+1] < numsCopy[j]; j-- {
				// Swap and send state after swap
				numsCopy[j], numsCopy[j+1] = numsCopy[j+1], numsCopy[j]
				currentState := make([]int, len(numsCopy))
				copy(currentState, numsCopy)
				numsCh <- stateChannel{
					nums: currentState,
					i:    i + 1,
					j:    j,
				}
			}
		}
		close(numsCh)
	}()

	return InsertionSortIterator{numsCh}
}

type InsertionSortMsg struct {
	Nums []int
	I    int
	J    int
	Done bool
}

func (it *InsertionSortIterator) NextCmd() tea.Cmd {
	return func() tea.Msg {
		if state, ok := <-it.ch; ok {
			return InsertionSortMsg{
				Nums: state.nums,
				I:    state.i,
				J:    state.j,
				Done: false,
			}
		} else {
			return InsertionSortMsg{Done: true}
		}
	}
}

func (it *InsertionSortIterator) Abort() {
	// Drain the channel to stop the goroutine
	for range it.ch {
	}
}
