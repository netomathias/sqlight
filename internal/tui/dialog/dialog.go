package dialog

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Kind represents the type of dialog.
type Kind int

const (
	KindNone Kind = iota
	KindEdit
	KindInsert
	KindConfirmDelete
)

// Model is the dialog component for editing, inserting, and confirming.
type Model struct {
	kind   Kind
	title  string
	fields []field
	cursor int

	// For delete confirmation
	deleteTable string
	deletePKCol string
	deletePKVal string

	width  int
	height int
	styles DialogStyles
}

// DialogStyles holds styles for the dialog component.
type DialogStyles struct {
	Border  lipgloss.Style
	Title   lipgloss.Style
	Input   lipgloss.Style
	Button  lipgloss.Style
	Overlay lipgloss.Style
}

type field struct {
	label string
	input textinput.Model
}

// New creates a new dialog model.
func New(styles DialogStyles) Model {
	return Model{
		styles: styles,
	}
}

// Active returns whether any dialog is shown.
func (m Model) Active() bool {
	return m.kind != KindNone
}

// Kind returns the current dialog kind.
func (m Model) Kind() Kind {
	return m.kind
}

// SetSize sets dialog dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// ShowEdit shows the edit cell dialog.
func (m *Model) ShowEdit(col, oldVal string) {
	m.kind = KindEdit
	m.title = fmt.Sprintf("Edit: %s", col)
	ti := textinput.New()
	ti.SetValue(oldVal)
	ti.Focus()
	ti.Width = 40
	m.fields = []field{{label: col, input: ti}}
	m.cursor = 0
}

// ShowInsert shows the insert row dialog with fields for each column.
func (m *Model) ShowInsert(columns []string) {
	m.kind = KindInsert
	m.title = "Insert Row"
	m.fields = make([]field, len(columns))
	for i, col := range columns {
		ti := textinput.New()
		ti.Placeholder = col
		ti.Width = 40
		if i == 0 {
			ti.Focus()
		}
		m.fields[i] = field{label: col, input: ti}
	}
	m.cursor = 0
}

// ShowConfirmDelete shows a delete confirmation dialog.
func (m *Model) ShowConfirmDelete(table, pkCol, pkVal string) {
	m.kind = KindConfirmDelete
	m.title = "Confirm Delete"
	m.deleteTable = table
	m.deletePKCol = pkCol
	m.deletePKVal = pkVal
	m.fields = nil
	m.cursor = 0
}

// Close closes the dialog.
func (m *Model) Close() {
	m.kind = KindNone
	m.fields = nil
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.Active() {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.Close()
			return m, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			return m, m.submit()

		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			if m.kind == KindInsert && len(m.fields) > 1 {
				m.fields[m.cursor].input.Blur()
				m.cursor = (m.cursor + 1) % len(m.fields)
				m.fields[m.cursor].input.Focus()
				return m, nil
			}

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			if m.kind == KindInsert && len(m.fields) > 1 {
				m.fields[m.cursor].input.Blur()
				m.cursor = (m.cursor - 1 + len(m.fields)) % len(m.fields)
				m.fields[m.cursor].input.Focus()
				return m, nil
			}
		}
	}

	// Update the focused text input
	if m.kind != KindConfirmDelete && len(m.fields) > 0 {
		var cmd tea.Cmd
		m.fields[m.cursor].input, cmd = m.fields[m.cursor].input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) submit() tea.Cmd {
	switch m.kind {
	case KindEdit:
		if len(m.fields) > 0 {
			val := m.fields[0].input.Value()
			return func() tea.Msg {
				return EditSubmitMsg{Value: val}
			}
		}
	case KindInsert:
		data := make(map[string]string)
		for _, f := range m.fields {
			val := f.input.Value()
			if val != "" {
				data[f.label] = val
			}
		}
		return func() tea.Msg {
			return InsertSubmitMsg{Data: data}
		}
	case KindConfirmDelete:
		return func() tea.Msg {
			return ConfirmDeleteMsg{
				Table: m.deleteTable,
				PKCol: m.deletePKCol,
				PKVal: m.deletePKVal,
			}
		}
	}
	return nil
}

// EditSubmitMsg carries the edited value.
type EditSubmitMsg struct {
	Value string
}

// InsertSubmitMsg carries insert field values.
type InsertSubmitMsg struct {
	Data map[string]string
}

// ConfirmDeleteMsg confirms a delete operation.
type ConfirmDeleteMsg struct {
	Table string
	PKCol string
	PKVal string
}

// View implements tea.Model.
func (m Model) View() string {
	if !m.Active() {
		return ""
	}

	var b strings.Builder

	title := m.styles.Title.Render(m.title)
	b.WriteString(title)
	b.WriteString("\n\n")

	switch m.kind {
	case KindEdit, KindInsert:
		for i, f := range m.fields {
			label := f.label
			if i == m.cursor {
				label = "▸ " + label
			} else {
				label = "  " + label
			}
			b.WriteString(lipgloss.NewStyle().Bold(true).Render(label))
			b.WriteString("\n")
			b.WriteString(f.input.View())
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("enter: submit  esc: cancel  tab: next field"))

	case KindConfirmDelete:
		b.WriteString(fmt.Sprintf("Delete row where %s = %s?\n\n", m.deletePKCol, m.deletePKVal))
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true).Render("This cannot be undone."))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("enter: confirm  esc: cancel"))
	}

	dialogWidth := 50
	if dialogWidth > m.width-4 {
		dialogWidth = m.width - 4
	}

	return m.styles.Border.Width(dialogWidth).Render(b.String())
}
