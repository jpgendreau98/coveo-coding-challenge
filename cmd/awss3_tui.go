package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AwsS3Tui struct {
	choices   map[string][]string
	nameInput string
	listInput string
	event     string
	input     textinput.Model
	cursor    int // which to-do list item our cursor is pointing at
	selected  map[int]struct{}
}

func (m AwsS3Tui) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

func AwsS3() AwsS3Tui {
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Filter By bucket Name"
	ti.CharLimit = 156
	ti.Width = 20
	ti.Focus()
	return AwsS3Tui{
		choices: map[string][]string{
			"Group By":                {"region"},
			"Order By (Inc)":          {"name", "price", "storage-class", "size"},
			"Order By (Dec)":          {"name", "price", "storage-class", "size"},
			"Filter By Storage Class": {"STANDARD", "REDUCED_REDUNDANCY", "GLACIER", "STANDARD_IA", "INTELLIGENT_TIERING", "DEEP_ARCHIVE", "GLACIER_IR"},
			"Returns Emtpy Buckets":   {"yes", "no"},
		},
		input:    ti,
		selected: make(map[int]struct{}),
	}
}

func (m AwsS3Tui) View() string {
	// The header
	s := "What services do you want to inspect?\n\n"

	// Iterate over our choices
	for key, list := range m.choices {
		s += key + "\n"
		for i, choice := range list {

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
	}

	b := &strings.Builder{}
	b.WriteString("Enter your event:\n")
	b.WriteString(m.input.View())
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (m AwsS3Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

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
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	case tea.WindowSizeMsg:
		_, _ = msg.Width, msg.Height
	}
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	m.nameInput = m.input.Value()

	// Return the updated AwsS3Tui to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, tea.Batch(cmds...)
}
