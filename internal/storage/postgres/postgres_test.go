package postgres

import (
	"database/sql"
	"errors"
	"main/internal/apperrors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func newMockPostgres(t *testing.T) (*PostgresDB, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}

	return &PostgresDB{db: db}, mock, func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet SQL expectations: %v", err)
		}
		db.Close()
	}
}

func TestGetShortURL_Found(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://example.com").
		WillReturnRows(sqlmock.NewRows([]string{"shortURL"}).AddRow("abc1234567"))

	got, err := p.GetShortURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc1234567" {
		t.Fatalf("expected abc1234567, got %q", got)
	}
}

func TestGetShortURL_NotFound(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://missing.example").
		WillReturnError(sql.ErrNoRows)

	_, err := p.GetShortURL("https://missing.example")
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetLongURL_Found(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT longURL FROM urls WHERE shortURL = \$1`).
		WithArgs("abc1234567").
		WillReturnRows(sqlmock.NewRows([]string{"longURL"}).AddRow("https://example.com"))

	got, err := p.GetLongURL("abc1234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://example.com" {
		t.Fatalf("expected https://example.com, got %q", got)
	}
}

func TestAddURL_ReturnsExistingShortURL(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://example.com").
		WillReturnRows(sqlmock.NewRows([]string{"shortURL"}).AddRow("abc1234567"))

	got, err := p.AddURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc1234567" {
		t.Fatalf("expected abc1234567, got %q", got)
	}
}

func TestAddURL_InsertsNewURL(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://example.com").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(`INSERT INTO urls \(longURL, shortURL\) VALUES \(\$1,\$2\)`).
		WithArgs("https://example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	got, err := p.AddURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 10 {
		t.Fatalf("expected short code length 10, got %q", got)
	}
}

func TestAddURL_LongURLConflictReturnsExistingShortURL(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://example.com").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(`INSERT INTO urls \(longURL, shortURL\) VALUES \(\$1,\$2\)`).
		WithArgs("https://example.com", sqlmock.AnyArg()).
		WillReturnError(&pq.Error{Code: "23505", Constraint: "urls_longurl_key"})
	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://example.com").
		WillReturnRows(sqlmock.NewRows([]string{"shortURL"}).AddRow("abc1234567"))

	got, err := p.AddURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc1234567" {
		t.Fatalf("expected abc1234567, got %q", got)
	}
}

func TestAddURL_ShortURLConflictRetries(t *testing.T) {
	p, mock, done := newMockPostgres(t)
	defer done()

	mock.ExpectQuery(`SELECT shortURL FROM urls WHERE longURL = \$1`).
		WithArgs("https://example.com").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec(`INSERT INTO urls \(longURL, shortURL\) VALUES \(\$1,\$2\)`).
		WithArgs("https://example.com", sqlmock.AnyArg()).
		WillReturnError(&pq.Error{Code: "23505", Constraint: "urls_shorturl_key"})
	mock.ExpectExec(`INSERT INTO urls \(longURL, shortURL\) VALUES \(\$1,\$2\)`).
		WithArgs("https://example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	got, err := p.AddURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 10 {
		t.Fatalf("expected short code length 10, got %q", got)
	}
}
