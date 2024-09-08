package service

import (
	"trello-backend/lib"
	"trello-backend/models"
	"trello-backend/pkg/repository"
)

type SettingsService struct {
	repo repository.Settings
}

func NewSettingsService(repo repository.Settings) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) GetSettings(userId int) (*models.Settings, *lib.WebError) {
	return s.repo.GetSettings(userId)
}

func (s *SettingsService) CreateSettings(userId int) *lib.WebError {
	return s.repo.CreateSettings(userId)
}

func (s *SettingsService) UpdateSettings(settings models.Settings) *lib.WebError {
	return s.repo.UpdateSettings(settings)
}
