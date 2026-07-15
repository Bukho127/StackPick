package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type rowKind int

const (
	rowHeader rowKind = iota
	rowLib
)

// row is a flattened line in the checklist: either a category header
// (not selectable) or a library entry (selectable, toggled with space).
type row struct {
	kind   rowKind
	header string
	catIdx int
	libIdx int
}

type sessionState int

const (
	stateChoice sessionState = iota
	stateSelecting
	stateInstalling
	stateDone
)

type Model struct {
	rows     []row
	cursor   int
	selected map[string]bool // key "catIdx-libIdx" -> selected?

	state      sessionState
	spinner    spinner.Model
	packages   []string
	installOut string
	installErr error
	choice     string
	catalog    []Category
}

func buildRows(cats []Category) []row {
	var rows []row
	for ci, cat := range cats {
		rows = append(rows, row{kind: rowHeader, header: cat.Name})
		for li := range cat.Libs {
			rows = append(rows, row{kind: rowLib, catIdx: ci, libIdx: li})
		}
	}
	return rows
}

func (m *Model) startSelection(choice string) {
	m.choice = choice
	m.catalog = FrontendCategories
	if choice == "backend" {
		m.catalog = BackendCategories
	}
	m.rows = buildRows(m.catalog)
	m.selected = make(map[string]bool)
	m.cursor = 0
	for i, r := range m.rows {
		if r.kind == rowLib {
			m.cursor = i
			break
		}
	}
	m.state = stateSelecting
}

// NewModel starts at the initial choice screen.
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
		m.startSelection("frontend")
		return m, nil
	case "b", "B", "backend", "Backend":
		m.startSelection("backend")
		return m, nil
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateSelecting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		m.moveCursor(-1)

	case "down", "j":
		m.moveCursor(1)

	case " ":
		r := m.rows[m.cursor]
		if r.kind == rowLib {
			key := keyOf(r.catIdx, r.libIdx)
			m.selected[key] = !m.selected[key]
		}

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

func (m *Model) moveCursor(dir int) {
	n := len(m.rows)
	i := m.cursor
	for {
		i += dir
		if i < 0 || i >= n {
			return
		}
		if m.rows[i].kind == rowLib {
			m.cursor = i
			return
		}
	}
}

func (m Model) collectPackages() []string {
	var pkgs []string
	for ci, cat := range m.catalog {
		for li := range cat.Libs {
			if m.selected[keyOf(ci, li)] {
				pkgs = append(pkgs, m.catalog[ci].Libs[li].Packages...)
			}
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
	b.WriteString(titleStyle.Render("⚛  StackPick") + "\n\n")
	b.WriteString(headerStyle.Render("Frontend or Backend ?") + "\n\n")
	b.WriteString("  [F] Frontend\n")
	b.WriteString("  [B] Backend\n\n")
	b.WriteString(helpStyle.Render("press f or b to choose"))
	return b.String()
}

func (m Model) viewSelecting() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("⚛  Frontend Framework Picker") + "\n\n")
	b.WriteString(headerStyle.Render("Category") + "  " + headerStyle.Render("Framework") + "  " + headerStyle.Render("Status") + "\n")
	b.WriteString(strings.Repeat("-", 70) + "\n")

	for i, r := range m.rows {
		if r.kind == rowHeader {
			b.WriteString("\n" + headerStyle.Render(r.header) + "\n")
			continue
		}
		lib := m.catalog[r.catIdx].Libs[r.libIdx]
		checked := "[ ]"
		if m.selected[keyOf(r.catIdx, r.libIdx)] {
			checked = "[x]"
		}
		line := fmt.Sprintf("%-16s  %-24s  %s", m.catalog[r.catIdx].Name, lib.Name, checked)
		if i == m.cursor {
			b.WriteString(cursorStyle.Render("› "+line) + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}

	count := 0
	for _, v := range m.selected {
		if v {
			count++
		}
	}
	b.WriteString("\n" + helpStyle.Render(fmt.Sprintf(
		"↑/↓ move · space select (%d selected) · enter install · q quit", count)))
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
