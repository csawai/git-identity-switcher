package main

import (
	"fmt"
	"os"

	"github.com/csawai/gitx/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI",
	Long:  "Launch a text-based user interface for managing identities and binding repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := launchTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

type model struct {
	identities []config.Identity
	selected   int
	quitting   bool
}

func initialModel() model {
	cfg, _ := config.LoadConfig()
	return model{
		identities: cfg.Identities,
		selected:   0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.identities)-1 {
				m.selected++
			}
		case "enter", " ":
			if len(m.identities) > 0 {
				// Bind to selected identity
				identity := m.identities[m.selected]
				if err := bindIdentity(identity.Alias); err == nil {
					return m, tea.Quit
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := "Select an identity to bind to this repository:\n\n"

	for i, id := range m.identities {
		cursor := " "
		if i == m.selected {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s (%s)\n", cursor, id.Alias, id.Email)
	}

	s += "\nPress q to quit, ↑/↓ to navigate, Enter to select"
	return s
}

func launchTUI() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if len(cfg.Identities) == 0 {
		fmt.Println("No identities configured. Use 'gitx add identity' first.")
		return nil
	}

	p := tea.NewProgram(initialModel())
	_, err = p.Run()
	return err
}

