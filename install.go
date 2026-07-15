package main

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

// installFinishedMsg is sent once the npm install process completes,
// carrying its combined stdout/stderr and any error.
type installFinishedMsg struct {
	output string
	err    error
}

// runInstall builds and runs `npm install <packages...>` in the current
// working directory (i.e. wherever the user launched the tool from).
func runInstall(packages []string) tea.Cmd {
	return func() tea.Msg {
		args := append([]string{"install"}, packages...)
		cmd := exec.Command("npm", args...)
		out, err := cmd.CombinedOutput()
		return installFinishedMsg{output: string(out), err: err}
	}
}
