CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    label TEXT NOT NULL,
    start_time INTEGER NOT NULL,
    end_time INTEGER NOT NULL,
    session_duration INTEGER NOT NULL
);

CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    priority INTEGER NOT NULL,
    status TEXT NOT NULL,
    due_date INTEGER,
    created_at INTEGER NOT NULL,
    completed_at INTEGER
);
