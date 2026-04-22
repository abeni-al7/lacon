package main

import (
	"bufio"
	"container/heap"
	"fmt"
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

func main() {
	path := filepath.Join(".", "test", "test.txt")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	frequencyTable, err := countCharacterOccurrences(file)
	if err != nil {
		log.Fatal(err)
	}

	tree := buildHuffmanTree(frequencyTable)
	fmt.Println(tree)
}
