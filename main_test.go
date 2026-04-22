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

func TestBuildHuffmanTree_MultipleLetters(t *testing.T) {
	frequencyTable := map[string]int{
		"a": 2,
		"b": 3,
		"c": 4,
	}

	tree := buildHuffmanTree(frequencyTable)

	if tree.Weight() != 9 {
		t.Fatalf("unexpected tree weight\nwant: %d\ngot:  %d", 9, tree.Weight())
	}

	if tree.IsLeaf() {
		t.Fatal("expected non-leaf root for multiple letters")
	}
}

func TestBuildHuffmanTree_SingleLetter(t *testing.T) {
	frequencyTable := map[string]int{
		"x": 7,
	}

	tree := buildHuffmanTree(frequencyTable)

	if tree.Weight() != 7 {
		t.Fatalf("unexpected tree weight\nwant: %d\ngot:  %d", 7, tree.Weight())
	}

	if !tree.IsLeaf() {
		t.Fatal("expected leaf root for single letter")
	}
}
