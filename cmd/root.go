package cmd

import (
	"fmt"
	"projet-devops-coveo/frontend"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "tool",
	}
	cmd.AddCommand(
		NewGuiCommand(),
		NewS3Command(),
	)
	return cmd
}

func NewGuiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gui",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := tea.NewProgram(frontend.InitialBase())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				return err
			}
			if frontend.RunCommand.Done {
				switch frontend.RunCommand.Command {
				case "AWSS3":
					fmt.Println(frontend.RunCommand.Options)
					err := RunS3Command(frontend.RunCommand.Options)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
	return cmd
}
