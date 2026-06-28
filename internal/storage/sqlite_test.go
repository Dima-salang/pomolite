package storage_test

import (
	"testing"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
)

func TestSQLiteSaveAndList(t *testing.T) {
	repo, err := storage.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite memory db: %v", err)
	}
	defer repo.Close()

	startTime := time.Date(2025, 9, 17, 12, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 9, 17, 12, 15, 0, 0, time.UTC)

	err = repo.Save("Test", startTime, endTime)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	sessions, err := repo.List(1)
	if err != nil {
		t.Fatalf("failed to list: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	s := sessions[0]
	if s.Label != "Test" {
		t.Errorf("expected label Test, got %s", s.Label)
	}
	if !s.StartTime.Equal(startTime) {
		t.Errorf("expected start %v, got %v", startTime, s.StartTime)
	}
	if !s.EndTime.Equal(endTime) {
		t.Errorf("expected end %v, got %v", endTime, s.EndTime)
	}
	if s.SessionDuration != 15*time.Minute {
		t.Errorf("expected duration 15m, got %v", s.SessionDuration)
	}
}
