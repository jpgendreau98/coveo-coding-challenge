package frontend

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type AwsS3OrderBy struct {
	choices  []string
	cursor   int // which to-do list item our cursor is pointing at
	selected map[int]struct{}
}

func (m AwsS3OrderBy) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func S3OrderBy() AwsS3OrderBy {
	return AwsS3OrderBy{
		choices:  []string{"name (INC)", "price (INC)", "storage-class (INC)", "size (INC)", "name (DEC)", "price (DEC)", "storage-class (DEC)", "size (DEC)"},
		selected: make(map[int]struct{}),
	}
}

func (m AwsS3OrderBy) View() string {
	// The header
	s := "I want to order by: _____ ?\n\n"

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

func (m AwsS3OrderBy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			switch m.choices[m.cursor] {
			case "name (INC)":
				RunCommand.Options.OrderByINC = "name"
			case "price (INC)":
				RunCommand.Options.OrderByINC = "price"
			case "storage-class (INC)":
				RunCommand.Options.OrderByINC = "storage-class"
			case "size (INC)":
				RunCommand.Options.OrderByINC = "size"
			case "name (DEC)":
				RunCommand.Options.OrderByDEC = "name"
			case "price (DEC)":
				RunCommand.Options.OrderByDEC = "price"
			case "storage-class (DEC)":
				RunCommand.Options.OrderByDEC = "storage-class"
			case "size (DEC)":
				RunCommand.Options.OrderByDEC = "size"
			}
			return S3FilterSC().Update(nil)

		}
	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	}

	// Return the updated AwsS3OrderBy to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
