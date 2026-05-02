package handler

import (
	"encoding/json"
	"errors"
	"main/internal/apperrors"
	"main/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockDB struct {
	addURLFn     func(url string) (string, error)
	getLongURLFn func(url string) (string, error)
	getShortFn   func(url string) (string, error)
}

func (m *mockDB) AddURL(url string) (string, error) {
	if m.addURLFn != nil {
		return m.addURLFn(url)
	}
	return "", errors.New("unexpected AddURL call")
}

func (m *mockDB) GetLongURL(url string) (string, error) {
	if m.getLongURLFn != nil {
		return m.getLongURLFn(url)
	}
	return "", errors.New("unexpected GetLongURL call")
}

func (m *mockDB) GetShortURL(url string) (string, error) {
	if m.getShortFn != nil {
		return m.getShortFn(url)
	}
	return "", errors.New("unexpected GetShortURL call")
}

func newTestHandler(db *mockDB) *Handler {
	return NewHandler(service.NewService(db))
}

func TestPostHandler_ValidURL(t *testing.T) {
	h := newTestHandler(&mockDB{
		addURLFn: func(url string) (string, error) {
			if url != "https://example.com" {
				t.Fatalf("expected URL to be passed to storage, got %q", url)
			}
			return "abc1234567", nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "http://short.test/url", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()

	h.PostHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp["shortURL"] != "http://short.test/abc1234567" {
		t.Fatalf("unexpected shortURL response: %q", resp["shortURL"])
	}
}

func TestPostHandler_EmptyURL(t *testing.T) {
	h := newTestHandler(&mockDB{
		addURLFn: func(url string) (string, error) {
			t.Fatal("AddURL should not be called for empty URL")
			return "", nil
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(`{"url":""}`))
	rw := httptest.NewRecorder()

	h.PostHandler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
	if strings.Contains(rw.Body.String(), "shortURL") {
		t.Fatalf("response body should not contain shortURL after error, got: %q", rw.Body.String())
	}
}

func TestPostHandler_BadJSON(t *testing.T) {
	h := newTestHandler(&mockDB{})
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader("not json at all"))
	rw := httptest.NewRecorder()

	h.PostHandler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestPostHandler_InternalError(t *testing.T) {
	h := newTestHandler(&mockDB{
		addURLFn: func(url string) (string, error) {
			return "", errors.New("storage down")
		},
	})
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(`{"url":"https://example.com"}`))
	rw := httptest.NewRecorder()

	h.PostHandler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rw.Code)
	}
}

func TestGetHandler_Found(t *testing.T) {
	h := newTestHandler(&mockDB{
		getLongURLFn: func(short string) (string, error) {
			if short != "abc1234567" {
				t.Fatalf("expected short code abc1234567, got %q", short)
			}
			return "https://example.com", nil
		},
	})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortURL}", h.GetHandler)
	req := httptest.NewRequest(http.MethodGet, "/abc1234567", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp["URL"] != "https://example.com" {
		t.Fatalf("expected URL=https://example.com, got %q", resp["URL"])
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	h := newTestHandler(&mockDB{
		getLongURLFn: func(short string) (string, error) {
			return "", apperrors.ErrNotFound
		},
	})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortURL}", h.GetHandler)
	req := httptest.NewRequest(http.MethodGet, "/unknown123", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
}

func TestGetHandler_InternalError(t *testing.T) {
	h := newTestHandler(&mockDB{
		getLongURLFn: func(short string) (string, error) {
			return "", errors.New("storage down")
		},
	})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortURL}", h.GetHandler)
	req := httptest.NewRequest(http.MethodGet, "/abc1234567", nil)
	rw := httptest.NewRecorder()

	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rw.Code)
	}
}
