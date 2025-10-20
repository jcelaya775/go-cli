package sort

import (
	"context"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"jcelaya775/go-cli/models"
	"jcelaya775/go-cli/services/sorting"
	"strconv"
	"time"
)

type Model struct {
	sortingAlgorithm models.SortingAlgorithm
	sortingIt        models.SortingAlgorithmIterator
	sortingItFactory sorting.SortingAlgorithmIteratorFactory
	nums             []int
	i                int
	j                int
	done             bool
	cancelled        bool
	pause            bool
	tickDuration     time.Duration
	lastTick         time.Time
	remainingTick    time.Duration
	spinner          spinner.Model
	ctxWithCancel    context.Context
	cancel           context.CancelFunc
}

func InitialModel(sortingAlgorithm models.SortingAlgorithm, nums []int, tickInterval time.Duration) (Model, error) {
	sortingItFactory, err := sorting.GetIteratorFactory(sortingAlgorithm)
	if err != nil {
		return Model{}, err
	}
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	sortingIt := sortingItFactory.New(ctxWithCancel, nums)

	s := spinner.New(
		spinner.WithSpinner(spinner.Monkey),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))),
	)

	return Model{
		sortingAlgorithm: sortingAlgorithm,
		sortingIt:        sortingIt,
		sortingItFactory: sortingItFactory,
		ctxWithCancel:    ctxWithCancel,
		cancel:           cancel,
		nums:             nums,
		i:                0,
		j:                1,
		done:             false,
		pause:            false,
		tickDuration:     tickInterval,
		lastTick:         time.Now(),
		remainingTick:    tickInterval,
		spinner:          s,
	}, nil
}

type TickMsg time.Time

func tickCmdWithDuration(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	m.lastTick = time.Now()
	m.remainingTick = m.tickDuration
	return tea.Batch(
		tickCmdWithDuration(m.tickDuration),
		m.spinner.Tick,
		m.sortingIt.NextCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case " ":
			m.pause = !m.pause
			if m.pause {
				m.remainingTick = m.tickDuration - time.Since(m.lastTick)
				if m.remainingTick < 0 {
					m.remainingTick = 0
				}
				return m, nil
			} else {
				m.lastTick = time.Now()
				return m, tea.Batch(
					tickCmdWithDuration(m.remainingTick),
					m.spinner.Tick,
				)
			}
		}
	case TickMsg:
		if m.pause {
			return m, nil
		}
		m.lastTick = time.Now()
		m.remainingTick = m.tickDuration
		return m, tea.Batch(
			tickCmdWithDuration(m.tickDuration),
			m.sortingIt.NextCmd(),
		)
	case models.SortingStateMsg:
		if msg.Done {
			m.done = true
			return m, tea.Quit
		}

		m.nums = msg.Nums
		m.i = msg.I
		m.j = msg.J
	case spinner.TickMsg:
		if m.pause {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) renderNums() string {
	s := ""

	if m.done {
		for _, num := range m.nums {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(" " + strconv.Itoa(num) + " ") // Green for sorted portion
		}
		return s
	}

	for i, num := range m.nums {
		if i == m.j {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Render(" " + strconv.Itoa(num) + " ") // Red for current element
		} else if i < m.i {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(" " + strconv.Itoa(num) + " ") // Green for sorted portion
		} else {
			s += lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Render(" " + strconv.Itoa(num) + " ") // Orange for unsorted portion
		}
	}
	return s
}

func (m Model) View() string {
	s := models.SortingAlgorithmNames[m.sortingAlgorithm] + "\n\n"
	if m.pause {
		s += "Paused " + m.spinner.View() + "\n\n"
	} else if m.done {
		s += "Complete!\n\n"
	} else {
		s += "Running " + m.spinner.View() + "\n\n"
	}

	s += m.renderNums()
	s += "\n\nPress space to pause/resume. Press ctrl+c to quit.\n"
	return s
}
