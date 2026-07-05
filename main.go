package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func countCharacterOccurrences(file *os.File) (map[string]int, error) {
	frequencyTable := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range line {
			frequencyTable[string(v)] += 1
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return frequencyTable, nil
}

func buildHuffmanTree(frequencyTable map[string]int) HuffmanTree {
	h := &MinHeap{}
	heap.Init(h)

	for letter, freq := range frequencyTable {
		node := NewLeafNode(letter, freq)
		current := NewHuffmanTree(node)

		h.Push(current)
	}

	for h.Len() > 1 {
		var node1 Node = h.Pop().(Node)
		var node2 Node = h.Pop().(Node)

		newNode := NewInternalNode(node1.Weight()+node2.Weight(), &node1, &node2)
		newTree := NewHuffmanTree(newNode)
		h.Push(newTree)
	}

	return h.Pop().(HuffmanTree)
}

func constructPrefixCodeTable(table map[string]string, node HuffmanTree, current string) {
	if s, ok := node.root.(LeafNode); ok {
		key := s.Letter()
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

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range line {
			err := bitWriter.Writebits(prefixCodeTable[string(v)])
			if err != nil {
				return err
			}
		}
	}
	err := bitWriter.Flush()
		if err != nil {
			return err
		}

	if err = scanner.Err(); err != nil {
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
}
