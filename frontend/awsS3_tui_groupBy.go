package frontend

import (
	"fmt"
	"projet-devops-coveo/pkg/util"

	tea "github.com/charmbracelet/bubbletea"
)

type AwsS3GroupBy struct {
	choices  []string
	cursor   int // which to-do list item our cursor is pointing at
	selected map[int]struct{}
}

func (m AwsS3GroupBy) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

var Options = util.CliOptions{}

func S3GroupBy() AwsS3GroupBy {
	return AwsS3GroupBy{
		choices:  []string{"region", "none"},
		selected: make(map[int]struct{}),
	}
}

func (m AwsS3GroupBy) View() string {
	// The header
	s := "I want to group by: _____ ?\n\n"

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

func (m AwsS3GroupBy) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.choices[m.cursor] == "region" {
				RunCommand.Options.OutputOptions.GroupBy = "region"
			}
			return S3OrderBy().Update(nil)
		}
	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	}

	// Return the updated AwsS3GroupBy to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
