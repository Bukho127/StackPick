package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const appVersion = "0.1.0"

func printBanner() {
	art := `  .   ____      _            __ _ _
 /\ / ___|__ _| | ___ _ __ / _(_) | ___
( ( ) |  _ / _' | |/ _ \ '__| |_| | |/ _ \
 \/| |_| | (_| | |  __/ | |  _| | |  __/
  '  \____|\__,_|_|\___|_| |_| |_|_|\___|`

	banner := lipgloss.JoinVertical(lipgloss.Center,
		art,
		fmt.Sprintf("StackPick v%s", appVersion),
	)

	boxed := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("212")).
		Padding(1, 3)

	fmt.Println(boxed.Render(bannerStyle.Render(banner)))
}

func main() {
	printBanner()

	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
