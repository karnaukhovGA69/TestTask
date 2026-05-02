package dbelg

import (
	"main/internal/apperrors"
	"main/shorturl"
	"sync"
)

const maxGenerateAttempts = 69

type DBelg struct {
	longToShort map[string]string
	shortToLong map[string]string
	mu          sync.RWMutex
}

func NewDBelg() *DBelg {
	return &DBelg{longToShort: make(map[string]string), shortToLong: make(map[string]string)}
}

func (db *DBelg) GetShortURL(longURL string) (string, error) {
	db.mu.RLock()
	shortURL, ok := db.longToShort[longURL]
	db.mu.RUnlock()
	if !ok {
		return "", apperrors.ErrNotFound
	}
	return shortURL, nil
}

func (db *DBelg) GetLongURL(shortURL string) (string, error) {
	db.mu.RLock()
	longURL, ok := db.shortToLong[shortURL]
	db.mu.RUnlock()
	if !ok {
		return "", apperrors.ErrNotFound
	}
	return longURL, nil
}

func (db *DBelg) AddURL(url string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if value, ok := db.longToShort[url]; ok {
		return value, nil
	}

	for i := 0; i < maxGenerateAttempts; i++ {
		shortURL, err := shorturl.MakeShortURL()
		if err != nil {
			return "", apperrors.ErrShortURLGeneration
		}

		if _, ok := db.shortToLong[shortURL]; ok {
			continue
		}

		db.longToShort[url] = shortURL
		db.shortToLong[shortURL] = url

		return shortURL, nil
	}

	return "", apperrors.ErrUniqueShortURLGeneration
}
