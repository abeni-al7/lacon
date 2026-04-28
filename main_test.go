package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConstructPrefixCodeTable_TwoLeafTree(t *testing.T) {
	leftLeaf := NewLeafNode("a", 1)
	rightLeaf := NewLeafNode("b", 1)

	var left Node = leftLeaf
	var right Node = rightLeaf
	root := NewInternalNode(2, &left, &right)
	tree := NewHuffmanTree(root)

	got := make(map[string]string)
	constructPrefixCodeTable(got, tree, "")

	want := map[string]string{
		"a": "0",
		"b": "1",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected prefix code table\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestConstructPrefixCodeTable_SingleLeafTree(t *testing.T) {
	tree := NewHuffmanTree(NewLeafNode("x", 7))

	got := make(map[string]string)
	constructPrefixCodeTable(got, tree, "")

	want := map[string]string{
		"x": "",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected prefix code table\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestConstructPrefixCodeTable_NestedTree(t *testing.T) {
	leftLeftLeaf := NewLeafNode("a", 1)
	leftRightLeaf := NewLeafNode("b", 1)
	rightLeaf := NewLeafNode("c", 2)

	var leftLeft Node = leftLeftLeaf
	var leftRight Node = leftRightLeaf
	leftRoot := NewInternalNode(2, &leftLeft, &leftRight)
	var leftTree Node = leftRoot
	var rightTree Node = rightLeaf
	root := NewInternalNode(4, &leftTree, &rightTree)
	tree := NewHuffmanTree(root)

	got := make(map[string]string)
	constructPrefixCodeTable(got, tree, "")

	want := map[string]string{
		"a": "00",
		"b": "01",
		"c": "1",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected prefix code table\nwant: %#v\ngot:  %#v", want, got)
	}
}

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
