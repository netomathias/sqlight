package db

import (
	"fmt"
	"strings"
)

// Insert inserts a row into the given table.
func (d *DB) Insert(table string, data map[string]string) error {
	if d.readOnly {
		return fmt.Errorf("database is opened in read-only mode")
	}
	if err := d.validateTable(table); err != nil {
		return err
	}
	for col := range data {
		if err := d.validateColumn(table, col); err != nil {
			return err
		}
	}

	cols := make([]string, 0, len(data))
	placeholders := make([]string, 0, len(data))
	vals := make([]any, 0, len(data))

	for col, val := range data {
		cols = append(cols, fmt.Sprintf(`"%s"`, quoteIdent(col)))
		placeholders = append(placeholders, "?")
		vals = append(vals, val)
	}

	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`,
		quoteIdent(table),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := d.rw.Exec(query, vals...)
	if err != nil {
		return fmt.Errorf("insert into %q: %w", table, err)
	}
	return nil
}

// Update updates a single row identified by pkCol=pkVal.
func (d *DB) Update(table, pkCol, pkVal, col, newVal string) error {
	if d.readOnly {
		return fmt.Errorf("database is opened in read-only mode")
	}
	if err := d.validateTable(table); err != nil {
		return err
	}
	if err := d.validateColumn(table, pkCol); err != nil {
		return err
	}
	if err := d.validateColumn(table, col); err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE "%s" SET "%s" = ? WHERE "%s" = ?`,
		quoteIdent(table), quoteIdent(col), quoteIdent(pkCol))

	result, err := d.rw.Exec(query, newVal, pkVal)
	if err != nil {
		return fmt.Errorf("update %q: %w", table, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no row found with %s = %s", pkCol, pkVal)
	}
	return nil
}

// Delete deletes a single row identified by pkCol=pkVal.
func (d *DB) Delete(table, pkCol, pkVal string) error {
	if d.readOnly {
		return fmt.Errorf("database is opened in read-only mode")
	}
	if err := d.validateTable(table); err != nil {
		return err
	}
	if err := d.validateColumn(table, pkCol); err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM "%s" WHERE "%s" = ?`,
		quoteIdent(table), quoteIdent(pkCol))

	result, err := d.rw.Exec(query, pkVal)
	if err != nil {
		return fmt.Errorf("delete from %q: %w", table, err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("no row found with %s = %s", pkCol, pkVal)
	}
	return nil
}
