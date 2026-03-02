//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:testdata/sample.db?mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		PRAGMA journal_mode=WAL;
		PRAGMA foreign_keys=ON;

		CREATE TABLE IF NOT EXISTS artists (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			country TEXT
		);

		CREATE TABLE IF NOT EXISTS albums (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			artist_id INTEGER NOT NULL,
			year INTEGER,
			FOREIGN KEY (artist_id) REFERENCES artists(id)
		);

		CREATE TABLE IF NOT EXISTS tracks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			album_id INTEGER NOT NULL,
			duration_seconds INTEGER,
			track_number INTEGER,
			FOREIGN KEY (album_id) REFERENCES albums(id)
		);

		CREATE VIEW IF NOT EXISTS album_details AS
			SELECT a.title AS album, ar.name AS artist, a.year,
				   COUNT(t.id) AS tracks, SUM(t.duration_seconds) AS total_seconds
			FROM albums a
			JOIN artists ar ON a.artist_id = ar.id
			LEFT JOIN tracks t ON t.album_id = a.id
			GROUP BY a.id;

		INSERT OR IGNORE INTO artists (id, name, country) VALUES
			(1, 'Pink Floyd', 'UK'),
			(2, 'Led Zeppelin', 'UK'),
			(3, 'Miles Davis', 'US'),
			(4, 'Radiohead', 'UK'),
			(5, 'Kraftwerk', 'DE');

		INSERT OR IGNORE INTO albums (id, title, artist_id, year) VALUES
			(1, 'The Dark Side of the Moon', 1, 1973),
			(2, 'Wish You Were Here', 1, 1975),
			(3, 'Led Zeppelin IV', 2, 1971),
			(4, 'Kind of Blue', 3, 1959),
			(5, 'OK Computer', 4, 1997),
			(6, 'Autobahn', 5, 1974);

		INSERT OR IGNORE INTO tracks (id, title, album_id, duration_seconds, track_number) VALUES
			(1, 'Speak to Me', 1, 68, 1),
			(2, 'Breathe', 1, 169, 2),
			(3, 'On the Run', 1, 225, 3),
			(4, 'Time', 1, 413, 4),
			(5, 'The Great Gig in the Sky', 1, 284, 5),
			(6, 'Money', 1, 382, 6),
			(7, 'Us and Them', 1, 469, 7),
			(8, 'Any Colour You Like', 1, 206, 8),
			(9, 'Brain Damage', 1, 228, 9),
			(10, 'Eclipse', 1, 131, 10),
			(11, 'Shine On You Crazy Diamond (Parts I-V)', 2, 810, 1),
			(12, 'Welcome to the Machine', 2, 450, 2),
			(13, 'Have a Cigar', 2, 307, 3),
			(14, 'Wish You Were Here', 2, 334, 4),
			(15, 'Shine On You Crazy Diamond (Parts VI-IX)', 2, 740, 5),
			(16, 'Black Dog', 3, 296, 1),
			(17, 'Rock and Roll', 3, 220, 2),
			(18, 'Stairway to Heaven', 3, 482, 4),
			(19, 'So What', 4, 562, 1),
			(20, 'Freddie Freeloader', 4, 588, 2),
			(21, 'Blue in Green', 4, 327, 3),
			(22, 'Paranoid Android', 5, 386, 2),
			(23, 'Lucky', 5, 270, 8),
			(24, 'Autobahn', 6, 1348, 1);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created testdata/sample.db")
}
