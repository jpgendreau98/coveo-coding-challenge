package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewS3Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "s3",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(AwsS3())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				return err
			}
			return nil
		},
	}
	return cmd
}
