package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBPath   string
	ReadOnly bool
	Version  bool
}

var version = "dev"

func Parse() (Config, error) {
	var cfg Config

	flag.BoolVar(&cfg.ReadOnly, "readonly", false, "open database in read-only mode")
	flag.BoolVar(&cfg.ReadOnly, "r", false, "open database in read-only mode (shorthand)")
	flag.BoolVar(&cfg.Version, "version", false, "print version and exit")
	flag.Parse()

	if cfg.Version {
		fmt.Println("sqlight", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		return cfg, fmt.Errorf("usage: sqlight [flags] <database.db>")
	}

	dbPath := args[0]

	info, err := os.Stat(dbPath)
	if err != nil {
		return cfg, fmt.Errorf("cannot access %q: %w", dbPath, err)
	}
	if !info.Mode().IsRegular() {
		return cfg, fmt.Errorf("%q is not a regular file", dbPath)
	}

	ext := filepath.Ext(dbPath)
	switch ext {
	case ".db", ".sqlite", ".sqlite3":
		// valid
	default:
		return cfg, fmt.Errorf("unsupported file extension %q (expected .db, .sqlite, or .sqlite3)", ext)
	}

	cfg.DBPath = dbPath
	return cfg, nil
}
