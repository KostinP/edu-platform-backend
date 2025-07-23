package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

type VisitorEventRepository interface {
	Create(ctx context.Context, event *entity.VisitorEvent) error
}

type PostgresVisitorEventRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresVisitorEventRepo(pool *pgxpool.Pool) *PostgresVisitorEventRepo {
	return &PostgresVisitorEventRepo{pool: pool}
}

func (r *PostgresVisitorEventRepo) Create(ctx context.Context, event *entity.VisitorEvent) error {
	eventJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO visitor_events (id, visitor_id, event_type, event_data) VALUES ($1, $2, $3, $4)`,
		event.ID, event.VisitorID, event.EventType, eventJSON,
	)
	return err
}
