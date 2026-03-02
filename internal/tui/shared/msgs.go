package shared

import "github.com/netomathias/sqlight/internal/db"

// Messages shared between TUI components.

// TableSelectedMsg is sent when a table is selected in the sidebar.
type TableSelectedMsg struct {
	Name string
	Type string // "table" or "view"
}

// TableDataMsg carries query results to display.
type TableDataMsg struct {
	Table   string
	Result  *db.QueryResult
	Page    int
	PerPage int
}

// SQLResultMsg carries the result of a raw SQL execution.
type SQLResultMsg struct {
	Result *db.QueryResult
	SQL    string
}

// ErrorMsg carries an error message.
type ErrorMsg struct {
	Err error
}

// StatusMsg updates the status bar text.
type StatusMsg struct {
	Text string
}

// RefreshTableMsg triggers a refresh of the current table data.
type RefreshTableMsg struct{}

// EditorToggleMsg toggles the SQL editor visibility.
type EditorToggleMsg struct{}

// CRUDMsg carries a CRUD operation result.
type CRUDMsg struct {
	Operation string // "insert", "update", "delete"
	Err       error
}

// EditCellMsg requests editing a cell.
type EditCellMsg struct {
	Table  string
	PKCol  string
	PKVal  string
	Col    string
	OldVal string
}

// InsertRowMsg requests inserting a new row.
type InsertRowMsg struct {
	Table string
}

// DeleteRowMsg requests deleting a row.
type DeleteRowMsg struct {
	Table string
	PKCol string
	PKVal string
}

// ConfirmDeleteMsg confirms a delete operation.
type ConfirmDeleteMsg struct {
	Table string
	PKCol string
	PKVal string
}
