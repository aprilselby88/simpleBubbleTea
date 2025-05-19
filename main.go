package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	inputs     []textinput.Model
	focusIndex int
	done       bool
}

func initialModel() model {
	inputs := make([]textinput.Model, 2)
	placeholders := []string{"First value", "Second value"}

	for i := range inputs {
		t := textinput.New()
		t.Placeholder = placeholders[i]
		t.Prompt = "> "
		t.Focus()
		inputs[i] = t
	}

	inputs[1].Blur() // Only focus on the first input at start

	return model{
		inputs:     inputs,
		focusIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.focusIndex < len(m.inputs)-1 {
				m.focusIndex++
				for i := range m.inputs {
					if i == m.focusIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			} else {
				m.done = true
			}
		case "q":
			return m, tea.Quit
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.done {
		first := m.inputs[0].Value()
		second := m.inputs[1].Value()

		table := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Padding(1, 2)
		header := fmt.Sprintf("%-20s %-20s\n", "Input 1", "Input 2")
		values := fmt.Sprintf("%-20s %-20s\n", first, second)

		paragraph := fmt.Sprintf(
			"\nYou entered '%s' and '%s' 🎉",
			first, second,
		)

		return table.Render(header + values) + paragraph + "\n"
	}

	s := "Enter two values:\n\n"
	for i := range m.inputs {
		s += m.inputs[i].View() + "\n"
	}
	s += "\n(press Enter to continue, Ctrl+C to quit)\n"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
