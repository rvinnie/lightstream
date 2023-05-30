package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	imagesCollection = "images"
)

type Images interface {
	GetById(ctx context.Context, id string) (string, error)
}

type ImagesPostgres struct {
	db *pgxpool.Pool
}

func NewImagesPostgres(db *pgxpool.Pool) *ImagesPostgres {
	return &ImagesPostgres{db: db}
}

func (r *ImagesPostgres) GetById(ctx context.Context, id string) (string, error) {
	query := fmt.Sprintf("SELECT path FROM %s WHERE id = $1", imagesCollection)
	row := r.db.QueryRow(ctx, query, id)

	var path string
	err := row.Scan(&path)
	return path, err
}
