package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/Dima-salang/pomolite/internal/storage/db"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db      *sql.DB
	queries *db.Queries
}

func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	var dbPath string
	if path == ":memory:" {
		dbPath = ":memory:"
	} else if path != "" {
		dbPath = path
	} else {
		userConfDir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		appPath := filepath.Join(userConfDir, "pomolite")
		if err := os.MkdirAll(appPath, 0755); err != nil {
			return nil, err
		}
		dbPath = filepath.Join(appPath, "pomolite.db")
	}

	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := initTable(dbConn); err != nil {
		dbConn.Close()
		return nil, err
	}
	return &SQLiteStorage{
		db:      dbConn,
		queries: db.New(dbConn),
	}, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func initTable(db *sql.DB) error {
	// Create sessions table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			label TEXT NOT NULL,
			start_time INTEGER NOT NULL,
			end_time INTEGER NOT NULL,
			session_duration INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			priority INTEGER NOT NULL,
			status TEXT NOT NULL,
			due_date INTEGER,
			created_at INTEGER NOT NULL,
			completed_at INTEGER
		)
	`)
	if err != nil {
		return err
	}

	// Migration for old session database structures
	rows, err := db.Query("PRAGMA table_info(sessions)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasSessionDuration := false
	hasTimeDuration := false
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltVal interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltVal, &pk); err != nil {
			return err
		}
		if name == "session_duration" {
			hasSessionDuration = true
		}
		if name == "time_duration" {
			hasTimeDuration = true
		}
	}

	if !hasSessionDuration {
		if hasTimeDuration {
			_, err = db.Exec("ALTER TABLE sessions RENAME COLUMN time_duration TO session_duration")
			if err != nil {
				_, _ = db.Exec("ALTER TABLE sessions ADD COLUMN session_duration INTEGER DEFAULT 0")
			}
		} else {
			_, err = db.Exec("ALTER TABLE sessions ADD COLUMN session_duration INTEGER DEFAULT 0")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
