-- name: SaveSession :exec
INSERT INTO sessions (label, start_time, end_time, session_duration)
VALUES (?, ?, ?, ?);

-- name: ListSessions :many
SELECT id, label, start_time, end_time, session_duration
FROM sessions
ORDER BY start_time DESC
LIMIT ?;

-- name: GetSessionsInTimeframe :many
SELECT id, label, start_time, end_time, session_duration
FROM sessions
WHERE start_time BETWEEN ? AND ?
ORDER BY start_time DESC;

-- name: AddTask :exec
INSERT INTO tasks (title, description, priority, status, due_date, created_at, completed_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateTask :exec
UPDATE tasks
SET title = ?, description = ?, priority = ?, status = ?, due_date = ?, completed_at = ?
WHERE id = ?;

-- name: ListTasks :many
SELECT id, title, description, priority, status, due_date, created_at, completed_at
FROM tasks
WHERE status IN (sqlc.slice('statuses'))
ORDER BY priority DESC, created_at DESC;

-- name: ListAllTasks :many
SELECT id, title, description, priority, status, due_date, created_at, completed_at
FROM tasks
ORDER BY priority DESC, created_at DESC;

-- name: GetTaskByID :one
SELECT id, title, description, priority, status, due_date, created_at, completed_at
FROM tasks
WHERE id = ?;
