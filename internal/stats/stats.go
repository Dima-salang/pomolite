package stats

import (
	"fmt"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
)

type PomoStats struct {
	TotalWorkDuration      time.Duration
	TotalSessions          int
	AverageSessionDuration time.Duration
	LongestSession         time.Duration
	ShortestSession        time.Duration
	HighestSessionLabel    map[string]time.Duration
	TimeSpentPerLabel      map[string]time.Duration
	PomosPerLabel          map[string]int
}

// Compute aggregates statistics from sessions found in the storage repository in the given timeframe.
func Compute(repo storage.Repository, timeframe string) (*PomoStats, error) {
	start, end, err := resolveTimeFrame(timeframe)
	if err != nil {
		return nil, err
	}

	sessions, err := repo.GetSessionsInTimeframe(start, end)
	if err != nil {
		return nil, err
	}

	stats := &PomoStats{
		HighestSessionLabel: make(map[string]time.Duration),
		TimeSpentPerLabel:   make(map[string]time.Duration),
		PomosPerLabel:       make(map[string]int),
	}

	stats.TotalSessions = len(sessions)
	if stats.TotalSessions == 0 {
		return stats, nil
	}

	var totalDuration time.Duration
	var longest time.Duration
	var shortest time.Duration = -1
	highestLabelLongest := make(map[string]time.Duration)

	for _, s := range sessions {
		duration := s.SessionDuration
		totalDuration += duration

		if duration > longest {
			longest = duration
		}
		if shortest == -1 || duration < shortest {
			shortest = duration
		}

		stats.TimeSpentPerLabel[s.Label] += duration
		stats.PomosPerLabel[s.Label]++

		if currentMax, exists := highestLabelLongest[s.Label]; !exists || duration > currentMax {
			highestLabelLongest[s.Label] = duration
		}
	}

	stats.TotalWorkDuration = totalDuration
	stats.AverageSessionDuration = totalDuration / time.Duration(stats.TotalSessions)
	stats.LongestSession = longest
	stats.ShortestSession = shortest

	var maxLabel string
	var maxLabelDur time.Duration
	for label, maxDur := range highestLabelLongest {
		if maxDur > maxLabelDur {
			maxLabel = label
			maxLabelDur = maxDur
		}
	}
	if maxLabel != "" {
		stats.HighestSessionLabel[maxLabel] = maxLabelDur
	}

	return stats, nil
}

func resolveTimeFrame(timeframe string) (time.Time, time.Time, error) {
	now := time.Now()
	var start, end time.Time

	switch timeframe {
	case "all":
		start = time.Time{}
		end = now

	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday → make it 7
			weekday = 7
		}
		startOfWeek := now.AddDate(0, 0, -(weekday - 1)) // Monday
		start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)

	case "month":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	case "year":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)

	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid timeframe: %s", timeframe)
	}

	return start, end, nil
}
