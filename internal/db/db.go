package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// DB wraps two connection pools: one for reads and one for writes.
type DB struct {
	rw       *sql.DB
	ro       *sql.DB
	readOnly bool
}

// Open opens a SQLite database with appropriate pragmas.
// If readOnly is true, the write pool is not created.
func Open(path string, readOnly bool) (*DB, error) {
	d := &DB{readOnly: readOnly}

	// Open rw first (if needed) so we can set journal_mode=WAL before ro connects.
	if !readOnly {
		rwDSN := fmt.Sprintf("file:%s?mode=rw", path)
		rwDB, err := sql.Open("sqlite", rwDSN)
		if err != nil {
			return nil, fmt.Errorf("open rw: %w", err)
		}
		rwDB.SetMaxOpenConns(1)
		if err := rwDB.Ping(); err != nil {
			rwDB.Close()
			return nil, fmt.Errorf("ping rw: %w", err)
		}
		if err := setWritePragmas(rwDB); err != nil {
			rwDB.Close()
			return nil, fmt.Errorf("pragmas rw: %w", err)
		}
		d.rw = rwDB
	}

	roDSN := fmt.Sprintf("file:%s?mode=ro", path)
	roDB, err := sql.Open("sqlite", roDSN)
	if err != nil {
		d.Close()
		return nil, fmt.Errorf("open ro: %w", err)
	}
	if err := roDB.Ping(); err != nil {
		roDB.Close()
		d.Close()
		return nil, fmt.Errorf("ping ro: %w", err)
	}
	if err := setReadPragmas(roDB); err != nil {
		roDB.Close()
		d.Close()
		return nil, fmt.Errorf("pragmas ro: %w", err)
	}
	d.ro = roDB

	return d, nil
}

// Close closes both connection pools.
func (d *DB) Close() error {
	var firstErr error
	if d.rw != nil {
		if err := d.rw.Close(); err != nil {
			firstErr = err
		}
	}
	if err := d.ro.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

// ReadOnly returns whether the database was opened in read-only mode.
func (d *DB) ReadOnly() bool {
	return d.readOnly
}

func setWritePragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
		"PRAGMA trusted_schema=OFF",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
	}
	return nil
}

func setReadPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA foreign_keys=ON",
		"PRAGMA busy_timeout=5000",
		"PRAGMA trusted_schema=OFF",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
	}
	return nil
}
