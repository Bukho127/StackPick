package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const appVersion = "0.1"

func printStartupBanner() {
	fmt.Println()
	fmt.Println(`  .   ____      _            __ _ _`)
	fmt.Println(` /\ / ___|__ _| | ___ _ __ / _(_) | ___`)
	fmt.Println(`( ( ) |  _ / _' | |/ _ \ '__| |_| | |/ _ \`)
	fmt.Println(` \/| |_| | (_| | |  __/ | |  _| | |  __/`)
	fmt.Println(`  '  \____|\__,_|_|\___|_| |_| |_|_|\___|`)
	fmt.Printf("        StackPick v%s\n", appVersion)
	fmt.Println()
	fmt.Println("Welcome to StackPick. Choose Frontend or Backend to get started.")
}

func main() {
	printStartupBanner()
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
