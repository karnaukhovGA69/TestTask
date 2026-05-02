package postgres

import (
	"database/sql"
	"errors"
	"main/internal/apperrors"
	"main/shorturl"

	"github.com/lib/pq"
)

const maxGenerateAttempts = 69

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) GetShortURL(url string) (string, error) {
	var shortURL string
	err := p.db.QueryRow(`SELECT shortURL FROM urls WHERE longURL = $1`, url).Scan(&shortURL)
	if err == nil {
		return shortURL, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrNotFound
	}
	return "", err
}

func (p *PostgresDB) GetLongURL(url string) (string, error) {

	var longURL string
	err := p.db.QueryRow(`SELECT longURL FROM urls WHERE shortURL = $1`, url).Scan(&longURL)

	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrNotFound
	}

	if err != nil {
		return "", err
	}

	return longURL, nil

}

func (p *PostgresDB) AddURL(url string) (string, error) {
	var oldShortURL string
	err := p.db.QueryRow(`SELECT shortURL FROM urls WHERE longURL = $1`, url).Scan(&oldShortURL)
	if err == nil {
		return oldShortURL, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	for i := 0; i < maxGenerateAttempts; i++ {
		newShortURL, err := shorturl.MakeShortURL()
		if err != nil {
			return "", apperrors.ErrShortURLGeneration
		}

		_, err = p.db.Exec(`INSERT INTO urls (longURL, shortURL) VALUES ($1,$2)`, url, newShortURL)
		if err == nil {
			return newShortURL, nil
		}
		if isUniqueViolation(err) {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Constraint == "urls_longurl_key" {
				var existing string
				selectErr := p.db.QueryRow(`SELECT shortURL FROM urls WHERE longURL = $1`, url).Scan(&existing)
				if selectErr == nil {
					return existing, nil
				}
			}
			continue
		}
		return "", err

	}
	return "", apperrors.ErrUniqueShortURLGeneration
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}

	return false
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}
