package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(InitialBase())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(
		NewS3Cmd(),
	)
	return cmd
}
