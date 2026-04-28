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

func constructPrefixCodeTable(table map[string]string, node HuffmanTree, current string) {
	fmt.Println(current)
	if s, ok := node.root.(LeafNode); ok {
		key := s.Letter()
		table[key] = current
	} else if s, ok := node.root.(InternalNode); ok {
		if l, ok := (s.Left()).(HuffmanTree); ok {
			constructPrefixCodeTable(table, l, current + "0")
		}
		if r, ok := (s.Right()).(HuffmanTree); ok {
			constructPrefixCodeTable(table, r, current + "1")
		}
	}
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

	prefixCodeTable := make(map[string]string)
	constructPrefixCodeTable(prefixCodeTable, tree, "")
	fmt.Println(prefixCodeTable)
}
