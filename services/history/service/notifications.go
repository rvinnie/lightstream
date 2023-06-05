package service

import (
	"context"
	"github.com/rvinnie/lightstream/services/history/repository"
)

type Notifications interface {
	Create(ctx context.Context, videoId string) (int, error)
}

type NotificationsService struct {
	repo repository.Notifications
}

func NewNotificationsService(repo repository.Notifications) *NotificationsService {
	return &NotificationsService{repo: repo}
}

func (s *NotificationsService) Create(ctx context.Context, videoId string) (int, error) {
	return s.repo.Create(ctx, videoId)
}
