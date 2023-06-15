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
	Create(ctx context.Context, path string) (int, error)
	GetAll(ctx context.Context) ([]string, error)
	GetById(ctx context.Context, id int) (string, error)
}

type ImagesPostgres struct {
	db *pgxpool.Pool
}

func NewImagesPostgres(db *pgxpool.Pool) *ImagesPostgres {
	return &ImagesPostgres{db: db}
}

func (r *ImagesPostgres) Create(ctx context.Context, path string) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (path) VALUES ($1) RETURNING id", imagesCollection)
	row := r.db.QueryRow(ctx, query, path)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ImagesPostgres) GetAll(ctx context.Context) ([]string, error) {
	query := fmt.Sprintf("SELECT path FROM %s", imagesCollection)
	rows, err := r.db.Query(ctx, query)
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		err = rows.Scan(&path)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func (r *ImagesPostgres) GetById(ctx context.Context, id int) (string, error) {
	query := fmt.Sprintf("SELECT path FROM %s WHERE id = $1", imagesCollection)
	row := r.db.QueryRow(ctx, query, id)

	var path string
	err := row.Scan(&path)

	return path, err
}
