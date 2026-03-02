package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/netomathias/sqlight/internal/db"
	"github.com/netomathias/sqlight/internal/tui/dialog"
	"github.com/netomathias/sqlight/internal/tui/editor"
	"github.com/netomathias/sqlight/internal/tui/shared"
	"github.com/netomathias/sqlight/internal/tui/sidebar"
	"github.com/netomathias/sqlight/internal/tui/statusbar"
	"github.com/netomathias/sqlight/internal/tui/table"
)

// Pane identifies which pane is focused.
type Pane int

const (
	PaneSidebar Pane = iota
	PaneTable
	PaneEditor
)

// App is the root TUI model.
type App struct {
	db       *db.DB
	dbPath   string
	readOnly bool

	sidebar   sidebar.Model
	table     table.Model
	editor    editor.Model
	statusbar statusbar.Model
	dialog    dialog.Model

	pane     Pane
	keys     shared.KeyMap
	styles   shared.Styles
	width    int
	height   int
	showHelp bool

	// State for pending edit
	editMsg shared.EditCellMsg
}

const sidebarWidth = 28

// NewApp creates a new root TUI model.
func NewApp(database *db.DB, dbPath string, readOnly bool) App {
	styles := shared.DefaultStyles()
	keys := shared.DefaultKeyMap()

	return App{
		db:        database,
		dbPath:    dbPath,
		readOnly:  readOnly,
		sidebar:   sidebar.New(styles.Sidebar, keys),
		table:     table.New(styles.Table, keys),
		editor:    editor.New(styles.Editor, keys),
		statusbar: statusbar.New(dbPath, readOnly, styles.StatusBar),
		dialog: dialog.New(dialog.DialogStyles{
			Border:  styles.Dialog.Border,
			Title:   styles.Dialog.Title,
			Input:   styles.Dialog.Input,
			Button:  styles.Dialog.Button,
			Overlay: styles.Dialog.Overlay,
		}),
		pane:   PaneSidebar,
		keys:   keys,
		styles: styles,
	}
}

// Init implements tea.Model.
func (a App) Init() tea.Cmd {
	return a.loadTables()
}

// Update implements tea.Model.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Dialog takes priority when active.
	if a.dialog.Active() {
		return a.updateDialog(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.layout()
		return a, nil

	case tea.KeyMsg:
		// Global keys
		switch {
		case key.Matches(msg, a.keys.Quit) && !a.editor.Focused():
			return a, tea.Quit
		case key.Matches(msg, a.keys.Help) && !a.editor.Focused():
			a.showHelp = !a.showHelp
			return a, nil
		case key.Matches(msg, a.keys.Tab) && !a.editor.Focused():
			a.nextPane()
			return a, nil
		case key.Matches(msg, a.keys.ShiftTab) && !a.editor.Focused():
			a.prevPane()
			return a, nil
		case key.Matches(msg, a.keys.ToggleEditor):
			a.toggleEditor()
			return a, nil
		case key.Matches(msg, a.keys.Escape):
			if a.editor.Focused() {
				a.pane = PaneTable
				a.updateFocus()
				return a, nil
			}
			if a.showHelp {
				a.showHelp = false
				return a, nil
			}
		}

	case tablesLoadedMsg:
		a.sidebar.SetTables(msg.tables)
		if len(msg.tables) > 0 {
			return a, func() tea.Msg {
				return shared.TableSelectedMsg{Name: msg.tables[0].Name, Type: msg.tables[0].Type}
			}
		}
		return a, nil

	case shared.TableSelectedMsg:
		a.statusbar.SetTable(msg.Name)
		a.statusbar.ClearMessage()
		return a, a.loadTableData(msg.Name, 0)

	case shared.TableDataMsg:
		a.table.SetData(msg.Table, msg.Result, msg.Page)
		return a, nil

	case shared.SQLResultMsg:
		a.table.SetData("query", msg.Result, 0)
		a.statusbar.SetTable("query result")
		a.statusbar.SetMessage(fmt.Sprintf("Query returned %d rows", len(msg.Result.Rows)), false)
		return a, nil

	case shared.RefreshTableMsg:
		name := a.table.TableName()
		if name != "" && name != "query" {
			return a, a.loadTableData(name, a.table.Page())
		}
		return a, nil

	case shared.ErrorMsg:
		a.statusbar.SetMessage(msg.Err.Error(), true)
		return a, nil

	case shared.EditCellMsg:
		if a.readOnly {
			a.statusbar.SetMessage("Read-only mode", true)
			return a, nil
		}
		if a.table.TableName() == "" {
			a.statusbar.SetMessage("No table selected", true)
			return a, nil
		}
		if a.table.TableName() == "query" {
			a.statusbar.SetMessage("Cannot edit query results — select a table first", true)
			return a, nil
		}
		a.editMsg = msg
		a.dialog.ShowEdit(msg.Col, msg.OldVal)
		return a, nil

	case shared.InsertRowMsg:
		if a.readOnly {
			a.statusbar.SetMessage("Read-only mode", true)
			return a, nil
		}
		if a.table.TableName() == "" || a.table.TableName() == "query" {
			a.statusbar.SetMessage("Cannot insert into query results — select a table first", true)
			return a, nil
		}
		cols := a.table.Columns()
		if len(cols) > 0 {
			a.dialog.ShowInsert(cols)
		}
		return a, nil

	case shared.DeleteRowMsg:
		if a.readOnly {
			a.statusbar.SetMessage("Read-only mode", true)
			return a, nil
		}
		if a.table.TableName() == "" || a.table.TableName() == "query" {
			a.statusbar.SetMessage("Cannot delete from query results — select a table first", true)
			return a, nil
		}
		a.dialog.ShowConfirmDelete(msg.Table, msg.PKCol, msg.PKVal)
		return a, nil

	case shared.CRUDMsg:
		if msg.Err != nil {
			a.statusbar.SetMessage(msg.Err.Error(), true)
		} else {
			a.statusbar.SetMessage(fmt.Sprintf("%s successful", msg.Operation), false)
		}
		if msg.Err == nil {
			return a, func() tea.Msg { return shared.RefreshTableMsg{} }
		}
		return a, nil
	}

	// Check for SQL execution from editor
	if sql, ok := editor.IsExecSQL(msg); ok {
		return a, a.execSQL(sql)
	}

	// Delegate to focused component
	var cmd tea.Cmd
	switch a.pane {
	case PaneSidebar:
		a.sidebar, cmd = a.sidebar.Update(msg)
	case PaneTable:
		a.table, cmd = a.table.Update(msg)
	case PaneEditor:
		a.editor, cmd = a.editor.Update(msg)
	}

	return a, cmd
}

func (a *App) updateDialog(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case dialog.EditSubmitMsg:
		a.dialog.Close()
		return a, a.updateCell(a.editMsg.Table, a.editMsg.PKCol, a.editMsg.PKVal, a.editMsg.Col, msg.Value)

	case dialog.InsertSubmitMsg:
		tableName := a.table.TableName()
		a.dialog.Close()
		return a, a.insertRow(tableName, msg.Data)

	case dialog.ConfirmDeleteMsg:
		a.dialog.Close()
		return a, a.deleteRow(msg.Table, msg.PKCol, msg.PKVal)
	}

	var cmd tea.Cmd
	a.dialog, cmd = a.dialog.Update(msg)
	return a, cmd
}

// View implements tea.Model.
func (a App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}

	if a.showHelp {
		return a.helpView()
	}

	sidebarView := a.sidebar.View()
	rightWidth := a.width - sidebarWidth

	var rightParts []string

	if a.editor.Visible() {
		editorHeight := 8
		tableHeight := a.height - editorHeight - 1
		if tableHeight < 5 {
			tableHeight = 5
		}

		tableModel := a.table
		tableModel.SetSize(rightWidth, tableHeight)
		rightParts = append(rightParts, tableModel.View())

		editorModel := a.editor
		editorModel.SetSize(rightWidth, editorHeight)
		rightParts = append(rightParts, editorModel.View())
	} else {
		tableHeight := a.height - 1
		tableModel := a.table
		tableModel.SetSize(rightWidth, tableHeight)
		rightParts = append(rightParts, tableModel.View())
	}

	rightView := lipgloss.JoinVertical(lipgloss.Left, rightParts...)
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, rightView)
	statusView := a.statusbar.View()
	fullView := lipgloss.JoinVertical(lipgloss.Left, mainView, statusView)

	if a.dialog.Active() {
		dialogView := a.dialog.View()
		fullView = a.overlayDialog(fullView, dialogView)
	}

	return fullView
}

func (a App) overlayDialog(bg, dlg string) string {
	bgLines := strings.Split(bg, "\n")
	dialogLines := strings.Split(dlg, "\n")

	startY := (len(bgLines) - len(dialogLines)) / 2
	if startY < 0 {
		startY = 0
	}
	startX := (a.width - lipgloss.Width(dlg)) / 2
	if startX < 0 {
		startX = 0
	}

	for i, dLine := range dialogLines {
		y := startY + i
		if y >= len(bgLines) {
			break
		}
		if startX < len(bgLines[y]) {
			bgLines[y] = fmt.Sprintf("%*s%s", startX, "", dLine)
		}
	}

	return strings.Join(bgLines, "\n")
}

func (a App) helpView() string {
	help := []struct{ key, desc string }{
		{"tab/shift+tab", "Switch pane"},
		{"↑/↓ or j/k", "Navigate rows"},
		{"←/→ or h/l", "Navigate columns"},
		{"n/p", "Next/prev page"},
		{"enter", "Select table"},
		{"ctrl+e", "Toggle SQL editor"},
		{"alt+enter/F5", "Execute SQL"},
		{"e", "Edit cell"},
		{"i", "Insert row"},
		{"d", "Delete row"},
		{"?", "Toggle help"},
		{"q/ctrl+c", "Quit"},
	}

	var b strings.Builder
	b.WriteString(a.styles.Help.Key.Render("  sqlight - Keyboard Shortcuts"))
	b.WriteString("\n\n")

	for _, h := range help {
		b.WriteString("  ")
		b.WriteString(a.styles.Help.Key.Width(20).Render(h.key))
		b.WriteString(a.styles.Help.Desc.Render(h.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(a.styles.Help.Desc.Render("  Press ? or esc to close"))

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Padding(2, 4).
		Render(b.String())
}

// Layout & focus management

func (a *App) layout() {
	tableHeight := a.height - 1
	if a.editor.Visible() {
		tableHeight = a.height - 9
	}

	a.sidebar.SetSize(sidebarWidth, a.height-1)
	a.table.SetSize(a.width-sidebarWidth, tableHeight)
	a.editor.SetSize(a.width-sidebarWidth, 8)
	a.statusbar.SetWidth(a.width)
	a.dialog.SetSize(a.width, a.height)
}

func (a *App) nextPane() {
	if a.editor.Visible() {
		switch a.pane {
		case PaneSidebar:
			a.pane = PaneTable
		case PaneTable:
			a.pane = PaneEditor
		case PaneEditor:
			a.pane = PaneSidebar
		}
	} else {
		switch a.pane {
		case PaneSidebar:
			a.pane = PaneTable
		case PaneTable:
			a.pane = PaneSidebar
		}
	}
	a.updateFocus()
}

func (a *App) prevPane() {
	if a.editor.Visible() {
		switch a.pane {
		case PaneSidebar:
			a.pane = PaneEditor
		case PaneTable:
			a.pane = PaneSidebar
		case PaneEditor:
			a.pane = PaneTable
		}
	} else {
		switch a.pane {
		case PaneSidebar:
			a.pane = PaneTable
		case PaneTable:
			a.pane = PaneSidebar
		}
	}
	a.updateFocus()
}

func (a *App) toggleEditor() {
	if a.editor.Visible() {
		a.editor.SetVisible(false)
		if a.pane == PaneEditor {
			a.pane = PaneTable
		}
	} else {
		a.editor.SetVisible(true)
		a.pane = PaneEditor
	}
	a.updateFocus()
	a.layout()
}

func (a *App) updateFocus() {
	a.sidebar.SetFocused(a.pane == PaneSidebar)
	a.table.SetFocused(a.pane == PaneTable)
	a.editor.SetFocused(a.pane == PaneEditor)
}

// Commands (side effects)

type tablesLoadedMsg struct {
	tables []db.TableInfo
}

func (a *App) loadTables() tea.Cmd {
	return func() tea.Msg {
		tables, err := a.db.Tables()
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}
		return tablesLoadedMsg{tables: tables}
	}
}

func (a *App) loadTableData(name string, page int) tea.Cmd {
	perPage := a.table.PerPage()
	return func() tea.Msg {
		result, err := a.db.QueryTable(name, perPage, page*perPage)
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}
		return shared.TableDataMsg{Table: name, Result: result, Page: page, PerPage: perPage}
	}
}

func (a *App) execSQL(sql string) tea.Cmd {
	return func() tea.Msg {
		result, err := a.db.ExecSQL(sql)
		if err != nil {
			return shared.ErrorMsg{Err: err}
		}
		return shared.SQLResultMsg{Result: result, SQL: sql}
	}
}

func (a *App) updateCell(tableName, pkCol, pkVal, col, newVal string) tea.Cmd {
	return func() tea.Msg {
		err := a.db.Update(tableName, pkCol, pkVal, col, newVal)
		return shared.CRUDMsg{Operation: "update", Err: err}
	}
}

func (a *App) insertRow(tableName string, data map[string]string) tea.Cmd {
	return func() tea.Msg {
		err := a.db.Insert(tableName, data)
		return shared.CRUDMsg{Operation: "insert", Err: err}
	}
}

func (a *App) deleteRow(tableName, pkCol, pkVal string) tea.Cmd {
	return func() tea.Msg {
		err := a.db.Delete(tableName, pkCol, pkVal)
		return shared.CRUDMsg{Operation: "delete", Err: err}
	}
}
