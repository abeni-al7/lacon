package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCountCharacterOccurrences_MainUseCase(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "input.txt")
	content := "hello\nworld"

	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test input file: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("failed to open test input file: %v", err)
	}
	defer file.Close()

	got, err := countCharacterOccurrences(file)
	if err != nil {
		t.Fatalf("countCharacterOccurrences returned an error: %v", err)
	}

	want := map[string]int{
		"h": 1,
		"e": 1,
		"l": 3,
		"o": 2,
		"w": 1,
		"r": 1,
		"d": 1,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected frequency table\nwant: %#v\ngot:  %#v", want, got)
	}
}
