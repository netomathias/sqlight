package db

import (
	"testing"
)

func TestInsert(t *testing.T) {
	d := testDB(t)

	err := d.Insert("users", map[string]string{
		"name":  "Dave",
		"email": "dave@test.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	result, err := d.QueryTable("users", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 4 {
		t.Fatalf("expected 4 rows after insert, got %d", result.Total)
	}
}

func TestInsertInvalidTable(t *testing.T) {
	d := testDB(t)
	err := d.Insert("nope", map[string]string{"x": "y"})
	if err == nil {
		t.Fatal("expected error for invalid table")
	}
}

func TestInsertInvalidColumn(t *testing.T) {
	d := testDB(t)
	err := d.Insert("users", map[string]string{"badcol": "y"})
	if err == nil {
		t.Fatal("expected error for invalid column")
	}
}

func TestInsertReadOnly(t *testing.T) {
	d := testReadOnlyDB(t)
	err := d.Insert("items", map[string]string{"val": "x"})
	if err == nil {
		t.Fatal("expected error in read-only mode")
	}
}

func TestUpdate(t *testing.T) {
	d := testDB(t)

	err := d.Update("users", "id", "1", "name", "Alice Updated")
	if err != nil {
		t.Fatal(err)
	}

	result, err := d.ExecSQL("SELECT name FROM users WHERE id = 1")
	if err != nil {
		t.Fatal(err)
	}
	if result.Rows[0][0] != "Alice Updated" {
		t.Fatalf("expected 'Alice Updated', got %q", result.Rows[0][0])
	}
}

func TestUpdateNonexistentRow(t *testing.T) {
	d := testDB(t)
	err := d.Update("users", "id", "999", "name", "Ghost")
	if err == nil {
		t.Fatal("expected error for nonexistent row")
	}
}

func TestUpdateReadOnly(t *testing.T) {
	d := testReadOnlyDB(t)
	err := d.Update("items", "id", "1", "val", "x")
	if err == nil {
		t.Fatal("expected error in read-only mode")
	}
}

func TestDelete(t *testing.T) {
	d := testDB(t)

	err := d.Delete("users", "id", "3")
	if err != nil {
		t.Fatal(err)
	}

	result, err := d.QueryTable("users", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 rows after delete, got %d", result.Total)
	}
}

func TestDeleteNonexistentRow(t *testing.T) {
	d := testDB(t)
	err := d.Delete("users", "id", "999")
	if err == nil {
		t.Fatal("expected error for nonexistent row")
	}
}

func TestDeleteReadOnly(t *testing.T) {
	d := testReadOnlyDB(t)
	err := d.Delete("items", "id", "1")
	if err == nil {
		t.Fatal("expected error in read-only mode")
	}
}
