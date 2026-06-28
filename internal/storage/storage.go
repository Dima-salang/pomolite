package storage

import (
	"time"
)

// Session represents a completed focus/pomodoro work session.
type Session struct {
	ID              int
	Label           string
	StartTime       time.Time
	EndTime         time.Time
	SessionDuration time.Duration
}

// Repository handles persistence operations for focus sessions.
type Repository interface {
	Save(label string, startTime time.Time, endTime time.Time) error
	List(limit int) ([]Session, error)
	GetSessionsInTimeframe(start, end time.Time) ([]Session, error)
	Close() error
}


// task repository
type TaskStatus string

const (
	Backlog      TaskStatus = "backlog"
	Todo         TaskStatus = "todo"
	InProgress   TaskStatus = "in_progress"
	Done         TaskStatus = "done"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Priority    int
	Status      TaskStatus
	DueDate     time.Time
	CreatedAt   time.Time
	CompletedAt time.Time
}


type TaskRepository interface {
	AddTask(task Task) error
	UpdateTask(task Task) error
	ListTasks(status []TaskStatus) ([]Task, error)
	GetTaskByID(id int) (Task, error)
	DeleteTask(id int) error
	Close() error
}
