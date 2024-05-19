package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type base struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func (m base) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func InitialBase() base {
	return base{
		choices: []string{"AWS S3", "Azure Blob (Work in progress)", "GCP (Work in progress)"},
	}
}

func (m base) View() string {
	// The header
	s := "What services do you want to inspect?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (m base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter":
			if m.choices[m.cursor] == "AWS S3" {
				return S3GroupBy().Update(msg)
			}
		}
	}

	// Return the updated base to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
