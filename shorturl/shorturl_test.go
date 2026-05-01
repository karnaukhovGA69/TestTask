package shorturl

import (
	"strings"
	"testing"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func TestMakeShortURL_Length(t *testing.T) {
	for i := 0; i < 100; i++ {
		s, err := MakeShortURL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(s) != 10 {
			t.Errorf("expected length 10, got %d: %q", len(s), s)
		}
	}
}

func TestMakeShortURL_Alphabet(t *testing.T) {
	for i := 0; i < 100; i++ {
		s, err := MakeShortURL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, c := range s {
			if !strings.ContainsRune(allowedChars, c) {
				t.Errorf("disallowed char %q in %q", c, s)
			}
		}
	}
}

func TestMakeShortURL_Uniqueness(t *testing.T) {
	seen := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		s, err := MakeShortURL()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if seen[s] {
			t.Errorf("duplicate short URL generated: %q", s)
		}
		seen[s] = true
	}
}
