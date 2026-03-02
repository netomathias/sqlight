package db

import (
	"testing"
)

func TestQueryTable(t *testing.T) {
	d := testDB(t)

	result, err := d.QueryTable("users", 10, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(result.Columns))
	}
	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}
	if len(result.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(result.Rows))
	}

	if result.Rows[2][2] != "NULL" {
		t.Errorf("expected NULL for Charlie's email, got %q", result.Rows[2][2])
	}
}

func TestQueryTablePagination(t *testing.T) {
	d := testDB(t)

	result, err := d.QueryTable("users", 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Rows))
	}
	if result.Rows[0][1] != "Bob" {
		t.Errorf("expected Bob, got %q", result.Rows[0][1])
	}
}

func TestQueryTableInvalidTable(t *testing.T) {
	d := testDB(t)

	_, err := d.QueryTable("nonexistent", 10, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent table")
	}
}

func TestQueryView(t *testing.T) {
	d := testDB(t)

	result, err := d.QueryTable("user_posts", 10, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Rows) != 2 {
		t.Fatalf("expected 2 rows from view, got %d", len(result.Rows))
	}
}

func TestExecSQLSelect(t *testing.T) {
	d := testDB(t)

	result, err := d.ExecSQL("SELECT name FROM users WHERE id = 1")
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Columns) != 1 || result.Columns[0] != "name" {
		t.Fatalf("unexpected columns: %v", result.Columns)
	}
	if len(result.Rows) != 1 || result.Rows[0][0] != "Alice" {
		t.Fatalf("unexpected rows: %v", result.Rows)
	}
}

func TestExecSQLInsert(t *testing.T) {
	d := testDB(t)

	result, err := d.ExecSQL("INSERT INTO users (name, email) VALUES ('Dave', 'dave@test.com')")
	if err != nil {
		t.Fatal(err)
	}

	if result.Rows[0][0] != "1" {
		t.Fatalf("expected 1 row affected, got %s", result.Rows[0][0])
	}
}

func TestExecSQLReadOnlyBlock(t *testing.T) {
	d := testReadOnlyDB(t)

	_, err := d.ExecSQL("INSERT INTO items (val) VALUES ('x')")
	if err == nil {
		t.Fatal("expected error for write in read-only mode")
	}
}
