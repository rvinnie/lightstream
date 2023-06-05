package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	notificationsCollection = "notifications"
)

type Notifications interface {
	Create(ctx context.Context, videoId string) (int, error)
}

type NotificationsPostgres struct {
	db *pgxpool.Pool
}

func NewNotificationsPostgres(db *pgxpool.Pool) *NotificationsPostgres {
	return &NotificationsPostgres{db: db}
}

func (r *NotificationsPostgres) Create(ctx context.Context, videoId string) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (videoId) VALUES ($1) RETURNING id", notificationsCollection)
	row := r.db.QueryRow(ctx, query, videoId)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
