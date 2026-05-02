package service

import (
	"main/internal/apperrors"
	"main/internal/storage"
	"strings"
)

type Service struct {
	database storage.DB
}

func NewService(db storage.DB) *Service {
	return &Service{database: db}
}

func (service *Service) CreateShortURL(originalURL string) (string, error) {
	if strings.TrimSpace(originalURL) == "" {
		return "", apperrors.ErrEmptyURL
	}

	return service.database.AddURL(originalURL)
}

func (service *Service) GetLongURL(shortURL string) (string, error) {
	if strings.TrimSpace(shortURL) == "" {
		return "", apperrors.ErrEmptyURL
	}
	return service.database.GetLongURL(shortURL)
}

func (service *Service) GetShortURL(longURL string) (string, error) {
	if strings.TrimSpace(longURL) == "" {
		return "", apperrors.ErrEmptyURL
	}
	return service.database.GetShortURL(longURL)
}
