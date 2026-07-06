package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteHeader_WritesFrequencyTableAndNewline(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.dat")

	outputFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}

	frequencyTable := map[string]int{
		"a": 2,
		"b": 1,
		"c": 3,
	}

	writeHeader(frequencyTable, outputFile)

	if err := outputFile.Close(); err != nil {
		t.Fatalf("failed to close output file: %v", err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var gotHeader map[string]int
	if err := json.Unmarshal(bytes.TrimSpace(got), &gotHeader); err != nil {
		t.Fatalf("failed to unmarshal header: %v", err)
	}

	if !reflect.DeepEqual(gotHeader, frequencyTable) {
		t.Fatalf("unexpected header frequencies\nwant: %#v\ngot:  %#v", frequencyTable, gotHeader)
	}
}

func TestWriteContents_WritesEncodedBits(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	outputPath := filepath.Join(tempDir, "output.dat")

	if err := os.WriteFile(inputPath, []byte("ab"), 0o600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}

	prefixCodeTable := map[string]string{
		"a": "0",
		"b": "1",
	}

	if err := writeContents(prefixCodeTable, inputFile, outputFile); err != nil {
		t.Fatalf("writeContents returned an error: %v", err)
	}

	if err := outputFile.Close(); err != nil {
		t.Fatalf("failed to close output file: %v", err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	want := []byte{0x40}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected encoded bytes\nwant: %v\ngot:  %v", want, got)
	}
}

func TestReadAndRebuildPrefixTable_ParsesHeader(t *testing.T) {
	header := map[string]int{
		"a": 2,
		"b": 1,
		"\n": 1,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}

	prefixTable, totalCount, err := ReadAndRebuildPrefixTable(headerBytes)
	if err != nil {
		t.Fatalf("ReadAndRebuildPrefixTable returned an error: %v", err)
	}

	if totalCount != 4 {
		t.Fatalf("unexpected count\nwant: %d\ngot:  %d", 4, totalCount)
	}

	want := map[string]string{
		"a": "1",
		"b": "01",
		"\n": "00",
	}

	if !reflect.DeepEqual(prefixTable, want) {
		t.Fatalf("unexpected prefix table\nwant: %#v\ngot:  %#v", want, prefixTable)
	}
}

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
		"x": "0",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected prefix code table\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestHuffmanRoundTrip_PreservesNewlines(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	compressedPath := filepath.Join(tempDir, "compressed.dat")
	decodedPath := filepath.Join(tempDir, "decoded.txt")

	original := []byte("hello\nworld\n")
	if err := os.WriteFile(inputPath, original, 0o600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		t.Fatalf("failed to create compressed file: %v", err)
	}
	defer compressedFile.Close()

	huffmanEncode(inputFile, compressedFile)

	if _, err := compressedFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to rewind compressed file: %v", err)
	}

	decodedFile, err := os.Create(decodedPath)
	if err != nil {
		t.Fatalf("failed to create decoded file: %v", err)
	}
	defer decodedFile.Close()

	if err := huffmanDecode(compressedFile, decodedFile); err != nil {
		t.Fatalf("huffmanDecode returned an error: %v", err)
	}

	got, err := os.ReadFile(decodedPath)
	if err != nil {
		t.Fatalf("failed to read decoded file: %v", err)
	}

	if !bytes.Equal(got, original) {
		t.Fatalf("round-trip mismatch\nwant: %q\ngot:  %q", original, got)
	}
}

func TestHuffmanRoundTrip_SingleSymbol(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.txt")
	compressedPath := filepath.Join(tempDir, "compressed.dat")
	decodedPath := filepath.Join(tempDir, "decoded.txt")

	original := []byte("aaaaaa")
	if err := os.WriteFile(inputPath, original, 0o600); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	compressedFile, err := os.Create(compressedPath)
	if err != nil {
		t.Fatalf("failed to create compressed file: %v", err)
	}
	defer compressedFile.Close()

	huffmanEncode(inputFile, compressedFile)

	if _, err := compressedFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to rewind compressed file: %v", err)
	}

	decodedFile, err := os.Create(decodedPath)
	if err != nil {
		t.Fatalf("failed to create decoded file: %v", err)
	}
	defer decodedFile.Close()

	if err := huffmanDecode(compressedFile, decodedFile); err != nil {
		t.Fatalf("huffmanDecode returned an error: %v", err)
	}

	got, err := os.ReadFile(decodedPath)
	if err != nil {
		t.Fatalf("failed to read decoded file: %v", err)
	}

	if !bytes.Equal(got, original) {
		t.Fatalf("round-trip mismatch\nwant: %q\ngot:  %q", original, got)
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
		"\n": 1,
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
