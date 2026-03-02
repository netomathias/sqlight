# sqlight

A terminal UI for SQLite databases. Browse tables, run queries, and edit data — all from your terminal.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO required).

## Install

### From release binaries

Download the latest binary from the [releases page](https://github.com/netomathias/sqlight/releases).

### With Go

```bash
go install github.com/netomathias/sqlight@latest
```

### Build from source

```bash
git clone https://github.com/netomathias/sqlight.git
cd sqlight
go build -o sqlight .
```

## Usage

```bash
sqlight path/to/database.db
```

Open in read-only mode (prevents any writes):

```bash
sqlight --readonly path/to/database.db
sqlight -r path/to/database.db
```

Supported file extensions: `.db`, `.sqlite`, `.sqlite3`

## Keyboard Shortcuts

| Key              | Action               |
|------------------|----------------------|
| `tab`/`shift+tab`| Switch pane          |
| `↑`/`↓` or `j`/`k` | Navigate rows     |
| `←`/`→` or `h`/`l` | Navigate columns  |
| `enter`          | Select table         |
| `n`/`p`          | Next/prev page       |
| `ctrl+e`         | Toggle SQL editor    |
| `alt+enter`      | Execute SQL query    |
| `e`              | Edit cell            |
| `i`              | Insert row           |
| `d`              | Delete row           |
| `?`              | Toggle help          |
| `q`/`ctrl+c`     | Quit                 |

## Features

- **Table browser** — sidebar lists all tables and views, with paginated data grid
- **SQL editor** — write and execute arbitrary SQL queries
- **CRUD operations** — edit cells, insert rows, delete rows with confirmation
- **Read-only mode** — enforced at the SQLite level via `?mode=ro`
- **Pure Go** — single binary, no CGO, cross-platform

## License

MIT
