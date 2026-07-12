package main

import (
	"bufio"
	"container/heap"
	"encoding/json"
	"io"
	"log"
	"os"
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