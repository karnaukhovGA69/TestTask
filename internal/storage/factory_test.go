package storage

import (
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
	if err == nil {
		t.Error("expected error for unknown DB type")
	}
}
