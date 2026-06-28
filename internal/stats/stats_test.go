package stats_test

import (
	"testing"
	"time"

	"github.com/Dima-salang/pomolite/internal/stats"
	"github.com/Dima-salang/pomolite/internal/storage"
)

func TestComputeStats(t *testing.T) {
	repo, err := storage.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatalf("failed to create memory db: %v", err)
	}
	defer repo.Close()

	now := time.Now()
	// session 1: Work - 30 minutes
	err = repo.Save("Work", now.Add(-40*time.Minute), now.Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}
	// session 2: Read - 15 minutes
	err = repo.Save("Read", now.Add(-50*time.Minute), now.Add(-35*time.Minute))
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	res, err := stats.Compute(repo, "all")
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}

	if res.TotalSessions != 2 {
		t.Errorf("expected 2 sessions, got %d", res.TotalSessions)
	}

	if res.TotalWorkDuration != 45*time.Minute {
		t.Errorf("expected 45m duration, got %v", res.TotalWorkDuration)
	}

	if res.LongestSession != 30*time.Minute {
		t.Errorf("expected 30m longest, got %v", res.LongestSession)
	}

	if res.ShortestSession != 15*time.Minute {
		t.Errorf("expected 15m shortest, got %v", res.ShortestSession)
	}

	if res.PomosPerLabel["Work"] != 1 || res.PomosPerLabel["Read"] != 1 {
		t.Errorf("unexpected pomos per label count: %v", res.PomosPerLabel)
	}
}
