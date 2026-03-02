package editor

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netomathias/sqlight/internal/tui/shared"
)

// Model is the SQL editor component.
type Model struct {
	textarea textarea.Model
	visible  bool
	focused  bool
	width    int
	height   int
	styles   shared.EditorStyles
	keys     shared.KeyMap
}

// New creates a new editor model.
func New(styles shared.EditorStyles, keys shared.KeyMap) Model {
	ta := textarea.New()
	ta.Placeholder = "Enter SQL query..."
	ta.ShowLineNumbers = false
	ta.SetHeight(3)
	ta.CharLimit = 4096

	return Model{
		textarea: ta,
		styles:   styles,
		keys:     keys,
	}
}

// SetSize sets the editor dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.textarea.SetWidth(w - 4)
	m.textarea.SetHeight(h - 4)
}

// SetVisible sets visibility.
func (m *Model) SetVisible(v bool) {
	m.visible = v
}

// Visible returns whether the editor is visible.
func (m Model) Visible() bool {
	return m.visible
}

// SetFocused sets focus state.
func (m *Model) SetFocused(f bool) {
	m.focused = f
	if f {
		m.textarea.Focus()
	} else {
		m.textarea.Blur()
	}
}

// Focused returns whether the editor is focused.
func (m Model) Focused() bool {
	return m.focused
}

// Value returns the current SQL text.
func (m Model) Value() string {
	return strings.TrimSpace(m.textarea.Value())
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible || !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.ExecSQL) {
			sql := m.Value()
			if sql != "" {
				return m, func() tea.Msg {
					return execSQLMsg{SQL: sql}
				}
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// execSQLMsg is an internal message to trigger SQL execution.
type execSQLMsg struct {
	SQL string
}

// ExecSQLMsg returns the public type for pattern matching in app.go.
func IsExecSQL(msg tea.Msg) (string, bool) {
	if m, ok := msg.(execSQLMsg); ok {
		return m.SQL, true
	}
	return "", false
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	title := m.styles.Title.Render("SQL Editor") + "  " +
		lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("alt+enter to execute")

	content := title + "\n" + m.textarea.View()

	borderStyle := m.styles.Border.Width(m.width - 2)
	if m.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("4"))
	}

	return borderStyle.Height(m.height - 2).Render(content)
}
