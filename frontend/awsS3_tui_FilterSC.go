package frontend

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type AwsS3FilterSC struct {
	choices  []string
	cursor   int // which to-do list item our cursor is pointing at
	selected map[int]struct{}
}

func (m AwsS3FilterSC) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func S3FilterSC() AwsS3FilterSC {
	return AwsS3FilterSC{
		choices:  []string{"STANDARD", "REDUCED_REDUNDANCY", "GLACIER", "STANDARD_IA", "INTELLIGENT_TIERING", "DEEP_ARCHIVE", "GLACIER_IR"},
		selected: make(map[int]struct{}),
	}
}

func (m AwsS3FilterSC) View() string {
	// The header
	s := "I want to filter by Storage-Class: _____ \n\nPress (x) to select.\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}
		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress Enter to continue.\n"
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (m AwsS3FilterSC) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "x":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			for k := range m.selected {
				RunCommand.Options.FilterByStorageClass = append(RunCommand.Options.FilterByStorageClass, m.choices[k])
			}
			return S3SelectRegion().Update(nil)

		}

	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	}

	// Return the updated AwsS3FilterSC to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
