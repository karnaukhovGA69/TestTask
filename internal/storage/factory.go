package storage

import (
	"errors"
	"fmt"
	"strings"

	"main/internal/storage/dbelg"
)

func MakeDB(name string) (DB, error) {
	name = strings.ToLower(strings.TrimSpace(name))

	switch name {
	case "dbelg":
		return dbelg.NewDBelg(), nil
	case "postgres":
		return nil, errors.New("TODO postgres")
	default:
		return nil, fmt.Errorf("unknown database type: %s", name)
	}
}
