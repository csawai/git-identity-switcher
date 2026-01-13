package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/csawai/git-identity-switcher/internal/config"
	"github.com/csawai/git-identity-switcher/internal/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var listIdentitiesCmd = &cobra.Command{
	Use:   "list identities",
	Short: "List all stored identities",
	Long:  "Display all configured identities (no secrets shown).",
	Run: func(cmd *cobra.Command, args []string) {
		if err := listIdentities(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func listIdentities() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if len(cfg.Identities) == 0 {
		fmt.Println(ui.WarningBox.Render("⚠️  No identities configured.\n\nUse 'gitx add identity' to add your first identity."))
		return nil
	}

	// Build table
	var rows []string
	header := ui.TableHeaderStyle.Render("Alias") + " │ " +
		ui.TableHeaderStyle.Render("Name") + " │ " +
		ui.TableHeaderStyle.Render("Email") + " │ " +
		ui.TableHeaderStyle.Render("Auth") + " │ " +
		ui.TableHeaderStyle.Render("Status")
	
	divider := strings.Repeat("─", lipgloss.Width(header))
	rows = append(rows, header)
	rows = append(rows, divider)

	for _, id := range cfg.Identities {
		authIcon := ui.GetAuthIcon(id.AuthMethod)
		row := ui.TableRowStyle.Render("🔹 "+id.Alias) + " │ " +
			ui.TableRowStyle.Render(id.Name) + " │ " +
			ui.TableRowStyle.Render(id.Email) + " │ " +
			ui.TableRowStyle.Render(authIcon+" "+strings.ToUpper(id.AuthMethod)) + " │ " +
			ui.TableRowStyle.Render(ui.StatusBound)
		rows = append(rows, row)
	}

	content := strings.Join(rows, "\n")
	box := ui.BoxStyle.Copy().
		Width(lipgloss.Width(header) + 4).
		Render("🔐 Configured Identities\n\n" + content)

	fmt.Println(box)
	return nil
}

