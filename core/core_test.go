package core

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestWriteHeader_WritesFrequencyTableAndNewline(t *testing.T) {
	var buf bytes.Buffer

	frequencyTable := map[string]int{
		"a": 2,
		"b": 1,
		"c": 3,
	}

	writeHeader(frequencyTable, &buf)

	got := buf.Bytes()

	var gotHeader map[string]int
	if err := json.Unmarshal(bytes.TrimSpace(got), &gotHeader); err != nil {
		t.Fatalf("failed to unmarshal header: %v", err)
	}

	if !reflect.DeepEqual(gotHeader, frequencyTable) {
		t.Fatalf("unexpected header frequencies\nwant: %#v\ngot:  %#v", frequencyTable, gotHeader)
	}
}

func TestWriteContents_WritesEncodedBits(t *testing.T) {
	var inputBuf bytes.Buffer
	inputBuf.WriteString("ab")

	var outputBuf bytes.Buffer

	prefixCodeTable := map[string]string{
		"a": "0",
		"b": "1",
	}

	if err := writeContents(prefixCodeTable, &inputBuf, &outputBuf); err != nil {
		t.Fatalf("writeContents returned an error: %v", err)
	}

	got := outputBuf.Bytes()
	want := []byte{0x40}
	if !bytes.Equal(got, want) {
		t.Fatalf("unexpected encoded bytes\nwant: %v\ngot:  %v", want, got)
	}
}

func TestReadAndRebuildPrefixTable_ParsesHeader(t *testing.T) {
	header := map[string]int{
		"a":   2,
		"b":   1,
		"\n": 1,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}

	prefixTable, totalCount, err := readAndRebuildPrefixTable(headerBytes)
	if err != nil {
		t.Fatalf("readAndRebuildPrefixTable returned an error: %v", err)
	}

	if totalCount != 4 {
		t.Fatalf("unexpected count\nwant: %d\ngot:  %d", 4, totalCount)
	}

	want := map[string]string{
		"a":  "1",
		"b":  "01",
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
	original := []byte("hello\nworld\n")

	var compressedBuf bytes.Buffer
	if err := Encode(bytes.NewReader(original), &compressedBuf); err != nil {
		t.Fatalf("Encode returned an error: %v", err)
	}

	var decodedBuf bytes.Buffer
	if err := Decode(bytes.NewReader(compressedBuf.Bytes()), &decodedBuf); err != nil {
		t.Fatalf("Decode returned an error: %v", err)
	}

	got := decodedBuf.Bytes()
	if !bytes.Equal(got, original) {
		t.Fatalf("round-trip mismatch\nwant: %q\ngot:  %q", original, got)
	}
}

func TestHuffmanRoundTrip_SingleSymbol(t *testing.T) {
	original := []byte("aaaaaa")

	var compressedBuf bytes.Buffer
	if err := Encode(bytes.NewReader(original), &compressedBuf); err != nil {
		t.Fatalf("Encode returned an error: %v", err)
	}

	var decodedBuf bytes.Buffer
	if err := Decode(bytes.NewReader(compressedBuf.Bytes()), &decodedBuf); err != nil {
		t.Fatalf("Decode returned an error: %v", err)
	}

	got := decodedBuf.Bytes()
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
	content := "hello\nworld"
	reader := bytes.NewReader([]byte(content))

	got, err := countCharacterOccurrences(reader)
	if err != nil {
		t.Fatalf("countCharacterOccurrences returned an error: %v", err)
	}

	want := map[string]int{
		"h":  1,
		"e":  1,
		"l":  3,
		"o":  2,
		"\n": 1,
		"w":  1,
		"r":  1,
		"d":  1,
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

// TestEncodeDecode_WithFiles tests the public API using actual files (integration-style)
func TestEncodeDecode_WithFiles(t *testing.T) {
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

	if err := Encode(inputFile, compressedFile); err != nil {
		t.Fatalf("Encode returned an error: %v", err)
	}
	compressedFile.Close()

	compressedFile, err = os.Open(compressedPath)
	if err != nil {
		t.Fatalf("failed to open compressed file: %v", err)
	}
	defer compressedFile.Close()

	decodedFile, err := os.Create(decodedPath)
	if err != nil {
		t.Fatalf("failed to create decoded file: %v", err)
	}
	defer decodedFile.Close()

	if err := Decode(compressedFile, decodedFile); err != nil {
		t.Fatalf("Decode returned an error: %v", err)
	}

	got, err := os.ReadFile(decodedPath)
	if err != nil {
		t.Fatalf("failed to read decoded file: %v", err)
	}

	if !bytes.Equal(got, original) {
		t.Fatalf("round-trip mismatch\nwant: %q\ngot:  %q", original, got)
	}
}