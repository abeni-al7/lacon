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
	fmt.Print(frequencyTable)

	h := &MinHeap{}
	heap.Init(h)

	
}