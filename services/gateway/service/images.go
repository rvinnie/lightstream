package service

import (
	"context"

	"github.com/rvinnie/lightstream/services/gateway/repository"
)

type ImagesService struct {
	repo repository.Images
}

func NewImagesService(repo repository.Images) *ImagesService {
	return &ImagesService{repo: repo}
}

func (s *ImagesService) GetById(ctx context.Context, id string) (string, error) {
	return s.repo.GetById(ctx, id)
}
