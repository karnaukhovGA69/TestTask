package service

import (
	"errors"
	"testing"
)

type mockDB struct {
	addURLFn   func(url string) (string, error)
	getLongFn  func(url string) (string, error)
	getShortFn func(url string) (string, error)
}

func (m *mockDB) AddURL(url string) (string, error)      { return m.addURLFn(url) }
func (m *mockDB) GetLongURL(url string) (string, error)  { return m.getLongFn(url) }
func (m *mockDB) GetShortURL(url string) (string, error) { return m.getShortFn(url) }

func TestService_CreateShortURL_Empty(t *testing.T) {
	svc := NewService(&mockDB{})
	_, err := svc.CreateShortURL("")
	if err == nil {
		t.Error("expected error for empty URL")
	}
}

func TestService_CreateShortURL_Whitespace(t *testing.T) {
	svc := NewService(&mockDB{})
	_, err := svc.CreateShortURL("   ")
	if err == nil {
		t.Error("expected error for whitespace-only URL")
	}
}

func TestService_CreateShortURL_Valid(t *testing.T) {
	called := false
	mock := &mockDB{
		addURLFn: func(url string) (string, error) {
			called = true
			return "abc1234567", nil
		},
	}
	svc := NewService(mock)
	result, err := svc.CreateShortURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected AddURL to be called")
	}
	if result != "abc1234567" {
		t.Errorf("expected abc1234567, got %q", result)
	}
}

func TestService_GetLongURL_Empty(t *testing.T) {
	svc := NewService(&mockDB{})
	_, err := svc.GetLongURL("")
	if err == nil {
		t.Error("expected error for empty short URL")
	}
}

func TestService_GetShortURL_Empty(t *testing.T) {
	svc := NewService(&mockDB{})
	_, err := svc.GetShortURL("")
	if err == nil {
		t.Error("expected error for empty long URL")
	}
}

func TestService_CreateShortURL_StorageError(t *testing.T) {
	mock := &mockDB{
		addURLFn: func(url string) (string, error) {
			return "", errors.New("storage error")
		},
	}
	svc := NewService(mock)
	_, err := svc.CreateShortURL("https://example.com")
	if err == nil {
		t.Error("expected error to propagate from storage")
	}
}
