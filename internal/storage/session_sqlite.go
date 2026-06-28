package storage

import (
	"context"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage/db"
)

func (s *SQLiteStorage) Save(label string, startTime time.Time, endTime time.Time) error {
	ctx := context.Background()
	return s.queries.SaveSession(ctx, db.SaveSessionParams{
		Label:           label,
		StartTime:       startTime.Unix(),
		EndTime:         endTime.Unix(),
		SessionDuration: int64(endTime.Sub(startTime).Seconds()),
	})
}

func (s *SQLiteStorage) List(limit int) ([]Session, error) {
	ctx := context.Background()
	lim := int64(limit)
	if lim <= 0 {
		lim = -1
	}

	dbSessions, err := s.queries.ListSessions(ctx, lim)
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(dbSessions))
	for i, ds := range dbSessions {
		sessions[i] = Session{
			ID:              int(ds.ID),
			Label:           ds.Label,
			StartTime:       time.Unix(ds.StartTime, 0),
			EndTime:         time.Unix(ds.EndTime, 0),
			SessionDuration: time.Duration(ds.SessionDuration) * time.Second,
		}
	}
	return sessions, nil
}

func (s *SQLiteStorage) GetSessionsInTimeframe(start, end time.Time) ([]Session, error) {
	ctx := context.Background()
	dbSessions, err := s.queries.GetSessionsInTimeframe(ctx, db.GetSessionsInTimeframeParams{
		FromStartTime: start.Unix(),
		ToStartTime:   end.Unix(),
	})
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, len(dbSessions))
	for i, ds := range dbSessions {
		sessions[i] = Session{
			ID:              int(ds.ID),
			Label:           ds.Label,
			StartTime:       time.Unix(ds.StartTime, 0),
			EndTime:         time.Unix(ds.EndTime, 0),
			SessionDuration: time.Duration(ds.SessionDuration) * time.Second,
		}
	}
	return sessions, nil
}
