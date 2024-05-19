package frontend

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type AwsS3ReturnEmptyBuckets struct {
	choices  []string
	cursor   int // which to-do list item our cursor is pointing at
	selected map[int]struct{}
}

func (m AwsS3ReturnEmptyBuckets) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func S3ReturnsEmptyBucket() AwsS3ReturnEmptyBuckets {
	return AwsS3ReturnEmptyBuckets{
		choices:  []string{"yes", "no"},
		selected: make(map[int]struct{}),
	}
}

func (m AwsS3ReturnEmptyBuckets) View() string {
	s := "Do you want to display buckets that are empty?\n\n"

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

func (m AwsS3ReturnEmptyBuckets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			if m.choices[m.cursor] == "yes" {
				RunCommand.Options.ReturnEmptyBuckets = true
			}
			RunCommand.Done = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	}

	// Return the updated AwsS3ReturnEmptyBuckets to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
