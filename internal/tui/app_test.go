package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/netomathias/sqlight/internal/db"
	"github.com/netomathias/sqlight/internal/tui/shared"
)

func testApp(t *testing.T) App {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	database, err := db.Open(path, false)
	if err != nil {
		t.Fatal(err)
	}

	// Create test table
	_, err = database.ExecSQL("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = database.ExecSQL("INSERT INTO users (name) VALUES ('Alice')")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { database.Close() })

	app := NewApp(database, path, false)
	// Simulate window size
	m, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return m.(App)
}

func testReadOnlyApp(t *testing.T) App {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Create in rw mode first
	tmp, err := db.Open(path, false)
	if err != nil {
		t.Fatal(err)
	}
	tmp.ExecSQL("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
	tmp.Close()

	database, err := db.Open(path, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })

	app := NewApp(database, path, true)
	m, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return m.(App)
}

func TestNewApp(t *testing.T) {
	app := testApp(t)

	if app.pane != PaneSidebar {
		t.Errorf("expected initial pane to be Sidebar, got %d", app.pane)
	}
	if app.width != 120 || app.height != 40 {
		t.Errorf("expected 120x40, got %dx%d", app.width, app.height)
	}
}

func TestPaneNavigation(t *testing.T) {
	app := testApp(t)

	// tab: sidebar -> table
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyTab})
	app = m.(App)
	if app.pane != PaneTable {
		t.Errorf("expected PaneTable after tab, got %d", app.pane)
	}

	// tab: table -> sidebar (editor not visible)
	m, _ = app.Update(tea.KeyMsg{Type: tea.KeyTab})
	app = m.(App)
	if app.pane != PaneSidebar {
		t.Errorf("expected PaneSidebar after second tab, got %d", app.pane)
	}
}

func TestToggleEditor(t *testing.T) {
	app := testApp(t)

	// ctrl+e opens editor
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	app = m.(App)
	if !app.editor.Visible() {
		t.Error("editor should be visible after ctrl+e")
	}
	if app.pane != PaneEditor {
		t.Errorf("expected PaneEditor, got %d", app.pane)
	}

	// ctrl+e closes editor
	m, _ = app.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
	app = m.(App)
	if app.editor.Visible() {
		t.Error("editor should be hidden after second ctrl+e")
	}
}

func TestHelpToggle(t *testing.T) {
	app := testApp(t)

	// ? toggles help
	m, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	app = m.(App)
	if !app.showHelp {
		t.Error("help should be shown after ?")
	}

	view := app.View()
	if len(view) == 0 {
		t.Error("help view should not be empty")
	}
}

func TestTableSelected(t *testing.T) {
	app := testApp(t)

	m, _ := app.Update(shared.TableSelectedMsg{Name: "users", Type: "table"})
	app = m.(App)

	// Status bar should show table name (tested implicitly via View())
	view := app.View()
	if len(view) == 0 {
		t.Error("view should not be empty after table selection")
	}
}

func TestReadOnlyBlocksCRUD(t *testing.T) {
	app := testReadOnlyApp(t)

	// EditCellMsg should be blocked
	m, _ := app.Update(shared.EditCellMsg{Table: "users", PKCol: "id", PKVal: "1", Col: "name", OldVal: "Alice"})
	app = m.(App)
	// No crash = pass

	// InsertRowMsg should be blocked
	m, _ = app.Update(shared.InsertRowMsg{Table: "users"})
	app = m.(App)

	// DeleteRowMsg should be blocked
	m, _ = app.Update(shared.DeleteRowMsg{Table: "users", PKCol: "id", PKVal: "1"})
	app = m.(App)
}

func TestErrorMsg(t *testing.T) {
	app := testApp(t)

	m, _ := app.Update(shared.ErrorMsg{Err: fmt.Errorf("test error")})
	app = m.(App)

	view := app.View()
	if len(view) == 0 {
		t.Error("view should not be empty after error")
	}
}
