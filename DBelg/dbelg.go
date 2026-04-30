package dbelg

import (
	"errors"
	"sync"
)

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
		return "", errors.New("Ссылка не найден")
	}
	return shortURL, nil
}

func (db *DBelg) GetLongURL(shortURL string) (string, error) {
	db.mu.RLock()
	longURL, ok := db.shortToLong[shortURL]
	db.mu.RUnlock()
	if !ok {
		return "", errors.New("Ссылка не найден")
	}
	return longURL, nil
}

func (db *DBelg) AddURL(url string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if value, ok := db.longToShort[url]; ok {
		return value, nil
	}
	shortURL := MakeShortURL(url)
	if _, ok := db.shortToLong[shortURL]; ok {
		for ok {
			shortURL = MakeShortURL(url)
			_, ok = db.shortToLong[shortURL]
		}
	}
	db.longToShort[url] = shortURL
	db.shortToLong[shortURL] = url
	return shortURL, nil
}
func MakeShortURL(url string) string {
	return ""
}
