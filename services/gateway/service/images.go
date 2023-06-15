package service

import (
	"context"

	"github.com/rvinnie/lightstream/services/gateway/repository"
)

type Images interface {
	Create(ctx context.Context, name string) (int, error)
	GetAll(ctx context.Context) ([]string, error)
	GetById(ctx context.Context, id int) (string, error)
}

type ImagesService struct {
	repo repository.Images
}

func NewImagesService(repo repository.Images) *ImagesService {
	return &ImagesService{repo: repo}
}

func (s *ImagesService) Create(ctx context.Context, path string) (int, error) {
	return s.repo.Create(ctx, path)
}

func (s *ImagesService) GetAll(ctx context.Context) ([]string, error) {
	return s.repo.GetAll(ctx)
}

func (s *ImagesService) GetById(ctx context.Context, id int) (string, error) {
	return s.repo.GetById(ctx, id)
}
