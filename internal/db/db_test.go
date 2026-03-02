package db

import (
	"os"
	"path/filepath"
	"testing"
)

// testDB creates a temporary SQLite file with test data and returns an opened DB.
func testDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	tmp, err := Open(path, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmp.rw.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT
		);
		INSERT INTO users (name, email) VALUES ('Alice', 'alice@test.com');
		INSERT INTO users (name, email) VALUES ('Bob', 'bob@test.com');
		INSERT INTO users (name, email) VALUES ('Charlie', NULL);

		CREATE TABLE posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
		INSERT INTO posts (user_id, title) VALUES (1, 'Hello World');
		INSERT INTO posts (user_id, title) VALUES (2, 'Second Post');

		CREATE VIEW user_posts AS
			SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id;
	`)
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	d, err := Open(path, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func testReadOnlyDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	tmp, err := Open(path, false)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmp.rw.Exec(`CREATE TABLE items (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	d, err := Open(path, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

func TestOpenClose(t *testing.T) {
	d := testDB(t)
	if d.ro == nil {
		t.Fatal("ro pool should not be nil")
	}
	if d.rw == nil {
		t.Fatal("rw pool should not be nil")
	}
	if d.ReadOnly() {
		t.Fatal("should not be read-only")
	}
}

func TestOpenReadOnly(t *testing.T) {
	d := testReadOnlyDB(t)
	if d.rw != nil {
		t.Fatal("rw pool should be nil in read-only mode")
	}
	if !d.ReadOnly() {
		t.Fatal("should be read-only")
	}
}
