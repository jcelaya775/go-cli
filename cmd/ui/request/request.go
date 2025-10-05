package request

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type Model struct {
	status int
	err    error
}

// `Cmd`s are functions that perform i/O and return a `tea.Msg` (like this)
func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}
	return statusMsg(res.StatusCode)
}

// Note: `Msg`s can be any type, even an empty struct
type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// We want to call the server right away, so return it from Init
// We return the `Cmd` without calling it (bubbletea runtime will take care of that)
func (m Model) Init() tea.Cmd {
	return checkServer
}

// Internally, `Cmd`s run async in a goroutine, and the `Msg` returned is collected and sent here
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit

	case errMsg:
		// Note: Since, errMsg implements Error, it implements error
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// Ignore any other messages
	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Keeps getting displayed until status state is updated (and appended to s), or err is received
	s := fmt.Sprintf("Checking %s...", url)

	if m.status > 0 {
		s += fmt.Sprintf(" %d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n\n"
}
