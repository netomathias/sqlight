package statusbar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netomathias/sqlight/internal/tui/shared"
)

// Model is the status bar component.
type Model struct {
	dbPath   string
	readOnly bool
	table    string
	message  string
	isError  bool
	width    int
	styles   shared.StatusBarStyles
}

// New creates a new status bar.
func New(dbPath string, readOnly bool, styles shared.StatusBarStyles) Model {
	return Model{
		dbPath:   dbPath,
		readOnly: readOnly,
		styles:   styles,
	}
}

// SetWidth sets the status bar width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetTable sets the currently displayed table name.
func (m *Model) SetTable(name string) {
	m.table = name
}

// SetMessage sets a temporary message.
func (m *Model) SetMessage(msg string, isError bool) {
	m.message = msg
	m.isError = isError
}

// ClearMessage clears the temporary message.
func (m *Model) ClearMessage() {
	m.message = ""
	m.isError = false
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	left := m.styles.Key.Render("DB: ") + m.styles.Value.Render(m.dbPath)

	if m.readOnly {
		left += "  " + m.styles.Error.Render("[RO]")
	}

	if m.table != "" {
		left += "  " + m.styles.Key.Render("Table: ") + m.styles.Value.Render(m.table)
	}

	right := ""
	if m.message != "" {
		if m.isError {
			right = m.styles.Error.Render(m.message)
		} else {
			right = m.styles.Success.Render(m.message)
		}
	} else {
		right = m.styles.Value.Render("? help  q quit")
	}

	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	gap := m.width - leftWidth - rightWidth - 2
	if gap < 1 {
		gap = 1
	}

	bar := fmt.Sprintf("%s%*s%s", left, gap, "", right)

	return m.styles.Bar.Width(m.width).Render(bar)
}
