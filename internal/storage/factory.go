package storage

import (
	"fmt"
	"main/internal/apperrors"
	"main/internal/storage/dbelg"
	"main/internal/storage/postgres"
	"os"
	"strings"
)

func MakeDB(name string) (DB, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case "dbelg":
		return dbelg.NewDBelg(), nil
	case "postgres":
		for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
			if strings.TrimSpace(os.Getenv(key)) == "" {
				return nil, fmt.Errorf("%w: %s", apperrors.ErrMissingConfig, key)
			}
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSLMODE"))

		return postgres.NewPostgresDB(dsn)
	default:
		return nil, fmt.Errorf("%w: %s", apperrors.ErrUnknownStorage, name)
	}
}
