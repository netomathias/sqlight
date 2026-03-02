package db

import (
	"fmt"
)

// TableInfo describes a table or view in the database.
type TableInfo struct {
	Name string
	Type string // "table" or "view"
}

// Column describes a column in a table.
type Column struct {
	CID          int
	Name         string
	Type         string
	NotNull      bool
	DefaultValue *string
	PK           int
}

// Tables returns all tables and views in the database.
func (d *DB) Tables() ([]TableInfo, error) {
	rows, err := d.ro.Query(
		`SELECT name, type FROM sqlite_master
		 WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%'
		 ORDER BY type, name`)
	if err != nil {
		return nil, fmt.Errorf("list tables: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.Name, &t.Type); err != nil {
			return nil, fmt.Errorf("scan table: %w", err)
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

// Columns returns column info for a table, validated against sqlite_master.
func (d *DB) Columns(table string) ([]Column, error) {
	if err := d.validateTable(table); err != nil {
		return nil, err
	}

	rows, err := d.ro.Query(fmt.Sprintf(`PRAGMA table_info("%s")`, quoteIdent(table)))
	if err != nil {
		return nil, fmt.Errorf("columns %q: %w", table, err)
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		var c Column
		if err := rows.Scan(&c.CID, &c.Name, &c.Type, &c.NotNull, &c.DefaultValue, &c.PK); err != nil {
			return nil, fmt.Errorf("scan column: %w", err)
		}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}

// validateTable checks that a table name exists in sqlite_master.
func (d *DB) validateTable(name string) error {
	var count int
	err := d.ro.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE name = ? AND type IN ('table', 'view')`,
		name,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("validate table: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("table %q does not exist", name)
	}
	return nil
}

// validateColumn checks that a column exists in a given table.
func (d *DB) validateColumn(table, column string) error {
	cols, err := d.Columns(table)
	if err != nil {
		return err
	}
	for _, c := range cols {
		if c.Name == column {
			return nil
		}
	}
	return fmt.Errorf("column %q does not exist in table %q", column, table)
}

// quoteIdent double-quotes an identifier, escaping internal double quotes.
func quoteIdent(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			out = append(out, '"')
		}
		out = append(out, s[i])
	}
	return string(out)
}
