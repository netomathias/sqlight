package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netomathias/sqlight/internal/db"
	"github.com/netomathias/sqlight/internal/tui/shared"
)

const defaultPerPage = 50

// Model is the data table component.
type Model struct {
	tableName string
	columns   []string
	rows      [][]string
	total     int64
	page      int
	perPage   int

	cursorRow int
	cursorCol int
	colOffset int // horizontal scroll offset

	focused bool
	width   int
	height  int
	styles  shared.TableStyles
	keys    shared.KeyMap
}

// New creates a new table model.
func New(styles shared.TableStyles, keys shared.KeyMap) Model {
	return Model{
		styles:  styles,
		keys:    keys,
		perPage: defaultPerPage,
	}
}

// SetData sets the table data from a query result.
func (m *Model) SetData(table string, result *db.QueryResult, page int) {
	m.tableName = table
	m.columns = result.Columns
	m.rows = result.Rows
	m.total = result.Total
	m.page = page
	m.cursorRow = 0
	m.cursorCol = 0
	m.colOffset = 0
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

// Focused returns whether the table is focused.
func (m Model) Focused() bool {
	return m.focused
}

// TableName returns the current table name.
func (m Model) TableName() string {
	return m.tableName
}

// Page returns the current page number.
func (m Model) Page() int {
	return m.page
}

// PerPage returns the page size.
func (m Model) PerPage() int {
	return m.perPage
}

// TotalPages returns the total number of pages.
func (m Model) TotalPages() int {
	if m.total <= 0 {
		return 1
	}
	return int((m.total + int64(m.perPage) - 1) / int64(m.perPage))
}

// SelectedRow returns the currently selected row index and data.
func (m Model) SelectedRow() (int, []string) {
	if m.cursorRow >= len(m.rows) {
		return -1, nil
	}
	return m.cursorRow, m.rows[m.cursorRow]
}

// SelectedCol returns the currently selected column index and name.
func (m Model) SelectedCol() (int, string) {
	if len(m.columns) == 0 {
		return -1, ""
	}
	return m.cursorCol, m.columns[m.cursorCol]
}

// Columns returns the column names.
func (m Model) Columns() []string {
	return m.columns
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
			if m.cursorRow > 0 {
				m.cursorRow--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursorRow < len(m.rows)-1 {
				m.cursorRow++
			}
		case key.Matches(msg, m.keys.Left):
			if m.cursorCol > 0 {
				m.cursorCol--
			}
		case key.Matches(msg, m.keys.Right):
			if m.cursorCol < len(m.columns)-1 {
				m.cursorCol++
			}
		case key.Matches(msg, m.keys.Home):
			m.cursorCol = 0
			m.colOffset = 0
		case key.Matches(msg, m.keys.End):
			if len(m.columns) > 0 {
				m.cursorCol = len(m.columns) - 1
			}
		case key.Matches(msg, m.keys.NextPage):
			if m.page < m.TotalPages()-1 {
				m.page++
				return m, func() tea.Msg { return shared.RefreshTableMsg{} }
			}
		case key.Matches(msg, m.keys.PrevPage):
			if m.page > 0 {
				m.page--
				return m, func() tea.Msg { return shared.RefreshTableMsg{} }
			}
		case key.Matches(msg, m.keys.EditCell):
			return m, m.editCellCmd()
		case key.Matches(msg, m.keys.InsertRow):
			if m.tableName != "" {
				return m, func() tea.Msg {
					return shared.InsertRowMsg{Table: m.tableName}
				}
			}
		case key.Matches(msg, m.keys.DeleteRow):
			return m, m.deleteRowCmd()
		}
	}

	return m, nil
}

func (m Model) editCellCmd() tea.Cmd {
	if len(m.rows) == 0 || len(m.columns) == 0 {
		return nil
	}
	// Find PK column (first column by convention)
	pkCol := m.columns[0]
	pkVal := m.rows[m.cursorRow][0]
	col := m.columns[m.cursorCol]
	oldVal := m.rows[m.cursorRow][m.cursorCol]
	table := m.tableName

	return func() tea.Msg {
		return shared.EditCellMsg{
			Table:  table,
			PKCol:  pkCol,
			PKVal:  pkVal,
			Col:    col,
			OldVal: oldVal,
		}
	}
}

func (m Model) deleteRowCmd() tea.Cmd {
	if len(m.rows) == 0 || len(m.columns) == 0 {
		return nil
	}
	pkCol := m.columns[0]
	pkVal := m.rows[m.cursorRow][0]
	table := m.tableName

	return func() tea.Msg {
		return shared.DeleteRowMsg{
			Table: table,
			PKCol: pkCol,
			PKVal: pkVal,
		}
	}
}

// View implements tea.Model.
func (m Model) View() string {
	if len(m.columns) == 0 {
		content := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true).
			Render("Select a table to view data")
		borderStyle := m.styles.Border.Width(m.width - 2)
		if m.focused {
			borderStyle = borderStyle.BorderForeground(lipgloss.Color("4"))
		}
		return borderStyle.Height(m.height - 2).Render(content)
	}

	// Calculate column widths
	colWidths := m.calcColumnWidths()

	var b strings.Builder

	// Header
	headerCells := m.renderRow(m.columns, colWidths, -1, true)
	b.WriteString(headerCells)
	b.WriteString("\n")

	// Separator
	var sepParts []string
	for _, cw := range colWidths {
		sepParts = append(sepParts, strings.Repeat("─", cw))
	}
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(strings.Join(sepParts, "┼")))
	b.WriteString("\n")

	// Data rows
	visibleRows := m.height - 6 // borders + header + separator + pagination
	if visibleRows < 1 {
		visibleRows = 1
	}

	rowStart := 0
	if m.cursorRow >= visibleRows {
		rowStart = m.cursorRow - visibleRows + 1
	}
	rowEnd := rowStart + visibleRows
	if rowEnd > len(m.rows) {
		rowEnd = len(m.rows)
	}

	for i := rowStart; i < rowEnd; i++ {
		line := m.renderRow(m.rows[i], colWidths, i, false)
		b.WriteString(line)
		if i < rowEnd-1 {
			b.WriteString("\n")
		}
	}

	// Pagination info
	if m.total >= 0 {
		pageInfo := fmt.Sprintf("\nPage %d/%d (%d rows)", m.page+1, m.TotalPages(), m.total)
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(pageInfo))
	}

	borderStyle := m.styles.Border.Width(m.width - 2)
	if m.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("4"))
	}

	return borderStyle.Height(m.height - 2).Render(b.String())
}

func (m Model) calcColumnWidths() []int {
	if len(m.columns) == 0 {
		return nil
	}

	widths := make([]int, len(m.columns))
	for i, col := range m.columns {
		widths[i] = len(col)
	}
	for _, row := range m.rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Cap widths and add padding
	maxColWidth := 30
	for i := range widths {
		if widths[i] > maxColWidth {
			widths[i] = maxColWidth
		}
		widths[i] += 2 // padding
	}

	// Determine visible columns based on available width
	available := m.width - 4 // borders
	visibleWidth := 0
	maxCols := 0
	for i := m.colOffset; i < len(widths); i++ {
		if visibleWidth+widths[i]+1 > available && maxCols > 0 {
			break
		}
		visibleWidth += widths[i] + 1 // +1 for separator
		maxCols++
	}

	// Adjust colOffset to keep cursor visible
	if m.cursorCol < m.colOffset {
		// This is handled but won't mutate; caller should manage
	}

	return widths
}

func (m Model) renderRow(cells []string, widths []int, rowIdx int, isHeader bool) string {
	var parts []string

	available := m.width - 4
	usedWidth := 0

	for i := 0; i < len(cells); i++ {
		if i >= len(widths) {
			break
		}

		w := widths[i]
		if usedWidth+w+1 > available && len(parts) > 0 {
			break
		}
		usedWidth += w + 1

		cell := cells[i]
		if len(cell) > w-2 {
			cell = cell[:w-3] + "…"
		}

		var style lipgloss.Style
		switch {
		case isHeader:
			style = m.styles.Header.Width(w)
		case rowIdx == m.cursorRow && i == m.cursorCol && m.focused:
			style = m.styles.Selected.Width(w)
		case cell == "NULL":
			style = m.styles.Null.Width(w)
		default:
			style = m.styles.Cell.Width(w)
		}

		parts = append(parts, style.Render(cell))
	}

	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("│")
	return strings.Join(parts, sep)
}
