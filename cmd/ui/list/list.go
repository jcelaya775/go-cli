package list

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type Output struct {
	value string
}

type Model struct {
	Choices  []string         // items on the to-do list
	Cursor   int              // which to-do list item our Cursor is pointing at
	Selected map[int]struct{} // which to-do items are Selected
}

func InitialModel() Model {
	return Model{
		Choices:  []string{"buy groceries", "walk the dog", "watch a movie"},
		Cursor:   0,
		Selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no i/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the Cursor up
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else if m.Cursor == 0 {
				m.Cursor = len(m.Choices) - 1
			}

		// The "down" and "j" keys move the Cursor down
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			} else if m.Cursor == len(m.Choices)-1 {
				m.Cursor = 0
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the Selected state for the item that the Cursor is pointing at.
		case "enter", " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{}
			}
		}
	case int:
		fmt.Println("Received int message:", msg)
		return m, nil
	}

	// Return the updated Model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our Choices
	for i, choice := range m.Choices {

		// Is the Cursor pointing at this choice?
		cursor := " " // no Cursor
		if m.Cursor == i {
			cursor = ">" // Cursor!
		}

		// Is this choice Selected?
		checked := " " // not Selected
		if _, ok := m.Selected[i]; ok {
			checked = "x" // Selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
