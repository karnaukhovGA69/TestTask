package dbelg

import (
	"sync"
	"testing"
)

func TestDBelg_AddURL_ReturnsShort(t *testing.T) {
	db := NewDBelg()
	short, err := db.AddURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if short == "" {
		t.Error("expected non-empty short URL")
	}
}

func TestDBelg_AddURL_Idempotent(t *testing.T) {
	db := NewDBelg()
	url := "https://example.com"
	s1, err := db.AddURL(url)
	if err != nil {
		t.Fatalf("first AddURL error: %v", err)
	}
	s2, err := db.AddURL(url)
	if err != nil {
		t.Fatalf("second AddURL error: %v", err)
	}
	if s1 != s2 {
		t.Errorf("expected same short URL, got %q and %q", s1, s2)
	}
}

func TestDBelg_GetLongURL_Found(t *testing.T) {
	db := NewDBelg()
	original := "https://example.com"
	short, _ := db.AddURL(original)
	got, err := db.GetLongURL(short)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != original {
		t.Errorf("expected %q, got %q", original, got)
	}
}

func TestDBelg_GetLongURL_NotFound(t *testing.T) {
	db := NewDBelg()
	_, err := db.GetLongURL("nonexistent")
	if err == nil {
		t.Error("expected error for unknown short URL")
	}
}

func TestDBelg_GetShortURL_Found(t *testing.T) {
	db := NewDBelg()
	original := "https://example.com"
	short, _ := db.AddURL(original)
	got, err := db.GetShortURL(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != short {
		t.Errorf("expected %q, got %q", short, got)
	}
}

func TestDBelg_GetShortURL_NotFound(t *testing.T) {
	db := NewDBelg()
	_, err := db.GetShortURL("https://notadded.com")
	if err == nil {
		t.Error("expected error for unknown long URL")
	}
}

func TestDBelg_Concurrent(t *testing.T) {
	db := NewDBelg()
	url := "https://concurrent.example.com"
	var wg sync.WaitGroup
	results := make([]string, 50)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			s, err := db.AddURL(url)
			if err != nil {
				t.Errorf("goroutine %d error: %v", idx, err)
				return
			}
			results[idx] = s
		}(i)
	}
	wg.Wait()
	first := results[0]
	for i, r := range results {
		if r != first {
			t.Errorf("goroutine %d returned %q, expected %q", i, r, first)
		}
	}
}
