package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type sessionState int

const (
	stateChoice sessionState = iota
	stateSelecting
	stateInstalling
	stateDone
)

type Model struct {
	selected map[string]bool
	items    []LibrarySelection
	cursor   int

	state      sessionState
	spinner    spinner.Model
	packages   []string
	installOut string
	installErr error
	choice     string
	catalog    []Category
}

type LibrarySelection struct {
	catIdx      int
	libIdx      int
	name        string
	description string
}

// NewModel starts with the initial frontend/backend choice screen.
func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return Model{
		state:   stateChoice,
		spinner: s,
		catalog: FrontendCategories,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func keyOf(catIdx, libIdx int) string {
	return fmt.Sprintf("%d-%d", catIdx, libIdx)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateChoice:
			return m.updateChoice(msg)
		case stateSelecting:
			return m.updateSelecting(msg)
		case stateDone:
			switch msg.String() {
			case "q", "ctrl+c", "enter":
				return m, tea.Quit
			}
		}
		return m, nil

	case spinner.TickMsg:
		if m.state == stateInstalling {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case installFinishedMsg:
		m.state = stateDone
		m.installOut = msg.output
		m.installErr = msg.err
		return m, nil
	}
	return m, nil
}

func (m Model) updateChoice(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "f", "F", "frontend", "Frontend":
		m.choice = "frontend"
		m.catalog = FrontendCategories
		m.loadItems()
		m.state = stateSelecting
		return m, nil
	case "b", "B", "backend", "Backend":
		m.choice = "backend"
		m.catalog = BackendCategories
		m.loadItems()
		m.state = stateSelecting
		return m, nil
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSelecting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyRunes {
		if len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if r >= '1' && r <= '9' {
				idx := int(r - '0')
				if idx <= len(m.items) {
					item := m.items[idx-1]
					key := keyOf(item.catIdx, item.libIdx)
					m.selected[key] = !m.selected[key]
					m.cursor = idx - 1
					return m, nil
				}
			}
		}
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		pkgs := m.collectPackages()
		if len(pkgs) == 0 {
			return m, nil
		}
		m.packages = pkgs
		m.state = stateInstalling
		return m, tea.Batch(m.spinner.Tick, runInstall(pkgs))
	}
	return m, nil
}

func (m *Model) loadItems() {
	m.items = nil
	for ci, cat := range m.catalog {
		for li := range cat.Libs {
			m.items = append(m.items, LibrarySelection{
				catIdx:      ci,
				libIdx:      li,
				name:        cat.Libs[li].Name,
				description: cat.Libs[li].Description,
			})
		}
	}
	m.selected = make(map[string]bool)
	m.cursor = 0
}

func (m *Model) moveCursor(dir int) {
	n := len(m.items)
	if n == 0 {
		return
	}
	i := m.cursor + dir
	if i < 0 {
		i = 0
	} else if i >= n {
		i = n - 1
	}
	m.cursor = i
}

func (m Model) collectPackages() []string {
	var pkgs []string
	for _, item := range m.items {
		key := keyOf(item.catIdx, item.libIdx)
		if m.selected[key] {
			pkgs = append(pkgs, m.catalog[item.catIdx].Libs[item.libIdx].Packages...)
		}
	}
	return pkgs
}

func (m Model) View() string {
	switch m.state {
	case stateChoice:
		return m.viewChoice()
	case stateInstalling:
		return m.viewInstalling()
	case stateDone:
		return m.viewDone()
	default:
		return m.viewSelecting()
	}
}

func (m Model) viewChoice() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("⚛  StackPick") + "\n")
	b.WriteString(headerStyle.Render(fmt.Sprintf("Version %s", appVersion)) + "\n\n")
	b.WriteString(headerStyle.Render("Frontend or Backend?") + "\n\n")
	b.WriteString("  [F] Frontend\n")
	b.WriteString("  [B] Backend\n\n")
	b.WriteString(helpStyle.Render("press f or b to choose"))
	return b.String()
}

func (m Model) viewSelecting() string {
	var rows [][]string
	for idx, item := range m.items {
		checked := "[ ]"
		if m.selected[keyOf(item.catIdx, item.libIdx)] {
			checked = "[x]"
		}
		label := fmt.Sprintf("%d. %s", idx+1, item.name)
		rows = append(rows, []string{checked, label, item.description, m.catalog[item.catIdx].Name})
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		Headers("", "Library", "Description", "Category").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
			case row == m.cursor+1:
				return lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
			case col == 0:
				if rows[row][0] == "[x]" {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
				}
				return lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
			default:
				return lipgloss.NewStyle()
			}
		})

	count := 0
	for _, v := range m.selected {
		if v {
			count++
		}
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("⚛  StackPick") + "\n\n")
	b.WriteString(t.Render())
	b.WriteString("\n" + helpStyle.Render(fmt.Sprintf(
		"1-9 select · ↑/↓ move · enter install · q quit (%d selected)", count)))
	return b.String()
}

func (m Model) viewInstalling() string {
	pkgList := strings.Join(m.packages, ", ")
	return fmt.Sprintf(
		"\n%s Installing with npm:\n\n  %s\n\n%s\n",
		m.spinner.View(), pkgList, helpStyle.Render("this may take a moment..."),
	)
}

func (m Model) viewDone() string {
	var b strings.Builder
	if m.installErr != nil {
		b.WriteString(errorStyle.Render("✗ Installation failed") + "\n\n")
	} else {
		b.WriteString(successStyle.Render("✓ Installation complete!") + "\n\n")
	}
	b.WriteString(m.installOut + "\n")
	b.WriteString(helpStyle.Render("\npress enter or q to exit"))
	return b.String()
}
