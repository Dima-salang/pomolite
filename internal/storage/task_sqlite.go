package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage/db"
)

func (s *SQLiteStorage) AddTask(task Task) error {
	ctx := context.Background()

	var desc sql.NullString
	if task.Description != "" {
		desc = sql.NullString{String: task.Description, Valid: true}
	}

	var dueDate sql.NullInt64
	if !task.DueDate.IsZero() {
		dueDate = sql.NullInt64{Int64: task.DueDate.Unix(), Valid: true}
	}

	var completedAt sql.NullInt64
	if !task.CompletedAt.IsZero() {
		completedAt = sql.NullInt64{Int64: task.CompletedAt.Unix(), Valid: true}
	}

	return s.queries.AddTask(ctx, db.AddTaskParams{
		Title:       task.Title,
		Description: desc,
		Priority:    int64(task.Priority),
		Status:      string(task.Status),
		DueDate:     dueDate,
		CreatedAt:   task.CreatedAt.Unix(),
		CompletedAt: completedAt,
	})
}

func (s *SQLiteStorage) UpdateTask(task Task) error {
	ctx := context.Background()

	var desc sql.NullString
	if task.Description != "" {
		desc = sql.NullString{String: task.Description, Valid: true}
	}

	var dueDate sql.NullInt64
	if !task.DueDate.IsZero() {
		dueDate = sql.NullInt64{Int64: task.DueDate.Unix(), Valid: true}
	}

	var completedAt sql.NullInt64
	if !task.CompletedAt.IsZero() {
		completedAt = sql.NullInt64{Int64: task.CompletedAt.Unix(), Valid: true}
	}

	return s.queries.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          int64(task.ID),
		Title:       task.Title,
		Description: desc,
		Priority:    int64(task.Priority),
		Status:      string(task.Status),
		DueDate:     dueDate,
		CompletedAt: completedAt,
	})
}

func (s *SQLiteStorage) ListTasks(statuses []TaskStatus) ([]Task, error) {
	ctx := context.Background()

	var dbTasks []db.Task
	var err error

	if len(statuses) == 0 {
		dbTasks, err = s.queries.ListAllTasks(ctx)
	} else {
		strStatuses := make([]string, len(statuses))
		for i, st := range statuses {
			strStatuses[i] = string(st)
		}
		dbTasks, err = s.queries.ListTasks(ctx, strStatuses)
	}
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, len(dbTasks))
	for i, dt := range dbTasks {
		tasks[i] = toTask(dt)
	}
	return tasks, nil
}

func (s *SQLiteStorage) GetTaskByID(id int) (Task, error) {
	ctx := context.Background()
	dt, err := s.queries.GetTaskByID(ctx, int64(id))
	if err != nil {
		return Task{}, err
	}
	return toTask(dt), nil
}

func toTask(dt db.Task) Task {
	var desc string
	if dt.Description.Valid {
		desc = dt.Description.String
	}

	var dueDate time.Time
	if dt.DueDate.Valid {
		dueDate = time.Unix(dt.DueDate.Int64, 0)
	}

	var completedAt time.Time
	if dt.CompletedAt.Valid {
		completedAt = time.Unix(dt.CompletedAt.Int64, 0)
	}

	return Task{
		ID:          int(dt.ID),
		Title:       dt.Title,
		Description: desc,
		Priority:    int(dt.Priority),
		Status:      TaskStatus(dt.Status),
		DueDate:     dueDate,
		CreatedAt:   time.Unix(dt.CreatedAt, 0),
		CompletedAt: completedAt,
	}
}

func (s *SQLiteStorage) DeleteTask(id int) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

