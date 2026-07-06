package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

func countCharacterOccurrences(file *os.File) (map[string]int, error) {
	frequencyTable := make(map[string]int)

	reader := bufio.NewReader(file)
	for {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		frequencyTable[string(char)] += 1
	}

	return frequencyTable, nil
}

func buildHuffmanTree(frequencyTable map[string]int) HuffmanTree {
	h := &MinHeap{}
	heap.Init(h)

	for letter, freq := range frequencyTable {
		node := NewLeafNode(letter, freq)
		current := NewHuffmanTree(node)

		heap.Push(h, current)
	}

	for h.Len() > 1 {
		leftTree := heap.Pop(h).(HuffmanTree)
		rightTree := heap.Pop(h).(HuffmanTree)

		var left Node = leftTree.root
		var right Node = rightTree.root
		newNode := NewInternalNode(left.Weight()+right.Weight(), &left, &right)
		newTree := NewHuffmanTree(newNode)
		heap.Push(h, newTree)
	}

	return h.Pop().(HuffmanTree)
}

func treeSortKey(node Node) string {
	switch current := node.(type) {
	case LeafNode:
		return current.Letter()
	case InternalNode:
		leftKey := treeSortKey(current.Left())
		rightKey := treeSortKey(current.Right())
		if leftKey < rightKey {
			return leftKey
		}
		return rightKey
	default:
		return ""
	}
}

func constructPrefixCodeTable(table map[string]string, node HuffmanTree, current string) {
	if s, ok := node.root.(LeafNode); ok {
		key := s.Letter()
		if current == "" {
			current = "0" // handle trees with single node
		}
		table[key] = current
	} else if s, ok := node.root.(InternalNode); ok {
		constructPrefixCodeTable(table, NewHuffmanTree(s.Left()), current+"0")
		constructPrefixCodeTable(table, NewHuffmanTree(s.Right()), current+"1")
	}
}

func writeHeader(freqTable map[string]int, outputFile *os.File) {
	jsonData, err := json.Marshal(freqTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = outputFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	_, err = outputFile.WriteString("\n")
	if err != nil {
		log.Fatal(err)
	}
}

func writeContents(prefixCodeTable map[string]string, inputFile *os.File, outputFile *os.File) error {
	bitWriter := NewBitWriter(outputFile)

	reader := bufio.NewReader(inputFile)
	for {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := bitWriter.Writebits(prefixCodeTable[string(char)]); err != nil {
			return err
		}
	}
	if err := bitWriter.Flush(); err != nil {
		return err
	}
	return nil
}

func huffmanEncode(input *os.File, output *os.File) {
	frequencyTable, err := countCharacterOccurrences(input)
	if err != nil {
		log.Fatal(err)
	}

	tree := buildHuffmanTree(frequencyTable)

	prefixCodeTable := make(map[string]string)
	constructPrefixCodeTable(prefixCodeTable, tree, "")

	writeHeader(frequencyTable, output)

	if _, err := input.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	if err := writeContents(prefixCodeTable, input, output); err != nil {
		log.Fatal(err)
	}
}

func huffmanDecode(inputFile *os.File, outputFile *os.File) error {
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	headerBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	prefixTable, totalCount, err := ReadAndRebuildPrefixTable(headerBytes)
	if err != nil {
		return err
	}

	decodingMap := make(map[string]string)
	for char, bitString := range prefixTable {
		decodingMap[bitString] = char
	}

	bitReader := NewBitReader(reader)
	outputWriter := bufio.NewWriter(outputFile)

	var currentBits []byte
	decodedCount := 0
	for {
		if decodedCount >= totalCount {
			break
		}

		bit, err := bitReader.ReadBit()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		currentBits = append(currentBits, byte(bit))

		if originalChar, exists := decodingMap[string(currentBits)]; exists {
			_, err := outputWriter.WriteString(originalChar)
			if err != nil {
				return err
			}
			decodedCount++
			currentBits = currentBits[:0]
		}
	}
	if err := outputWriter.Flush(); err != nil {
		return err
	}
	return nil
}

func ReadAndRebuildPrefixTable(headerBytes []byte) (map[string]string, int, error) {
	var frequencyTable map[string]int
	err := json.Unmarshal(headerBytes, &frequencyTable)
	if err != nil {
		return nil, 0, err
	}

	totalCount := 0
	for _, frequency := range frequencyTable {
		totalCount += frequency
	}

	tree := buildHuffmanTree(frequencyTable)
	prefixCodeTable := make(map[string]string)
	constructPrefixCodeTable(prefixCodeTable, tree, "")

	return prefixCodeTable, totalCount, nil
}

func main() {
	inputPath := filepath.Join(".", "test", "test.txt")
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputPath := filepath.Join(".", "test", "test_output.txt")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	huffmanEncode(inputFile, outputFile)

	if _, err := outputFile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	decoded, err := os.Create(filepath.Join(".", "test", "decoded.txt"))
	if err != nil {
		log.Fatal(err)
	}
	defer decoded.Close()

	huffmanDecode(outputFile, decoded)
}
