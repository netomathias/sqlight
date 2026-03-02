package db

import (
	"fmt"
	"strings"
)

// QueryResult holds the result of executing a SQL query.
type QueryResult struct {
	Columns []string
	Rows    [][]string
	Total   int64 // -1 if unknown (e.g., raw SQL)
}

// QueryTable queries a table with pagination, validated against sqlite_master.
func (d *DB) QueryTable(table string, limit, offset int) (*QueryResult, error) {
	if err := d.validateTable(table); err != nil {
		return nil, err
	}

	// Get total count
	var total int64
	err := d.ro.QueryRow(
		fmt.Sprintf(`SELECT COUNT(*) FROM "%s"`, quoteIdent(table)),
	).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count %q: %w", table, err)
	}

	// Query with pagination
	query := fmt.Sprintf(`SELECT * FROM "%s" LIMIT ? OFFSET ?`, quoteIdent(table))
	rows, err := d.ro.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query %q: %w", table, err)
	}
	defer rows.Close()

	return scanResult(rows, total)
}

// ExecSQL executes a raw SQL statement.
// SELECT/PRAGMA/EXPLAIN returns rows; others return affected count.
func (d *DB) ExecSQL(sql string) (*QueryResult, error) {
	trimmed := strings.TrimSpace(sql)
	upper := strings.ToUpper(trimmed)

	isQuery := strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "PRAGMA") ||
		strings.HasPrefix(upper, "EXPLAIN") ||
		strings.HasPrefix(upper, "WITH")

	if isQuery {
		rows, err := d.ro.Query(trimmed)
		if err != nil {
			return nil, fmt.Errorf("exec query: %w", err)
		}
		defer rows.Close()
		return scanResult(rows, -1)
	}

	// Write operation
	if d.readOnly {
		return nil, fmt.Errorf("database is opened in read-only mode")
	}

	result, err := d.rw.Exec(trimmed)
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	affected, _ := result.RowsAffected()
	return &QueryResult{
		Columns: []string{"rows_affected"},
		Rows:    [][]string{{fmt.Sprintf("%d", affected)}},
		Total:   1,
	}, nil
}

func scanResult(rows interface{ Next() bool; Columns() ([]string, error); Scan(...any) error; Err() error }, total int64) (*QueryResult, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("columns: %w", err)
	}

	var data [][]string
	scanArgs := make([]any, len(cols))
	scanVals := make([]*string, len(cols))
	for i := range scanArgs {
		scanVals[i] = new(string)
		scanArgs[i] = &nullString{s: scanVals[i]}
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		row := make([]string, len(cols))
		for i, v := range scanVals {
			row[i] = *v
		}
		data = append(data, row)
	}

	return &QueryResult{Columns: cols, Rows: data, Total: total}, rows.Err()
}

// nullString implements sql.Scanner and converts NULL to "NULL".
type nullString struct {
	s *string
}

func (n *nullString) Scan(value any) error {
	if value == nil {
		*n.s = "NULL"
		return nil
	}
	*n.s = fmt.Sprintf("%v", value)
	return nil
}
