package internal

import (
	"errors"
	"strings"
)

type Service struct {
	database DB
}

func NewService(db DB) *Service {
	return &Service{database: db}
}

func (service *Service) CreateShortURL(originalURL string) (string, error) {
	if strings.TrimSpace(originalURL) == "" {
		return "", errors.New("Пустой URL")
	}

	return service.database.AddURL(originalURL)
}

func (service *Service) GetLongURL(shortURL string) (string, error) {
	if strings.TrimSpace(shortURL) == "" {
		return "", errors.New("Пустой URL")
	}
	return service.database.GetLongURL(shortURL)
}

func (service *Service) GetShortURL(longURL string) (string, error) {
	if strings.TrimSpace(longURL) == "" {
		return "", errors.New("Пустой URL")
	}
	return service.database.GetShortURL(longURL)
}
