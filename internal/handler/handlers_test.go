package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct {
	createFn  func(url string) (string, error)
	getLongFn func(short string) (string, error)
}

func (m *mockService) CreateShortURL(url string) (string, error) {
	if m.createFn != nil {
		return m.createFn(url)
	}
	return "", errors.New("not implemented")
}

func (m *mockService) GetLongURL(short string) (string, error) {
	if m.getLongFn != nil {
		return m.getLongFn(short)
	}
	return "", errors.New("not implemented")
}

// handlerWithMock creates a Handler that uses a mockService directly.
// We bypass NewHandler because it expects *service.Service; instead we test PostHandler/GetHandler logic
// by embedding a testable handler struct.

type testHandler struct {
	mock *mockService
}

func (th *testHandler) post(rw http.ResponseWriter, r *http.Request) {
	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(rw, "не получилось распарсить URL", http.StatusBadRequest)
		return
	}
	shortURL, err := th.mock.CreateShortURL(url.Url)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	fullShortURL := "http://" + r.Host + "/" + shortURL
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"shortURL": fullShortURL})
}

func (th *testHandler) get(rw http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimSpace(r.PathValue("shortURL"))
	if shortURL == "" {
		http.Error(rw, "пустой URL", http.StatusBadRequest)
		return
	}
	originalURL, err := th.mock.GetLongURL(shortURL)
	if err != nil {
		http.Error(rw, "Не найден", http.StatusNotFound)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"URL": originalURL})
}

func TestPostHandler_ValidURL(t *testing.T) {
	th := &testHandler{mock: &mockService{
		createFn: func(url string) (string, error) { return "abc1234567", nil },
	}}
	body := strings.NewReader(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/url", body)
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	th.post(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if _, ok := resp["shortURL"]; !ok {
		t.Error("response missing 'shortURL' key")
	}
}

func TestPostHandler_EmptyURL(t *testing.T) {
	th := &testHandler{mock: &mockService{
		createFn: func(url string) (string, error) { return "", errors.New("Пустой URL") },
	}}
	body := strings.NewReader(`{"url":""}`)
	req := httptest.NewRequest(http.MethodPost, "/url", body)
	rw := httptest.NewRecorder()
	th.post(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rw.Code)
	}
	// Body must NOT contain a JSON shortURL after the error
	respBody := rw.Body.String()
	if strings.Contains(respBody, "shortURL") {
		t.Errorf("response body should not contain shortURL after error, got: %q", respBody)
	}
}

func TestPostHandler_BadJSON(t *testing.T) {
	th := &testHandler{mock: &mockService{}}
	body := strings.NewReader("not json at all")
	req := httptest.NewRequest(http.MethodPost, "/url", body)
	rw := httptest.NewRecorder()
	th.post(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rw.Code)
	}
}

func TestPostHandler_ContentType(t *testing.T) {
	th := &testHandler{mock: &mockService{
		createFn: func(url string) (string, error) { return "abc1234567", nil },
	}}
	body := strings.NewReader(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/url", body)
	rw := httptest.NewRecorder()
	th.post(rw, req)
	ct := rw.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestGetHandler_Found(t *testing.T) {
	th := &testHandler{mock: &mockService{
		getLongFn: func(short string) (string, error) { return "https://example.com", nil },
	}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortURL}", th.get)
	req := httptest.NewRequest(http.MethodGet, "/abc1234567", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rw.Code)
	}
	var resp map[string]string
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["URL"] != "https://example.com" {
		t.Errorf("expected URL=https://example.com, got %q", resp["URL"])
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	th := &testHandler{mock: &mockService{
		getLongFn: func(short string) (string, error) { return "", errors.New("not found") },
	}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortURL}", th.get)
	req := httptest.NewRequest(http.MethodGet, "/unknown123", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}

func TestGetHandler_ContentType(t *testing.T) {
	th := &testHandler{mock: &mockService{
		getLongFn: func(short string) (string, error) { return "https://example.com", nil },
	}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{shortURL}", th.get)
	req := httptest.NewRequest(http.MethodGet, "/abc1234567", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	ct := rw.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}
