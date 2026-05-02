package storage

import (
	"errors"
	"main/internal/apperrors"
	"testing"
)

func TestMakeDB_Dbelg(t *testing.T) {
	db, err := MakeDB("dbelg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if db == nil {
		t.Error("expected non-nil DB")
	}
}

func TestMakeDB_DbelgCaseInsensitive(t *testing.T) {
	for _, name := range []string{"DBELG", "Dbelg", "DBelg"} {
		db, err := MakeDB(name)
		if err != nil {
			t.Fatalf("MakeDB(%q) unexpected error: %v", name, err)
		}
		if db == nil {
			t.Errorf("MakeDB(%q) returned nil", name)
		}
	}
}

func TestMakeDB_Unknown(t *testing.T) {
	_, err := MakeDB("mongodb")
	if !errors.Is(err, apperrors.ErrUnknownStorage) {
		t.Fatalf("expected ErrUnknownStorage, got %v", err)
	}
}

func TestMakeDB_PostgresMissingConfig(t *testing.T) {
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
		t.Setenv(key, "")
	}

	_, err := MakeDB("postgres")
	if !errors.Is(err, apperrors.ErrMissingConfig) {
		t.Fatalf("expected ErrMissingConfig, got %v", err)
	}
}
