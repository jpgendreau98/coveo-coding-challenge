package frontend

import (
	"fmt"

	"projet-devops-coveo/pkg/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Command struct {
	Command string
	Options *util.CliOptions
	Done    bool
}

var RunCommand = &Command{
	Options: &util.CliOptions{
		OutputOptions: &util.OutputOptions{},
	},
}

type base struct {
	choices []string
	cursor  int
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
			if m.choices[m.cursor] == "AWS S3" {
				RunCommand.Command = "AWSS3"
				return S3GroupBy().Update(nil)
			}
		}
	}
	return m, nil
}
