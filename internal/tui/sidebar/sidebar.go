package sidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netomathias/sqlight/internal/db"
	"github.com/netomathias/sqlight/internal/tui/shared"
)

// Model is the sidebar component showing tables and views.
type Model struct {
	tables   []db.TableInfo
	cursor   int
	focused  bool
	width    int
	height   int
	styles   shared.SidebarStyles
	keys     shared.KeyMap
}

// New creates a new sidebar model.
func New(styles shared.SidebarStyles, keys shared.KeyMap) Model {
	return Model{
		styles: styles,
		keys:   keys,
	}
}

// SetTables updates the list of tables.
func (m *Model) SetTables(tables []db.TableInfo) {
	m.tables = tables
	if m.cursor >= len(tables) {
		m.cursor = max(0, len(tables)-1)
	}
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// SetFocused sets focus state.
func (m *Model) SetFocused(f bool) {
	m.focused = f
}

// Focused returns whether the sidebar is focused.
func (m Model) Focused() bool {
	return m.focused
}

// SelectedTable returns the currently selected table name, or empty if none.
func (m Model) SelectedTable() string {
	if len(m.tables) == 0 {
		return ""
	}
	return m.tables[m.cursor].Name
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.tables)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Enter):
			if len(m.tables) > 0 {
				t := m.tables[m.cursor]
				return m, func() tea.Msg {
					return shared.TableSelectedMsg{Name: t.Name, Type: t.Type}
				}
			}
		}
	}

	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	var b strings.Builder

	title := m.styles.Title.Render("Tables")
	b.WriteString(title)
	b.WriteString("\n")

	contentHeight := m.height - 4 // border + title + padding
	if contentHeight < 1 {
		contentHeight = 1
	}

	// Calculate visible window
	start := 0
	if m.cursor >= contentHeight {
		start = m.cursor - contentHeight + 1
	}
	end := start + contentHeight
	if end > len(m.tables) {
		end = len(m.tables)
	}

	for i := start; i < end; i++ {
		t := m.tables[i]
		badge := m.styles.TypeBadge.Render(fmt.Sprintf("[%s]", t.Type[:1]))
		name := t.Name

		var line string
		if i == m.cursor && m.focused {
			line = m.styles.SelectedItem.Render(fmt.Sprintf("%s %s", badge, name))
		} else {
			line = m.styles.Item.Render(fmt.Sprintf("%s %s", badge, name))
		}
		b.WriteString(line)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	borderStyle := m.styles.Border.Width(m.width - 2)
	if m.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("4"))
	}

	return borderStyle.Height(m.height - 2).Render(b.String())
}
