package db

import (
	"testing"
)

func TestTables(t *testing.T) {
	d := testDB(t)

	tables, err := d.Tables()
	if err != nil {
		t.Fatal(err)
	}

	// Should have 2 tables + 1 view
	if len(tables) != 3 {
		t.Fatalf("expected 3 entries, got %d: %v", len(tables), tables)
	}

	found := map[string]string{}
	for _, ti := range tables {
		found[ti.Name] = ti.Type
	}

	if found["users"] != "table" {
		t.Error("missing users table")
	}
	if found["posts"] != "table" {
		t.Error("missing posts table")
	}
	if found["user_posts"] != "view" {
		t.Error("missing user_posts view")
	}
}

func TestColumns(t *testing.T) {
	d := testDB(t)

	cols, err := d.Columns("users")
	if err != nil {
		t.Fatal(err)
	}

	if len(cols) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(cols))
	}

	// Check first column is id and is PK
	if cols[0].Name != "id" || cols[0].PK != 1 {
		t.Errorf("expected id as PK, got %+v", cols[0])
	}
	if cols[1].Name != "name" || !cols[1].NotNull {
		t.Errorf("expected name NOT NULL, got %+v", cols[1])
	}
}

func TestColumnsInvalidTable(t *testing.T) {
	d := testDB(t)

	_, err := d.Columns("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent table")
	}
}

func TestValidateTable(t *testing.T) {
	d := testDB(t)

	if err := d.validateTable("users"); err != nil {
		t.Errorf("users should be valid: %v", err)
	}
	if err := d.validateTable("user_posts"); err != nil {
		t.Errorf("user_posts view should be valid: %v", err)
	}
	if err := d.validateTable("nope"); err == nil {
		t.Error("nonexistent table should fail validation")
	}
}

func TestQuoteIdent(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"users", "users"},
		{`my"table`, `my""table`},
		{"normal_name", "normal_name"},
	}
	for _, tt := range tests {
		got := quoteIdent(tt.in)
		if got != tt.want {
			t.Errorf("quoteIdent(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
